import { useEffect, useState, useCallback } from 'react'
import { Settings as SettingsIcon, Palette, RefreshCw, Info, ClipboardList, Trash2, ChevronDown, ChevronRight } from 'lucide-react'
import { api } from '../lib/api'
import { useI18n } from '../lib/i18n'
import { useTheme } from '../lib/theme'
import { toast } from '../lib/toast'

type SyncMode = 'symlink' | 'copy'
type ThemeOption = 'light' | 'dark' | 'system'
type TextSize = 'small' | 'default' | 'large'
type UpdateInterval = 'off' | '1h' | '6h' | '24h'

interface AuditEntry {
  id: number
  action: string
  target: string
  detail: string
  created_at: string
}

const actionColors: Record<string, string> = {
  install: 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400',
  delete: 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400',
  enable: 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400',
  disable: 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400',
  sync: 'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-400',
}

function getActionColor(action: string): string {
  if (actionColors[action]) return actionColors[action]
  if (action.includes('tag')) return 'bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-400'
  if (action.includes('group')) return 'bg-teal-100 text-teal-800 dark:bg-teal-900/30 dark:text-teal-400'
  return 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300'
}

function formatTime(iso: string): string {
  const date = new Date(iso)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffMin = Math.floor(diffMs / 60000)
  const diffHour = Math.floor(diffMs / 3600000)
  const diffDay = Math.floor(diffMs / 86400000)
  if (diffMin < 1) return 'just now'
  if (diffMin < 60) return `${diffMin}m ago`
  if (diffHour < 24) return `${diffHour}h ago`
  if (diffDay < 7) return `${diffDay}d ago`
  return date.toLocaleDateString(undefined, { month: 'short', day: 'numeric' })
}

export default function Settings() {
  const { t, locale, setLocale } = useI18n()
  const { theme, setTheme } = useTheme()
  const [settings, setSettings] = useState<Record<string, string>>({})
  const [loading, setLoading] = useState(true)
  const [auditEntries, setAuditEntries] = useState<AuditEntry[]>([])
  const [auditOpen, setAuditOpen] = useState(false)

  const fetchAudit = useCallback(() => {
    api.audit.list().then(setAuditEntries).catch(() => {})
  }, [])

  useEffect(() => {
    api.settings.get()
      .then(setSettings)
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [])

  const saveSetting = useCallback(async (key: string, value: string) => {
    try {
      await api.settings.update({ [key]: value })
      setSettings(prev => ({ ...prev, [key]: value }))
      toast.success(t('settings.saved'))
    } catch {
      toast.error(t('toast.error'))
    }
  }, [t])

  const syncMode = (settings.sync_mode as SyncMode) || 'symlink'
  const textSize = (settings.text_size as TextSize) || 'default'
  const autoUpdate = (settings.auto_update_interval as UpdateInterval) || 'off'

  const themeValue: ThemeOption = settings.theme === 'system'
    ? 'system'
    : theme

  const handleThemeChange = useCallback(async (value: ThemeOption) => {
    if (value === 'system') {
      const systemDark = window.matchMedia('(prefers-color-scheme: dark)').matches
      setTheme(systemDark ? 'dark' : 'light')
    } else {
      setTheme(value)
    }
    await saveSetting('theme', value)
  }, [setTheme, saveSetting])

  const handleLanguageChange = useCallback(async (value: string) => {
    if (value === 'en' || value === 'zh') {
      setLocale(value)
      await saveSetting('language', value)
    }
  }, [setLocale, saveSetting])

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500" />
      </div>
    )
  }

  return (
    <div>
      <div className="flex items-center gap-3 mb-6">
        <SettingsIcon className="w-6 h-6 text-slate-400" />
        <h2 className="text-2xl font-bold text-slate-900 dark:text-slate-100">{t('settings.title')}</h2>
      </div>

      <div className="space-y-6 max-w-2xl">
        {/* Storage Section */}
        <section className="bg-white dark:bg-slate-800 rounded-xl border border-slate-200 dark:border-slate-700 overflow-hidden">
          <div className="flex items-center gap-2 px-5 py-3 bg-slate-50 dark:bg-slate-800/50 border-b border-slate-200 dark:border-slate-700">
            <SettingsIcon className="w-4 h-4 text-slate-500 dark:text-slate-400" />
            <h3 className="font-semibold text-slate-700 dark:text-slate-300 text-sm">{t('settings.storage')}</h3>
          </div>
          <div className="p-5 space-y-4">
            <ReadOnlyField label="skills_dir" value={settings.skills_dir} />
            <ReadOnlyField label="cache_dir" value={settings.cache_dir} />
            <div>
              <label className="text-xs font-medium text-slate-500 dark:text-slate-400 uppercase tracking-wider block mb-2">
                {t('settings.syncMode')}
              </label>
              <div className="flex gap-2">
                <SegmentButton
                  active={syncMode === 'symlink'}
                  onClick={() => saveSetting('sync_mode', 'symlink')}
                  label={t('settings.symlink')}
                />
                <SegmentButton
                  active={syncMode === 'copy'}
                  onClick={() => saveSetting('sync_mode', 'copy')}
                  label={t('settings.copy')}
                />
              </div>
            </div>
          </div>
        </section>

        {/* Appearance Section */}
        <section className="bg-white dark:bg-slate-800 rounded-xl border border-slate-200 dark:border-slate-700 overflow-hidden">
          <div className="flex items-center gap-2 px-5 py-3 bg-slate-50 dark:bg-slate-800/50 border-b border-slate-200 dark:border-slate-700">
            <Palette className="w-4 h-4 text-slate-500 dark:text-slate-400" />
            <h3 className="font-semibold text-slate-700 dark:text-slate-300 text-sm">{t('settings.appearance')}</h3>
          </div>
          <div className="p-5 space-y-4">
            <SelectField
              label={t('settings.theme')}
              value={themeValue}
              onChange={(v) => handleThemeChange(v as ThemeOption)}
              options={[
                { value: 'light', label: 'Light' },
                { value: 'dark', label: 'Dark' },
                { value: 'system', label: 'System' },
              ]}
            />
            <SelectField
              label={t('settings.language')}
              value={locale}
              onChange={handleLanguageChange}
              options={[
                { value: 'en', label: 'English' },
                { value: 'zh', label: '中文' },
              ]}
            />
            <SelectField
              label={t('settings.textSize')}
              value={textSize}
              onChange={(v) => saveSetting('text_size', v)}
              options={[
                { value: 'small', label: 'Small' },
                { value: 'default', label: 'Default' },
                { value: 'large', label: 'Large' },
              ]}
            />
          </div>
        </section>

        {/* Updates Section */}
        <section className="bg-white dark:bg-slate-800 rounded-xl border border-slate-200 dark:border-slate-700 overflow-hidden">
          <div className="flex items-center gap-2 px-5 py-3 bg-slate-50 dark:bg-slate-800/50 border-b border-slate-200 dark:border-slate-700">
            <RefreshCw className="w-4 h-4 text-slate-500 dark:text-slate-400" />
            <h3 className="font-semibold text-slate-700 dark:text-slate-300 text-sm">{t('settings.updates')}</h3>
          </div>
          <div className="p-5">
            <SelectField
              label={t('settings.autoUpdate')}
              value={autoUpdate}
              onChange={(v) => saveSetting('auto_update_interval', v)}
              options={[
                { value: 'off', label: 'Off' },
                { value: '1h', label: '1 hour' },
                { value: '6h', label: '6 hours' },
                { value: '24h', label: '24 hours' },
              ]}
            />
          </div>
        </section>

        {/* About Section */}
        <section className="bg-white dark:bg-slate-800 rounded-xl border border-slate-200 dark:border-slate-700 overflow-hidden">
          <div className="flex items-center gap-2 px-5 py-3 bg-slate-50 dark:bg-slate-800/50 border-b border-slate-200 dark:border-slate-700">
            <Info className="w-4 h-4 text-slate-500 dark:text-slate-400" />
            <h3 className="font-semibold text-slate-700 dark:text-slate-300 text-sm">{t('settings.about')}</h3>
          </div>
          <div className="p-5 space-y-4">
            <ReadOnlyField label={t('settings.version')} value={settings.version || '0.1.0'} />
            <ReadOnlyField label="skills_dir" value={settings.skills_dir} />
            <ReadOnlyField label="cache_dir" value={settings.cache_dir} />
          </div>
        </section>

        {/* Audit Log Section */}
        <section className="bg-white dark:bg-slate-800 rounded-xl border border-slate-200 dark:border-slate-700 overflow-hidden">
          <button
            onClick={() => { setAuditOpen(o => !o); if (!auditOpen) fetchAudit() }}
            className="flex items-center justify-between w-full px-5 py-3 bg-slate-50 dark:bg-slate-800/50 border-b border-slate-200 dark:border-slate-700"
          >
            <div className="flex items-center gap-2">
              <ClipboardList className="w-4 h-4 text-slate-500 dark:text-slate-400" />
              <h3 className="font-semibold text-slate-700 dark:text-slate-300 text-sm">{t('audit.title')}</h3>
            </div>
            {auditOpen ? <ChevronDown className="w-4 h-4 text-slate-400" /> : <ChevronRight className="w-4 h-4 text-slate-400" />}
          </button>
          {auditOpen && (
            <div className="p-5">
              <div className="flex justify-end mb-3">
                <button
                  onClick={async () => {
                    try { await api.audit.prune(); toast.success(t('audit.pruned')); fetchAudit() }
                    catch { toast.error(t('toast.error')) }
                  }}
                  className="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg hover:bg-red-100 dark:hover:bg-red-900/30 transition-colors"
                >
                  <Trash2 className="w-3.5 h-3.5" />
                  {t('audit.prune')}
                </button>
              </div>
              {auditEntries.length === 0 ? (
                <div className="text-center py-8 text-slate-400 dark:text-slate-500">
                  <ClipboardList className="w-10 h-10 mx-auto mb-2 opacity-50" />
                  <p className="text-sm">{t('audit.noEntries')}</p>
                </div>
              ) : (
                <div className="overflow-x-auto">
                  <table className="w-full text-sm" data-testid="audit-table">
                    <thead>
                      <tr className="border-b border-slate-200 dark:border-slate-700">
                        <th className="text-left px-3 py-2 font-medium text-slate-500 dark:text-slate-400 text-xs">{t('audit.time')}</th>
                        <th className="text-left px-3 py-2 font-medium text-slate-500 dark:text-slate-400 text-xs">{t('audit.action')}</th>
                        <th className="text-left px-3 py-2 font-medium text-slate-500 dark:text-slate-400 text-xs">{t('audit.target')}</th>
                        <th className="text-left px-3 py-2 font-medium text-slate-500 dark:text-slate-400 text-xs">{t('audit.detail')}</th>
                      </tr>
                    </thead>
                    <tbody>
                      {auditEntries.map(entry => (
                        <tr key={entry.id} className="border-b border-slate-100 dark:border-slate-700/50 last:border-0">
                          <td className="px-3 py-2 text-slate-500 dark:text-slate-400 whitespace-nowrap text-xs">{formatTime(entry.created_at)}</td>
                          <td className="px-3 py-2">
                            <span className={`inline-block px-2 py-0.5 rounded-full text-xs font-medium ${getActionColor(entry.action)}`}>{entry.action}</span>
                          </td>
                          <td className="px-3 py-2 text-slate-700 dark:text-slate-300 font-mono text-xs">{entry.target}</td>
                          <td className="px-3 py-2 text-slate-600 dark:text-slate-400 text-xs">{entry.detail}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </div>
          )}
        </section>
      </div>
    </div>
  )
}

function ReadOnlyField({ label, value }: { label: string; value?: string }) {
  return (
    <div>
      <label className="text-xs font-medium text-slate-500 dark:text-slate-400 uppercase tracking-wider block mb-1">
        {label}
      </label>
      <div className="px-3 py-2 bg-slate-50 dark:bg-slate-800/50 rounded-lg text-sm text-slate-500 dark:text-slate-400 font-mono">
        {value || '—'}
      </div>
    </div>
  )
}

function SegmentButton({ active, onClick, label }: { active: boolean; onClick: () => void; label: string }) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
        active
          ? 'bg-blue-500 text-white'
          : 'bg-slate-100 dark:bg-slate-700 text-slate-600 dark:text-slate-300 hover:bg-slate-200 dark:hover:bg-slate-600'
      }`}
    >
      {label}
    </button>
  )
}

function SelectField({
  label,
  value,
  onChange,
  options,
}: {
  label: string
  value: string
  onChange: (value: string) => void
  options: { value: string; label: string }[]
}) {
  return (
    <div>
      <label className="text-xs font-medium text-slate-500 dark:text-slate-400 uppercase tracking-wider block mb-1">
        {label}
      </label>
      <select
        value={value}
        onChange={(e) => onChange(e.target.value)}
        className="w-full px-3 py-2 bg-slate-50 dark:bg-slate-700 border border-slate-200 dark:border-slate-600 rounded-lg text-sm text-slate-700 dark:text-slate-300 focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        {options.map((opt) => (
          <option key={opt.value} value={opt.value}>
            {opt.label}
          </option>
        ))}
      </select>
    </div>
  )
}
