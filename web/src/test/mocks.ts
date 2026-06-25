import type { Skill, Group, Agent, Target } from '../lib/api'

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
