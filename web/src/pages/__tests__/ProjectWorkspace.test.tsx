import { describe, it, expect, vi } from 'vitest'
import { screen } from '@testing-library/react'
import { Routes, Route } from 'react-router-dom'
import ProjectWorkspace from '../ProjectWorkspace'
import { renderWithRouter } from '../../test/utils'

vi.mock('../../lib/api', () => ({
  api: {
    projects: {
      list: vi.fn().mockResolvedValue([
        { id: 'p1', name: 'frontend-app', path: '/home/user/frontend-app', created_at: '2026-01-15T10:00:00Z' },
        { id: 'p2', name: 'backend-svc', path: '/home/user/backend-svc', created_at: '2026-02-20T08:00:00Z' },
      ]),
      skills: vi.fn().mockResolvedValue([
        { agent: 'claude', agent_display: 'Claude Code', skill_name: 'react-helper', description: 'React utils', skill_path: '/skills/react-helper', enabled: true },
      ]),
    },
    skills: { list: vi.fn().mockResolvedValue([]) },
    agents: { list: vi.fn().mockResolvedValue([]) },
  },
}))

vi.mock('../../lib/toast', () => ({
  toast: { success: vi.fn(), error: vi.fn() },
}))

describe('ProjectWorkspace', () => {
  it('renders project list with projects', async () => {
    renderWithRouter(<ProjectWorkspace />)
    expect(screen.getByText('Project Workspace')).toBeInTheDocument()
    expect(await screen.findByText('frontend-app')).toBeInTheDocument()
    expect(screen.getByText('backend-svc')).toBeInTheDocument()
  })

  it('shows empty state when no projects', async () => {
    const { api } = await import('../../lib/api')
    vi.mocked(api.projects.list).mockResolvedValueOnce([])
    renderWithRouter(<ProjectWorkspace />)
    expect(await screen.findByText('No projects registered')).toBeInTheDocument()
    expect(screen.getByText('Add a project directory to manage its skills')).toBeInTheDocument()
  })

  it('renders project detail heading', async () => {
    renderWithRouter(
      <Routes>
        <Route path="/projects/:id" element={<ProjectWorkspace />} />
      </Routes>,
      { route: '/projects/p1' },
    )
    expect(await screen.findByText('frontend-app')).toBeInTheDocument()
    expect(screen.getByText('Back to Projects')).toBeInTheDocument()
  })
})
