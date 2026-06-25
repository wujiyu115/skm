import { describe, it, expect, vi } from 'vitest'
import { screen } from '@testing-library/react'
import Settings from '../Settings'
import { renderWithRouter } from '../../test/utils'

vi.mock('../../lib/api', () => ({
  api: {
    settings: {
      get: vi.fn().mockResolvedValue({
        skills_dir: '/home/user/.skm/skills',
        cache_dir: '/home/user/.skm/cache',
        sync_mode: 'symlink',
      }),
    },
  },
}))

describe('Settings', () => {
  it('renders settings page title', () => {
    renderWithRouter(<Settings />)
    expect(screen.getByText('Settings')).toBeInTheDocument()
  })

  it('displays settings values', async () => {
    renderWithRouter(<Settings />)
    expect(await screen.findByText('/home/user/.skm/skills')).toBeInTheDocument()
    expect(screen.getByText('/home/user/.skm/cache')).toBeInTheDocument()
  })

  it('groups dir settings under Storage', async () => {
    renderWithRouter(<Settings />)
    expect(await screen.findByText('Storage')).toBeInTheDocument()
  })

  it('shows other settings in separate section', async () => {
    renderWithRouter(<Settings />)
    expect(await screen.findByText('Other')).toBeInTheDocument()
    expect(screen.getByText('symlink')).toBeInTheDocument()
  })
})
