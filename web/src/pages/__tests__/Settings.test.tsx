import { describe, it, expect, vi, beforeEach } from 'vitest'
import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import Settings from '../Settings'
import { renderWithRouter } from '../../test/utils'

const mockUpdate = vi.fn().mockResolvedValue({})

vi.mock('../../lib/api', () => ({
  api: {
    settings: {
      get: vi.fn().mockResolvedValue({
        skills_dir: '/home/user/.skm/skills',
        cache_dir: '/home/user/.skm/cache',
        sync_mode: 'symlink',
        theme: 'light',
        text_size: 'default',
        auto_update_interval: 'off',
      }),
      update: (...args: unknown[]) => mockUpdate(...args),
    },
  },
}))

vi.mock('../../lib/theme', () => ({
  useTheme: () => ({
    theme: 'light',
    setTheme: vi.fn(),
    toggle: vi.fn(),
  }),
}))

vi.mock('../../lib/toast', () => ({
  toast: {
    success: vi.fn(),
    error: vi.fn(),
  },
}))

beforeEach(() => {
  mockUpdate.mockClear()
})

describe('Settings', () => {
  it('renders settings page title', async () => {
    renderWithRouter(<Settings />)
    expect(await screen.findByText('Settings')).toBeInTheDocument()
  })

  it('renders all sections', async () => {
    renderWithRouter(<Settings />)
    expect(await screen.findByText('Storage')).toBeInTheDocument()
    expect(screen.getByText('Appearance')).toBeInTheDocument()
    expect(screen.getByText('Updates')).toBeInTheDocument()
    expect(screen.getByText('About')).toBeInTheDocument()
  })

  it('displays read-only settings values', async () => {
    renderWithRouter(<Settings />)
    const skillsDirs = await screen.findAllByText('/home/user/.skm/skills')
    expect(skillsDirs.length).toBeGreaterThanOrEqual(1)
    const cacheDirs = screen.getAllByText('/home/user/.skm/cache')
    expect(cacheDirs.length).toBeGreaterThanOrEqual(1)
  })

  it('shows sync mode buttons', async () => {
    renderWithRouter(<Settings />)
    expect(await screen.findByText('Symlink')).toBeInTheDocument()
    expect(screen.getByText('Copy')).toBeInTheDocument()
  })

  it('calls API when sync mode is toggled', async () => {
    const user = userEvent.setup()
    renderWithRouter(<Settings />)
    const copyBtn = await screen.findByText('Copy')
    await user.click(copyBtn)
    await waitFor(() => {
      expect(mockUpdate).toHaveBeenCalledWith({ sync_mode: 'copy' })
    })
  })

  it('shows theme select with options', async () => {
    renderWithRouter(<Settings />)
    await screen.findByText('Storage')
    const themeSelect = screen.getByDisplayValue('Light')
    expect(themeSelect).toBeInTheDocument()
  })

  it('shows language select', async () => {
    renderWithRouter(<Settings />)
    await screen.findByText('Storage')
    const langSelect = screen.getByDisplayValue('English')
    expect(langSelect).toBeInTheDocument()
  })

  it('shows auto update interval select', async () => {
    renderWithRouter(<Settings />)
    await screen.findByText('Storage')
    const updateSelect = screen.getByDisplayValue('Off')
    expect(updateSelect).toBeInTheDocument()
  })

  it('shows version in about section', async () => {
    renderWithRouter(<Settings />)
    expect(await screen.findByText('0.1.0')).toBeInTheDocument()
  })
})
