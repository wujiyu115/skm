const BASE = '/api'

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error(err.error || res.statusText)
  }
  return res.json()
}

export interface Skill {
  ID: string
  Name: string
  Description: string
  SourceType: string
  SourceRef: string
  CentralPath: string
  ContentHash: string
  Enabled: boolean
  targets?: Target[]
  tags?: string[]
}

export interface Target {
  skill_id: string
  agent: string
  target_path: string
  mode: string
}

export interface Group {
  id: string
  name: string
  description: string
  skill_count?: number
}

export interface Agent {
  name: string
  display_name: string
  project_dir: string
  global_dir: string
  detected: boolean
}

export const api = {
  skills: {
    list: () => request<Skill[]>('/skills').then(r => r ?? []),
    get: (id: string) => request<{ skill: Skill; targets: Target[] }>(`/skills/${id}`),
    install: (source: string, agents: string[], global: boolean) =>
      request<{ installed: string[] }>('/skills/install', {
        method: 'POST',
        body: JSON.stringify({ source, agents, global }),
      }),
    remove: (id: string) => request(`/skills/${id}`, { method: 'DELETE' }),
    setEnabled: (id: string, enabled: boolean) =>
      request<{ ok: boolean }>(`/skills/${id}/enable`, {
        method: 'PUT',
        body: JSON.stringify({ enabled }),
      }),
    content: (id: string) => request<{ content: string }>(`/skills/${id}/content`),
    sync: (id: string, agents: string[]) =>
      request(`/skills/${id}/sync`, {
        method: 'POST',
        body: JSON.stringify({ agents }),
      }),
  },
  groups: {
    list: () => request<Group[]>('/groups').then(r => r ?? []),
    get: (id: string) => request<{ group: Group; skills: Skill[] }>(`/groups/${id}`),
    create: (name: string, description: string) =>
      request<{ id: string }>('/groups', {
        method: 'POST',
        body: JSON.stringify({ name, description }),
      }),
    update: (id: string, name: string, description: string) =>
      request(`/groups/${id}`, {
        method: 'PUT',
        body: JSON.stringify({ name, description }),
      }),
    remove: (id: string) => request(`/groups/${id}`, { method: 'DELETE' }),
    addSkills: (id: string, skillIds: string[]) =>
      request(`/groups/${id}/skills`, {
        method: 'POST',
        body: JSON.stringify({ skill_ids: skillIds }),
      }),
    removeSkill: (id: string, skillId: string) =>
      request(`/groups/${id}/skills/${skillId}`, { method: 'DELETE' }),
    install: (id: string, agents: string[]) =>
      request(`/groups/${id}/install`, {
        method: 'POST',
        body: JSON.stringify({ agents }),
      }),
  },
  agents: {
    list: () => request<Agent[]>('/agents').then(r => r ?? []),
  },
  sync: {
    status: () => request<{ total: number; synced: number; stale: number }>('/sync/status'),
    trigger: (agents: string[]) =>
      request('/sync', {
        method: 'POST',
        body: JSON.stringify({ agents }),
      }),
  },
  batch: {
    delete: (ids: string[]) =>
      request<{ ok: boolean; processed: number; errors: string[] }>('/skills/batch/delete', {
        method: 'POST',
        body: JSON.stringify({ ids }),
      }),
    enable: (ids: string[], enabled: boolean) =>
      request<{ ok: boolean; processed: number; errors: string[] }>('/skills/batch/enable', {
        method: 'POST',
        body: JSON.stringify({ ids, enabled }),
      }),
    tag: (ids: string[], tags: string[], action: 'add' | 'remove') =>
      request<{ ok: boolean; processed: number; errors: string[] }>('/skills/batch/tag', {
        method: 'POST',
        body: JSON.stringify({ ids, tags, action }),
      }),
    sync: (ids: string[], agents: string[]) =>
      request<{ ok: boolean; processed: number; errors: string[] }>('/skills/batch/sync', {
        method: 'POST',
        body: JSON.stringify({ ids, agents }),
      }),
  },
  tags: {
    list: () => request<string[]>('/tags').then(r => r ?? []),
    getForSkill: (skillId: string) => request<string[]>(`/skills/${skillId}/tags`).then(r => r ?? []),
    setForSkill: (skillId: string, tags: string[]) =>
      request(`/skills/${skillId}/tags`, { method: 'PUT', body: JSON.stringify({ tags }) }),
  },
  settings: {
    get: () => request<Record<string, string>>('/settings'),
    update: (settings: Record<string, string>) =>
      request('/settings', { method: 'PUT', body: JSON.stringify(settings) }),
  },
  audit: {
    list: (limit = 100) =>
      request<Array<{ id: number; action: string; target: string; detail: string; created_at: string }>>(
        `/audit?limit=${limit}`,
      ).then(r => r ?? []),
    prune: () => request<{ ok: boolean; pruned: boolean }>('/audit', { method: 'DELETE' }),
  },
}
