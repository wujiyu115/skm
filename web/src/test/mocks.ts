import type { Skill, Group, Agent, Target, Project, ProjectSkill } from '../lib/api'

export function mockSkill(overrides?: Partial<Skill>): Skill {
  return {
    ID: 'skill-1',
    Name: 'test-skill',
    Description: 'A test skill for unit testing',
    SourceType: 'git',
    SourceRef: 'https://github.com/test/repo',
    CentralPath: '/home/user/.skm/skills/test-skill',
    ContentHash: 'abc123',
    Enabled: true,
    targets: [],
    ...overrides,
  }
}

export function mockTarget(overrides?: Partial<Target>): Target {
  return {
    skill_id: 'skill-1',
    agent: 'claude',
    target_path: '/home/user/.claude/skills/test-skill',
    mode: 'symlink',
    ...overrides,
  }
}

export function mockGroup(overrides?: Partial<Group>): Group {
  return {
    id: 'group-1',
    name: 'test-group',
    description: 'A test group',
    skill_count: 3,
    ...overrides,
  }
}

export function mockAgent(overrides?: Partial<Agent>): Agent {
  return {
    name: 'claude',
    display_name: 'Claude Code',
    project_dir: '.claude/skills',
    global_dir: '.claude/skills',
    detected: true,
    ...overrides,
  }
}

export function mockProject(overrides?: Partial<Project>): Project {
  return {
    id: 'proj-1',
    name: 'my-project',
    path: '/home/user/projects/my-project',
    created_at: '2026-01-15T10:00:00Z',
    ...overrides,
  }
}

export function mockProjectSkill(overrides?: Partial<ProjectSkill>): ProjectSkill {
  return {
    agent: 'claude',
    agent_display: 'Claude Code',
    skill_name: 'test-skill',
    description: 'A test skill',
    skill_path: '/home/user/projects/my-project/.claude/skills/test-skill',
    enabled: true,
    ...overrides,
  }
}
