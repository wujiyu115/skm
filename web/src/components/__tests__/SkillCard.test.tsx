import { describe, it, expect, vi } from 'vitest'
import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import SkillCard from '../SkillCard'
import { mockSkill, mockTarget } from '../../test/mocks'
import { renderWithRouter } from '../../test/utils'

vi.mock('../../lib/api', () => ({
  api: {
    skills: {
      sync: vi.fn().mockResolvedValue({}),
      unsync: vi.fn().mockResolvedValue({}),
    },
  },
}))

const detectedAgents = [
  { name: 'claude', display_name: 'Claude Code', detected: true, global_dir: '.claude/skills', project_dir: '.claude/skills', detect_paths: ['.claude'], is_builtin: true, category: 'coding' },
  { name: 'cursor', display_name: 'Cursor', detected: true, global_dir: '.cursor/skills', project_dir: '.cursor/skills', detect_paths: ['.cursor'], is_builtin: true, category: 'coding' },
]

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

  it('shows agent badges with synced styling for targets', () => {
    const skill = mockSkill({
      targets: [mockTarget({ agent: 'claude' })],
    })
    renderWithRouter(<SkillCard skill={skill} detectedAgents={detectedAgents} />)
    expect(screen.getByText('Claude Code')).toBeInTheDocument()
    expect(screen.getByText('Cursor')).toBeInTheDocument()
  })

  it('calls onRemove when remove clicked', async () => {
    const onRemove = vi.fn()
    renderWithRouter(<SkillCard skill={mockSkill()} onRemove={onRemove} />)
    await userEvent.click(screen.getByText('remove'))
    expect(onRemove).toHaveBeenCalledWith('skill-1')
  })

  it('calls onToggleEnabled with toggled value when badge is clicked', async () => {
    const onToggleEnabled = vi.fn()
    renderWithRouter(<SkillCard skill={mockSkill({ Enabled: true })} onToggleEnabled={onToggleEnabled} />)
    await userEvent.click(screen.getByText('Enabled'))
    expect(onToggleEnabled).toHaveBeenCalledWith('skill-1', false)
  })

  it('calls onToggleEnabled to enable when disabled badge is clicked', async () => {
    const onToggleEnabled = vi.fn()
    renderWithRouter(<SkillCard skill={mockSkill({ Enabled: false })} onToggleEnabled={onToggleEnabled} />)
    await userEvent.click(screen.getByText('Disabled'))
    expect(onToggleEnabled).toHaveBeenCalledWith('skill-1', true)
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

  it('renders tag pills when tags provided', () => {
    renderWithRouter(<SkillCard skill={mockSkill()} tags={['frontend', 'react']} />)
    expect(screen.getByText('frontend')).toBeInTheDocument()
    expect(screen.getByText('react')).toBeInTheDocument()
  })

  it('does not render tag section when tags is empty', () => {
    renderWithRouter(<SkillCard skill={mockSkill()} tags={[]} />)
    expect(screen.queryByText('frontend')).not.toBeInTheDocument()
  })

  it('does not render tag section when tags is undefined', () => {
    renderWithRouter(<SkillCard skill={mockSkill()} />)
    expect(screen.getByText('test-skill')).toBeInTheDocument()
  })
})
