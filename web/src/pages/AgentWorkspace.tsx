import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Terminal, MousePointer, Code, ArrowLeft, CheckCircle, XCircle, Plus, Trash2 } from 'lucide-react'
import { api, type Agent, type Skill, type ProjectSkill } from '../lib/api'
import { useI18n } from '../lib/i18n'
import { toast } from '../lib/toast'

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
  const navigate = useNavigate()

  useEffect(() => {
    api.agents.list().then(a => setAgents(a ?? [])).catch(() => {})
  }, [])

  return (
    <div>
      <h2 className="text-2xl font-bold text-slate-900 dark:text-slate-100 mb-6">{t('agents.title')}</h2>
      <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
        {[...agents].sort((a, b) => Number(b.detected) - Number(a.detected)).map(a => {
          const Icon = agentIcons[a.name] ?? Terminal
          const colorClass = agentColors[a.name] ?? 'bg-slate-100 dark:bg-slate-700 text-slate-600 dark:text-slate-400'
          return (
            <div
              key={a.name}
              onClick={() => navigate(`/agents/${a.name}`)}
              className="bg-white dark:bg-slate-800 rounded-xl border border-slate-200 dark:border-slate-700 p-5 hover:shadow-md transition-shadow cursor-pointer"
            >
              <div className="flex items-center gap-3 mb-3">
                <div className={`w-10 h-10 rounded-lg ${colorClass} flex items-center justify-center`}>
                  <Icon className="w-5 h-5" />
                </div>
                <div>
                  <h3 className="font-semibold text-slate-900 dark:text-slate-100">{a.display_name}</h3>
                  <div className="flex items-center gap-1.5 mt-0.5">
                    {a.detected ? (
                      <>
                        <CheckCircle className="w-3.5 h-3.5 text-primary-500" />
                        <span className="text-xs text-primary-600 font-medium">{t('agents.active')}</span>
                      </>
                    ) : (
                      <>
                        <XCircle className="w-3.5 h-3.5 text-slate-400" />
                        <span className="text-xs text-slate-500 dark:text-slate-400">{t('agents.notDetected')}</span>
                      </>
                    )}
                  </div>
                </div>
              </div>
              <div className="text-xs text-slate-500 dark:text-slate-400 space-y-1">
                <div>Global: <code className="bg-slate-100 dark:bg-slate-700 px-1.5 py-0.5 rounded">~/{a.global_dir}</code></div>
              </div>
              <div className="mt-3 flex items-center justify-between">
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
  const [skills, setSkills] = useState<ProjectSkill[]>([])
  const [showAdd, setShowAdd] = useState(false)
  const [librarySkills, setLibrarySkills] = useState<Skill[]>([])
  const navigate = useNavigate()

  const load = () => {
    Promise.all([api.agents.list(), api.agents.skills(name)])
      .then(([agents, sk]) => {
        const found = agents?.find((a: Agent) => a.name === name)
        if (found) setAgent(found)
        else navigate('/agents')
        setSkills(sk ?? [])
      })
      .catch(() => navigate('/agents'))
  }

  useEffect(() => { load() }, [name])

  const openAddForm = () => {
    setShowAdd(true)
    api.skills.list().then(s => setLibrarySkills(s ?? [])).catch(() => {})
  }

  const addSkill = async (skillId: string) => {
    try {
      await api.agents.addSkill(name, skillId)
      toast.success(t('toast.skillAddedToProject'))
      load()
    } catch (e) {
      toast.error(e instanceof Error ? e.message : t('toast.error'))
    }
  }

  const toggleSkill = async (skillName: string, currentEnabled: boolean) => {
    try {
      await api.agents.toggleSkill(name, skillName, !currentEnabled)
      toast.success(t('toast.skillToggled'))
      load()
    } catch {
      toast.error(t('toast.error'))
    }
  }

  const removeSkill = async (skillPath: string) => {
    if (!window.confirm(t('projects.confirmRemove'))) return
    try {
      await api.agents.removeSkill(name, skillPath)
      toast.success(t('toast.skillRemovedFromProject'))
      load()
    } catch {
      toast.error(t('toast.error'))
    }
  }

  if (!agent) return null

  const Icon = agentIcons[name] ?? Terminal
  const colorClass = agentColors[name] ?? 'bg-slate-100 dark:bg-slate-700 text-slate-600 dark:text-slate-400'

  return (
    <div>
      <button
        onClick={() => navigate('/agents')}
        className="flex items-center gap-1.5 text-sm text-slate-500 dark:text-slate-400 hover:text-slate-700 dark:hover:text-slate-300 mb-4"
      >
        <ArrowLeft className="w-4 h-4" /> {t('agents.back')}
      </button>

      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <div className={`w-12 h-12 rounded-lg ${colorClass} flex items-center justify-center`}>
            <Icon className="w-6 h-6" />
          </div>
          <div>
            <h2 className="text-2xl font-bold text-slate-900 dark:text-slate-100">{agent.display_name}</h2>
            <div className="flex items-center gap-1.5">
              {agent.detected ? (
                <><CheckCircle className="w-3.5 h-3.5 text-primary-500" /><span className="text-sm text-primary-600">{t('agents.active')}</span></>
              ) : (
                <><XCircle className="w-3.5 h-3.5 text-slate-400" /><span className="text-sm text-slate-500 dark:text-slate-400">{t('agents.notDetected')}</span></>
              )}
            </div>
            <code className="text-xs text-slate-500 dark:text-slate-400 bg-slate-100 dark:bg-slate-700 px-1.5 py-0.5 rounded">~/{agent.global_dir}</code>
          </div>
        </div>
        <button
          onClick={showAdd ? () => setShowAdd(false) : openAddForm}
          className="flex items-center gap-2 px-4 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 transition-colors"
        >
          <Plus className="w-4 h-4" /> {t('projects.addFromLibrary')}
        </button>
      </div>

      {showAdd && (
        <div className="bg-white dark:bg-slate-800 rounded-xl border border-slate-200 dark:border-slate-700 p-5 mb-6">
          <h3 className="font-semibold text-slate-900 dark:text-slate-100 mb-3">{t('projects.addFromLibrary')}</h3>
          {librarySkills.length === 0 ? (
            <p className="text-sm text-slate-500 dark:text-slate-400 py-4 text-center">{t('skills.noSkills')}</p>
          ) : (
            <div className="space-y-1 max-h-80 overflow-y-auto">
              {librarySkills.map(sk => (
                <div key={sk.ID} className="flex items-center justify-between px-3 py-2 rounded-lg hover:bg-slate-50 dark:hover:bg-slate-700 transition-colors">
                  <div>
                    <span className="font-medium text-slate-900 dark:text-slate-100 text-sm">{sk.Name}</span>
                    <span className="text-xs text-slate-500 dark:text-slate-400 ml-2">{sk.Description}</span>
                  </div>
                  <button
                    onClick={() => addSkill(sk.ID)}
                    className="px-3 py-1 bg-primary-600 text-white rounded text-xs font-medium hover:bg-primary-700 transition-colors"
                  >
                    {t('groups.add')}
                  </button>
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      <h3 className="text-lg font-semibold text-slate-900 dark:text-slate-100 mb-3">{t('agents.syncedSkills')} ({skills.length})</h3>

      {skills.length === 0 && !showAdd ? (
        <div className="text-center py-8 text-slate-500 dark:text-slate-400">
          <p>{t('agents.noSkills')}</p>
          <button
            onClick={openAddForm}
            className="mt-3 px-4 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 transition-colors"
          >
            <Plus className="w-4 h-4 inline mr-1" />{t('projects.addFromLibrary')}
          </button>
        </div>
      ) : (
        <div className="space-y-2">
          {skills.map(sk => (
            <div key={sk.skill_path} className="flex items-center justify-between bg-white dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700 px-4 py-3">
              <div className="min-w-0 flex-1">
                <span className="font-medium text-slate-900 dark:text-slate-100">{sk.skill_name}</span>
                {sk.description && (
                  <span className="text-sm text-slate-500 dark:text-slate-400 ml-2">{sk.description}</span>
                )}
              </div>
              <div className="flex items-center gap-2 ml-4">
                <button
                  onClick={() => toggleSkill(sk.skill_name, sk.enabled)}
                  className={`px-3 py-1 rounded text-xs font-medium transition-colors ${
                    sk.enabled
                      ? 'bg-green-100 dark:bg-green-900 text-green-700 dark:text-green-300 hover:bg-green-200 dark:hover:bg-green-800'
                      : 'bg-slate-100 dark:bg-slate-700 text-slate-500 dark:text-slate-400 hover:bg-slate-200 dark:hover:bg-slate-600'
                  }`}
                >
                  {sk.enabled ? t('projects.disable') : t('projects.enable')}
                </button>
                <button
                  onClick={() => removeSkill(sk.skill_path)}
                  className="text-slate-400 hover:text-red-500 transition-colors"
                >
                  <Trash2 className="w-4 h-4" />
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
