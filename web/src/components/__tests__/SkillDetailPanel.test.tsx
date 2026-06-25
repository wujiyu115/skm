import { describe, it, expect, vi, beforeEach } from 'vitest'
import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import SkillDetailPanel from '../SkillDetailPanel'
import { renderWithRouter } from '../../test/utils'
import { api } from '../../lib/api'

vi.mock('../../lib/api', async () => {
  const actual = await vi.importActual('../../lib/api')
  return {
    ...actual,
    api: {
      skills: {
        get: vi.fn(),
        content: vi.fn(),
        setEnabled: vi.fn(),
      },
      tags: {
        getForSkill: vi.fn(),
      },
    },
  }
})

const mockSkillResponse = {
  skill: {
    ID: 'skill-1',
    Name: 'test-skill',
    Description: 'A test skill',
    SourceType: 'git',
    SourceRef: 'https://github.com/test/repo',
    CentralPath: '/path/to/skill',
    ContentHash: 'abc123',
    Enabled: true,
  },
  targets: [
    { skill_id: 'skill-1', agent: 'claude', target_path: '/path', mode: 'symlink' },
  ],
}

beforeEach(() => {
  vi.clearAllMocks()
  vi.mocked(api.skills.get).mockResolvedValue(mockSkillResponse)
  vi.mocked(api.skills.content).mockResolvedValue({ content: '# Hello\nSome content here.' })
  vi.mocked(api.tags.getForSkill).mockResolvedValue(['frontend', 'react'])
})

describe('SkillDetailPanel', () => {
  it('does not show panel content when skillId is null', () => {
    renderWithRouter(<SkillDetailPanel skillId={null} onClose={() => {}} />)
    expect(screen.queryByText('test-skill')).not.toBeInTheDocument()
  })

  it('renders panel with skill details when skillId is set', async () => {
    renderWithRouter(<SkillDetailPanel skillId="skill-1" onClose={() => {}} />)

    await waitFor(() => {
      expect(screen.getByText('test-skill')).toBeInTheDocument()
    })
    expect(screen.getByText('A test skill')).toBeInTheDocument()
    expect(screen.getByText('Enabled')).toBeInTheDocument()
    expect(screen.getByText('git')).toBeInTheDocument()
    expect(screen.getByText('frontend')).toBeInTheDocument()
    expect(screen.getByText('react')).toBeInTheDocument()
    expect(screen.getByText('claude')).toBeInTheDocument()
  })

  it('renders markdown content', async () => {
    renderWithRouter(<SkillDetailPanel skillId="skill-1" onClose={() => {}} />)

    await waitFor(() => {
      expect(screen.getByText('Hello')).toBeInTheDocument()
    })
    expect(screen.getByText('Some content here.')).toBeInTheDocument()
  })

  it('shows "No content available" when content is empty', async () => {
    vi.mocked(api.skills.content).mockResolvedValue({ content: '' })

    renderWithRouter(<SkillDetailPanel skillId="skill-1" onClose={() => {}} />)

    await waitFor(() => {
      expect(screen.getByText('No content available')).toBeInTheDocument()
    })
  })

  it('shows loading spinner then content', async () => {
    let resolveGet: (v: typeof mockSkillResponse) => void
    vi.mocked(api.skills.get).mockReturnValue(
      new Promise(r => { resolveGet = r })
    )

    renderWithRouter(<SkillDetailPanel skillId="skill-1" onClose={() => {}} />)
    expect(screen.getByTestId('loading-spinner')).toBeInTheDocument()

    resolveGet!(mockSkillResponse)

    await waitFor(() => {
      expect(screen.getByText('test-skill')).toBeInTheDocument()
    })
    expect(screen.queryByTestId('loading-spinner')).not.toBeInTheDocument()
  })

  it('calls onClose when close button is clicked', async () => {
    const onClose = vi.fn()
    renderWithRouter(<SkillDetailPanel skillId="skill-1" onClose={onClose} />)

    await waitFor(() => {
      expect(screen.getByText('test-skill')).toBeInTheDocument()
    })

    await userEvent.click(screen.getByTestId('close-button'))
    expect(onClose).toHaveBeenCalledOnce()
  })

  it('calls onClose when overlay is clicked', async () => {
    const onClose = vi.fn()
    renderWithRouter(<SkillDetailPanel skillId="skill-1" onClose={onClose} />)

    await waitFor(() => {
      expect(screen.getByText('test-skill')).toBeInTheDocument()
    })

    await userEvent.click(screen.getByTestId('detail-overlay'))
    expect(onClose).toHaveBeenCalledOnce()
  })

  it('calls onClose when Escape key is pressed', async () => {
    const onClose = vi.fn()
    renderWithRouter(<SkillDetailPanel skillId="skill-1" onClose={onClose} />)

    await waitFor(() => {
      expect(screen.getByText('test-skill')).toBeInTheDocument()
    })

    await userEvent.keyboard('{Escape}')
    expect(onClose).toHaveBeenCalledOnce()
  })
})
