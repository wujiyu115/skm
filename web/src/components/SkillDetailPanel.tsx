import { useEffect, useState, useCallback } from 'react'
import { X } from 'lucide-react'
import { api, type Skill, type Target } from '../lib/api'
import { useI18n } from '../lib/i18n'
import { toast } from '../lib/toast'
import { getTagColor } from '../lib/tagColors'
import SkillMarkdown from './SkillMarkdown'

interface SkillDetailPanelProps {
  skillId: string | null
  onClose: () => void
}

const agentColors: Record<string, string> = {
  claude: 'bg-orange-100 text-orange-700 dark:bg-orange-900 dark:text-orange-300',
  cursor: 'bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-300',
  codex: 'bg-purple-100 text-purple-700 dark:bg-purple-900 dark:text-purple-300',
}

export default function SkillDetailPanel({ skillId, onClose }: SkillDetailPanelProps) {
  const { t } = useI18n()
  const [skill, setSkill] = useState<Skill | null>(null)
  const [targets, setTargets] = useState<Target[]>([])
  const [content, setContent] = useState<string>('')
  const [tags, setTags] = useState<string[]>([])
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (!skillId) {
      setSkill(null)
      setTargets([])
      setContent('')
      setTags([])
      return
    }

    setLoading(true)
    Promise.all([
      api.skills.get(skillId),
      api.skills.content(skillId).catch(() => ({ content: '' })),
      api.tags.getForSkill(skillId).catch(() => [] as string[]),
    ])
      .then(([detail, contentRes, tagsRes]) => {
        setSkill(detail.skill)
        setTargets(detail.targets ?? [])
        setContent(contentRes.content ?? '')
        setTags(tagsRes)
      })
      .catch(() => {
        toast.error(t('toast.error'))
      })
      .finally(() => setLoading(false))
  }, [skillId, t])

  const handleToggleEnabled = async () => {
    if (!skill) return
    try {
      await api.skills.setEnabled(skill.ID, !skill.Enabled)
      setSkill({ ...skill, Enabled: !skill.Enabled })
      toast.success(t(skill.Enabled ? 'toast.skillDisabled' : 'toast.skillEnabled'))
    } catch {
      toast.error(t('toast.error'))
    }
  }

  const handleKeyDown = useCallback((e: KeyboardEvent) => {
    if (e.key === 'Escape') onClose()
  }, [onClose])

  useEffect(() => {
    if (skillId) {
      document.addEventListener('keydown', handleKeyDown)
      return () => document.removeEventListener('keydown', handleKeyDown)
    }
  }, [skillId, handleKeyDown])

  const isOpen = skillId !== null

  return (
    <>
      {/* Overlay */}
      <div
        className={`fixed inset-0 bg-black/40 z-40 transition-opacity duration-300 ${
          isOpen ? 'opacity-100' : 'opacity-0 pointer-events-none'
        }`}
        onClick={onClose}
        data-testid="detail-overlay"
      />

      {/* Panel */}
      <div
        className={`fixed top-0 right-0 h-full w-full max-w-2xl bg-white dark:bg-slate-900 shadow-2xl z-50 transform transition-transform duration-300 ease-in-out ${
          isOpen ? 'translate-x-0' : 'translate-x-full'
        }`}
        data-testid="detail-panel"
      >
        {loading ? (
          <div className="flex items-center justify-center h-full">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600" data-testid="loading-spinner" />
          </div>
        ) : skill ? (
          <div className="h-full flex flex-col overflow-hidden">
            {/* Header */}
            <div className="flex items-center justify-between px-6 py-4 border-b border-slate-200 dark:border-slate-700">
              <h2 className="text-xl font-bold text-slate-900 dark:text-slate-100 truncate">
                {skill.Name}
              </h2>
              <button
                onClick={onClose}
                className="p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-500 dark:text-slate-400 transition-colors"
                aria-label={t('detail.close')}
                data-testid="close-button"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            {/* Body */}
            <div className="flex-1 overflow-y-auto px-6 py-5 space-y-5">
              {/* Description + Enabled toggle */}
              <div className="flex items-start justify-between gap-3">
                <p className="text-sm text-slate-600 dark:text-slate-400">{skill.Description}</p>
                <button
                  type="button"
                  onClick={handleToggleEnabled}
                  className={`shrink-0 px-3 py-1 rounded-full text-xs font-medium cursor-pointer transition-colors ${
                    skill.Enabled
                      ? 'bg-green-100 text-green-700 hover:bg-green-200 dark:bg-green-900 dark:text-green-300 dark:hover:bg-green-800'
                      : 'bg-slate-100 dark:bg-slate-700 text-slate-500 dark:text-slate-400 hover:bg-slate-200 dark:hover:bg-slate-600'
                  }`}
                >
                  {skill.Enabled ? t('detail.enabled') : t('detail.disabled')}
                </button>
              </div>

              {/* Tags */}
              {tags.length > 0 && (
                <div>
                  <h3 className="text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wider mb-2">
                    {t('detail.tags')}
                  </h3>
                  <div className="flex flex-wrap gap-1.5">
                    {tags.map(tag => (
                      <span
                        key={tag}
                        className={`px-2 py-0.5 rounded-full text-xs font-medium ${getTagColor(tag)}`}
                      >
                        {tag}
                      </span>
                    ))}
                  </div>
                </div>
              )}

              {/* Source metadata */}
              <div>
                <h3 className="text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wider mb-2">
                  {t('detail.metadata')}
                </h3>
                <div className="grid grid-cols-2 gap-2 text-sm">
                  <div className="bg-slate-50 dark:bg-slate-800 rounded-lg px-3 py-2">
                    <span className="text-slate-500 dark:text-slate-400 text-xs">{t('detail.sourceType')}</span>
                    <p className="font-medium text-slate-900 dark:text-slate-100">{skill.SourceType}</p>
                  </div>
                  <div className="bg-slate-50 dark:bg-slate-800 rounded-lg px-3 py-2">
                    <span className="text-slate-500 dark:text-slate-400 text-xs">{t('detail.sourceRef')}</span>
                    <p className="font-medium text-slate-900 dark:text-slate-100 truncate">{skill.SourceRef}</p>
                  </div>
                </div>
              </div>

              {/* Synced agents */}
              {targets.length > 0 && (
                <div>
                  <h3 className="text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wider mb-2">
                    {t('detail.agents')}
                  </h3>
                  <div className="flex flex-wrap gap-2">
                    {targets.map(target => (
                      <span
                        key={target.agent}
                        className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg text-xs font-medium ${
                          agentColors[target.agent] ?? 'bg-slate-100 dark:bg-slate-700 text-slate-600 dark:text-slate-400'
                        }`}
                      >
                        <span className="w-1.5 h-1.5 rounded-full bg-green-500" />
                        {target.agent}
                      </span>
                    ))}
                  </div>
                </div>
              )}

              {/* Content */}
              <div>
                <h3 className="text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wider mb-2">
                  {t('detail.content')}
                </h3>
                {content ? (
                  <SkillMarkdown content={content} />
                ) : (
                  <p className="text-sm text-slate-400 dark:text-slate-500 italic">
                    {t('detail.noContent')}
                  </p>
                )}
              </div>
            </div>
          </div>
        ) : null}
      </div>
    </>
  )
}
