import { useEffect, useRef } from 'react'
import { EditorView } from '@codemirror/view'
import { EditorState } from '@codemirror/state'
import { baseExtensions, languageForPath } from './cmSetup'

interface CodeViewProps {
  content: string
  path: string
}

export default function CodeView({ content, path }: CodeViewProps) {
  const containerRef = useRef<HTMLDivElement>(null)
  const viewRef = useRef<EditorView | null>(null)

  useEffect(() => {
    if (!containerRef.current) return

    const view = new EditorView({
      state: EditorState.create({
        doc: content,
        extensions: [
          ...baseExtensions(true),
          ...languageForPath(path),
        ],
      }),
      parent: containerRef.current,
    })
    viewRef.current = view

    return () => {
      view.destroy()
      viewRef.current = null
    }
  }, [content, path])

  return (
    <div
      ref={containerRef}
      className="h-full overflow-hidden"
      style={{ background: 'var(--bg-primary)' }}
    />
  )
}
