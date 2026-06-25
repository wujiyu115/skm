import { describe, it, expect, vi } from 'vitest'
import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import Sidebar from '../Sidebar'
import { renderWithRouter } from '../../test/utils'

vi.mock('../../lib/api', () => ({
  api: {
    groups: {
      list: vi.fn().mockResolvedValue([
        { id: 'g1', name: 'Frontend', description: '', skill_count: 5 },
        { id: 'g2', name: 'DevOps', description: '', skill_count: 2 },
      ]),
    },
    agents: {
      list: vi.fn().mockResolvedValue([
        { name: 'claude', display_name: 'Claude Code', project_dir: '.claude/skills', global_dir: '.claude/skills', detected: true },
        { name: 'cursor', display_name: 'Cursor', project_dir: '.cursor/skills', global_dir: '.cursor/skills', detected: false },
      ]),
    },
    skills: {
      list: vi.fn().mockResolvedValue([
        { ID: 's1', Name: 'sk1', targets: [{ agent: 'claude', skill_id: 's1', target_path: '', mode: 'symlink' }] },
        { ID: 's2', Name: 'sk2', targets: [{ agent: 'claude', skill_id: 's2', target_path: '', mode: 'symlink' }] },
      ]),
    },
    projects: {
      list: vi.fn().mockResolvedValue([
        { id: 'p1', name: 'my-app', path: '/home/user/my-app', created_at: '2026-01-15T10:00:00Z' },
      ]),
    },
  },
}))

describe('Sidebar', () => {
  it('renders logo and Skills Manager text', async () => {
    renderWithRouter(<Sidebar />)
    expect(screen.getByText('Skills Manager')).toBeInTheDocument()
  })

  it('renders main nav links', () => {
    renderWithRouter(<Sidebar />)
    expect(screen.getByText('Dashboard')).toBeInTheDocument()
    expect(screen.getByText('Library')).toBeInTheDocument()
    expect(screen.getByText('Install Skills')).toBeInTheDocument()
    expect(screen.getByText('Settings')).toBeInTheDocument()
  })

  it('highlights active link based on route', () => {
    renderWithRouter(<Sidebar />, { route: '/skills' })
    const libraryLink = screen.getByText('Library').closest('a')
    expect(libraryLink?.className).toContain('bg-primary-600')
  })

  it('renders groups section', async () => {
    renderWithRouter(<Sidebar />)
    expect(screen.getByText('Presets')).toBeInTheDocument()
    expect(await screen.findByText('Frontend')).toBeInTheDocument()
    expect(await screen.findByText('DevOps')).toBeInTheDocument()
  })

  it('shows skill counts next to groups', async () => {
    renderWithRouter(<Sidebar />)
    expect(await screen.findByText('5')).toBeInTheDocument()
    const twos = await screen.findAllByText('2')
    expect(twos.length).toBeGreaterThanOrEqual(1)
  })

  it('renders agents section with only detected agents', async () => {
    renderWithRouter(<Sidebar />)
    expect(screen.getByText('Global Workspace')).toBeInTheDocument()
    expect(await screen.findByText('Claude Code')).toBeInTheDocument()
    expect(screen.queryByText('Cursor')).not.toBeInTheDocument()
  })

  it('collapses presets section on click', async () => {
    renderWithRouter(<Sidebar />)
    await screen.findByText('Frontend')
    await userEvent.click(screen.getByText('Presets'))
    expect(screen.queryByText('Frontend')).not.toBeInTheDocument()
  })

  it('collapses agents section on click', async () => {
    renderWithRouter(<Sidebar />)
    const claude = await screen.findByText('Claude Code')
    expect(claude).toBeInTheDocument()
    await userEvent.click(screen.getByText('Global Workspace'))
    expect(screen.queryByText('Claude Code')).not.toBeInTheDocument()
  })

  it('renders project workspace section with projects', async () => {
    renderWithRouter(<Sidebar />)
    expect(screen.getByText('Project Workspace')).toBeInTheDocument()
    expect(await screen.findByText('my-app')).toBeInTheDocument()
    expect(screen.getByText('+ Add Project')).toBeInTheDocument()
  })
})
