interface TagFilterProps {
  tags: string[]
  activeTags: string[]
  onToggle: (tag: string) => void
}

export default function TagFilter({ tags, activeTags, onToggle }: TagFilterProps) {
  if (tags.length === 0) return null

  return (
    <div className="flex flex-wrap gap-2">
      {tags.map(tag => {
        const active = activeTags.includes(tag)
        return (
          <button
            key={tag}
            onClick={() => onToggle(tag)}
            className={`px-3 py-1 rounded-full text-xs font-medium transition-colors ${
              active
                ? 'bg-primary-500 text-white'
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
