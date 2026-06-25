import { useEffect, useState } from 'react'
import { Link, useLocation } from 'react-router-dom'
import {
  LayoutDashboard, BookOpen, Download, Settings,
  ChevronDown, ChevronRight, FolderOpen, Globe, Terminal, Code, MousePointer, Sparkles, Languages,
  Sun, Moon,
} from 'lucide-react'
import { api, type Group, type Agent, type Skill } from '../lib/api'
import { useI18n } from '../lib/i18n'
import { useTheme } from '../lib/theme'

const agentIcons: Record<string, typeof Terminal> = {
  claude: Terminal,
  cursor: MousePointer,
  codex: Code,
}

export default function Sidebar() {
  const location = useLocation()
  const { t, locale, setLocale } = useI18n()
  const { theme, toggle: toggleTheme } = useTheme()
  const [groups, setGroups] = useState<Group[]>([])
  const [agents, setAgents] = useState<Agent[]>([])
  const [skills, setSkills] = useState<Skill[]>([])
  const [collapsed, setCollapsed] = useState<Record<string, boolean>>({})

  useEffect(() => {
    Promise.all([
      api.groups.list(),
      api.agents.list(),
      api.skills.list(),
    ]).then(([g, a, s]) => {
      setGroups(g ?? [])
      setAgents(a ?? [])
      setSkills(s ?? [])
    }).catch(() => {})
  }, [location.pathname])

  const toggle = (section: string) => {
    setCollapsed(prev => ({ ...prev, [section]: !prev[section] }))
  }

  const isActive = (path: string) => location.pathname === path

  const agentSkillCount = (agentName: string) =>
    skills.filter(s => s.targets?.some(t => t.agent === agentName)).length

  const mainNav = [
    { path: '/', label: t('nav.dashboard'), icon: LayoutDashboard },
    { path: '/skills', label: t('nav.library'), icon: BookOpen },
    { path: '/install', label: t('nav.install'), icon: Download },
  ]

  return (
    <nav className="w-64 bg-sidebar text-slate-300 flex flex-col h-screen" data-testid="sidebar">
      <div className="p-5 flex items-center gap-2 border-b border-slate-700">
        <Sparkles className="w-6 h-6 text-primary-400" />
        <span className="font-bold text-white text-lg">{t('app.title')}</span>
      </div>

      <div className="flex-1 overflow-y-auto py-3 px-3 space-y-1">
        {mainNav.map(item => {
          const Icon = item.icon
          return (
            <Link
              key={item.path}
              to={item.path}
              className={`flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors ${
                isActive(item.path)
                  ? 'bg-primary-600 text-white font-medium'
                  : 'hover:bg-sidebar-hover'
              }`}
            >
              <Icon className="w-4 h-4" />
              {item.label}
            </Link>
          )
        })}

        <div className="pt-4">
          <button
            onClick={() => toggle('presets')}
            className="flex items-center gap-2 px-3 py-1.5 text-xs font-semibold uppercase text-slate-500 w-full hover:text-slate-300"
          >
            {collapsed.presets ? <ChevronRight className="w-3 h-3" /> : <ChevronDown className="w-3 h-3" />}
            <FolderOpen className="w-3 h-3" />
            {t('nav.presets')}
          </button>
          {!collapsed.presets && (
            <div className="ml-4 space-y-0.5">
              {groups.map(g => (
                <Link
                  key={g.id}
                  to={`/groups/${g.id}`}
                  className={`flex items-center justify-between px-3 py-1.5 rounded-md text-sm transition-colors ${
                    isActive(`/groups/${g.id}`)
                      ? 'bg-primary-600 text-white'
                      : 'hover:bg-sidebar-hover'
                  }`}
                >
                  <span>{g.name}</span>
                  <span className="text-xs text-slate-500">{g.skill_count ?? 0}</span>
                </Link>
              ))}
              {groups.length === 0 && (
                <div className="px-3 py-1.5 text-xs text-slate-500 italic">
                  {locale === 'zh' ? '暂无分组' : 'No groups'}
                </div>
              )}
              <Link
                to="/groups"
                className="flex items-center px-3 py-1.5 text-xs text-slate-500 hover:text-slate-300"
              >
                {t('nav.newGroup')}
              </Link>
            </div>
          )}
        </div>

        <div className="pt-2">
          <button
            onClick={() => toggle('agents')}
            className="flex items-center gap-2 px-3 py-1.5 text-xs font-semibold uppercase text-slate-500 w-full hover:text-slate-300"
          >
            {collapsed.agents ? <ChevronRight className="w-3 h-3" /> : <ChevronDown className="w-3 h-3" />}
            <Globe className="w-3 h-3" />
            {t('nav.workspace')}
          </button>
          {!collapsed.agents && (
            <div className="ml-4 space-y-0.5">
              <Link
                to="/agents"
                className={`flex items-center justify-between px-3 py-1.5 rounded-md text-sm transition-colors ${
                  location.pathname === '/agents' && !location.pathname.includes('/agents/')
                    ? 'bg-primary-600 text-white'
                    : 'hover:bg-sidebar-hover'
                }`}
              >
                <span>{t('nav.allAgents')}</span>
              </Link>
              {agents.map(a => {
                const Icon = agentIcons[a.name] ?? Terminal
                const count = agentSkillCount(a.name)
                return (
                  <Link
                    key={a.name}
                    to={`/agents/${a.name}`}
                    className={`flex items-center justify-between px-3 py-1.5 rounded-md text-sm transition-colors ${
                      isActive(`/agents/${a.name}`)
                        ? 'bg-primary-600 text-white'
                        : 'hover:bg-sidebar-hover'
                    }`}
                  >
                    <span className="flex items-center gap-2">
                      <Icon className="w-3.5 h-3.5" />
                      {a.display_name}
                    </span>
                    <span className="flex items-center gap-1.5">
                      {a.detected && <span className="w-1.5 h-1.5 rounded-full bg-primary-400" />}
                      <span className="text-xs text-slate-500">{count}</span>
                    </span>
                  </Link>
                )
              })}
            </div>
          )}
        </div>
      </div>

      <div className="p-3 border-t border-slate-700 space-y-1">
        <button
          onClick={toggleTheme}
          className="flex items-center gap-3 px-3 py-2 rounded-lg text-sm w-full hover:bg-sidebar-hover transition-colors"
        >
          {theme === 'light' ? <Moon className="w-4 h-4" /> : <Sun className="w-4 h-4" />}
          {theme === 'light' ? (locale === 'zh' ? '深色模式' : 'Dark Mode') : (locale === 'zh' ? '浅色模式' : 'Light Mode')}
        </button>
        <button
          onClick={() => setLocale(locale === 'en' ? 'zh' : 'en')}
          className="flex items-center gap-3 px-3 py-2 rounded-lg text-sm w-full hover:bg-sidebar-hover transition-colors"
        >
          <Languages className="w-4 h-4" />
          {locale === 'en' ? '中文' : 'English'}
        </button>
        <Link
          to="/settings"
          className={`flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors ${
            isActive('/settings')
              ? 'bg-primary-600 text-white font-medium'
              : 'hover:bg-sidebar-hover'
          }`}
        >
          <Settings className="w-4 h-4" />
          {t('nav.settings')}
        </Link>
      </div>
    </nav>
  )
}
