import { GitBranch, FolderOpen, Trash2, RefreshCw } from 'lucide-react'
import type { Skill } from '../lib/api'
import { useI18n } from '../lib/i18n'
import { getTagColor } from '../lib/tagColors'

interface SkillCardProps {
  skill: Skill
  tags?: string[]
  onRemove?: (id: string) => void
  onSync?: (id: string) => void
  onToggleEnabled?: (id: string, enabled: boolean) => void
  selected?: boolean
  onSelect?: (id: string) => void
  onClick?: (id: string) => void
}

const agentColors: Record<string, string> = {
  claude: 'bg-orange-100 text-orange-700',
  cursor: 'bg-blue-100 text-blue-700',
  codex: 'bg-purple-100 text-purple-700',
}

export default function SkillCard({ skill, tags, onRemove, onSync, onToggleEnabled, selected, onSelect, onClick }: SkillCardProps) {
  const { t } = useI18n()
  const sourceIcon = skill.SourceType === 'git'
    ? <GitBranch className="w-3 h-3" />
    : <FolderOpen className="w-3 h-3" />

  return (
    <div
      className="bg-white dark:bg-slate-800 rounded-xl border border-slate-200 dark:border-slate-700 p-5 hover:shadow-md transition-shadow cursor-pointer"
      onClick={() => onClick?.(skill.ID)}
      data-testid="skill-card"
    >
      <div className="flex items-start justify-between">
        <div className="flex items-center gap-2">
          {onSelect && (
            <input
              type="checkbox"
              checked={selected}
              onChange={() => onSelect(skill.ID)}
              onClick={e => e.stopPropagation()}
              className="rounded border-slate-300"
            />
          )}
          <h3 className="font-semibold text-slate-900 dark:text-slate-100">{skill.Name}</h3>
        </div>
        <button
          type="button"
          onClick={e => { e.stopPropagation(); onToggleEnabled?.(skill.ID, !skill.Enabled) }}
          className={`px-2 py-0.5 rounded-full text-xs font-medium cursor-pointer transition-colors ${
            skill.Enabled
              ? 'bg-green-100 text-green-700 hover:bg-green-200 dark:bg-green-900 dark:text-green-300 dark:hover:bg-green-800'
              : 'bg-slate-100 dark:bg-slate-700 text-slate-500 dark:text-slate-400 hover:bg-slate-200 dark:hover:bg-slate-600'
          }`}
        >
          {skill.Enabled ? t('detail.enabled') : t('detail.disabled')}
        </button>
      </div>

      <p className="text-sm text-slate-500 dark:text-slate-400 mt-2 line-clamp-2">{skill.Description}</p>

      {tags && tags.length > 0 && (
        <div className="flex flex-wrap gap-1 mt-2">
          {tags.map(tag => (
            <span
              key={tag}
              className={`px-1.5 py-0.5 rounded-full text-[10px] font-medium ${getTagColor(tag)}`}
            >
              {tag}
            </span>
          ))}
        </div>
      )}

      <div className="flex items-center gap-2 mt-3">
        {onSync && (
          <button
            onClick={e => { e.stopPropagation(); onSync(skill.ID) }}
            className="text-xs text-primary-600 hover:text-primary-700 font-medium flex items-center gap-1"
          >
            <RefreshCw className="w-3 h-3" /> {t('skills.sync')}
          </button>
        )}
        {onRemove && (
          <button
            onClick={e => { e.stopPropagation(); onRemove(skill.ID) }}
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
