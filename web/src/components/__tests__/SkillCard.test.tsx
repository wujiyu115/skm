import { describe, it, expect, vi } from 'vitest'
import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import SkillCard from '../SkillCard'
import { mockSkill, mockTarget } from '../../test/mocks'
import { renderWithRouter } from '../../test/utils'

describe('SkillCard', () => {
  it('renders skill name and description', () => {
    renderWithRouter(<SkillCard skill={mockSkill()} />)
    expect(screen.getByText('test-skill')).toBeInTheDocument()
    expect(screen.getByText('A test skill for unit testing')).toBeInTheDocument()
  })

  it('shows Enabled badge when skill is enabled', () => {
    renderWithRouter(<SkillCard skill={mockSkill({ Enabled: true })} />)
    expect(screen.getByText('Enabled')).toBeInTheDocument()
  })

  it('shows Disabled badge when skill is disabled', () => {
    renderWithRouter(<SkillCard skill={mockSkill({ Enabled: false })} />)
    expect(screen.getByText('Disabled')).toBeInTheDocument()
  })

  it('shows source type badge', () => {
    renderWithRouter(<SkillCard skill={mockSkill({ SourceType: 'git' })} />)
    expect(screen.getByText('git')).toBeInTheDocument()
  })

  it('shows agent badges for targets', () => {
    const skill = mockSkill({
      targets: [mockTarget({ agent: 'claude' }), mockTarget({ agent: 'cursor' })],
    })
    renderWithRouter(<SkillCard skill={skill} />)
    expect(screen.getByText('claude')).toBeInTheDocument()
    expect(screen.getByText('cursor')).toBeInTheDocument()
  })

  it('calls onRemove when remove clicked', async () => {
    const onRemove = vi.fn()
    renderWithRouter(<SkillCard skill={mockSkill()} onRemove={onRemove} />)
    await userEvent.click(screen.getByText('remove'))
    expect(onRemove).toHaveBeenCalledWith('skill-1')
  })

  it('calls onSync when sync clicked', async () => {
    const onSync = vi.fn()
    renderWithRouter(<SkillCard skill={mockSkill()} onSync={onSync} />)
    await userEvent.click(screen.getByText('Sync'))
    expect(onSync).toHaveBeenCalledWith('skill-1')
  })

  it('shows checkbox when onSelect provided', () => {
    renderWithRouter(<SkillCard skill={mockSkill()} onSelect={() => {}} selected={false} />)
    expect(screen.getByRole('checkbox')).toBeInTheDocument()
  })

  it('calls onSelect when checkbox clicked', async () => {
    const onSelect = vi.fn()
    renderWithRouter(<SkillCard skill={mockSkill()} onSelect={onSelect} selected={false} />)
    await userEvent.click(screen.getByRole('checkbox'))
    expect(onSelect).toHaveBeenCalledWith('skill-1')
  })
})
