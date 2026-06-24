import { Routes, Route, Link, useLocation } from 'react-router-dom'
import Dashboard from './pages/Dashboard'
import SkillsLibrary from './pages/SkillsLibrary'
import Groups from './pages/Groups'
import AgentWorkspace from './pages/AgentWorkspace'
import Settings from './pages/Settings'

const navItems = [
  { path: '/', label: 'Dashboard' },
  { path: '/skills', label: 'Skills' },
  { path: '/groups', label: 'Groups' },
  { path: '/agents', label: 'Agents' },
  { path: '/settings', label: 'Settings' },
]

export default function App() {
  const location = useLocation()

  return (
    <div className="flex h-screen bg-gray-50">
      <nav className="w-56 bg-white border-r border-gray-200 p-4">
        <h1 className="text-xl font-bold mb-6 text-gray-900">SKM</h1>
        <ul className="space-y-1">
          {navItems.map(item => (
            <li key={item.path}>
              <Link
                to={item.path}
                className={`block px-3 py-2 rounded-md text-sm ${
                  location.pathname === item.path
                    ? 'bg-blue-50 text-blue-700 font-medium'
                    : 'text-gray-600 hover:bg-gray-50'
                }`}
              >
                {item.label}
              </Link>
            </li>
          ))}
        </ul>
      </nav>
      <main className="flex-1 overflow-auto p-6">
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/skills" element={<SkillsLibrary />} />
          <Route path="/groups" element={<Groups />} />
          <Route path="/agents" element={<AgentWorkspace />} />
          <Route path="/settings" element={<Settings />} />
        </Routes>
      </main>
    </div>
  )
}
