import { useI18n } from '../lib/i18n'
import { getTagColor } from '../lib/tagColors'

interface TagFilterProps {
  tags: string[]
  activeTags: string[]
  onToggle: (tag: string) => void
}

export default function TagFilter({ tags, activeTags, onToggle }: TagFilterProps) {
  const { t } = useI18n()

  if (tags.length === 0) return null

  const isAllActive = activeTags.length === 0

  const handleAllClick = () => {
    // Clear all active tags — parent handles this by detecting empty array
    for (const tag of activeTags) {
      onToggle(tag)
    }
  }

  return (
    <div className="flex flex-wrap gap-2">
      <button
        onClick={handleAllClick}
        className={`px-3 py-1 rounded-full text-xs font-medium transition-colors ${
          isAllActive
            ? 'bg-primary-500 text-white'
            : 'bg-slate-100 dark:bg-slate-700 text-slate-600 dark:text-slate-400 hover:bg-slate-200 dark:hover:bg-slate-600'
        }`}
      >
        {t('skills.all')}
      </button>
      <button
        onClick={() => onToggle('__untagged__')}
        className={`px-3 py-1 rounded-full text-xs font-medium transition-colors ${
          activeTags.includes('__untagged__')
            ? 'bg-slate-600 text-white dark:bg-slate-300 dark:text-slate-900'
            : 'bg-slate-100 dark:bg-slate-700 text-slate-600 dark:text-slate-400 hover:bg-slate-200 dark:hover:bg-slate-600'
        }`}
      >
        {t('tags.untagged')}
      </button>
      {tags.map(tag => {
        const active = activeTags.includes(tag)
        return (
          <button
            key={tag}
            onClick={() => onToggle(tag)}
            className={`px-3 py-1 rounded-full text-xs font-medium transition-colors ${
              active
                ? getTagColor(tag)
                : 'bg-slate-100 dark:bg-slate-700 text-slate-600 dark:text-slate-400 hover:bg-slate-200 dark:hover:bg-slate-600'
            }`}
          >
            {tag}
          </button>
        )
      })}
    </div>
  )
}
