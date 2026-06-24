import { useEffect, useState } from 'react'
import { api } from '../lib/api'

export default function Settings() {
  const [settings, setSettings] = useState<Record<string, string>>({})
  useEffect(() => { api.settings.get().then(setSettings).catch(() => {}) }, [])

  return (
    <div>
      <h2 className="text-2xl font-bold text-gray-900 mb-6">Settings</h2>
      <div className="bg-white rounded-lg border p-6 max-w-lg">
        {Object.entries(settings).map(([key, value]) => (
          <div key={key} className="mb-4">
            <label className="text-sm font-medium text-gray-700 block">{key}</label>
            <input value={value} readOnly className="mt-1 px-3 py-2 border rounded-md text-sm w-full bg-gray-50" />
          </div>
        ))}
      </div>
    </div>
  )
}
