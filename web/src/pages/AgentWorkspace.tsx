import { useEffect, useState } from 'react'
import { api, type Agent } from '../lib/api'

export default function AgentWorkspace() {
  const [agents, setAgents] = useState<Agent[]>([])
  useEffect(() => { api.agents.list().then(setAgents).catch(() => {}) }, [])

  return (
    <div>
      <h2 className="text-2xl font-bold text-gray-900 mb-6">Agent Workspaces</h2>
      <div className="grid grid-cols-3 gap-4">
        {agents.map(a => (
          <div key={a.name} className="bg-white rounded-lg border p-4">
            <div className="flex items-center gap-2">
              <div className={`w-2 h-2 rounded-full ${a.detected ? 'bg-green-500' : 'bg-gray-300'}`} />
              <span className="font-medium text-gray-900">{a.display_name}</span>
            </div>
            <div className="text-xs text-gray-500 mt-2">
              <div>Project: {a.project_dir}</div>
              <div>Global: ~/{a.global_dir}</div>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
