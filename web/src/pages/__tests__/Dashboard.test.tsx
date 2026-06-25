import { describe, it, expect, vi } from 'vitest'
import { screen } from '@testing-library/react'
import Dashboard from '../Dashboard'
import { renderWithRouter } from '../../test/utils'

vi.mock('../../lib/api', () => ({
  api: {
    skills: { list: vi.fn().mockResolvedValue([{ ID: '1' }, { ID: '2' }, { ID: '3' }]) },
    groups: { list: vi.fn().mockResolvedValue([{ id: '1' }]) },
    agents: { list: vi.fn().mockResolvedValue([{ detected: true }, { detected: false }]) },
    sync: { status: vi.fn().mockResolvedValue({ synced: 5, stale: 1, total: 6 }) },
  },
}))

describe('Dashboard', () => {
  it('renders all stat cards', () => {
    renderWithRouter(<Dashboard />)
    expect(screen.getByText('Skills')).toBeInTheDocument()
    expect(screen.getByText('Groups')).toBeInTheDocument()
    expect(screen.getByText('Agents')).toBeInTheDocument()
    expect(screen.getByText('Synced')).toBeInTheDocument()
    expect(screen.getByText('Stale')).toBeInTheDocument()
  })

  it('displays correct counts from API', async () => {
    renderWithRouter(<Dashboard />)
    expect(await screen.findByText('3')).toBeInTheDocument()
    expect(await screen.findByText('5')).toBeInTheDocument()
  })

  it('renders quick action buttons', () => {
    renderWithRouter(<Dashboard />)
    expect(screen.getByText('Sync All')).toBeInTheDocument()
    expect(screen.getByText('Install Skill')).toBeInTheDocument()
    expect(screen.getByText('Create Group')).toBeInTheDocument()
  })

  it('renders page title', () => {
    renderWithRouter(<Dashboard />)
    expect(screen.getByText('Dashboard')).toBeInTheDocument()
  })
})
