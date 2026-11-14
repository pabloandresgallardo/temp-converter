import { useState } from 'react'
import { ConversionForm } from '@/features/conversion'
import { HistoryList } from '@/features/history'

export default function App() {
  const [refreshKey, setRefreshKey] = useState(0)

  return (
    <main className="min-h-screen bg-gray-50 p-8">
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Temp Converter</h1>
      <ConversionForm onConverted={() => setRefreshKey(k => k + 1)} />
      <HistoryList refreshKey={refreshKey} />
    </main>
  )
}
