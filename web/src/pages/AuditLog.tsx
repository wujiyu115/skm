import { useEffect, useState } from 'react'
import { ClipboardList, Trash2 } from 'lucide-react'
import { api } from '../lib/api'
import { useI18n } from '../lib/i18n'
import { toast } from '../lib/toast'

interface AuditEntry {
  id: number
  action: string
  target: string
  detail: string
  created_at: string
}

const actionColors: Record<string, string> = {
  install:
    'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400',
  delete: 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400',
  enable: 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400',
  disable: 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400',
  sync: 'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-400',
  tag: 'bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-400',
  untag: 'bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-400',
  'rename-tag':
    'bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-400',
  'delete-tag':
    'bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-400',
  'group-create':
    'bg-teal-100 text-teal-800 dark:bg-teal-900/30 dark:text-teal-400',
  'group-update':
    'bg-teal-100 text-teal-800 dark:bg-teal-900/30 dark:text-teal-400',
  'group-delete':
    'bg-teal-100 text-teal-800 dark:bg-teal-900/30 dark:text-teal-400',
  'group-add-skill':
    'bg-teal-100 text-teal-800 dark:bg-teal-900/30 dark:text-teal-400',
  'group-remove-skill':
    'bg-teal-100 text-teal-800 dark:bg-teal-900/30 dark:text-teal-400',
}

function getActionColor(action: string): string {
  if (actionColors[action]) return actionColors[action]
  if (action.startsWith('tag') || action.includes('tag'))
    return 'bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-400'
  if (action.startsWith('group') || action.includes('group'))
    return 'bg-teal-100 text-teal-800 dark:bg-teal-900/30 dark:text-teal-400'
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
  return date.toLocaleDateString(undefined, {
    month: 'short',
    day: 'numeric',
    year: date.getFullYear() !== now.getFullYear() ? 'numeric' : undefined,
  })
}

export default function AuditLog() {
  const { t } = useI18n()
  const [entries, setEntries] = useState<AuditEntry[]>([])
  const [loading, setLoading] = useState(true)

  const fetchEntries = () => {
    setLoading(true)
    api.audit
      .list()
      .then(setEntries)
      .catch(() => {})
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    fetchEntries()
  }, [])

  const handlePrune = async () => {
    try {
      await api.audit.prune()
      toast.success(t('audit.pruned'))
      fetchEntries()
    } catch {
      toast.error(t('toast.error'))
    }
  }

  return (
    <div className="max-w-5xl mx-auto">
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <ClipboardList className="w-7 h-7 text-primary-500 dark:text-primary-400" />
          <h1 className="text-2xl font-bold text-slate-900 dark:text-white">
            {t('audit.title')}
          </h1>
        </div>
        <button
          onClick={handlePrune}
          className="flex items-center gap-2 px-4 py-2 text-sm font-medium text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg hover:bg-red-100 dark:hover:bg-red-900/30 transition-colors"
        >
          <Trash2 className="w-4 h-4" />
          {t('audit.prune')}
        </button>
      </div>

      {loading ? (
        <div className="text-center py-12 text-slate-400 dark:text-slate-500">
          Loading...
        </div>
      ) : entries.length === 0 ? (
        <div className="text-center py-16" data-testid="audit-empty">
          <ClipboardList className="w-12 h-12 mx-auto text-slate-300 dark:text-slate-600 mb-3" />
          <p className="text-slate-500 dark:text-slate-400">
            {t('audit.noEntries')}
          </p>
        </div>
      ) : (
        <div className="bg-white dark:bg-slate-800 rounded-xl border border-slate-200 dark:border-slate-700 overflow-hidden">
          <table className="w-full text-sm" data-testid="audit-table">
            <thead>
              <tr className="border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-800/50">
                <th className="text-left px-4 py-3 font-medium text-slate-500 dark:text-slate-400">
                  {t('audit.time')}
                </th>
                <th className="text-left px-4 py-3 font-medium text-slate-500 dark:text-slate-400">
                  {t('audit.action')}
                </th>
                <th className="text-left px-4 py-3 font-medium text-slate-500 dark:text-slate-400">
                  {t('audit.target')}
                </th>
                <th className="text-left px-4 py-3 font-medium text-slate-500 dark:text-slate-400">
                  {t('audit.detail')}
                </th>
              </tr>
            </thead>
            <tbody>
              {entries.map(entry => (
                <tr
                  key={entry.id}
                  className="border-b border-slate-100 dark:border-slate-700/50 last:border-0 hover:bg-slate-50 dark:hover:bg-slate-700/30 transition-colors"
                >
                  <td className="px-4 py-3 text-slate-500 dark:text-slate-400 whitespace-nowrap">
                    {formatTime(entry.created_at)}
                  </td>
                  <td className="px-4 py-3">
                    <span
                      className={`inline-block px-2 py-0.5 rounded-full text-xs font-medium ${getActionColor(entry.action)}`}
                    >
                      {entry.action}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-slate-700 dark:text-slate-300 font-mono text-xs">
                    {entry.target}
                  </td>
                  <td className="px-4 py-3 text-slate-600 dark:text-slate-400">
                    {entry.detail}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
