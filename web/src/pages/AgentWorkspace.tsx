import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Terminal, MousePointer, Code, ArrowLeft, CheckCircle, XCircle } from 'lucide-react'
import { api, type Agent, type Skill } from '../lib/api'
import { useI18n } from '../lib/i18n'

const agentIcons: Record<string, typeof Terminal> = {
  claude: Terminal,
  cursor: MousePointer,
  codex: Code,
}

const agentColors: Record<string, string> = {
  claude: 'bg-orange-100 text-orange-600',
  cursor: 'bg-blue-100 text-blue-600',
  codex: 'bg-purple-100 text-purple-600',
}

export default function AgentWorkspace() {
  const { name } = useParams<{ name: string }>()
  return name ? <AgentDetail name={name} /> : <AgentList />
}

function AgentList() {
  const { t } = useI18n()
  const [agents, setAgents] = useState<Agent[]>([])
  const [skills, setSkills] = useState<Skill[]>([])
  const navigate = useNavigate()

  useEffect(() => {
    Promise.all([api.agents.list(), api.skills.list()])
      .then(([a, s]) => { setAgents(a ?? []); setSkills(s ?? []) })
      .catch(() => {})
  }, [])

  const skillCount = (agentName: string) =>
    skills.filter(s => s.targets?.some(t => t.agent === agentName)).length

  return (
    <div>
      <h2 className="text-2xl font-bold text-slate-900 mb-6">{t('agents.title')}</h2>
      <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
        {agents.map(a => {
          const Icon = agentIcons[a.name] ?? Terminal
          const colorClass = agentColors[a.name] ?? 'bg-slate-100 text-slate-600'
          const count = skillCount(a.name)
          return (
            <div
              key={a.name}
              onClick={() => navigate(`/agents/${a.name}`)}
              className="bg-white rounded-xl border border-slate-200 p-5 hover:shadow-md transition-shadow cursor-pointer"
            >
              <div className="flex items-center gap-3 mb-3">
                <div className={`w-10 h-10 rounded-lg ${colorClass} flex items-center justify-center`}>
                  <Icon className="w-5 h-5" />
                </div>
                <div>
                  <h3 className="font-semibold text-slate-900">{a.display_name}</h3>
                  <div className="flex items-center gap-1.5 mt-0.5">
                    {a.detected ? (
                      <>
                        <CheckCircle className="w-3.5 h-3.5 text-primary-500" />
                        <span className="text-xs text-primary-600 font-medium">{t('agents.active')}</span>
                      </>
                    ) : (
                      <>
                        <XCircle className="w-3.5 h-3.5 text-slate-400" />
                        <span className="text-xs text-slate-500">{t('agents.notDetected')}</span>
                      </>
                    )}
                  </div>
                </div>
              </div>
              <div className="text-xs text-slate-500 space-y-1">
                <div>Project: <code className="bg-slate-100 px-1.5 py-0.5 rounded">{a.project_dir}</code></div>
                <div>Global: <code className="bg-slate-100 px-1.5 py-0.5 rounded">~/{a.global_dir}</code></div>
              </div>
              <div className="mt-3 flex items-center justify-between">
                <span className="text-sm text-slate-500">{count} {t('agents.skillsSynced')}</span>
                <span className="text-xs text-primary-600 font-medium">{t('agents.view')}</span>
              </div>
            </div>
          )
        })}
      </div>
    </div>
  )
}

function AgentDetail({ name }: { name: string }) {
  const { t } = useI18n()
  const [agent, setAgent] = useState<Agent | null>(null)
  const [skills, setSkills] = useState<Skill[]>([])
  const navigate = useNavigate()

  useEffect(() => {
    Promise.all([api.agents.list(), api.skills.list()])
      .then(([agents, allSkills]) => {
        const found = agents?.find((a: Agent) => a.name === name)
        if (found) setAgent(found)
        setSkills((allSkills ?? []).filter((s: Skill) => s.targets?.some(t => t.agent === name)))
      })
      .catch(() => navigate('/agents'))
  }, [name, navigate])

  if (!agent) return null

  const Icon = agentIcons[name] ?? Terminal
  const colorClass = agentColors[name] ?? 'bg-slate-100 text-slate-600'

  return (
    <div>
      <button
        onClick={() => navigate('/agents')}
        className="flex items-center gap-1.5 text-sm text-slate-500 hover:text-slate-700 mb-4"
      >
        <ArrowLeft className="w-4 h-4" /> {t('agents.back')}
      </button>

      <div className="flex items-center gap-3 mb-6">
        <div className={`w-12 h-12 rounded-lg ${colorClass} flex items-center justify-center`}>
          <Icon className="w-6 h-6" />
        </div>
        <div>
          <h2 className="text-2xl font-bold text-slate-900">{agent.display_name}</h2>
          <div className="flex items-center gap-1.5">
            {agent.detected ? (
              <><CheckCircle className="w-3.5 h-3.5 text-primary-500" /><span className="text-sm text-primary-600">{t('agents.active')}</span></>
            ) : (
              <><XCircle className="w-3.5 h-3.5 text-slate-400" /><span className="text-sm text-slate-500">{t('agents.notDetected')}</span></>
            )}
          </div>
        </div>
      </div>

      <h3 className="text-lg font-semibold text-slate-900 mb-3">{t('agents.syncedSkills')} ({skills.length})</h3>

      {skills.length === 0 ? (
        <div className="text-center py-8 text-slate-500">
          <p>{t('agents.noSkills')}</p>
        </div>
      ) : (
        <div className="space-y-2">
          {skills.map(sk => (
            <div key={sk.ID} className="flex items-center justify-between bg-white rounded-lg border border-slate-200 px-4 py-3">
              <div>
                <span className="font-medium text-slate-900">{sk.Name}</span>
                <span className="text-sm text-slate-500 ml-2">{sk.Description}</span>
              </div>
              <span className={`px-2 py-0.5 rounded text-xs font-medium ${sk.Enabled ? 'bg-primary-100 text-primary-700' : 'bg-slate-100 text-slate-500'}`}>
                {sk.Enabled ? t('skills.enabled') : t('skills.disabled')}
              </span>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
