import { Routes, Route } from 'react-router-dom'
import Sidebar from './components/Sidebar'
import Dashboard from './pages/Dashboard'
import SkillsLibrary from './pages/SkillsLibrary'
import Install from './pages/Install'
import Groups from './pages/Groups'
import AgentWorkspace from './pages/AgentWorkspace'
import Settings from './pages/Settings'

export default function App() {
  return (
    <div className="flex h-screen">
      <Sidebar />
      <main className="flex-1 overflow-auto bg-slate-50 dark:bg-slate-900 p-8">
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/skills" element={<SkillsLibrary />} />
          <Route path="/install" element={<Install />} />
          <Route path="/groups" element={<Groups />} />
          <Route path="/groups/:id" element={<Groups />} />
          <Route path="/agents" element={<AgentWorkspace />} />
          <Route path="/agents/:name" element={<AgentWorkspace />} />
          <Route path="/settings" element={<Settings />} />
        </Routes>
      </main>
    </div>
  )
}
