import { Trash2, ToggleLeft, ToggleRight, RefreshCw, X } from 'lucide-react'
import { useI18n } from '../lib/i18n'

interface BatchToolbarProps {
  selectedIds: string[]
  onClear: () => void
  onAction: (action: string) => void
}

export default function BatchToolbar({ selectedIds, onClear, onAction }: BatchToolbarProps) {
  const { t } = useI18n()

  if (selectedIds.length === 0) return null

  const handleDelete = () => {
    if (window.confirm(t('batch.confirm'))) {
      onAction('delete')
    }
  }

  return (
    <div
      data-testid="batch-toolbar"
      className="fixed bottom-0 left-0 right-0 z-50 transform transition-transform duration-200 ease-out"
      style={{ transform: selectedIds.length > 0 ? 'translateY(0)' : 'translateY(100%)' }}
    >
      <div className="mx-auto max-w-4xl px-4 pb-4">
        <div className="flex items-center justify-between gap-3 rounded-xl bg-slate-900 dark:bg-slate-800 px-5 py-3 shadow-2xl border border-slate-700">
          <div className="flex items-center gap-3">
            <span className="text-sm font-medium text-white">
              {selectedIds.length} {t('batch.selected')}
            </span>
            <button
              onClick={onClear}
              className="p-1 text-slate-400 hover:text-white rounded transition-colors"
              aria-label={t('batch.cancel')}
            >
              <X className="w-4 h-4" />
            </button>
          </div>

          <div className="flex items-center gap-2">
            <button
              onClick={() => onAction('enable')}
              className="flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium text-slate-200 hover:bg-slate-700 dark:hover:bg-slate-600 rounded-lg transition-colors"
            >
              <ToggleRight className="w-4 h-4" />
              {t('batch.enable')}
            </button>
            <button
              onClick={() => onAction('disable')}
              className="flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium text-slate-200 hover:bg-slate-700 dark:hover:bg-slate-600 rounded-lg transition-colors"
            >
              <ToggleLeft className="w-4 h-4" />
              {t('batch.disable')}
            </button>
            <button
              onClick={() => onAction('sync')}
              className="flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium text-slate-200 hover:bg-slate-700 dark:hover:bg-slate-600 rounded-lg transition-colors"
            >
              <RefreshCw className="w-4 h-4" />
              {t('batch.sync')}
            </button>
            <div className="w-px h-6 bg-slate-600 mx-1" />
            <button
              onClick={handleDelete}
              className="flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium text-red-400 hover:bg-red-900/30 rounded-lg transition-colors"
            >
              <Trash2 className="w-4 h-4" />
              {t('batch.delete')}
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
