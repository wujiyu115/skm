import { useEffect, useState } from 'react'
import { Settings as SettingsIcon, FolderOpen, Database, HardDrive } from 'lucide-react'
import { api } from '../lib/api'
import { useI18n } from '../lib/i18n'

export default function Settings() {
  const { t } = useI18n()
  const [settings, setSettings] = useState<Record<string, string>>({})
  useEffect(() => { api.settings.get().then(setSettings).catch(() => {}) }, [])

  const sections = [
    {
      title: t('settings.storage'),
      icon: Database,
      items: Object.entries(settings).filter(([k]) => k.includes('dir') || k.includes('path')),
    },
    {
      title: t('settings.other'),
      icon: HardDrive,
      items: Object.entries(settings).filter(([k]) => !k.includes('dir') && !k.includes('path')),
    },
  ].filter(s => s.items.length > 0)

  return (
    <div>
      <div className="flex items-center gap-3 mb-6">
        <SettingsIcon className="w-6 h-6 text-slate-400" />
        <h2 className="text-2xl font-bold text-slate-900 dark:text-slate-100">{t('settings.title')}</h2>
      </div>

      {sections.length === 0 ? (
        <div className="bg-white dark:bg-slate-800 rounded-xl border border-slate-200 dark:border-slate-700 p-8 text-center text-slate-500 dark:text-slate-400 max-w-lg">
          <FolderOpen className="w-10 h-10 mx-auto mb-3 text-slate-300" />
          <p>{t('settings.noSettings')}</p>
        </div>
      ) : (
        <div className="space-y-6 max-w-2xl">
          {sections.map(section => {
            const Icon = section.icon
            return (
              <div key={section.title} className="bg-white dark:bg-slate-800 rounded-xl border border-slate-200 dark:border-slate-700 overflow-hidden">
                <div className="flex items-center gap-2 px-5 py-3 bg-slate-50 dark:bg-slate-800/50 border-b border-slate-200 dark:border-slate-700">
                  <Icon className="w-4 h-4 text-slate-500 dark:text-slate-400" />
                  <h3 className="font-semibold text-slate-700 dark:text-slate-300 text-sm">{section.title}</h3>
                </div>
                <div className="p-5 space-y-4">
                  {section.items.map(([key, value]) => (
                    <div key={key}>
                      <label className="text-xs font-medium text-slate-500 dark:text-slate-400 uppercase tracking-wider block mb-1">{key}</label>
                      <div className="px-3 py-2 bg-slate-50 dark:bg-slate-800/50 rounded-lg text-sm text-slate-700 dark:text-slate-300 font-mono">
                        {value || '—'}
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )
          })}
        </div>
      )}
    </div>
  )
}
