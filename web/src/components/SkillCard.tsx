import { GitBranch, FolderOpen, Trash2, RefreshCw } from 'lucide-react'
import type { Skill } from '../lib/api'
import { useI18n } from '../lib/i18n'

interface SkillCardProps {
  skill: Skill
  onRemove?: (id: string) => void
  onSync?: (id: string) => void
  selected?: boolean
  onSelect?: (id: string) => void
}

const agentColors: Record<string, string> = {
  claude: 'bg-orange-100 text-orange-700',
  cursor: 'bg-blue-100 text-blue-700',
  codex: 'bg-purple-100 text-purple-700',
}

export default function SkillCard({ skill, onRemove, onSync, selected, onSelect }: SkillCardProps) {
  const { t } = useI18n()
  const sourceIcon = skill.SourceType === 'git'
    ? <GitBranch className="w-3 h-3" />
    : <FolderOpen className="w-3 h-3" />

  return (
    <div className="bg-white dark:bg-slate-800 rounded-xl border border-slate-200 dark:border-slate-700 p-5 hover:shadow-md transition-shadow">
      <div className="flex items-start justify-between">
        <div className="flex items-center gap-2">
          {onSelect && (
            <input
              type="checkbox"
              checked={selected}
              onChange={() => onSelect(skill.ID)}
              className="rounded border-slate-300"
            />
          )}
          <h3 className="font-semibold text-slate-900 dark:text-slate-100">{skill.Name}</h3>
        </div>
        <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${
          skill.Enabled
            ? 'bg-primary-100 text-primary-700'
            : 'bg-slate-100 dark:bg-slate-700 text-slate-500 dark:text-slate-400'
        }`}>
          {skill.Enabled ? t('skills.enabled') : t('skills.disabled')}
        </span>
      </div>

      <p className="text-sm text-slate-500 dark:text-slate-400 mt-2 line-clamp-2">{skill.Description}</p>

      <div className="flex items-center gap-2 mt-3">
        {onSync && (
          <button
            onClick={() => onSync(skill.ID)}
            className="text-xs text-primary-600 hover:text-primary-700 font-medium flex items-center gap-1"
          >
            <RefreshCw className="w-3 h-3" /> {t('skills.sync')}
          </button>
        )}
        {onRemove && (
          <button
            onClick={() => onRemove(skill.ID)}
            className="text-xs text-red-500 hover:text-red-600 font-medium flex items-center gap-1"
          >
            <Trash2 className="w-3 h-3" /> {t('skills.remove')}
          </button>
        )}
      </div>

      <div className="flex flex-wrap items-center gap-1.5 mt-3">
        <span className="inline-flex items-center gap-1 px-2 py-0.5 bg-slate-100 dark:bg-slate-700 text-slate-600 dark:text-slate-400 rounded text-xs">
          {sourceIcon} {skill.SourceType}
        </span>
        {(skill.targets ?? []).map(t => (
          <span key={t.agent} className={`px-2 py-0.5 rounded text-xs font-medium ${agentColors[t.agent] ?? 'bg-slate-100 dark:bg-slate-700 text-slate-600 dark:text-slate-400'}`}>
            {t.agent}
          </span>
        ))}
      </div>
    </div>
  )
}
