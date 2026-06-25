import { describe, it, expect, vi } from 'vitest'
import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import SkillsLibrary from '../SkillsLibrary'
import { renderWithRouter } from '../../test/utils'

vi.mock('../../lib/api', () => ({
  api: {
    skills: {
      list: vi.fn().mockResolvedValue([
        { ID: '1', Name: 'react-helper', Description: 'React utilities', SourceType: 'git', SourceRef: '', CentralPath: '', ContentHash: '', Enabled: true, targets: [] },
        { ID: '2', Name: 'go-patterns', Description: 'Go design patterns', SourceType: 'local', SourceRef: '', CentralPath: '', ContentHash: '', Enabled: false, targets: [] },
        { ID: '3', Name: 'test-skill', Description: 'Testing utilities', SourceType: 'git', SourceRef: '', CentralPath: '', ContentHash: '', Enabled: true, targets: [] },
      ]),
      install: vi.fn().mockResolvedValue({ installed: ['test'] }),
      remove: vi.fn().mockResolvedValue({}),
      sync: vi.fn().mockResolvedValue({}),
      unsync: vi.fn().mockResolvedValue({}),
    },
    tags: {
      list: vi.fn().mockResolvedValue(['frontend', 'backend']),
      getForSkill: vi.fn().mockImplementation((id: string) => {
        const map: Record<string, string[]> = {
          '1': ['frontend'],
          '2': ['backend'],
          '3': [],
        }
        return Promise.resolve(map[id] ?? [])
      }),
    },
    agents: {
      list: vi.fn().mockResolvedValue([
        { name: 'claude', display_name: 'Claude Code', detected: true, global_dir: '.claude/skills' },
      ]),
    },
    sync: { trigger: vi.fn().mockResolvedValue({}) },
  },
}))

describe('SkillsLibrary', () => {
  it('renders page title with skill count', async () => {
    renderWithRouter(<SkillsLibrary />)
    expect(screen.getByText('Library')).toBeInTheDocument()
    expect(await screen.findByText('3')).toBeInTheDocument()
  })

  it('renders skill cards', async () => {
    renderWithRouter(<SkillsLibrary />)
    expect(await screen.findByText('react-helper')).toBeInTheDocument()
    expect(screen.getByText('go-patterns')).toBeInTheDocument()
    expect(screen.getByText('test-skill')).toBeInTheDocument()
  })

  it('filters skills by search text', async () => {
    renderWithRouter(<SkillsLibrary />)
    await screen.findByText('react-helper')
    const searchInput = screen.getByPlaceholderText('Search skills in the control library...')
    await userEvent.type(searchInput, 'react')
    expect(screen.getByText('react-helper')).toBeInTheDocument()
    expect(screen.queryByText('go-patterns')).not.toBeInTheDocument()
  })

  it('filters by Enabled tab', async () => {
    renderWithRouter(<SkillsLibrary />)
    await screen.findByText('react-helper')
    // The tab button is the first "Enabled" button (before the card badges)
    const enabledButtons = screen.getAllByRole('button', { name: 'Enabled' })
    await userEvent.click(enabledButtons[0])
    expect(screen.getByText('react-helper')).toBeInTheDocument()
    expect(screen.queryByText('go-patterns')).not.toBeInTheDocument()
  })

  it('filters by Available tab', async () => {
    renderWithRouter(<SkillsLibrary />)
    await screen.findByText('go-patterns')
    await userEvent.click(screen.getByRole('button', { name: 'Available' }))
    expect(screen.getByText('go-patterns')).toBeInTheDocument()
    expect(screen.queryByText('react-helper')).not.toBeInTheDocument()
  })

  it('renders action buttons', () => {
    renderWithRouter(<SkillsLibrary />)
    expect(screen.getByText('Sync')).toBeInTheDocument()
    expect(screen.getByText('Check All')).toBeInTheDocument()
  })

  it('renders tag filter with tags from API', async () => {
    renderWithRouter(<SkillsLibrary />)
    // "frontend" and "backend" appear both as TagFilter pills and as tag pills on skill cards
    const frontendElements = await screen.findAllByText('frontend')
    expect(frontendElements.length).toBeGreaterThanOrEqual(1)
    const backendElements = screen.getAllByText('backend')
    expect(backendElements.length).toBeGreaterThanOrEqual(1)
    expect(screen.getByText('Untagged')).toBeInTheDocument()
  })
})
