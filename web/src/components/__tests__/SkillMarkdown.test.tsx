import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import SkillMarkdown from '../SkillMarkdown'

describe('SkillMarkdown', () => {
  it('renders markdown content', () => {
    render(<SkillMarkdown content="# Hello World" />)
    expect(screen.getByText('Hello World')).toBeInTheDocument()
    expect(screen.getByRole('heading', { level: 1 })).toBeInTheDocument()
  })

  it('renders paragraphs and links', () => {
    render(<SkillMarkdown content="This is a [link](https://example.com)." />)
    const link = screen.getByRole('link', { name: 'link' })
    expect(link).toHaveAttribute('href', 'https://example.com')
  })

  it('strips YAML frontmatter', () => {
    const content = `---
title: My Skill
version: 1.0
---
# Actual Content

This is the body.`
    render(<SkillMarkdown content={content} />)
    expect(screen.getByText('Actual Content')).toBeInTheDocument()
    expect(screen.queryByText('title: My Skill')).not.toBeInTheDocument()
    expect(screen.queryByText('version: 1.0')).not.toBeInTheDocument()
  })

  it('renders content without frontmatter unchanged', () => {
    render(<SkillMarkdown content="Just plain text" />)
    expect(screen.getByText('Just plain text')).toBeInTheDocument()
  })

  it('renders code blocks', () => {
    const content = '```\nconsole.log("hello")\n```'
    render(<SkillMarkdown content={content} />)
    expect(screen.getByText('console.log("hello")')).toBeInTheDocument()
  })

  it('renders GFM tables', () => {
    const content = `| Col A | Col B |
| --- | --- |
| val1 | val2 |`
    render(<SkillMarkdown content={content} />)
    expect(screen.getByText('Col A')).toBeInTheDocument()
    expect(screen.getByText('val1')).toBeInTheDocument()
  })

  it('renders lists', () => {
    const content = `- item one
- item two
- item three`
    render(<SkillMarkdown content={content} />)
    expect(screen.getByText('item one')).toBeInTheDocument()
    expect(screen.getByText('item three')).toBeInTheDocument()
  })
})
