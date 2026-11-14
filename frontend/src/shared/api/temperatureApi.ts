export async function postConvert(
  value: number,
  from: string,
  to: string,
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

