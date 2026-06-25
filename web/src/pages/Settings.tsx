import { useEffect, useState } from 'react'
import { Settings as SettingsIcon, FolderOpen, Database, HardDrive } from 'lucide-react'
import { api } from '../lib/api'

export default function Settings() {
  const [settings, setSettings] = useState<Record<string, string>>({})
  useEffect(() => { api.settings.get().then(setSettings).catch(() => {}) }, [])

  const sections = [
    {
      title: 'Storage',
      icon: Database,
      items: Object.entries(settings).filter(([k]) => k.includes('dir') || k.includes('path')),
    },
    {
      title: 'Other',
      icon: HardDrive,
      items: Object.entries(settings).filter(([k]) => !k.includes('dir') && !k.includes('path')),
    },
  ].filter(s => s.items.length > 0)

  return (
    <div>
      <div className="flex items-center gap-3 mb-6">
        <SettingsIcon className="w-6 h-6 text-slate-400" />
        <h2 className="text-2xl font-bold text-slate-900">Settings</h2>
      </div>

      {sections.length === 0 ? (
        <div className="bg-white rounded-xl border border-slate-200 p-8 text-center text-slate-500 max-w-lg">
          <FolderOpen className="w-10 h-10 mx-auto mb-3 text-slate-300" />
          <p>No settings configured</p>
        </div>
      ) : (
        <div className="space-y-6 max-w-2xl">
          {sections.map(section => {
            const Icon = section.icon
            return (
              <div key={section.title} className="bg-white rounded-xl border border-slate-200 overflow-hidden">
                <div className="flex items-center gap-2 px-5 py-3 bg-slate-50 border-b border-slate-200">
                  <Icon className="w-4 h-4 text-slate-500" />
                  <h3 className="font-semibold text-slate-700 text-sm">{section.title}</h3>
                </div>
                <div className="p-5 space-y-4">
                  {section.items.map(([key, value]) => (
                    <div key={key}>
                      <label className="text-xs font-medium text-slate-500 uppercase tracking-wider block mb-1">{key}</label>
                      <div className="px-3 py-2 bg-slate-50 rounded-lg text-sm text-slate-700 font-mono">
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
