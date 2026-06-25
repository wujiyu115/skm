import { createContext, useContext, useState, useCallback, useEffect, type ReactNode } from 'react'

export type Theme = 'light' | 'dark'

interface ThemeContextType {
  theme: Theme
  setTheme: (theme: Theme) => void
  toggle: () => void
}

const ThemeContext = createContext<ThemeContextType>({
  theme: 'light',
  setTheme: () => {},
  toggle: () => {},
})

export function ThemeProvider({ children }: { children: ReactNode }) {
  const [theme, setThemeState] = useState<Theme>(() => {
    const saved = localStorage.getItem('skm-theme')
    if (saved === 'dark' || saved === 'light') return saved
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
  })

  useEffect(() => {
    document.documentElement.classList.toggle('dark', theme === 'dark')
  }, [theme])

  const setTheme = useCallback((t: Theme) => {
    setThemeState(t)
    localStorage.setItem('skm-theme', t)
  }, [])

  const toggle = useCallback(() => {
    setTheme(theme === 'light' ? 'dark' : 'light')
  }, [theme, setTheme])

  return (
    <ThemeContext.Provider value={{ theme, setTheme, toggle }}>
      {children}
    </ThemeContext.Provider>
  )
}

export function useTheme() {
  return useContext(ThemeContext)
}
