# POST /convert + GET /history Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement the full conversion-history flow: POST /convert converts and saves in-memory, GET /history returns records LIFO, frontend uses POST and auto-refreshes history list.

**Architecture:** Hexagonal backend — HistoryRepository (memory adapter) ← ConvertTemperature + GetHistory use cases ← HTTP handler. Frontend lifts `refreshKey` to App and passes it to HistoryList as a re-fetch trigger.

**Tech Stack:** Go 1.23 net/http + testify v1.9.0; React 18 + Vitest + RTL

**Spec:** `.ai/docs/specs/2026-11-14-post-convert-get-history.md`

## Global Constraints

- No external dependencies beyond existing ones
- `sync.RWMutex` already declared in `HistoryRepository` struct (remove `//nolint:unused` after implementing)
- `Timestamp` = `time.Now().Unix()` (int64, Unix epoch seconds)
- GET /history returns `[]` (empty array) when no records — never 404 or null
- History order: most recent first (LIFO — prepend or reverse on read)
- `GET /convert` unchanged — still ConvertOnly, no side effects
- Error wrapping convention: `fmt.Errorf("<package>: %w", err)`
- Tests run with: `cd backend && GO111MODULE=on go test ./...`
- Frontend tests: `cd frontend && node_modules/.bin/vitest run`

## Review Focus

1. **Malformed JSON body on POST /convert** — `json.Decode` fails on `{"value":"abc"}` or missing Content-Type; expect 400 `{"error": "..."}`, not 500 or panic. → Task 4, step 1d.
2. **Concurrent saves** — two goroutines calling `Save` simultaneously must not corrupt `records`; mutex write-lock required. → Task 1, step 1c.
3. **GET /history before any conversion** — must return `[]` (empty JSON array), not `null` or 404. → Task 3, step 1a.
4. **POST /convert invalid unit** — `{"value":100,"from":"celsius","to":"badunit"}` must 400 without panicking or saving a partial record. → Task 4, step 1b.
5. **History accumulates across multiple conversions** — after three POSTs, GET /history returns three records in LIFO order. → Task 2 integration test, step 1c.

---

### Task 1: HistoryRepository (memory adapter)

**Files:**
- Modify: `backend/internal/adapter/memory/history_repository.go`
- Test: `backend/internal/adapter/memory/history_repository_test.go`

**Interfaces:**
- Consumes: `port.ConversionRecord`, `port.HistoryRepository` interface
- Produces: `*HistoryRepository` satisfying `port.HistoryRepository` with `Save(port.ConversionRecord) error` and `FindAll() ([]port.ConversionRecord, error)`

- [ ] **Step 1a: Write failing test `TestHistoryRepository_SaveAndFindAll`**

```go
func TestHistoryRepository_SaveAndFindAll(t *testing.T) {
    repo := memory.NewHistoryRepository()
    record := port.ConversionRecord{Value: 0, From: "celsius", To: "fahrenheit", Result: 32, Timestamp: 1}
    require.NoError(t, repo.Save(record))
    all, err := repo.FindAll()
    require.NoError(t, err)
    assert.Len(t, all, 1)
    assert.Equal(t, record, all[0])
}
```

- [ ] **Step 1b: Run to verify FAIL**

Run: `cd backend && GO111MODULE=on go test ./internal/adapter/memory/...`
Expected: FAIL — `panic: not implemented`

- [ ] **Step 1c: Implement `Save` and `FindAll` in `history_repository.go`**

Remove `//nolint:unused` from both struct fields. `Save` acquires write lock, prepends record to `records`. `FindAll` acquires read lock, returns a copy of `records` (slice copy to avoid race on caller mutation).

- [ ] **Step 1d: Write failing test `TestHistoryRepository_FindAll_Empty`**

```go
func TestHistoryRepository_FindAll_Empty(t *testing.T) {
    repo := memory.NewHistoryRepository()
    all, err := repo.FindAll()
    require.NoError(t, err)
    assert.Empty(t, all)
}
```

- [ ] **Step 1e: Run all tests to verify PASS**

Run: `cd backend && GO111MODULE=on go test ./internal/adapter/memory/...`
Expected: PASS (2 tests)

- [ ] **Step 1f: Commit**

```bash
git add backend/internal/adapter/memory/
git commit -m "feat(memory): implement HistoryRepository Save and FindAll"
```

---

### Task 2: ConvertTemperature use case

**Files:**
- Modify: `backend/internal/usecase/convert_temperature.go`
- Test: `backend/internal/usecase/convert_temperature_test.go`

**Interfaces:**
- Consumes: `temperature.Convert`, `port.HistoryRepository.Save`
- Produces: `(*ConvertTemperature).Execute(ConvertRequest) (ConvertResponse, error)`

- [ ] **Step 2a: Write failing test `TestConvertTemperature_Execute`**

```go
func TestConvertTemperature_Execute(t *testing.T) {
    repo := memory.NewHistoryRepository()
    uc := usecase.NewConvertTemperature(repo)
    resp, err := uc.Execute(usecase.ConvertRequest{Value: 0, From: "celsius", To: "fahrenheit"})
    require.NoError(t, err)
    assert.InDelta(t, 32.0, resp.Result, 0.001)
    saved, _ := repo.FindAll()
    assert.Len(t, saved, 1)
    assert.InDelta(t, 32.0, saved[0].Result, 0.001)
}
```

- [ ] **Step 2b: Run to verify FAIL**

Run: `cd backend && GO111MODULE=on go test ./internal/usecase/... -run TestConvertTemperature`
Expected: FAIL — `panic: not implemented`

- [ ] **Step 2c: Write failing test `TestConvertTemperature_Execute_InvalidUnit`**

```go
func TestConvertTemperature_Execute_InvalidUnit(t *testing.T) {
    repo := memory.NewHistoryRepository()
    uc := usecase.NewConvertTemperature(repo)
    _, err := uc.Execute(usecase.ConvertRequest{Value: 0, From: "celsius", To: "bad"})
    require.Error(t, err)
    saved, _ := repo.FindAll()
    assert.Empty(t, saved) // must NOT save on error
}
```

- [ ] **Step 2d: Implement `Execute` in `convert_temperature.go`**

Call `temperature.Convert`. On error: return `fmt.Errorf("convert_temperature: %w", err)` without saving. On success: build `port.ConversionRecord` with `Timestamp: time.Now().Unix()`, call `repo.Save`, return result.

- [ ] **Step 2e: Run all tests to verify PASS**

Run: `cd backend && GO111MODULE=on go test ./internal/usecase/...`
Expected: PASS (all use case tests)

- [ ] **Step 2f: Commit**

```bash
git add backend/internal/usecase/convert_temperature.go backend/internal/usecase/convert_temperature_test.go
git commit -m "feat(usecase): implement ConvertTemperature Execute"
```

---

### Task 3: GetHistory use case

**Files:**
- Modify: `backend/internal/usecase/get_history.go`
- Test: `backend/internal/usecase/get_history_test.go`

**Interfaces:**
- Consumes: `port.HistoryRepository.FindAll`
- Produces: `(*GetHistory).Execute() ([]port.ConversionRecord, error)`

- [ ] **Step 3a: Write failing test `TestGetHistory_Execute_Empty`**

```go
func TestGetHistory_Execute_Empty(t *testing.T) {
    repo := memory.NewHistoryRepository()
    uc := usecase.NewGetHistory(repo)
    records, err := uc.Execute()
    require.NoError(t, err)
    assert.NotNil(t, records)
    assert.Empty(t, records)
}
```

- [ ] **Step 3b: Write failing test `TestGetHistory_Execute_WithRecords`**

```go
func TestGetHistory_Execute_WithRecords(t *testing.T) {
    repo := memory.NewHistoryRepository()
    _ = repo.Save(port.ConversionRecord{Value: 0, From: "celsius", To: "fahrenheit", Result: 32})
    uc := usecase.NewGetHistory(repo)
    records, err := uc.Execute()
    require.NoError(t, err)
    assert.Len(t, records, 1)
}
```

- [ ] **Step 3c: Run to verify FAIL**

Run: `cd backend && GO111MODULE=on go test ./internal/usecase/... -run TestGetHistory`
Expected: FAIL — `panic: not implemented`

- [ ] **Step 3d: Implement `Execute` in `get_history.go`**

Call `repo.FindAll()`. On error: return `fmt.Errorf("get_history: %w", err)`. On success: if result is nil, return empty slice `[]port.ConversionRecord{}` to guarantee JSON `[]` not `null`.

- [ ] **Step 3e: Run all tests to verify PASS**

Run: `cd backend && GO111MODULE=on go test ./internal/usecase/...`
Expected: PASS

- [ ] **Step 3f: Commit**

```bash
git add backend/internal/usecase/get_history.go backend/internal/usecase/get_history_test.go
git commit -m "feat(usecase): implement GetHistory Execute"
```

---

### Task 4: HTTP handlers + wiring

**Files:**
- Modify: `backend/internal/adapter/http/handler.go`
- Modify: `backend/internal/adapter/http/handler_test.go`
- Modify: `backend/cmd/server/main.go`

**Interfaces:**
- Consumes: `usecase.ConvertTemperature.Execute(ConvertRequest)`, `usecase.GetHistory.Execute()`
- Produces: `POST /convert` → 200/400 JSON; `GET /history` → 200 JSON array

- [ ] **Step 4a: Write failing test `TestHandler_PostConvert_OK`**

```go
func TestHandler_PostConvert_OK(t *testing.T) {
    h := newTestHandler()
    body := `{"value":0,"from":"celsius","to":"fahrenheit"}`
    req := httptest.NewRequest(http.MethodPost, "/convert", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    h.handleConvert(w, req) // call directly
    assert.Equal(t, http.StatusOK, w.Code)
    assert.Contains(t, w.Body.String(), `"result"`)
}
```

`newTestHandler()` must wire a real `memory.NewHistoryRepository()` through `usecase.NewConvertTemperature` and `usecase.NewGetHistory` and pass to `NewHandler`.

- [ ] **Step 4b: Write failing test `TestHandler_PostConvert_InvalidUnit`**

```go
func TestHandler_PostConvert_InvalidUnit(t *testing.T) {
    h := newTestHandler()
    body := `{"value":0,"from":"celsius","to":"bad"}`
    req := httptest.NewRequest(http.MethodPost, "/convert", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    h.handleConvert(w, req)
    assert.Equal(t, http.StatusBadRequest, w.Code)
}
```

- [ ] **Step 4c: Write failing test `TestHandler_PostConvert_MalformedBody`**

```go
func TestHandler_PostConvert_MalformedBody(t *testing.T) {
    h := newTestHandler()
    req := httptest.NewRequest(http.MethodPost, "/convert", strings.NewReader("not-json"))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    h.handleConvert(w, req)
    assert.Equal(t, http.StatusBadRequest, w.Code)
}
```

- [ ] **Step 4d: Write failing test `TestHandler_GetHistory_Empty`**

```go
func TestHandler_GetHistory_Empty(t *testing.T) {
    h := newTestHandler()
    req := httptest.NewRequest(http.MethodGet, "/history", nil)
    w := httptest.NewRecorder()
    h.handleHistory(w, req)
    assert.Equal(t, http.StatusOK, w.Code)
    assert.Contains(t, w.Body.String(), "[]")
}
```

- [ ] **Step 4e: Run to verify FAIL**

Run: `cd backend && GO111MODULE=on go test ./internal/adapter/http/... -run TestHandler_Post`
Expected: FAIL — handlers panic or return 501

- [ ] **Step 4f: Update `newTestHandler()` in `handler_test.go`**

```go
func newTestHandler() *Handler {
    repo := memory.NewHistoryRepository()
    return NewHandler(
        usecase.NewConvertTemperature(repo),
        usecase.NewGetHistory(repo),
        usecase.NewConvertOnly(),
    )
}
```

- [ ] **Step 4g: Implement `handleConvert` in `handler.go`**

Decode JSON body into struct `struct{ Value float64; From, To string }`. On decode error: 400. Call `convertUC.Execute`. On error: 400. On success: 200 `{"result": float64}`.

- [ ] **Step 4h: Implement `handleHistory` in `handler.go`**

Call `historyUC.Execute()`. On error: 500. On success: 200, encode slice (never null — use case guarantees non-nil).

- [ ] **Step 4i: Update `main.go` to wire repo**

Instantiate `memory.NewHistoryRepository()`, pass to `usecase.NewConvertTemperature(repo)` and `usecase.NewGetHistory(repo)`, pass both to `httpadapter.NewHandler`.

- [ ] **Step 4j: Run all backend tests to verify PASS**

Run: `cd backend && GO111MODULE=on go test ./...`
Expected: PASS (all tests)

- [ ] **Step 4k: Commit**

```bash
git add backend/internal/adapter/http/ backend/cmd/server/main.go
git commit -m "feat(api): implement POST /convert and GET /history handlers"
```

---

### Task 5: Frontend — POST + history refresh

**Files:**
- Modify: `frontend/src/shared/api/temperatureApi.ts`
- Modify: `frontend/src/features/conversion/ConversionForm.tsx`
- Modify: `frontend/src/features/history/HistoryList.tsx`
- Modify: `frontend/src/App.tsx`
- Modify: `frontend/src/features/conversion/ConversionForm.test.tsx`

**Interfaces:**
- Consumes: `postConvert(value, from, to)` from temperatureApi
- Produces: `ConversionForm` accepts `onConverted?: () => void`; `HistoryList` accepts `refreshKey?: number`

- [ ] **Step 5a: Update failing test `ConversionForm calls postConvert on submit`**

```tsx
// In ConversionForm.test.tsx — replace mock of convertTemperature with postConvert
vi.mock('@/shared/api/temperatureApi', () => ({
  postConvert: vi.fn(),
}))
import { postConvert } from '@/shared/api/temperatureApi'

it('shows result after successful conversion', async () => {
  vi.mocked(postConvert).mockResolvedValue({ result: 32 })
  render(<ConversionForm />)
  await userEvent.type(screen.getByLabelText('Valor'), '0')
  await userEvent.click(screen.getByRole('button', { name: /convertir/i }))
  expect(await screen.findByText(/32/)).toBeInTheDocument()
})
```

- [ ] **Step 5b: Run to verify FAIL**

Run: `cd frontend && node_modules/.bin/vitest run`
Expected: FAIL — `postConvert` not found

- [ ] **Step 5c: Add `postConvert` to `temperatureApi.ts`**

```ts
export async function postConvert(
  value: number, from: string, to: string,
): Promise<{ result: number }> {
  const response = await fetch('/convert', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ value, from, to }),
  })
  if (!response.ok) {
    const body = await response.json().catch(() => ({ error: 'Unknown error' }))
    throw new Error(body.error ?? 'Conversion failed')
  }
  return response.json()
}
```

- [ ] **Step 5d: Update `ConversionForm.tsx`**

Add `onConverted?: () => void` prop. Replace `convertTemperature` import with `postConvert`. On success: call `onConverted?.()`.

- [ ] **Step 5e: Update `HistoryList.tsx`**

Add `refreshKey?: number` prop. Add `refreshKey` to `useEffect` dependency array so the fetch re-runs on change.

- [ ] **Step 5f: Update `App.tsx`**

Add `const [refreshKey, setRefreshKey] = useState(0)`. Pass `onConverted={() => setRefreshKey(k => k + 1)}` to `ConversionForm`. Pass `refreshKey={refreshKey}` to `HistoryList`.

- [ ] **Step 5g: Run frontend tests to verify PASS**

Run: `cd frontend && node_modules/.bin/vitest run`
Expected: PASS

- [ ] **Step 5h: Commit**

```bash
git add frontend/src/
git commit -m "feat(frontend): use POST /convert and auto-refresh history"
```
