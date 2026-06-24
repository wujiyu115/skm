import { useEffect, useState } from 'react'
import { api } from '../lib/api'

export default function Dashboard() {
  const [stats, setStats] = useState({ skills: 0, groups: 0, agents: 0, synced: 0, stale: 0 })

  useEffect(() => {
    Promise.all([api.skills.list(), api.groups.list(), api.agents.list(), api.sync.status()])
      .then(([skills, groups, agents, sync]) => {
        setStats({
          skills: skills?.length ?? 0,
          groups: groups?.length ?? 0,
          agents: agents?.filter(a => a.detected).length ?? 0,
          synced: sync.synced,
          stale: sync.stale,
        })
      })
      .catch(() => {
        // API may not be running during dev — leave defaults
      })
  }, [])

  const cards = [
    { label: 'Skills', value: stats.skills, color: 'bg-blue-500' },
    { label: 'Groups', value: stats.groups, color: 'bg-purple-500' },
    { label: 'Agents', value: stats.agents, color: 'bg-green-500' },
    { label: 'Synced', value: stats.synced, color: 'bg-teal-500' },
    { label: 'Stale', value: stats.stale, color: 'bg-amber-500' },
  ]

  return (
    <div>
      <h2 className="text-2xl font-bold text-gray-900 mb-6">Dashboard</h2>
      <div className="grid grid-cols-5 gap-4">
        {cards.map(c => (
          <div key={c.label} className="bg-white rounded-lg border p-4">
            <div className="text-xs font-medium text-gray-500 uppercase">{c.label}</div>
            <div className="text-3xl font-bold text-gray-900 mt-1">{c.value}</div>
          </div>
        ))}
      </div>
    </div>
  )
}
