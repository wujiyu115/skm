import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import TagFilter from '../TagFilter'
import { I18nProvider } from '../../lib/i18n'
import { getTagColor } from '../../lib/tagColors'

function renderTagFilter(props: { tags: string[]; activeTags: string[]; onToggle: (tag: string) => void }) {
  return render(
    <I18nProvider>
      <TagFilter {...props} />
    </I18nProvider>
  )
}

describe('TagFilter', () => {
  const tags = ['frontend', 'backend', 'devops']

  it('renders all tags plus All and Untagged buttons', () => {
    renderTagFilter({ tags, activeTags: [], onToggle: () => {} })
    expect(screen.getByText('All')).toBeInTheDocument()
    expect(screen.getByText('Untagged')).toBeInTheDocument()
    expect(screen.getByText('frontend')).toBeInTheDocument()
    expect(screen.getByText('backend')).toBeInTheDocument()
    expect(screen.getByText('devops')).toBeInTheDocument()
  })

  it('highlights All when no tags are active', () => {
    renderTagFilter({ tags, activeTags: [], onToggle: () => {} })
    const allBtn = screen.getByText('All')
    expect(allBtn.className).toContain('bg-primary-500')
  })

  it('All is not highlighted when tags are active', () => {
    renderTagFilter({ tags, activeTags: ['frontend'], onToggle: () => {} })
    const allBtn = screen.getByText('All')
    expect(allBtn.className).toContain('bg-slate-100')
  })

  it('active tags get deterministic color', () => {
    renderTagFilter({ tags, activeTags: ['frontend'], onToggle: () => {} })
    const frontendBtn = screen.getByText('frontend')
    const expectedColor = getTagColor('frontend')
    // The active tag should have one of the color classes from getTagColor
    for (const cls of expectedColor.split(' ')) {
      expect(frontendBtn.className).toContain(cls)
    }
  })

  it('inactive tags have gray styling', () => {
    renderTagFilter({ tags, activeTags: ['frontend'], onToggle: () => {} })
    const backendBtn = screen.getByText('backend')
    expect(backendBtn.className).toContain('bg-slate-100')
  })

  it('calls onToggle when tag clicked', async () => {
    const onToggle = vi.fn()
    renderTagFilter({ tags, activeTags: [], onToggle })
    await userEvent.click(screen.getByText('backend'))
    expect(onToggle).toHaveBeenCalledWith('backend')
  })

  it('calls onToggle when Untagged clicked', async () => {
    const onToggle = vi.fn()
    renderTagFilter({ tags, activeTags: [], onToggle })
    await userEvent.click(screen.getByText('Untagged'))
    expect(onToggle).toHaveBeenCalledWith('__untagged__')
  })

  it('clicking All deselects all active tags', async () => {
    const onToggle = vi.fn()
    renderTagFilter({ tags, activeTags: ['frontend', 'backend'], onToggle })
    await userEvent.click(screen.getByText('All'))
    expect(onToggle).toHaveBeenCalledTimes(2)
    expect(onToggle).toHaveBeenCalledWith('frontend')
    expect(onToggle).toHaveBeenCalledWith('backend')
  })

  it('renders nothing for empty tags', () => {
    const { container } = renderTagFilter({ tags: [], activeTags: [], onToggle: () => {} })
    expect(container.innerHTML).toBe('')
  })
})
