import { render, type RenderOptions } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { I18nProvider } from '../lib/i18n'
import type { ReactElement } from 'react'

export function renderWithRouter(
  ui: ReactElement,
  { route = '/', ...options }: RenderOptions & { route?: string } = {},
) {
  return render(ui, {
    wrapper: ({ children }) => (
      <I18nProvider>
        <MemoryRouter initialEntries={[route]}>{children}</MemoryRouter>
      </I18nProvider>
    ),
    ...options,
  })
}
