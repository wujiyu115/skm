import { useEffect, useState } from 'react'
import { Search, LayoutGrid, List, RefreshCw, CheckSquare, Download } from 'lucide-react'
import { api, type Skill } from '../lib/api'
import { useI18n } from '../lib/i18n'
import { toast } from '../lib/toast'
import SkillCard from '../components/SkillCard'
import TagFilter from '../components/TagFilter'

type Tab = 'all' | 'enabled' | 'available'

export default function SkillsLibrary() {
  const { t } = useI18n()
  const [skills, setSkills] = useState<Skill[]>([])
  const [search, setSearch] = useState('')
  const [activeTab, setActiveTab] = useState<Tab>('all')
  const [activeTags, setActiveTags] = useState<string[]>([])
  const [selectedSkills, setSelectedSkills] = useState<Set<string>>(new Set())
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid')
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

  const tags = [...new Set(skills.map(s => s.SourceType).filter(Boolean))]

  const toggleTag = (tag: string) => {
    setActiveTags(prev =>
      prev.includes(tag) ? prev.filter(t => t !== tag) : [...prev, tag]
    )
  }

  const filtered = skills.filter(sk => {
    if (search && !sk.Name.toLowerCase().includes(search.toLowerCase()) &&
        !sk.Description.toLowerCase().includes(search.toLowerCase())) return false
    if (activeTab === 'enabled' && !sk.Enabled) return false
    if (activeTab === 'available' && sk.Enabled) return false
    if (activeTags.length > 0 && !activeTags.includes(sk.SourceType)) return false
    return true
  })

  const toggleSelect = (id: string) => {
    setSelectedSkills(prev => {
      const next = new Set(prev)
      if (next.has(id)) next.delete(id)
      else next.add(id)
      return next
    })
  }

  const checkAll = () => {
    if (selectedSkills.size === filtered.length) {
      setSelectedSkills(new Set())
    } else {
      setSelectedSkills(new Set(filtered.map(s => s.ID)))
    }
  }

  const syncAll = async () => {
    try {
      await api.sync.trigger([])
      load()
    } catch { /* ignore */ }
  }

  const removeSkill = async (id: string) => {
    try {
      await api.skills.remove(id)
      load()
    } catch { /* ignore */ }
  }

  const handleToggleEnabled = async (id: string, enabled: boolean) => {
    try {
      await api.skills.setEnabled(id, enabled)
      toast.success(t(enabled ? 'toast.skillEnabled' : 'toast.skillDisabled'))
      load()
    } catch {
      toast.error(t('toast.error'))
    }
  }

  const tabs: { key: Tab; label: string }[] = [
    { key: 'all', label: t('skills.all') },
    { key: 'enabled', label: t('skills.enabled') },
    { key: 'available', label: t('skills.available') },
  ]

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <h2 className="text-2xl font-bold text-slate-900 dark:text-slate-100">{t('skills.title')}</h2>
          <span className="px-2.5 py-0.5 bg-primary-100 text-primary-700 rounded-full text-sm font-medium">
            {skills.length}
          </span>
        </div>
        <div className="flex gap-2">
          <input
            value={source}
            onChange={e => setSource(e.target.value)}
            placeholder={t('skills.placeholder')}
            className="px-3 py-2 border border-slate-200 dark:border-slate-700 rounded-lg text-sm w-72 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            onKeyDown={e => e.key === 'Enter' && install()}
          />
          <button
            onClick={install}
            disabled={installing}
            className="flex items-center gap-2 px-4 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 disabled:opacity-50 transition-colors"
          >
            <Download className="w-4 h-4" />
            {installing ? t('skills.installing') : t('skills.install')}
          </button>
        </div>
      </div>

      <div className="relative mb-4">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
        <input
          value={search}
          onChange={e => setSearch(e.target.value)}
          placeholder={t('skills.search')}
          className="w-full pl-10 pr-4 py-2.5 border border-slate-200 dark:border-slate-700 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
        />
      </div>

      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-1">
          {tabs.map(tab => (
            <button
              key={tab.key}
              onClick={() => setActiveTab(tab.key)}
              className={`px-4 py-1.5 rounded-lg text-sm font-medium transition-colors ${
                activeTab === tab.key
                  ? 'bg-primary-600 text-white'
                  : 'text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-700'
              }`}
            >
              {tab.label}
            </button>
          ))}
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={syncAll}
            className="flex items-center gap-1.5 px-3 py-1.5 text-sm text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-700 rounded-lg transition-colors"
          >
            <RefreshCw className="w-3.5 h-3.5" /> {t('skills.sync')}
          </button>
          <button
            onClick={checkAll}
            className="flex items-center gap-1.5 px-3 py-1.5 text-sm text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-700 rounded-lg transition-colors"
          >
            <CheckSquare className="w-3.5 h-3.5" /> {t('skills.checkAll')}
          </button>
          <div className="flex border border-slate-200 dark:border-slate-700 rounded-lg overflow-hidden ml-2">
            <button
              onClick={() => setViewMode('grid')}
              className={`p-1.5 ${viewMode === 'grid' ? 'bg-primary-600 text-white' : 'text-slate-500 dark:text-slate-400 hover:bg-slate-50 dark:hover:bg-slate-700'}`}
            >
              <LayoutGrid className="w-4 h-4" />
            </button>
            <button
              onClick={() => setViewMode('list')}
              className={`p-1.5 ${viewMode === 'list' ? 'bg-primary-600 text-white' : 'text-slate-500 dark:text-slate-400 hover:bg-slate-50 dark:hover:bg-slate-700'}`}
            >
              <List className="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>

      {tags.length > 0 && (
        <div className="mb-5">
          <TagFilter tags={tags} activeTags={activeTags} onToggle={toggleTag} />
        </div>
      )}

      {filtered.length === 0 ? (
        <div className="text-center py-12 text-slate-500 dark:text-slate-400">
          <p className="text-lg">{t('skills.noSkills')}</p>
          <p className="text-sm mt-1">{t('skills.noSkillsHint')}</p>
        </div>
      ) : (
        <div className={viewMode === 'grid'
          ? 'grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4'
          : 'space-y-3'
        }>
          {filtered.map(sk => (
            <SkillCard
              key={sk.ID}
              skill={sk}
              selected={selectedSkills.has(sk.ID)}
              onSelect={toggleSelect}
              onRemove={removeSkill}
              onSync={id => api.skills.sync(id, []).then(load)}
              onToggleEnabled={handleToggleEnabled}
            />
          ))}
        </div>
      )}
    </div>
  )
}
