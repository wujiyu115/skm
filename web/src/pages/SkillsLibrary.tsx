import { useEffect, useState } from 'react'
import { api, type Skill } from '../lib/api'

export default function SkillsLibrary() {
  const [skills, setSkills] = useState<Skill[]>([])
  const [source, setSource] = useState('')
  const [installing, setInstalling] = useState(false)

  const load = () => {
    api.skills.list().then(setSkills).catch(() => {})
  }
  useEffect(() => { load() }, [])

  const install = async () => {
    if (!source.trim()) return
    setInstalling(true)
    try {
      await api.skills.install(source, [], false)
      setSource('')
      load()
    } finally {
      setInstalling(false)
    }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold text-gray-900">Skills Library</h2>
        <div className="flex gap-2">
          <input
            value={source}
            onChange={e => setSource(e.target.value)}
            placeholder="owner/repo or URL..."
            className="px-3 py-2 border rounded-md text-sm w-80"
            onKeyDown={e => e.key === 'Enter' && install()}
          />
          <button
            onClick={install}
            disabled={installing}
            className="px-4 py-2 bg-blue-600 text-white rounded-md text-sm hover:bg-blue-700 disabled:opacity-50"
          >
            {installing ? 'Installing...' : 'Install'}
          </button>
        </div>
      </div>
      {skills.length === 0 ? (
        <p className="text-gray-500">No skills installed yet.</p>
      ) : (
        <div className="grid grid-cols-3 gap-4">
          {skills.map(sk => (
            <div key={sk.ID} className="bg-white rounded-lg border p-4">
              <div className="font-medium text-gray-900">{sk.Name}</div>
              <div className="text-sm text-gray-500 mt-1">{sk.Description}</div>
              <div className="flex gap-1 mt-3">
                {(sk.targets ?? []).map(t => (
                  <span key={t.agent} className="px-2 py-0.5 bg-green-100 text-green-700 rounded text-xs">
                    {t.agent}
                  </span>
                ))}
              </div>
              <div className="text-xs text-gray-400 mt-2">{sk.SourceType}: {sk.SourceRef}</div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
