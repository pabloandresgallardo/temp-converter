import { render, screen, waitFor } from '@testing-library/react'
import { vi } from 'vitest'
import { HistoryList } from './HistoryList'

global.fetch = vi.fn()

afterEach(() => vi.resetAllMocks())

it('shows unavailable when fetch fails', async () => {
  vi.mocked(fetch).mockRejectedValue(new Error('network error'))
  render(<HistoryList />)
  expect(await screen.findByText(/no disponible/i)).toBeInTheDocument()
})

it('shows unavailable when response is not ok', async () => {
  vi.mocked(fetch).mockResolvedValue({ ok: false } as Response)
  render(<HistoryList />)
  expect(await screen.findByText(/no disponible/i)).toBeInTheDocument()
})

it('shows empty message when history is empty', async () => {
  vi.mocked(fetch).mockResolvedValue({
    ok: true,
    json: async () => [],
  } as Response)
  render(<HistoryList />)
  expect(await screen.findByText(/sin conversiones/i)).toBeInTheDocument()
})

it('shows records when history has entries', async () => {
  vi.mocked(fetch).mockResolvedValue({
    ok: true,
    json: async () => [
      { value: 0, from: 'celsius', to: 'fahrenheit', result: 32, timestamp: 1 },
    ],
  } as Response)
  render(<HistoryList />)
  expect(await screen.findByText(/fahrenheit/)).toBeInTheDocument()
  expect(screen.getByText(/32/)).toBeInTheDocument()
})

it('re-fetches when refreshKey changes', async () => {
  vi.mocked(fetch).mockResolvedValue({
    ok: true,
    json: async () => [],
  } as Response)
  const { rerender } = render(<HistoryList refreshKey={0} />)
  await waitFor(() => expect(fetch).toHaveBeenCalledTimes(1))
  rerender(<HistoryList refreshKey={1} />)
  await waitFor(() => expect(fetch).toHaveBeenCalledTimes(2))
})
