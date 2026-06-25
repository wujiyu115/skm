import { useState } from 'react'
import { Download, Search, FolderOpen, GitBranch, ChevronDown, ChevronUp, Check, Package } from 'lucide-react'
import { api } from '../lib/api'
import { useI18n } from '../lib/i18n'

type Tab = 'registry' | 'local' | 'git'
type Filter = 'all' | 'official' | 'community'

interface RegistrySkill {
  name: string
  description: string
  source: string
  repo: string
  category: Filter
  icon: string
}

const registry: RegistrySkill[] = [
  { name: 'find-skills', description: 'Discover and search for AI agent skills', source: '@anthropics/skills', repo: 'anthropics/skills', category: 'official', icon: '🔍' },
  { name: 'frontend-design', description: 'Frontend UI/UX design patterns and components', source: '@anthropics/skills', repo: 'anthropics/skills', category: 'official', icon: '🎨' },
  { name: 'agent-browser', description: 'Browser automation and web scraping skills', source: '@vercel-labs/agent-skills', repo: 'vercel-labs/agent-skills', category: 'official', icon: '🌐' },
  { name: 'vercel-react-best-practices', description: 'React best practices from Vercel team', source: '@vercel-labs/agent-skills', repo: 'vercel-labs/agent-skills', category: 'official', icon: '⚛️' },
  { name: 'web-design-guidelines', description: 'Web design guidelines and accessibility', source: '@vercel-labs/agent-skills', repo: 'vercel-labs/agent-skills', category: 'official', icon: '📐' },
  { name: 'microsoft-foundry', description: 'Microsoft Foundry integration skills', source: '@anthropics/skills', repo: 'anthropics/skills', category: 'official', icon: '🏗️' },
  { name: 'azure-ai', description: 'Azure AI services integration', source: '@microsoft/azure-skills', repo: 'microsoft/azure-skills', category: 'official', icon: '🤖' },
  { name: 'azure-deploy', description: 'Azure deployment and CI/CD automation', source: '@microsoft/azure-skills', repo: 'microsoft/azure-skills', category: 'official', icon: '🚀' },
  { name: 'azure-diagnostics', description: 'Azure diagnostics and monitoring', source: '@microsoft/azure-skills', repo: 'microsoft/azure-skills', category: 'official', icon: '🔧' },
  { name: 'azure-prepare', description: 'Azure environment preparation and setup', source: '@microsoft/azure-skills', repo: 'microsoft/azure-skills', category: 'official', icon: '📦' },
  { name: 'azure-storage', description: 'Azure storage management skills', source: '@microsoft/azure-skills', repo: 'microsoft/azure-skills', category: 'official', icon: '💾' },
  { name: 'azure-validate', description: 'Azure configuration validation', source: '@microsoft/azure-skills', repo: 'microsoft/azure-skills', category: 'official', icon: '✅' },
  { name: 'remoetions-dev-skills', description: 'Remote development workflow skills', source: '@remoetions/dev-skills', repo: 'remoetions/dev-skills', category: 'community', icon: '💻' },
]

const sources = [...new Set(registry.map(s => s.source))]

export default function Install() {
  const { t } = useI18n()
  const [tab, setTab] = useState<Tab>('registry')
  const [search, setSearch] = useState('')
  const [filter, setFilter] = useState<Filter>('all')
  const [sourceFilter, setSourceFilter] = useState('')
  const [installing, setInstalling] = useState<string | null>(null)
  const [installed, setInstalled] = useState<Set<string>>(new Set())
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null)

  const [localPath, setLocalPath] = useState('')
  const [gitUrl, setGitUrl] = useState('')
  const [showAdvanced, setShowAdvanced] = useState(false)
  const [globalInstall, setGlobalInstall] = useState(false)

  const filteredRegistry = registry.filter(sk => {
    if (search && !sk.name.toLowerCase().includes(search.toLowerCase()) &&
        !sk.description.toLowerCase().includes(search.toLowerCase())) return false
    if (filter !== 'all' && sk.category !== filter) return false
    if (sourceFilter && sk.source !== sourceFilter) return false
    return true
  })

  const doInstall = async (source: string, key?: string) => {
    const installKey = key ?? source
    setInstalling(installKey)
    setMessage(null)
    try {
      await api.skills.install(source, [], globalInstall)
      setInstalled(prev => new Set(prev).add(installKey))
      setMessage({ type: 'success', text: t('install.success') })
    } catch {
      setMessage({ type: 'error', text: t('install.error') })
    } finally {
      setInstalling(null)
    }
  }

  const tabs: { key: Tab; label: string; icon: typeof Download }[] = [
    { key: 'registry', label: t('install.tabRegistry'), icon: Package },
    { key: 'local', label: t('install.tabLocal'), icon: FolderOpen },
    { key: 'git', label: t('install.tabGit'), icon: GitBranch },
  ]

  const filters: { key: Filter; label: string }[] = [
    { key: 'all', label: t('install.filterAll') },
    { key: 'official', label: t('install.filterOfficial') },
    { key: 'community', label: t('install.filterCommunity') },
  ]

  return (
    <div>
      <div className="flex items-center gap-3 mb-6">
        <Download className="w-6 h-6 text-primary-600" />
        <h2 className="text-2xl font-bold text-slate-900">{t('install.title')}</h2>
      </div>

      {message && (
        <div className={`mb-4 px-4 py-3 rounded-lg text-sm font-medium ${
          message.type === 'success' ? 'bg-primary-50 text-primary-700 border border-primary-200' : 'bg-red-50 text-red-700 border border-red-200'
        }`}>
          {message.text}
        </div>
      )}

      <div className="flex items-center gap-1 mb-6 border-b border-slate-200">
        {tabs.map(t => {
          const Icon = t.icon
          return (
            <button
              key={t.key}
              onClick={() => setTab(t.key)}
              className={`flex items-center gap-2 px-4 py-2.5 text-sm font-medium border-b-2 transition-colors ${
                tab === t.key
                  ? 'border-primary-600 text-primary-600'
                  : 'border-transparent text-slate-500 hover:text-slate-700'
              }`}
            >
              <Icon className="w-4 h-4" /> {t.label}
            </button>
          )
        })}
      </div>

      {tab === 'registry' && (
        <div>
          <div className="flex items-center gap-3 mb-4">
            <div className="flex items-center gap-1">
              {filters.map(f => (
                <button
                  key={f.key}
                  onClick={() => setFilter(f.key)}
                  className={`px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
                    filter === f.key ? 'bg-primary-600 text-white' : 'text-slate-600 hover:bg-slate-100'
                  }`}
                >
                  {f.label}
                </button>
              ))}
            </div>
            <select
              value={sourceFilter}
              onChange={e => setSourceFilter(e.target.value)}
              className="px-3 py-1.5 border border-slate-200 rounded-lg text-sm text-slate-600 focus:outline-none focus:ring-2 focus:ring-primary-500"
            >
              <option value="">{t('install.filterSource')}</option>
              {sources.map(s => <option key={s} value={s}>{s}</option>)}
            </select>
            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
              <input
                value={search}
                onChange={e => setSearch(e.target.value)}
                placeholder={t('install.search')}
                className="w-full pl-10 pr-4 py-2 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
              />
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4 gap-4">
            {filteredRegistry.map(sk => {
              const isInstalled = installed.has(sk.name)
              const isInstalling = installing === sk.name
              return (
                <div key={sk.name} className="bg-white rounded-xl border border-slate-200 p-4 hover:shadow-md transition-shadow">
                  <div className="flex items-start justify-between">
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 rounded-lg bg-primary-50 flex items-center justify-center text-lg">
                        {sk.icon}
                      </div>
                      <div>
                        <h3 className="font-semibold text-slate-900 text-sm">{sk.name}</h3>
                        <span className="text-xs text-slate-400">{sk.source}</span>
                      </div>
                    </div>
                    <button
                      onClick={() => doInstall(sk.repo, sk.name)}
                      disabled={isInstalling || isInstalled}
                      className={`flex items-center gap-1 px-3 py-1.5 rounded-lg text-xs font-medium transition-colors ${
                        isInstalled
                          ? 'bg-primary-100 text-primary-700'
                          : 'bg-primary-600 text-white hover:bg-primary-700 disabled:opacity-50'
                      }`}
                    >
                      {isInstalled ? (
                        <><Check className="w-3 h-3" /> {t('install.installed')}</>
                      ) : isInstalling ? (
                        t('install.installing')
                      ) : (
                        <><Download className="w-3 h-3" /> {t('install.installBtn')}</>
                      )}
                    </button>
                  </div>
                  <p className="text-xs text-slate-500 mt-2 line-clamp-2">{sk.description}</p>
                </div>
              )
            })}
          </div>
        </div>
      )}

      {tab === 'local' && (
        <div className="max-w-xl">
          <div className="bg-white rounded-xl border border-slate-200 p-6">
            <div className="flex items-center gap-3 mb-4">
              <FolderOpen className="w-8 h-8 text-slate-400" />
              <div>
                <h3 className="font-semibold text-slate-900">{t('install.localTitle')}</h3>
                <p className="text-sm text-slate-500">{t('install.localDesc')}</p>
              </div>
            </div>
            <div className="space-y-4">
              <input
                value={localPath}
                onChange={e => setLocalPath(e.target.value)}
                placeholder={t('install.localPlaceholder')}
                className="w-full px-4 py-2.5 border border-slate-200 rounded-lg text-sm font-mono focus:outline-none focus:ring-2 focus:ring-primary-500"
                onKeyDown={e => e.key === 'Enter' && localPath.trim() && doInstall(localPath)}
              />
              <AdvancedOptions
                show={showAdvanced}
                onToggle={() => setShowAdvanced(!showAdvanced)}
                globalInstall={globalInstall}
                setGlobalInstall={setGlobalInstall}
              />
              <button
                onClick={() => localPath.trim() && doInstall(localPath)}
                disabled={!localPath.trim() || installing !== null}
                className="flex items-center gap-2 px-4 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 disabled:opacity-50 transition-colors"
              >
                <Download className="w-4 h-4" />
                {installing ? t('install.installing') : t('install.installBtn')}
              </button>
            </div>
          </div>
        </div>
      )}

      {tab === 'git' && (
        <div className="max-w-xl">
          <div className="bg-white rounded-xl border border-slate-200 p-6">
            <div className="flex items-center gap-3 mb-4">
              <GitBranch className="w-8 h-8 text-slate-400" />
              <div>
                <h3 className="font-semibold text-slate-900">{t('install.gitTitle')}</h3>
                <p className="text-sm text-slate-500">{t('install.gitDesc')}</p>
              </div>
            </div>
            <div className="space-y-4">
              <div>
                <input
                  value={gitUrl}
                  onChange={e => setGitUrl(e.target.value)}
                  placeholder={t('install.gitPlaceholder')}
                  className="w-full px-4 py-2.5 border border-slate-200 rounded-lg text-sm font-mono focus:outline-none focus:ring-2 focus:ring-primary-500"
                  onKeyDown={e => e.key === 'Enter' && gitUrl.trim() && doInstall(gitUrl)}
                />
                <p className="text-xs text-slate-400 mt-1.5">
                  https://github.com/user/repo.git or github.com/user/repo
                </p>
              </div>
              <AdvancedOptions
                show={showAdvanced}
                onToggle={() => setShowAdvanced(!showAdvanced)}
                globalInstall={globalInstall}
                setGlobalInstall={setGlobalInstall}
              />
              <button
                onClick={() => gitUrl.trim() && doInstall(gitUrl)}
                disabled={!gitUrl.trim() || installing !== null}
                className="flex items-center gap-2 px-4 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 disabled:opacity-50 transition-colors"
              >
                <Download className="w-4 h-4" />
                {installing ? t('install.installing') : t('install.installBtn')}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

function AdvancedOptions({ show, onToggle, globalInstall, setGlobalInstall }: {
  show: boolean
  onToggle: () => void
  globalInstall: boolean
  setGlobalInstall: (v: boolean) => void
}) {
  const { t } = useI18n()
  return (
    <div>
      <button
        onClick={onToggle}
        className="flex items-center gap-1.5 text-sm text-slate-500 hover:text-slate-700"
      >
        {show ? <ChevronUp className="w-4 h-4" /> : <ChevronDown className="w-4 h-4" />}
        {t('install.advancedOptions')}
      </button>
      {show && (
        <div className="mt-3 p-4 bg-slate-50 rounded-lg space-y-3">
          <label className="flex items-center gap-2 text-sm text-slate-700 cursor-pointer">
            <input
              type="checkbox"
              checked={globalInstall}
              onChange={e => setGlobalInstall(e.target.checked)}
              className="rounded border-slate-300"
            />
            {t('install.global')}
          </label>
        </div>
      )}
    </div>
  )
}
