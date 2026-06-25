import { describe, it, expect, vi } from 'vitest'
import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import BatchToolbar from '../BatchToolbar'
import { renderWithRouter } from '../../test/utils'

describe('BatchToolbar', () => {
  it('renders nothing when no items are selected', () => {
    const { container } = renderWithRouter(
      <BatchToolbar selectedIds={[]} onClear={vi.fn()} onAction={vi.fn()} />
    )
    expect(container.innerHTML).toBe('')
  })

  it('shows toolbar when items are selected', () => {
    renderWithRouter(
      <BatchToolbar selectedIds={['a', 'b']} onClear={vi.fn()} onAction={vi.fn()} />
    )
    expect(screen.getByTestId('batch-toolbar')).toBeInTheDocument()
  })

  it('displays correct count of selected items', () => {
    renderWithRouter(
      <BatchToolbar selectedIds={['a', 'b', 'c']} onClear={vi.fn()} onAction={vi.fn()} />
    )
    expect(screen.getByText(/3\s+selected/)).toBeInTheDocument()
  })

  it('calls onClear when close button is clicked', async () => {
    const onClear = vi.fn()
    renderWithRouter(
      <BatchToolbar selectedIds={['a']} onClear={onClear} onAction={vi.fn()} />
    )
    await userEvent.click(screen.getByLabelText('Cancel'))
    expect(onClear).toHaveBeenCalledOnce()
  })

  it('calls onAction with "enable" when enable button clicked', async () => {
    const onAction = vi.fn()
    renderWithRouter(
      <BatchToolbar selectedIds={['a']} onClear={vi.fn()} onAction={onAction} />
    )
    await userEvent.click(screen.getByText('Enable Selected'))
    expect(onAction).toHaveBeenCalledWith('enable')
  })

  it('calls onAction with "disable" when disable button clicked', async () => {
    const onAction = vi.fn()
    renderWithRouter(
      <BatchToolbar selectedIds={['a']} onClear={vi.fn()} onAction={onAction} />
    )
    await userEvent.click(screen.getByText('Disable Selected'))
    expect(onAction).toHaveBeenCalledWith('disable')
  })

  it('calls onAction with "sync" when sync button clicked', async () => {
    const onAction = vi.fn()
    renderWithRouter(
      <BatchToolbar selectedIds={['a']} onClear={vi.fn()} onAction={onAction} />
    )
    await userEvent.click(screen.getByText('Sync Selected'))
    expect(onAction).toHaveBeenCalledWith('sync')
  })

  it('calls onAction with "delete" after confirmation', async () => {
    const onAction = vi.fn()
    vi.spyOn(window, 'confirm').mockReturnValue(true)
    renderWithRouter(
      <BatchToolbar selectedIds={['a']} onClear={vi.fn()} onAction={onAction} />
    )
    await userEvent.click(screen.getByText('Delete Selected'))
    expect(window.confirm).toHaveBeenCalled()
    expect(onAction).toHaveBeenCalledWith('delete')
    vi.restoreAllMocks()
  })

  it('does not call onAction for delete when confirmation is cancelled', async () => {
    const onAction = vi.fn()
    vi.spyOn(window, 'confirm').mockReturnValue(false)
    renderWithRouter(
      <BatchToolbar selectedIds={['a']} onClear={vi.fn()} onAction={onAction} />
    )
    await userEvent.click(screen.getByText('Delete Selected'))
    expect(window.confirm).toHaveBeenCalled()
    expect(onAction).not.toHaveBeenCalled()
    vi.restoreAllMocks()
  })
})
