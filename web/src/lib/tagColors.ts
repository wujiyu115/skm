const tagColors = [
  'bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-300',
  'bg-purple-100 text-purple-700 dark:bg-purple-900 dark:text-purple-300',
  'bg-pink-100 text-pink-700 dark:bg-pink-900 dark:text-pink-300',
  'bg-amber-100 text-amber-700 dark:bg-amber-900 dark:text-amber-300',
  'bg-teal-100 text-teal-700 dark:bg-teal-900 dark:text-teal-300',
]

export function getTagColor(tag: string): string {
  const hash = [...tag].reduce((h, c) => (h << 5) - h + c.charCodeAt(0), 0)
  return tagColors[Math.abs(hash) % tagColors.length]
}
