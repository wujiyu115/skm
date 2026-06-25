import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import TagFilter from '../TagFilter'

describe('TagFilter', () => {
  const tags = ['git', 'local', 'frontend']

  it('renders all tags', () => {
    render(<TagFilter tags={tags} activeTags={[]} onToggle={() => {}} />)
    expect(screen.getByText('git')).toBeInTheDocument()
    expect(screen.getByText('local')).toBeInTheDocument()
    expect(screen.getByText('frontend')).toBeInTheDocument()
  })

  it('highlights active tags', () => {
    render(<TagFilter tags={tags} activeTags={['git']} onToggle={() => {}} />)
    const gitBtn = screen.getByText('git')
    expect(gitBtn.className).toContain('bg-primary-500')
    const localBtn = screen.getByText('local')
    expect(localBtn.className).toContain('bg-slate-100')
  })

  it('calls onToggle when tag clicked', async () => {
    const onToggle = vi.fn()
    render(<TagFilter tags={tags} activeTags={[]} onToggle={onToggle} />)
    await userEvent.click(screen.getByText('local'))
    expect(onToggle).toHaveBeenCalledWith('local')
  })

  it('renders nothing for empty tags', () => {
    const { container } = render(<TagFilter tags={[]} activeTags={[]} onToggle={() => {}} />)
    expect(container.innerHTML).toBe('')
  })
})
