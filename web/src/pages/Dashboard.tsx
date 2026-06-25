import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Puzzle, FolderOpen, Bot, RefreshCw, AlertTriangle, ArrowRight } from 'lucide-react'
import { api } from '../lib/api'

interface Stats {
  skills: number
  groups: number
  agents: number
  synced: number
  stale: number
}

const cards = [
  { key: 'skills' as const, label: 'Skills', icon: Puzzle, color: 'bg-primary-100 text-primary-600', accent: 'border-primary-500' },
  { key: 'groups' as const, label: 'Groups', icon: FolderOpen, color: 'bg-purple-100 text-purple-600', accent: 'border-purple-500' },
  { key: 'agents' as const, label: 'Agents', icon: Bot, color: 'bg-blue-100 text-blue-600', accent: 'border-blue-500' },
  { key: 'synced' as const, label: 'Synced', icon: RefreshCw, color: 'bg-teal-100 text-teal-600', accent: 'border-teal-500' },
  { key: 'stale' as const, label: 'Stale', icon: AlertTriangle, color: 'bg-amber-100 text-amber-600', accent: 'border-amber-500' },
]

export default function Dashboard() {
  const [stats, setStats] = useState<Stats>({ skills: 0, groups: 0, agents: 0, synced: 0, stale: 0 })
  const navigate = useNavigate()

  useEffect(() => {
    Promise.all([api.skills.list(), api.groups.list(), api.agents.list(), api.sync.status()])
      .then(([skills, groups, agents, sync]) => {
        setStats({
          skills: skills?.length ?? 0,
          groups: groups?.length ?? 0,
          agents: agents?.filter((a: { detected: boolean }) => a.detected).length ?? 0,
          synced: sync.synced,
          stale: sync.stale,
        })
      })
      .catch(() => {})
  }, [])

  return (
    <div>
      <div className="mb-8">
        <h2 className="text-2xl font-bold text-slate-900">Dashboard</h2>
        <p className="text-slate-500 mt-1">Overview of your skills library</p>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-5">
        {cards.map(c => {
          const Icon = c.icon
          return (
            <div key={c.key} className={`bg-white rounded-xl border-b-2 ${c.accent} border border-slate-200 p-5 shadow-sm`}>
              <div className="flex items-center gap-3">
                <div className={`w-10 h-10 rounded-lg ${c.color} flex items-center justify-center`}>
                  <Icon className="w-5 h-5" />
                </div>
                <div>
                  <div className="text-xs font-medium text-slate-500 uppercase">{c.label}</div>
                  <div className="text-2xl font-bold text-slate-900">{stats[c.key]}</div>
                </div>
              </div>
            </div>
          )
        })}
      </div>

      <div className="mt-8">
        <h3 className="text-lg font-semibold text-slate-900 mb-4">Quick Actions</h3>
        <div className="flex flex-wrap gap-3">
          <button
            onClick={() => navigate('/skills')}
            className="flex items-center gap-2 px-4 py-2.5 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 transition-colors"
          >
            <RefreshCw className="w-4 h-4" /> Sync All
          </button>
          <button
            onClick={() => navigate('/install')}
            className="flex items-center gap-2 px-4 py-2.5 bg-white border border-slate-200 text-slate-700 rounded-lg text-sm font-medium hover:bg-slate-50 transition-colors"
          >
            <ArrowRight className="w-4 h-4" /> Install Skill
          </button>
          <button
            onClick={() => navigate('/groups')}
            className="flex items-center gap-2 px-4 py-2.5 bg-white border border-slate-200 text-slate-700 rounded-lg text-sm font-medium hover:bg-slate-50 transition-colors"
          >
            <FolderOpen className="w-4 h-4" /> Create Group
          </button>
        </div>
      </div>
    </div>
  )
}
