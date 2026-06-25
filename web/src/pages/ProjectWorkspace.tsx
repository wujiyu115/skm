import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { FolderOpen, Plus, ArrowLeft, Trash2, Briefcase } from 'lucide-react'
import { api, type Project, type ProjectSkill, type Agent } from '../lib/api'
import { useI18n } from '../lib/i18n'
import { toast } from '../lib/toast'
import AddFromLibraryModal from '../components/AddFromLibraryModal'

export default function ProjectWorkspace() {
  const { id } = useParams<{ id: string }>()
  return id ? <ProjectDetail id={id} /> : <ProjectList />
}

function ProjectList() {
  const { t } = useI18n()
  const [projects, setProjects] = useState<Project[]>([])
  const [path, setPath] = useState('')
  const navigate = useNavigate()

  const load = () => {
    api.projects.list().then(setProjects).catch(() => {})
  }
  useEffect(() => { load() }, [])

  const addProject = async () => {
    if (!path.trim()) return
    try {
      await api.projects.create(path.trim())
      toast.success(t('toast.projectAdded'))
      setPath('')
      load()
    } catch (e) {
      toast.error(e instanceof Error ? e.message : t('toast.error'))
    }
  }

  const removeProject = async (id: string) => {
    if (!window.confirm(t('projects.confirmRemoveProject'))) return
    try {
      await api.projects.remove(id)
      toast.success(t('toast.projectRemoved'))
      load()
    } catch {
      toast.error(t('toast.error'))
    }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <h2 className="text-2xl font-bold text-slate-900 dark:text-slate-100">{t('projects.title')}</h2>
          <span className="px-2.5 py-0.5 bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-300 rounded-full text-sm font-medium">
            {projects.length}
          </span>
        </div>
      </div>

      <div className="bg-white dark:bg-slate-800 rounded-xl border border-slate-200 dark:border-slate-700 p-5 mb-6">
        <h3 className="font-semibold text-slate-900 dark:text-slate-100 mb-3">{t('projects.addProject')}</h3>
        <div className="flex gap-2">
          <input
            value={path}
            onChange={e => setPath(e.target.value)}
            placeholder={t('projects.pathPlaceholder')}
            className="flex-1 px-3 py-2 border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-slate-100 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
            onKeyDown={e => e.key === 'Enter' && addProject()}
          />
          <button
            onClick={addProject}
            className="flex items-center gap-2 px-4 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 transition-colors"
          >
            <Plus className="w-4 h-4" /> {t('projects.addProject')}
          </button>
        </div>
      </div>

      {projects.length === 0 ? (
        <div className="text-center py-12 text-slate-500 dark:text-slate-400">
          <FolderOpen className="w-12 h-12 mx-auto mb-3 text-slate-300 dark:text-slate-600" />
          <p className="text-lg">{t('projects.noProjects')}</p>
          <p className="text-sm mt-1">{t('projects.noProjectsHint')}</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
          {projects.map(p => (
            <div
              key={p.id}
              className="bg-white dark:bg-slate-800 rounded-xl border border-slate-200 dark:border-slate-700 p-5 hover:shadow-md transition-shadow cursor-pointer"
              onClick={() => navigate(`/projects/${p.id}`)}
            >
              <div className="flex items-start justify-between">
                <div className="flex items-center gap-2">
                  <Briefcase className="w-5 h-5 text-blue-500" />
                  <h3 className="font-semibold text-slate-900 dark:text-slate-100">{p.name}</h3>
                </div>
                <button
                  onClick={e => { e.stopPropagation(); removeProject(p.id) }}
                  className="text-slate-400 hover:text-red-500 transition-colors"
                >
                  <Trash2 className="w-4 h-4" />
                </button>
              </div>
              <p className="text-sm text-slate-500 dark:text-slate-400 mt-2 truncate">
                <code className="bg-slate-100 dark:bg-slate-700 px-1.5 py-0.5 rounded text-xs">{p.path}</code>
              </p>
              <div className="mt-3">
                <span className="text-xs text-slate-400 dark:text-slate-500">
                  {new Date(p.created_at).toLocaleDateString()}
                </span>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

function ProjectDetail({ id }: { id: string }) {
  const { t } = useI18n()
  const [project, setProject] = useState<Project | null>(null)
  const [skills, setSkills] = useState<ProjectSkill[]>([])
  const [showAddModal, setShowAddModal] = useState(false)
  const [agents, setAgents] = useState<Agent[]>([])
  const navigate = useNavigate()

  const load = () => {
    Promise.all([api.projects.list(), api.projects.skills(id)])
      .then(([projectsList, sk]) => {
        const found = projectsList.find((p: Project) => p.id === id)
        if (found) setProject(found)
        else navigate('/projects')
        setSkills(sk ?? [])
      })
      .catch(() => navigate('/projects'))
  }

  useEffect(() => { load() }, [id])

  useEffect(() => {
    api.agents.list().then(a => setAgents(a ?? [])).catch(() => {})
  }, [])

  const toggleSkill = async (agent: string, skillName: string, enabled: boolean) => {
    try {
      await api.projects.toggleSkill(id, agent, skillName, !enabled)
      toast.success(t('toast.skillToggled'))
      load()
    } catch {
      toast.error(t('toast.error'))
    }
  }

  const removeSkill = async (skillPath: string) => {
    if (!window.confirm(t('projects.confirmRemove'))) return
    try {
      await api.projects.removeSkill(id, skillPath)
      toast.success(t('toast.skillRemovedFromProject'))
      load()
    } catch {
      toast.error(t('toast.error'))
    }
  }

  if (!project) return null

  // Group skills by agent
  const groupedSkills: Record<string, ProjectSkill[]> = {}
  for (const sk of skills) {
    if (!groupedSkills[sk.agent_display || sk.agent]) {
      groupedSkills[sk.agent_display || sk.agent] = []
    }
    groupedSkills[sk.agent_display || sk.agent].push(sk)
  }

  return (
    <div>
      <button
        onClick={() => navigate('/projects')}
        className="flex items-center gap-1.5 text-sm text-slate-500 dark:text-slate-400 hover:text-slate-700 dark:hover:text-slate-300 mb-4"
      >
        <ArrowLeft className="w-4 h-4" /> {t('projects.back')}
      </button>

      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <Briefcase className="w-6 h-6 text-blue-500" />
          <div>
            <h2 className="text-2xl font-bold text-slate-900 dark:text-slate-100">{project.name}</h2>
            <code className="text-sm text-slate-500 dark:text-slate-400 bg-slate-100 dark:bg-slate-700 px-1.5 py-0.5 rounded">{project.path}</code>
          </div>
        </div>
        <button
          onClick={() => setShowAddModal(true)}
          className="flex items-center gap-2 px-4 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 transition-colors"
        >
          <Plus className="w-4 h-4" /> {t('projects.addFromLibrary')}
        </button>
      </div>

      <AddFromLibraryModal
        open={showAddModal}
        onClose={() => setShowAddModal(false)}
        mode="project"
        projectId={id}
        agents={agents}
        existingSkillNames={skills.map(sk => sk.skill_name)}
        onSuccess={load}
      />

      {skills.length === 0 ? (
        <div className="text-center py-12 text-slate-500 dark:text-slate-400">
          <FolderOpen className="w-12 h-12 mx-auto mb-3 text-slate-300 dark:text-slate-600" />
          <p className="text-lg">{t('projects.noSkills')}</p>
          <p className="text-sm mt-2">{t('projects.noSkillsHint')}</p>
          <button
            onClick={() => setShowAddModal(true)}
            className="mt-4 px-4 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 transition-colors"
          >
            <Plus className="w-4 h-4 inline mr-1" />{t('projects.addFromLibrary')}
          </button>
        </div>
      ) : (
        <div className="space-y-6">
          {Object.entries(groupedSkills).map(([agentName, agentSkills]) => (
            <div key={agentName}>
              <h3 className="text-lg font-semibold text-slate-900 dark:text-slate-100 mb-3">
                {agentName}
                <span className="ml-2 text-sm font-normal text-slate-500 dark:text-slate-400">
                  ({agentSkills.length})
                </span>
              </h3>
              <div className="space-y-2">
                {agentSkills.map(sk => (
                  <div key={sk.skill_path} className="flex items-center justify-between bg-white dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700 px-4 py-3">
                    <div className="min-w-0 flex-1">
                      <span className="font-medium text-slate-900 dark:text-slate-100">{sk.skill_name}</span>
                      {sk.description && (
                        <span className="text-sm text-slate-500 dark:text-slate-400 ml-2">{sk.description}</span>
                      )}
                    </div>
                    <div className="flex items-center gap-2 ml-4">
                      <button
                        onClick={() => toggleSkill(sk.agent, sk.skill_name, sk.enabled)}
                        className={`px-3 py-1 rounded text-xs font-medium transition-colors ${
                          sk.enabled
                            ? 'bg-primary-100 dark:bg-primary-900 text-primary-700 dark:text-primary-300 hover:bg-primary-200 dark:hover:bg-primary-800'
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
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
