import { createContext, useContext, useState, useCallback, type ReactNode } from 'react'

export type Locale = 'en' | 'zh'

const translations = {
  en: {
    'app.title': 'Skills Manager',
    'nav.dashboard': 'Dashboard',
    'nav.library': 'Library',
    'nav.install': 'Install Skills',
    'nav.settings': 'Settings',
    'nav.presets': 'Presets',
    'nav.workspace': 'Global Workspace',
    'nav.allAgents': 'All Agents',
    'nav.newGroup': '+ New Group',

    'dashboard.title': 'Dashboard',
    'dashboard.subtitle': 'Overview of your skills library',
    'dashboard.skills': 'Skills',
    'dashboard.groups': 'Groups',
    'dashboard.agents': 'Agents',
    'dashboard.synced': 'Synced',
    'dashboard.stale': 'Stale',
    'dashboard.quickActions': 'Quick Actions',
    'dashboard.syncAll': 'Sync All',
    'dashboard.installSkill': 'Install Skill',
    'dashboard.createGroup': 'Create Group',

    'skills.title': 'Library',
    'skills.search': 'Search skills in the control library...',
    'skills.install': 'Install',
    'skills.installing': 'Installing...',
    'skills.placeholder': 'owner/repo or URL...',
    'skills.all': 'All',
    'skills.enabled': 'Enabled',
    'skills.available': 'Available',
    'skills.disabled': 'Disabled',
    'skills.sync': 'Sync',
    'skills.checkAll': 'Check All',
    'skills.noSkills': 'No skills found',
    'skills.noSkillsHint': 'Install a skill using the input above',
    'skills.remove': 'remove',

    'groups.title': 'Skill Groups',
    'groups.create': 'Create Group',
    'groups.new': 'New Group',
    'groups.name': 'Group name',
    'groups.desc': 'Description (optional)',
    'groups.cancel': 'Cancel',
    'groups.noGroups': 'No groups yet',
    'groups.noGroupsHint': 'Create a group to organize your skills',
    'groups.skills': 'skills',
    'groups.back': 'Back to Groups',
    'groups.noSkillsInGroup': 'No skills in this group',
    'groups.addHint': 'Add skills using the CLI: skm group add',

    'agents.title': 'Agent Workspaces',
    'agents.active': 'Active',
    'agents.notDetected': 'Not detected',
    'agents.skillsSynced': 'skills synced',
    'agents.view': 'View →',
    'agents.back': 'Back to Agents',
    'agents.syncedSkills': 'Synced Skills',
    'agents.noSkills': 'No skills synced to this agent',

    'settings.title': 'Settings',
    'settings.storage': 'Storage',
    'settings.other': 'Other',
    'settings.noSettings': 'No settings configured',

    'install.title': 'Install Skills',
    'install.tabRegistry': 'Remote Registry',
    'install.tabLocal': 'Local Path',
    'install.tabGit': 'Git Install',
    'install.search': 'Search skills...',
    'install.filterAll': 'All',
    'install.filterOfficial': 'Official',
    'install.filterCommunity': 'Community',
    'install.filterSource': 'Source',
    'install.installBtn': 'Install',
    'install.installing': 'Installing...',
    'install.installed': 'Installed',
    'install.localTitle': 'Install from Local Path',
    'install.localDesc': 'Enter a local directory path containing skill files.',
    'install.localPlaceholder': '/path/to/skills/directory',
    'install.gitTitle': 'Install from Git Repository',
    'install.gitDesc': 'Enter a Git repository URL. Supports any Git-hosted repository containing skill definitions.',
    'install.gitPlaceholder': 'https://github.com/user/repo.git',
    'install.advancedOptions': 'Advanced Options',
    'install.agents': 'Target Agents',
    'install.global': 'Install Globally',
    'install.success': 'Skills installed successfully!',
    'install.error': 'Installation failed',
  },
  zh: {
    'app.title': '技能管理器',
    'nav.dashboard': '仪表盘',
    'nav.library': '技能库',
    'nav.install': '安装技能',
    'nav.settings': '设置',
    'nav.presets': '技能组',
    'nav.workspace': '全局工作区',
    'nav.allAgents': '所有智能体',
    'nav.newGroup': '+ 新建组',

    'dashboard.title': '仪表盘',
    'dashboard.subtitle': '技能库概览',
    'dashboard.skills': '技能',
    'dashboard.groups': '技能组',
    'dashboard.agents': '智能体',
    'dashboard.synced': '已同步',
    'dashboard.stale': '待更新',
    'dashboard.quickActions': '快捷操作',
    'dashboard.syncAll': '全部同步',
    'dashboard.installSkill': '安装技能',
    'dashboard.createGroup': '创建技能组',

    'skills.title': '技能库',
    'skills.search': '搜索技能...',
    'skills.install': '安装',
    'skills.installing': '安装中...',
    'skills.placeholder': 'owner/repo 或 URL...',
    'skills.all': '全部',
    'skills.enabled': '已启用',
    'skills.available': '可用',
    'skills.disabled': '已禁用',
    'skills.sync': '同步',
    'skills.checkAll': '全选',
    'skills.noSkills': '未找到技能',
    'skills.noSkillsHint': '使用上方输入框安装技能',
    'skills.remove': '删除',

    'groups.title': '技能组',
    'groups.create': '创建技能组',
    'groups.new': '新建技能组',
    'groups.name': '技能组名称',
    'groups.desc': '描述（可选）',
    'groups.cancel': '取消',
    'groups.noGroups': '暂无技能组',
    'groups.noGroupsHint': '创建技能组来组织你的技能',
    'groups.skills': '个技能',
    'groups.back': '返回技能组',
    'groups.noSkillsInGroup': '该组暂无技能',
    'groups.addHint': '使用 CLI 添加技能: skm group add',

    'agents.title': '智能体工作区',
    'agents.active': '已检测',
    'agents.notDetected': '未检测到',
    'agents.skillsSynced': '个技能已同步',
    'agents.view': '查看 →',
    'agents.back': '返回智能体',
    'agents.syncedSkills': '已同步技能',
    'agents.noSkills': '该智能体暂无同步技能',

    'settings.title': '设置',
    'settings.storage': '存储',
    'settings.other': '其他',
    'settings.noSettings': '暂无设置',

    'install.title': '安装 Skills',
    'install.tabRegistry': '远程仓库',
    'install.tabLocal': '本地路径安装',
    'install.tabGit': 'Git 安装',
    'install.search': '搜索技能...',
    'install.filterAll': '全部',
    'install.filterOfficial': '官方推荐',
    'install.filterCommunity': '社区',
    'install.filterSource': '来源',
    'install.installBtn': '安装',
    'install.installing': '安装中...',
    'install.installed': '已安装',
    'install.localTitle': '从本地路径安装',
    'install.localDesc': '输入包含技能文件的本地目录路径。',
    'install.localPlaceholder': '/path/to/skills/directory',
    'install.gitTitle': '从 Git 仓库安装',
    'install.gitDesc': '输入 Git 仓库 URL，支持任何包含技能定义的 Git 托管仓库。',
    'install.gitPlaceholder': 'https://github.com/user/repo.git',
    'install.advancedOptions': '高级开关选项',
    'install.agents': '目标智能体',
    'install.global': '全局安装',
    'install.success': '技能安装成功！',
    'install.error': '安装失败',
  },
} as const

type TranslationKey = keyof typeof translations.en

interface I18nContextType {
  locale: Locale
  setLocale: (locale: Locale) => void
  t: (key: TranslationKey) => string
}

const I18nContext = createContext<I18nContextType>({
  locale: 'en',
  setLocale: () => {},
  t: (key) => key,
})

export function I18nProvider({ children }: { children: ReactNode }) {
  const [locale, setLocaleState] = useState<Locale>(() => {
    const saved = localStorage.getItem('skm-locale')
    if (saved === 'zh' || saved === 'en') return saved
    return navigator.language.startsWith('zh') ? 'zh' : 'en'
  })

  const setLocale = useCallback((l: Locale) => {
    setLocaleState(l)
    localStorage.setItem('skm-locale', l)
  }, [])

  const t = useCallback((key: TranslationKey): string => {
    return translations[locale][key] ?? key
  }, [locale])

  return (
    <I18nContext.Provider value={{ locale, setLocale, t }}>
      {children}
    </I18nContext.Provider>
  )
}

export function useI18n() {
  return useContext(I18nContext)
}
