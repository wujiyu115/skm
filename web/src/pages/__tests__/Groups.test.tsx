import { describe, it, expect, vi } from 'vitest'
import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import Groups from '../Groups'
import { renderWithRouter } from '../../test/utils'

vi.mock('../../lib/api', () => ({
  api: {
    groups: {
      list: vi.fn().mockResolvedValue([
        { id: 'g1', name: 'Frontend', description: 'React skills', skill_count: 3 },
        { id: 'g2', name: 'Backend', description: '', skill_count: 1 },
      ]),
      get: vi.fn().mockResolvedValue({
        group: { id: 'g1', name: 'Frontend', description: 'React skills' },
        skills: [
          { ID: 's1', Name: 'react-helper', Description: 'React utils' },
        ],
      }),
      create: vi.fn().mockResolvedValue({ id: 'g3' }),
      remove: vi.fn().mockResolvedValue({}),
      removeSkill: vi.fn().mockResolvedValue({}),
    },
  },
}))

describe('Groups', () => {
  it('renders groups list', async () => {
    renderWithRouter(<Groups />)
    expect(screen.getByText('Skill Groups')).toBeInTheDocument()
    expect(await screen.findByText('Frontend')).toBeInTheDocument()
    expect(screen.getByText('Backend')).toBeInTheDocument()
  })

  it('shows skill count per group', async () => {
    renderWithRouter(<Groups />)
    expect(await screen.findByText('3 skills')).toBeInTheDocument()
    expect(screen.getByText('1 skills')).toBeInTheDocument()
  })

  it('shows create form on button click', async () => {
    renderWithRouter(<Groups />)
    await userEvent.click(screen.getByText('Create Group'))
    expect(screen.getByPlaceholderText('Group name')).toBeInTheDocument()
  })

  it('renders group count badge', async () => {
    renderWithRouter(<Groups />)
    expect(await screen.findByText('2')).toBeInTheDocument()
  })
})
