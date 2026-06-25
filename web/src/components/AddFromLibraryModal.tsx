import { useEffect, useState, useCallback } from 'react'
import { X, Search, Plus, Layers } from 'lucide-react'
import { api, type Skill, type Agent, type Group } from '../lib/api'
import { useI18n } from '../lib/i18n'
import { toast } from '../lib/toast'
import TagFilter from './TagFilter'

interface AddFromLibraryModalProps {
  open: boolean
  onClose: () => void
  mode: 'agent' | 'project'
  agentName?: string
  agentDisplayName?: string
  projectId?: string
  agents?: Agent[]
  existingSkillNames: string[]
  onSuccess: () => void
}

type ModalTab = 'skills' | 'groups'

export default function AddFromLibraryModal({
  open,
  onClose,
  mode,
  agentName,
  agentDisplayName,
  projectId,
  agents,
  existingSkillNames,
  onSuccess,
}: AddFromLibraryModalProps) {
  const { t } = useI18n()
  const [activeTab, setActiveTab] = useState<ModalTab>('skills')
  const [search, setSearch] = useState('')
  const [allSkills, setAllSkills] = useState<Skill[]>([])
  const [allGroups, setAllGroups] = useState<Group[]>([])
  const [allTags, setAllTags] = useState<string[]>([])
  const [skillTags, setSkillTags] = useState<Record<string, string[]>>({})
  const [activeTags, setActiveTags] = useState<string[]>([])
  const [adding, setAdding] = useState<string | null>(null)
  const [selectedAgents, setSelectedAgents] = useState<string[]>([])

  useEffect(() => {
    if (!open) return
    setSearch('')
    setActiveTags([])
    setActiveTab('skills')
    setAdding(null)

    if (mode === 'project' && agents?.length) {
      setSelectedAgents(agents.filter(a => a.detected).map(a => a.name))
    }

    Promise.all([
      api.skills.list(),
      api.groups.list(),
      api.tags.list(),
    ]).then(([skills, groups, tags]) => {
      setAllSkills(skills ?? [])
      setAllGroups(groups ?? [])
      setAllTags(tags ?? [])
      Promise.all(
        (skills ?? []).map(sk =>
          api.tags.getForSkill(sk.ID)
            .then(t => ({ id: sk.ID, tags: t }))
            .catch(() => ({ id: sk.ID, tags: [] as string[] }))
        )
      ).then(results => {
        const map: Record<string, string[]> = {}
        for (const r of results) map[r.id] = r.tags
        setSkillTags(map)
      })
    }).catch(() => {})
  }, [open])

  const handleKeyDown = useCallback((e: KeyboardEvent) => {
    if (e.key === 'Escape') onClose()
  }, [onClose])

  useEffect(() => {
    if (open) {
      document.addEventListener('keydown', handleKeyDown)
      return () => document.removeEventListener('keydown', handleKeyDown)
    }
  }, [open, handleKeyDown])

  if (!open) return null

  const toggleTag = (tag: string) => {
    setActiveTags(prev =>
      prev.includes(tag) ? prev.filter(t => t !== tag) : [...prev, tag]
    )
  }

  const filteredSkills = allSkills.filter(sk => {
    if (search && !sk.Name.toLowerCase().includes(search.toLowerCase()) &&
      !sk.Description.toLowerCase().includes(search.toLowerCase())) return false
    if (activeTags.length > 0) {
      const tags = skillTags[sk.ID] ?? []
      if (activeTags.includes('__untagged__')) {
        const otherTags = activeTags.filter(t => t !== '__untagged__')
        if (tags.length === 0) return true
        if (otherTags.length > 0 && otherTags.some(t => tags.includes(t))) return true
        return false
      }
      if (!activeTags.some(t => tags.includes(t))) return false
    }
    return true
  })

  const filteredGroups = allGroups.filter(g => {
    if (search && !g.name.toLowerCase().includes(search.toLowerCase()) &&
      !(g.description ?? '').toLowerCase().includes(search.toLowerCase())) return false
    return true
  })

  const isAdded = (name: string) => existingSkillNames.includes(name)

  const addSingleSkill = async (skillId: string) => {
    setAdding(skillId)
    try {
      if (mode === 'agent' && agentName) {
        await api.agents.addSkill(agentName, skillId)
      } else if (mode === 'project' && projectId) {
        if (selectedAgents.length === 0) {
          toast.error(t('projects.selectAgents'))
          setAdding(null)
          return
        }
        await api.projects.addSkill(projectId, skillId, selectedAgents)
      }
      toast.success(t('toast.skillAddedToProject'))
      onSuccess()
    } catch (e) {
      toast.error(e instanceof Error ? e.message : t('toast.error'))
    } finally {
      setAdding(null)
    }
  }

  const addGroup = async (groupId: string) => {
    setAdding(groupId)
    try {
      const data = await api.groups.get(groupId)
      const skills = data?.skills ?? []
      for (const sk of skills) {
        if (isAdded(sk.Name)) continue
        if (mode === 'agent' && agentName) {
          await api.agents.addSkill(agentName, sk.ID)
        } else if (mode === 'project' && projectId) {
          if (selectedAgents.length === 0) break
          await api.projects.addSkill(projectId, sk.ID, selectedAgents)
        }
      }
      toast.success(t('toast.skillAddedToProject'))
      onSuccess()
    } catch (e) {
      toast.error(e instanceof Error ? e.message : t('toast.error'))
    } finally {
      setAdding(null)
    }
  }

  const targetLabel = mode === 'agent' ? agentDisplayName : undefined

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50" onClick={onClose}>
      <div
        className="w-full max-w-2xl bg-white dark:bg-slate-800 rounded-xl shadow-2xl border border-slate-200 dark:border-slate-700 overflow-hidden max-h-[80vh] flex flex-col"
        onClick={e => e.stopPropagation()}
      >
        {/* Header */}
        <div className="flex items-center justify-between px-6 py-4 border-b border-slate-200 dark:border-slate-700">
          <h2 className="text-lg font-bold text-slate-900 dark:text-slate-100">{t('modal.addFromLibrary')}</h2>
          <button onClick={onClose} className="p-1 text-slate-400 hover:text-slate-600 dark:hover:text-slate-300 rounded transition-colors">
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Target + Agent selection for project mode */}
        <div className="px-6 pt-4 space-y-3">
          {targetLabel && (
            <div className="flex items-center gap-2">
              <span className="text-sm text-slate-500 dark:text-slate-400">{t('modal.target')}:</span>
              <span className="px-2.5 py-1 bg-primary-100 dark:bg-primary-900 text-primary-700 dark:text-primary-300 rounded-lg text-sm font-medium">{targetLabel}</span>
            </div>
          )}
          {mode === 'project' && agents && agents.length > 0 && (
            <div className="flex items-center gap-2 flex-wrap">
              <span className="text-sm text-slate-500 dark:text-slate-400">{t('projects.selectAgents')}:</span>
              {agents.filter(a => a.detected).map(a => (
                <label key={a.name} className="flex items-center gap-1.5 text-sm text-slate-700 dark:text-slate-300">
                  <input
                    type="checkbox"
                    checked={selectedAgents.includes(a.name)}
                    onChange={() => setSelectedAgents(prev =>
                      prev.includes(a.name) ? prev.filter(n => n !== a.name) : [...prev, a.name]
                    )}
                    className="rounded border-slate-300 dark:border-slate-600"
                  />
                  {a.display_name}
                </label>
              ))}
            </div>
          )}
        </div>

        {/* Search */}
        <div className="px-6 pt-3">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
            <input
              value={search}
              onChange={e => setSearch(e.target.value)}
              placeholder={t('modal.searchPlaceholder')}
              className="w-full pl-10 pr-4 py-2.5 border border-slate-200 dark:border-slate-700 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent bg-white dark:bg-slate-900"
              autoFocus
            />
          </div>
        </div>

        {/* Tabs */}
        <div className="flex items-center gap-1 px-6 pt-3">
          <button
            onClick={() => setActiveTab('skills')}
            className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
              activeTab === 'skills'
                ? 'bg-primary-600 text-white'
                : 'text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-700'
            }`}
          >
            <Plus className="w-3.5 h-3.5" /> {t('modal.tabSkills')}
          </button>
          <button
            onClick={() => setActiveTab('groups')}
            className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
              activeTab === 'groups'
                ? 'bg-primary-600 text-white'
                : 'text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-700'
            }`}
          >
            <Layers className="w-3.5 h-3.5" /> {t('modal.tabGroups')}
          </button>
        </div>

        {/* Tag filter (skills tab only) */}
        {activeTab === 'skills' && allTags.length > 0 && (
          <div className="px-6 pt-3">
            <TagFilter tags={allTags} activeTags={activeTags} onToggle={toggleTag} />
          </div>
        )}

        {/* Scrollable list */}
        <div className="flex-1 overflow-y-auto px-6 py-3 min-h-0">
          {activeTab === 'skills' && (
            filteredSkills.length === 0 ? (
              <p className="text-center text-sm text-slate-500 dark:text-slate-400 py-8">{t('modal.noResults')}</p>
            ) : (
              <div className="space-y-1">
                {filteredSkills.map(sk => {
                  const added = isAdded(sk.Name)
                  const isLoading = adding === sk.ID
                  return (
                    <div key={sk.ID} className="flex items-center justify-between px-3 py-2.5 rounded-lg hover:bg-slate-50 dark:hover:bg-slate-700/50 transition-colors">
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2">
                          <span className="font-medium text-slate-900 dark:text-slate-100 text-sm">{sk.Name}</span>
                          {sk.SourceType && (
                            <span className="text-xs px-1.5 py-0.5 bg-slate-100 dark:bg-slate-700 text-slate-500 dark:text-slate-400 rounded">
                              {sk.SourceType}
                            </span>
                          )}
                        </div>
                        {sk.Description && (
                          <p className="text-xs text-slate-500 dark:text-slate-400 mt-0.5 truncate">{sk.Description}</p>
                        )}
                      </div>
                      {added ? (
                        <span className="px-2.5 py-1 bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400 rounded text-xs font-medium">
                          {t('modal.alreadyAdded')}
                        </span>
                      ) : (
                        <button
                          onClick={() => addSingleSkill(sk.ID)}
                          disabled={isLoading}
                          className="px-3 py-1.5 bg-primary-600 text-white rounded-lg text-xs font-medium hover:bg-primary-700 disabled:opacity-50 transition-colors"
                        >
                          {isLoading ? t('modal.adding') : t('groups.add')}
                        </button>
                      )}
                    </div>
                  )
                })}
              </div>
            )
          )}

          {activeTab === 'groups' && (
            filteredGroups.length === 0 ? (
              <p className="text-center text-sm text-slate-500 dark:text-slate-400 py-8">{t('groups.noGroups')}</p>
            ) : (
              <div className="space-y-1">
                {filteredGroups.map(g => {
                  const isLoading = adding === g.id
                  return (
                    <div key={g.id} className="flex items-center justify-between px-3 py-2.5 rounded-lg hover:bg-slate-50 dark:hover:bg-slate-700/50 transition-colors">
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2">
                          <span className="font-medium text-slate-900 dark:text-slate-100 text-sm">{g.name}</span>
                          <span className="text-xs px-1.5 py-0.5 bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 rounded">
                            {g.skill_count ?? 0} {t('groups.skills')}
                          </span>
                        </div>
                        {g.description && (
                          <p className="text-xs text-slate-500 dark:text-slate-400 mt-0.5 truncate">{g.description}</p>
                        )}
                      </div>
                      <button
                        onClick={() => addGroup(g.id)}
                        disabled={isLoading}
                        className="px-3 py-1.5 bg-primary-600 text-white rounded-lg text-xs font-medium hover:bg-primary-700 disabled:opacity-50 transition-colors"
                      >
                        {isLoading ? t('modal.adding') : t('groups.add')}
                      </button>
                    </div>
                  )
                })}
              </div>
            )
          )}
        </div>
      </div>
    </div>
  )
}
