import { describe, it, expect, vi } from 'vitest'
import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import CommandPalette from '../CommandPalette'
import { renderWithRouter } from '../../test/utils'

vi.mock('../../lib/api', () => ({
  api: {
    skills: {
      list: vi.fn().mockResolvedValue([
        { ID: 's1', Name: 'React Helper', Description: 'Helps with React', Enabled: true },
        { ID: 's2', Name: 'Git Workflow', Description: 'Git best practices', Enabled: true },
      ]),
    },
    groups: {
      list: vi.fn().mockResolvedValue([
        { id: 'g1', name: 'Frontend', description: 'Frontend tools', skill_count: 3 },
        { id: 'g2', name: 'DevOps', description: 'DevOps tools', skill_count: 1 },
      ]),
    },
  },
}))

describe('CommandPalette', () => {
  it('renders nothing initially (palette is closed)', () => {
    renderWithRouter(<CommandPalette />)
    expect(screen.queryByTestId('command-palette-overlay')).not.toBeInTheDocument()
  })

  it('opens on Ctrl+K keypress', async () => {
    renderWithRouter(<CommandPalette />)
    await userEvent.keyboard('{Control>}k{/Control}')
    expect(screen.getByTestId('command-palette-overlay')).toBeInTheDocument()
  })

  it('shows search input when open', async () => {
    renderWithRouter(<CommandPalette />)
    await userEvent.keyboard('{Control>}k{/Control}')
    expect(screen.getByTestId('command-palette-input')).toBeInTheDocument()
    expect(screen.getByPlaceholderText('Search skills, groups, pages...')).toBeInTheDocument()
  })

  it('closes on Escape', async () => {
    renderWithRouter(<CommandPalette />)
    await userEvent.keyboard('{Control>}k{/Control}')
    expect(screen.getByTestId('command-palette-overlay')).toBeInTheDocument()
    await userEvent.keyboard('{Escape}')
    expect(screen.queryByTestId('command-palette-overlay')).not.toBeInTheDocument()
  })

  it('shows navigation pages in results', async () => {
    renderWithRouter(<CommandPalette />)
    await userEvent.keyboard('{Control>}k{/Control}')
    await waitFor(() => {
      expect(screen.getByText('Dashboard')).toBeInTheDocument()
      expect(screen.getByText('Library')).toBeInTheDocument()
      expect(screen.getByText('Install')).toBeInTheDocument()
      expect(screen.getByText('Settings')).toBeInTheDocument()
      expect(screen.getByText('Audit Log')).toBeInTheDocument()
    })
  })

  it('filters results based on search input', async () => {
    const user = userEvent.setup()
    renderWithRouter(<CommandPalette />)
    await user.keyboard('{Control>}k{/Control}')

    // Wait for data to load
    await waitFor(() => {
      expect(screen.getByText('React Helper')).toBeInTheDocument()
    })

    // Type a filter query
    const input = screen.getByTestId('command-palette-input')
    await user.type(input, 'react')

    // React Helper skill should be visible, Git Workflow should not
    expect(screen.getByText('React Helper')).toBeInTheDocument()
    expect(screen.queryByText('Git Workflow')).not.toBeInTheDocument()
    // Page names that don't match should be hidden
    expect(screen.queryByText('Dashboard')).not.toBeInTheDocument()
  })

  it('shows skills and groups from API', async () => {
    renderWithRouter(<CommandPalette />)
    await userEvent.keyboard('{Control>}k{/Control}')

    await waitFor(() => {
      expect(screen.getByText('React Helper')).toBeInTheDocument()
      expect(screen.getByText('Git Workflow')).toBeInTheDocument()
      expect(screen.getByText('Frontend')).toBeInTheDocument()
      expect(screen.getByText('DevOps')).toBeInTheDocument()
    })
  })

  it('closes when overlay is clicked', async () => {
    renderWithRouter(<CommandPalette />)
    await userEvent.keyboard('{Control>}k{/Control}')
    expect(screen.getByTestId('command-palette-overlay')).toBeInTheDocument()

    await userEvent.click(screen.getByTestId('command-palette-overlay'))
    expect(screen.queryByTestId('command-palette-overlay')).not.toBeInTheDocument()
  })
})
