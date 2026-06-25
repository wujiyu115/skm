import { describe, it, expect, vi, beforeEach } from 'vitest'
import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import AuditLog from '../AuditLog'
import { renderWithRouter } from '../../test/utils'

const mockList = vi.fn()
const mockPrune = vi.fn()

vi.mock('../../lib/api', () => ({
  api: {
    audit: {
      list: (...args: unknown[]) => mockList(...args),
      prune: (...args: unknown[]) => mockPrune(...args),
    },
  },
}))

vi.mock('../../lib/toast', () => ({
  toast: {
    success: vi.fn(),
    error: vi.fn(),
  },
}))

const sampleEntries = [
  { id: 1, action: 'install', target: 'my-skill', detail: 'Installed from git', created_at: new Date().toISOString() },
  { id: 2, action: 'delete', target: 'old-skill', detail: 'Removed by user', created_at: new Date(Date.now() - 3600000).toISOString() },
  { id: 3, action: 'enable', target: 'another-skill', detail: 'Enabled', created_at: new Date(Date.now() - 86400000).toISOString() },
]

beforeEach(() => {
  mockList.mockReset()
  mockPrune.mockReset()
  mockList.mockResolvedValue(sampleEntries)
  mockPrune.mockResolvedValue({ ok: true, pruned: true })
})

describe('AuditLog', () => {
  it('renders page title', async () => {
    renderWithRouter(<AuditLog />)
    expect(await screen.findByText('Audit Log')).toBeInTheDocument()
  })

  it('renders audit entries in table', async () => {
    renderWithRouter(<AuditLog />)
    expect(await screen.findByText('my-skill')).toBeInTheDocument()
    expect(screen.getByText('old-skill')).toBeInTheDocument()
    expect(screen.getByText('another-skill')).toBeInTheDocument()
    expect(screen.getByText('install')).toBeInTheDocument()
    expect(screen.getByText('delete')).toBeInTheDocument()
    expect(screen.getByText('enable')).toBeInTheDocument()
    expect(screen.getByTestId('audit-table')).toBeInTheDocument()
  })

  it('shows empty state when no entries', async () => {
    mockList.mockResolvedValue([])
    renderWithRouter(<AuditLog />)
    expect(await screen.findByText('No audit entries')).toBeInTheDocument()
    expect(screen.getByTestId('audit-empty')).toBeInTheDocument()
  })

  it('prune button calls API and refreshes', async () => {
    const user = userEvent.setup()
    renderWithRouter(<AuditLog />)
    const pruneBtn = await screen.findByText('Prune Old Entries')
    await user.click(pruneBtn)
    await waitFor(() => {
      expect(mockPrune).toHaveBeenCalled()
    })
    // list is called once on mount and once after prune
    expect(mockList).toHaveBeenCalledTimes(2)
  })
})
