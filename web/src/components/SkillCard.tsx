import { GitBranch, FolderOpen, Trash2 } from 'lucide-react'
import { api, type Skill, type Agent } from '../lib/api'
import { useI18n } from '../lib/i18n'
import { getTagColor } from '../lib/tagColors'
import { toast } from '../lib/toast'

interface SkillCardProps {
  skill: Skill
  tags?: string[]
  detectedAgents?: Agent[]
  onRemove?: (id: string) => void
  onToggleEnabled?: (id: string, enabled: boolean) => void
  onSyncChange?: () => void
  selected?: boolean
  onSelect?: (id: string) => void
  onClick?: (id: string) => void
}

const agentSyncedColors: Record<string, string> = {
  claude: 'bg-orange-100 text-orange-700 border-orange-300',
  cursor: 'bg-blue-100 text-blue-700 border-blue-300',
  codex: 'bg-purple-100 text-purple-700 border-purple-300',
  cline: 'bg-green-100 text-green-700 border-green-300',
  windsurf: 'bg-teal-100 text-teal-700 border-teal-300',
  github_copilot: 'bg-indigo-100 text-indigo-700 border-indigo-300',
}

const unsyncedStyle = 'bg-slate-50 dark:bg-slate-700 text-slate-400 dark:text-slate-500 border-slate-200 dark:border-slate-600'

export default function SkillCard({ skill, tags, detectedAgents, onRemove, onToggleEnabled, onSyncChange, selected, onSelect, onClick }: SkillCardProps) {
  const { t } = useI18n()
  const sourceIcon = skill.SourceType === 'git'
    ? <GitBranch className="w-3 h-3" />
    : <FolderOpen className="w-3 h-3" />

  const syncedAgents = new Set((skill.targets ?? []).map(t => t.agent))

  const handleAgentClick = async (e: React.MouseEvent, agentName: string) => {
    e.stopPropagation()
    try {
      if (syncedAgents.has(agentName)) {
        await api.skills.unsync(skill.ID, agentName)
      } else {
        await api.skills.sync(skill.ID, [agentName])
      }
      onSyncChange?.()
    } catch (err: any) {
      toast.error(err?.message || t('toast.error'))
    }
  }

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
        {(detectedAgents ?? []).map(a => {
          const isSynced = syncedAgents.has(a.name)
          const colorClass = isSynced
            ? (agentSyncedColors[a.name] ?? 'bg-primary-100 text-primary-700 border-primary-300')
            : unsyncedStyle
          const tip = isSynced ? `${t('skills.unsync')} ${a.display_name}` : `${t('skills.sync')} → ${a.display_name}`
          return (
            <span key={a.name} className="relative group/tip">
              <button
                onClick={e => handleAgentClick(e, a.name)}
                className={`px-2 py-0.5 rounded text-xs font-medium border transition-all hover:scale-105 ${colorClass}`}
              >
                {a.display_name}
              </button>
              <span className="absolute bottom-full left-1/2 -translate-x-1/2 mb-1.5 px-2 py-1 bg-slate-800 dark:bg-slate-200 text-white dark:text-slate-800 text-[11px] rounded whitespace-nowrap opacity-0 group-hover/tip:opacity-100 transition-opacity pointer-events-none z-10">
                {tip}
              </span>
            </span>
          )
        })}
      </div>
    </div>
  )
}
