import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'

interface SkillMarkdownProps {
  content: string
}

/** Strip YAML frontmatter (content between --- markers at the start) */
function stripFrontmatter(text: string): string {
  const match = text.match(/^---\r?\n[\s\S]*?\r?\n---\r?\n?/)
  return match ? text.slice(match[0].length) : text
}

export default function SkillMarkdown({ content }: SkillMarkdownProps) {
  const body = stripFrontmatter(content)

  return (
    <div className="skill-markdown prose prose-slate dark:prose-invert max-w-none
      [&_h1]:text-2xl [&_h1]:font-bold [&_h1]:mb-3 [&_h1]:text-slate-900 [&_h1]:dark:text-slate-100
      [&_h2]:text-xl [&_h2]:font-semibold [&_h2]:mb-2 [&_h2]:mt-5 [&_h2]:text-slate-800 [&_h2]:dark:text-slate-200
      [&_h3]:text-lg [&_h3]:font-semibold [&_h3]:mb-2 [&_h3]:mt-4 [&_h3]:text-slate-800 [&_h3]:dark:text-slate-200
      [&_p]:mb-3 [&_p]:text-slate-700 [&_p]:dark:text-slate-300 [&_p]:leading-relaxed
      [&_ul]:list-disc [&_ul]:pl-6 [&_ul]:mb-3
      [&_ol]:list-decimal [&_ol]:pl-6 [&_ol]:mb-3
      [&_li]:mb-1 [&_li]:text-slate-700 [&_li]:dark:text-slate-300
      [&_a]:text-primary-600 [&_a]:dark:text-primary-400 [&_a]:underline [&_a]:hover:text-primary-700
      [&_code]:bg-slate-100 [&_code]:dark:bg-slate-700 [&_code]:px-1.5 [&_code]:py-0.5 [&_code]:rounded [&_code]:text-sm [&_code]:text-pink-600 [&_code]:dark:text-pink-400
      [&_pre]:bg-slate-100 [&_pre]:dark:bg-slate-800 [&_pre]:rounded-lg [&_pre]:p-4 [&_pre]:mb-3 [&_pre]:overflow-x-auto
      [&_pre_code]:bg-transparent [&_pre_code]:p-0 [&_pre_code]:text-slate-800 [&_pre_code]:dark:text-slate-200
      [&_blockquote]:border-l-4 [&_blockquote]:border-slate-300 [&_blockquote]:dark:border-slate-600 [&_blockquote]:pl-4 [&_blockquote]:italic [&_blockquote]:text-slate-600 [&_blockquote]:dark:text-slate-400
      [&_table]:w-full [&_table]:border-collapse [&_table]:mb-3
      [&_th]:border [&_th]:border-slate-300 [&_th]:dark:border-slate-600 [&_th]:px-3 [&_th]:py-1.5 [&_th]:bg-slate-50 [&_th]:dark:bg-slate-700 [&_th]:text-left [&_th]:font-medium
      [&_td]:border [&_td]:border-slate-300 [&_td]:dark:border-slate-600 [&_td]:px-3 [&_td]:py-1.5
      [&_hr]:border-slate-200 [&_hr]:dark:border-slate-700 [&_hr]:my-4
    ">
      <ReactMarkdown remarkPlugins={[remarkGfm]}>{body}</ReactMarkdown>
    </div>
  )
}
