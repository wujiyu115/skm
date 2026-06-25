import { useEffect, useState, useCallback, useRef } from 'react'
import { useNavigate } from 'react-router-dom'
import {
  Search, BookOpen, FolderOpen, LayoutDashboard, Download,
  Settings, ClipboardList, FileText,
} from 'lucide-react'
import { api, type Skill, type Group } from '../lib/api'
import { useI18n } from '../lib/i18n'

interface SearchResult {
  id: string
  name: string
  description?: string
  category: 'skill' | 'group' | 'page'
  path: string
  icon: typeof Search
}

const PAGE_TARGETS: Omit<SearchResult, 'id'>[] = [
  { name: 'Dashboard', path: '/', category: 'page', icon: LayoutDashboard },
  { name: 'Library', path: '/skills', category: 'page', icon: BookOpen },
  { name: 'Install', path: '/install', category: 'page', icon: Download },
  { name: 'Settings', path: '/settings', category: 'page', icon: Settings },
  { name: 'Audit Log', path: '/audit', category: 'page', icon: ClipboardList },
]

export default function CommandPalette() {
  const [open, setOpen] = useState(false)
  const [query, setQuery] = useState('')
  const [skills, setSkills] = useState<Skill[]>([])
  const [groups, setGroups] = useState<Group[]>([])
  const [activeIndex, setActiveIndex] = useState(0)
  const inputRef = useRef<HTMLInputElement>(null)
  const listRef = useRef<HTMLDivElement>(null)
  const navigate = useNavigate()
  const { t } = useI18n()

  // Global keyboard shortcut
  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
        e.preventDefault()
        setOpen(prev => !prev)
      }
    }
    document.addEventListener('keydown', handler)
    return () => document.removeEventListener('keydown', handler)
  }, [])

  // Fetch data when opening
  useEffect(() => {
    if (!open) return
    setQuery('')
    setActiveIndex(0)
    Promise.all([api.skills.list(), api.groups.list()])
      .then(([s, g]) => {
        setSkills(s ?? [])
        setGroups(g ?? [])
      })
      .catch(() => {})
    // Focus input after render
    requestAnimationFrame(() => inputRef.current?.focus())
  }, [open])

  // Build results
  const results = useCallback((): SearchResult[] => {
    const q = query.toLowerCase().trim()

    const skillResults: SearchResult[] = skills
      .filter(s => {
        if (!q) return true
        return s.Name.toLowerCase().includes(q) ||
          (s.Description && s.Description.toLowerCase().includes(q))
      })
      .map(s => ({
        id: `skill-${s.ID}`,
        name: s.Name,
        description: s.Description,
        category: 'skill' as const,
        path: '/skills',
        icon: FileText,
      }))

    const groupResults: SearchResult[] = groups
      .filter(g => {
        if (!q) return true
        return g.name.toLowerCase().includes(q)
      })
      .map(g => ({
        id: `group-${g.id}`,
        name: g.name,
        description: g.description,
        category: 'group' as const,
        path: `/groups/${g.id}`,
        icon: FolderOpen,
      }))

    const pageResults: SearchResult[] = PAGE_TARGETS
      .filter(p => {
        if (!q) return true
        return p.name.toLowerCase().includes(q)
      })
      .map((p, i) => ({
        ...p,
        id: `page-${i}`,
      }))

    return [...skillResults, ...groupResults, ...pageResults]
  }, [query, skills, groups])

  const allResults = results()

  // Keyboard navigation inside palette
  useEffect(() => {
    if (!open) return
    const handler = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        e.preventDefault()
        setOpen(false)
      } else if (e.key === 'ArrowDown') {
        e.preventDefault()
        setActiveIndex(prev => Math.min(prev + 1, allResults.length - 1))
      } else if (e.key === 'ArrowUp') {
        e.preventDefault()
        setActiveIndex(prev => Math.max(prev - 1, 0))
      } else if (e.key === 'Enter') {
        e.preventDefault()
        const item = allResults[activeIndex]
        if (item) {
          navigate(item.path)
          setOpen(false)
        }
      }
    }
    document.addEventListener('keydown', handler)
    return () => document.removeEventListener('keydown', handler)
  }, [open, activeIndex, allResults, navigate])

  // Scroll active item into view
  useEffect(() => {
    if (!open || !listRef.current) return
    const active = listRef.current.querySelector('[data-active="true"]')
    if (active && typeof active.scrollIntoView === 'function') {
      active.scrollIntoView({ block: 'nearest' })
    }
  }, [activeIndex, open])

  // Reset active index when query changes
  useEffect(() => {
    setActiveIndex(0)
  }, [query])

  if (!open) return null

  const categoryLabel = (cat: SearchResult['category']) => {
    switch (cat) {
      case 'skill': return t('search.skills')
      case 'group': return t('search.groups')
      case 'page': return t('search.pages')
    }
  }

  const categoryBadgeColor = (cat: SearchResult['category']) => {
    switch (cat) {
      case 'skill': return 'bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-300'
      case 'group': return 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300'
      case 'page': return 'bg-purple-100 text-purple-700 dark:bg-purple-900 dark:text-purple-300'
    }
  }

  return (
    <div
      className="fixed inset-0 z-50 flex items-start justify-center pt-[20vh] bg-black/50"
      data-testid="command-palette-overlay"
      onClick={() => setOpen(false)}
    >
      <div
        className="w-full max-w-lg bg-white dark:bg-slate-800 rounded-xl shadow-2xl border border-slate-200 dark:border-slate-700 overflow-hidden"
        onClick={e => e.stopPropagation()}
      >
        {/* Search input */}
        <div className="flex items-center gap-3 px-4 py-3 border-b border-slate-200 dark:border-slate-700">
          <Search className="w-5 h-5 text-slate-400" />
          <input
            ref={inputRef}
            type="text"
            data-testid="command-palette-input"
            className="flex-1 bg-transparent text-slate-900 dark:text-white placeholder-slate-400 outline-none text-sm"
            placeholder={t('search.placeholder')}
            value={query}
            onChange={e => setQuery(e.target.value)}
          />
          <kbd className="hidden sm:inline-block px-1.5 py-0.5 text-xs text-slate-400 bg-slate-100 dark:bg-slate-700 rounded">
            ESC
          </kbd>
        </div>

        {/* Results */}
        <div ref={listRef} className="max-h-80 overflow-y-auto py-2" data-testid="command-palette-results">
          {allResults.length === 0 ? (
            <div className="px-4 py-8 text-center text-sm text-slate-400">
              {t('search.noResults')}
            </div>
          ) : (
            allResults.map((item, index) => {
              const Icon = item.icon
              return (
                <button
                  key={item.id}
                  data-active={index === activeIndex}
                  className={`w-full flex items-center gap-3 px-4 py-2 text-sm text-left transition-colors ${
                    index === activeIndex
                      ? 'bg-primary-50 dark:bg-primary-900/30 text-primary-700 dark:text-primary-300'
                      : 'text-slate-700 dark:text-slate-300 hover:bg-slate-50 dark:hover:bg-slate-700/50'
                  }`}
                  onClick={() => {
                    navigate(item.path)
                    setOpen(false)
                  }}
                >
                  <Icon className="w-4 h-4 shrink-0" />
                  <span className="flex-1 truncate">{item.name}</span>
                  <span className={`px-1.5 py-0.5 text-xs rounded ${categoryBadgeColor(item.category)}`}>
                    {categoryLabel(item.category)}
                  </span>
                </button>
              )
            })
          )}
        </div>
      </div>
    </div>
  )
}
