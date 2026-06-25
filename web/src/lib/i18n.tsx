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
    'groups.addSkills': 'Add Skills',
    'groups.searchSkills': 'Search skills...',
    'groups.noAvailable': 'No available skills to add',
    'groups.add': 'Add',

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

    // Tags
    'tags.title': 'Tags',
    'tags.add': 'Add Tag',
    'tags.remove': 'Remove',
    'tags.rename': 'Rename Tag',
    'tags.delete': 'Delete Tag',
    'tags.search': 'Search tags...',
    'tags.noTags': 'No tags',
    'tags.untagged': 'Untagged',
    'tags.placeholder': 'Enter tag name...',

    // Skill detail
    'detail.title': 'Skill Details',
    'detail.content': 'Content',
    'detail.source': 'Source',
    'detail.metadata': 'Metadata',
    'detail.sourceType': 'Source Type',
    'detail.sourceRef': 'Source Ref',
    'detail.enabled': 'Enabled',
    'detail.disabled': 'Disabled',
    'detail.close': 'Close',
    'detail.tags': 'Tags',
    'detail.agents': 'Synced Agents',
    'detail.noContent': 'No content available',

    // Audit
    'audit.title': 'Audit Log',
    'audit.action': 'Action',
    'audit.target': 'Target',
    'audit.detail': 'Detail',
    'audit.time': 'Time',
    'audit.prune': 'Prune Old Entries',
    'audit.noEntries': 'No audit entries',
    'audit.pruned': 'Old entries pruned',

    // Settings (expanded)
    'settings.agents': 'Agent Management',
    'settings.agentToggle': 'Toggle agent',
    'settings.syncMode': 'Sync Mode',
    'settings.symlink': 'Symlink',
    'settings.copy': 'Copy',
    'settings.appearance': 'Appearance',
    'settings.theme': 'Theme',
    'settings.language': 'Language',
    'settings.textSize': 'Text Size',
    'settings.updates': 'Updates',
    'settings.autoUpdate': 'Auto Update Interval',
    'settings.about': 'About',
    'settings.version': 'Version',
    'settings.saved': 'Settings saved',

    // Toast
    'toast.skillEnabled': 'Skill enabled',
    'toast.skillDisabled': 'Skill disabled',
    'toast.tagAdded': 'Tag added',
    'toast.tagRemoved': 'Tag removed',
    'toast.error': 'Operation failed',
    'toast.synced': 'Sync complete',
    'toast.deleted': 'Deleted successfully',
    'toast.copied': 'Copied to clipboard',

    // Batch
    'batch.selected': 'selected',
    'batch.delete': 'Delete Selected',
    'batch.enable': 'Enable Selected',
    'batch.disable': 'Disable Selected',
    'batch.tag': 'Tag Selected',
    'batch.sync': 'Sync Selected',
    'batch.confirm': 'Are you sure?',
    'batch.cancel': 'Cancel',

    // Search / Command Palette
    'search.placeholder': 'Search skills, groups, pages...',
    'search.noResults': 'No results found',
    'search.skills': 'Skills',
    'search.groups': 'Groups',
    'search.pages': 'Pages',
    'search.hint': '⌘K to search',

    // Nav additions
    'nav.audit': 'Audit Log',
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
    'groups.addSkills': '添加技能',
    'groups.searchSkills': '搜索技能...',
    'groups.noAvailable': '没有可添加的技能',
    'groups.add': '添加',

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

    // Tags
    'tags.title': '标签',
    'tags.add': '添加标签',
    'tags.remove': '移除',
    'tags.rename': '重命名标签',
    'tags.delete': '删除标签',
    'tags.search': '搜索标签...',
    'tags.noTags': '暂无标签',
    'tags.untagged': '未标记',
    'tags.placeholder': '输入标签名...',

    // Skill detail
    'detail.title': '技能详情',
    'detail.content': '内容',
    'detail.source': '来源',
    'detail.metadata': '元数据',
    'detail.sourceType': '来源类型',
    'detail.sourceRef': '来源引用',
    'detail.enabled': '已启用',
    'detail.disabled': '已禁用',
    'detail.close': '关闭',
    'detail.tags': '标签',
    'detail.agents': '已同步智能体',
    'detail.noContent': '暂无内容',

    // Audit
    'audit.title': '操作日志',
    'audit.action': '操作',
    'audit.target': '目标',
    'audit.detail': '详情',
    'audit.time': '时间',
    'audit.prune': '清理旧记录',
    'audit.noEntries': '暂无操作记录',
    'audit.pruned': '旧记录已清理',

    // Settings (expanded)
    'settings.agents': '智能体管理',
    'settings.agentToggle': '切换智能体',
    'settings.syncMode': '同步模式',
    'settings.symlink': '符号链接',
    'settings.copy': '复制',
    'settings.appearance': '外观',
    'settings.theme': '主题',
    'settings.language': '语言',
    'settings.textSize': '字体大小',
    'settings.updates': '更新',
    'settings.autoUpdate': '自动更新间隔',
    'settings.about': '关于',
    'settings.version': '版本',
    'settings.saved': '设置已保存',

    // Toast
    'toast.skillEnabled': '技能已启用',
    'toast.skillDisabled': '技能已禁用',
    'toast.tagAdded': '标签已添加',
    'toast.tagRemoved': '标签已移除',
    'toast.error': '操作失败',
    'toast.synced': '同步完成',
    'toast.deleted': '删除成功',
    'toast.copied': '已复制到剪贴板',

    // Batch
    'batch.selected': '已选中',
    'batch.delete': '删除选中',
    'batch.enable': '启用选中',
    'batch.disable': '禁用选中',
    'batch.tag': '标记选中',
    'batch.sync': '同步选中',
    'batch.confirm': '确定要执行此操作吗？',
    'batch.cancel': '取消',

    // Search / Command Palette
    'search.placeholder': '搜索技能、分组、页面...',
    'search.noResults': '未找到结果',
    'search.skills': '技能',
    'search.groups': '分组',
    'search.pages': '页面',
    'search.hint': '⌘K 搜索',

    // Nav additions
    'nav.audit': '操作日志',
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
