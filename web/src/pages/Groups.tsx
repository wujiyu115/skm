import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { FolderOpen, Plus, ArrowLeft, Trash2, X } from 'lucide-react'
import { api, type Group, type Skill } from '../lib/api'
import { useI18n } from '../lib/i18n'

export default function Groups() {
  const { id } = useParams<{ id: string }>()
  return id ? <GroupDetail id={id} /> : <GroupList />
}

function GroupList() {
  const { t } = useI18n()
  const [groups, setGroups] = useState<Group[]>([])
  const [name, setName] = useState('')
  const [desc, setDesc] = useState('')
  const [showForm, setShowForm] = useState(false)
  const navigate = useNavigate()

  const load = () => {
    api.groups.list().then(setGroups).catch(() => {})
  }
  useEffect(() => { load() }, [])

  const create = async () => {
    if (!name.trim()) return
    await api.groups.create(name, desc)
    setName('')
    setDesc('')
    setShowForm(false)
    load()
  }

  const remove = async (id: string) => {
    await api.groups.remove(id)
    load()
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <h2 className="text-2xl font-bold text-slate-900">{t('groups.title')}</h2>
          <span className="px-2.5 py-0.5 bg-purple-100 text-purple-700 rounded-full text-sm font-medium">
            {groups.length}
          </span>
        </div>
        <button
          onClick={() => setShowForm(!showForm)}
          className="flex items-center gap-2 px-4 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 transition-colors"
        >
          <Plus className="w-4 h-4" /> {t('groups.create')}
        </button>
      </div>

      {showForm && (
        <div className="bg-white rounded-xl border border-slate-200 p-5 mb-6">
          <h3 className="font-semibold text-slate-900 mb-3">{t('groups.new')}</h3>
          <div className="space-y-3">
            <input
              value={name}
              onChange={e => setName(e.target.value)}
              placeholder={t('groups.name')}
              className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
              onKeyDown={e => e.key === 'Enter' && create()}
            />
            <input
              value={desc}
              onChange={e => setDesc(e.target.value)}
              placeholder={t('groups.desc')}
              className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
            />
            <div className="flex gap-2">
              <button onClick={create} className="px-4 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700">
                {t('groups.create')}
              </button>
              <button onClick={() => setShowForm(false)} className="px-4 py-2 text-slate-600 hover:bg-slate-100 rounded-lg text-sm">
                {t('groups.cancel')}
              </button>
            </div>
          </div>
        </div>
      )}

      {groups.length === 0 ? (
        <div className="text-center py-12 text-slate-500">
          <FolderOpen className="w-12 h-12 mx-auto mb-3 text-slate-300" />
          <p className="text-lg">{t('groups.noGroups')}</p>
          <p className="text-sm mt-1">{t('groups.noGroupsHint')}</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
          {groups.map(g => (
            <div
              key={g.id}
              className="bg-white rounded-xl border border-slate-200 p-5 hover:shadow-md transition-shadow cursor-pointer"
              onClick={() => navigate(`/groups/${g.id}`)}
            >
              <div className="flex items-start justify-between">
                <div className="flex items-center gap-2">
                  <FolderOpen className="w-5 h-5 text-purple-500" />
                  <h3 className="font-semibold text-slate-900">{g.name}</h3>
                </div>
                <button
                  onClick={e => { e.stopPropagation(); remove(g.id) }}
                  className="text-slate-400 hover:text-red-500 transition-colors"
                >
                  <Trash2 className="w-4 h-4" />
                </button>
              </div>
              {g.description && (
                <p className="text-sm text-slate-500 mt-2">{g.description}</p>
              )}
              <div className="mt-3">
                <span className="text-xs text-slate-500">{g.skill_count ?? 0} {t('groups.skills')}</span>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

function GroupDetail({ id }: { id: string }) {
  const { t } = useI18n()
  const [group, setGroup] = useState<Group | null>(null)
  const [skills, setSkills] = useState<Skill[]>([])
  const navigate = useNavigate()

  useEffect(() => {
    api.groups.get(id).then(data => {
      setGroup(data.group)
      setSkills(data.skills ?? [])
    }).catch(() => navigate('/groups'))
  }, [id, navigate])

  const removeSkill = async (skillId: string) => {
    await api.groups.removeSkill(id, skillId)
    const data = await api.groups.get(id)
    setSkills(data.skills ?? [])
  }

  if (!group) return null

  return (
    <div>
      <button
        onClick={() => navigate('/groups')}
        className="flex items-center gap-1.5 text-sm text-slate-500 hover:text-slate-700 mb-4"
      >
        <ArrowLeft className="w-4 h-4" /> {t('groups.back')}
      </button>

      <div className="flex items-center gap-3 mb-6">
        <FolderOpen className="w-6 h-6 text-purple-500" />
        <h2 className="text-2xl font-bold text-slate-900">{group.name}</h2>
        <span className="px-2.5 py-0.5 bg-purple-100 text-purple-700 rounded-full text-sm font-medium">
          {skills.length} {t('groups.skills')}
        </span>
      </div>

      {group.description && (
        <p className="text-slate-500 mb-6">{group.description}</p>
      )}

      {skills.length === 0 ? (
        <div className="text-center py-12 text-slate-500">
          <p>{t('groups.noSkillsInGroup')}</p>
          <p className="text-sm mt-1">{t('groups.addHint')} {group.name} &lt;skill&gt;</p>
        </div>
      ) : (
        <div className="space-y-2">
          {skills.map(sk => (
            <div key={sk.ID} className="flex items-center justify-between bg-white rounded-lg border border-slate-200 px-4 py-3">
              <div>
                <span className="font-medium text-slate-900">{sk.Name}</span>
                <span className="text-sm text-slate-500 ml-2">{sk.Description}</span>
              </div>
              <button
                onClick={() => removeSkill(sk.ID)}
                className="text-slate-400 hover:text-red-500 transition-colors"
              >
                <X className="w-4 h-4" />
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
