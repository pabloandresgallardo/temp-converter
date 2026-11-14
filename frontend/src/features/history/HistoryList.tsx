import { useEffect, useState } from 'react'

interface ConversionRecord {
  value: number
  from: string
  to: string
  result: number
  timestamp: number
}

interface Props {
  refreshKey?: number
}

export function HistoryList({ refreshKey }: Props) {
  const [records, setRecords] = useState<ConversionRecord[]>([])
  const [status, setStatus] = useState<'loading' | 'unavailable' | 'ok'>('loading')

  useEffect(() => {
    fetch('/history')
      .then(r => {
        if (!r.ok) {
          setStatus('unavailable')
          return
        }
        return r.json()
      })
      .then(data => {
        if (data) {
          setRecords(data)
          setStatus('ok')
        }
      })
      .catch(() => setStatus('unavailable'))
  }, [refreshKey])

  if (status === 'loading') return null

  if (status === 'unavailable') {
    return (
      <p className="text-sm text-gray-400 mt-6">Historial no disponible aún.</p>
    )
  }

  if (records.length === 0) {
    return (
      <p className="text-sm text-gray-400 mt-6">Sin conversiones recientes.</p>
    )
  }

  return (
    <section className="mt-6 max-w-sm">
      <h2 className="text-sm font-semibold text-gray-700 mb-2">Historial</h2>
      <ul className="space-y-1">
        {records.map((r, i) => (
          <li key={`${r.timestamp ?? i}-${r.value}`} className="text-sm text-gray-600">
            {r.value} {r.from} → {r.result} {r.to}
          </li>
        ))}
      </ul>
    </section>
  )
}
