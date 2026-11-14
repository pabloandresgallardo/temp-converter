import { useState } from 'react'
import { postConvert } from '@/shared/api/temperatureApi'

const UNITS = ['celsius', 'fahrenheit', 'kelvin'] as const

interface Props {
  onConverted?: () => void
}

export function ConversionForm({ onConverted }: Props) {
  const [value, setValue] = useState('')
  const [from, setFrom] = useState<string>('celsius')
  const [to, setTo] = useState<string>('fahrenheit')
  const [result, setResult] = useState<number | null>(null)
  const [error, setError] = useState<string | null>(null)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError(null)
    setResult(null)
    try {
      const resp = await postConvert(Number(value), from, to)
      setResult(resp.result)
      onConverted?.()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error desconocido')
    }
  }

  return (
    <form aria-label="form" onSubmit={handleSubmit} className="space-y-4 max-w-sm">
      <div>
        <label htmlFor="value" className="block text-sm font-medium text-gray-700">
          Valor
        </label>
        <input
          id="value"
          type="number"
          value={value}
          onChange={e => setValue(e.target.value)}
          className="mt-1 block w-full rounded border-gray-300 shadow-sm"
          required
        />
      </div>

      <div className="flex gap-2">
        <div className="flex-1">
          <label htmlFor="from" className="block text-sm font-medium text-gray-700">De</label>
          <select
            id="from"
            value={from}
            onChange={e => setFrom(e.target.value)}
            className="mt-1 block w-full rounded border-gray-300 shadow-sm"
          >
            {UNITS.map(u => <option key={u} value={u}>{u}</option>)}
          </select>
        </div>
        <div className="flex-1">
          <label htmlFor="to" className="block text-sm font-medium text-gray-700">A</label>
          <select
            id="to"
            value={to}
            onChange={e => setTo(e.target.value)}
            className="mt-1 block w-full rounded border-gray-300 shadow-sm"
          >
            {UNITS.map(u => <option key={u} value={u}>{u}</option>)}
          </select>
        </div>
      </div>

      <button type="submit" className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700">
        Convertir
      </button>

      {result !== null && (
        <p className="text-lg font-semibold text-gray-900">Resultado: {result}</p>
      )}
      {error && (
        <div role="alert" className="text-red-600 text-sm">{error}</div>
      )}
    </form>
  )
}
