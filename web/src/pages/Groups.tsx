import { useEffect, useState } from 'react'
import { api, type Group } from '../lib/api'

export default function Groups() {
  const [groups, setGroups] = useState<Group[]>([])
  const [name, setName] = useState('')

  const load = () => {
    api.groups.list().then(setGroups).catch(() => {})
  }
  useEffect(() => { load() }, [])

  const create = async () => {
    if (!name.trim()) return
    await api.groups.create(name, '')
    setName('')
    load()
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold text-gray-900">Skill Groups</h2>
        <div className="flex gap-2">
          <input
            value={name}
            onChange={e => setName(e.target.value)}
            placeholder="New group name..."
            className="px-3 py-2 border rounded-md text-sm"
            onKeyDown={e => e.key === 'Enter' && create()}
          />
          <button onClick={create} className="px-4 py-2 bg-purple-600 text-white rounded-md text-sm hover:bg-purple-700">
            Create
          </button>
        </div>
      </div>
      {groups.length === 0 ? (
        <p className="text-gray-500">No groups yet.</p>
      ) : (
        <div className="grid grid-cols-3 gap-4">
          {groups.map(g => (
            <div key={g.id} className="bg-white rounded-lg border p-4">
              <div className="font-medium text-gray-900">{g.name}</div>
              <div className="text-sm text-gray-500 mt-1">{g.skill_count} skills</div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
