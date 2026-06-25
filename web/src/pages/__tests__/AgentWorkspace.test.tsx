import { describe, it, expect, vi } from 'vitest'
import { screen } from '@testing-library/react'
import AgentWorkspace from '../AgentWorkspace'
import { renderWithRouter } from '../../test/utils'

vi.mock('../../lib/api', () => ({
  api: {
    agents: {
      list: vi.fn().mockResolvedValue([
        { name: 'claude', display_name: 'Claude Code', project_dir: '.claude/skills', global_dir: '.claude/skills', detected: true },
        { name: 'cursor', display_name: 'Cursor', project_dir: '.cursor/skills', global_dir: '.cursor/skills', detected: false },
      ]),
      skills: vi.fn().mockResolvedValue([]),
    },
    skills: {
      list: vi.fn().mockResolvedValue([]),
    },
  },
}))

describe('AgentWorkspace', () => {
  it('renders agent cards', async () => {
    renderWithRouter(<AgentWorkspace />)
    expect(screen.getByText('Agent Workspaces')).toBeInTheDocument()
    expect(await screen.findByText('Claude Code')).toBeInTheDocument()
    expect(screen.getByText('Cursor')).toBeInTheDocument()
  })

  it('shows detected status', async () => {
    renderWithRouter(<AgentWorkspace />)
    expect(await screen.findByText('Active')).toBeInTheDocument()
    expect(screen.getByText('Not detected')).toBeInTheDocument()
  })

  it('shows global dir', async () => {
    renderWithRouter(<AgentWorkspace />)
    await screen.findByText('Claude Code')
    const codeEls = screen.getAllByText((_, el) => el?.tagName === 'CODE' && (el?.textContent?.includes('.claude/skills') ?? false))
    expect(codeEls.length).toBeGreaterThan(0)
  })
})
