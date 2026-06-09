import { useEffect, useRef, useState } from 'react'
import { EditorView, ViewUpdate } from '@codemirror/view'
import { EditorState } from '@codemirror/state'
import { baseExtensions, languageForPath } from './cmSetup'

interface FileEditorProps {
  initialContent: string
  path: string
  onDirtyChange: (dirty: boolean) => void
  onSave: (content: string) => void
  editorRef?: React.MutableRefObject<{ getContent: () => string } | null>
}

export default function FileEditor({ initialContent, path, onDirtyChange, onSave, editorRef }: FileEditorProps) {
  const containerRef = useRef<HTMLDivElement>(null)
  const viewRef = useRef<EditorView | null>(null)
  const [ready, setReady] = useState(false)

  useEffect(() => {
    if (!containerRef.current) return

    const handleSave = () => {
      const content = viewRef.current?.state.doc.toString() ?? ''
      onSave(content)
    }

    const view = new EditorView({
      state: EditorState.create({
        doc: initialContent,
        extensions: [
          ...baseExtensions(false, handleSave),
          ...languageForPath(path),
          EditorView.updateListener.of((update: ViewUpdate) => {
            if (update.docChanged) {
              const current = update.state.doc.toString()
              onDirtyChange(current !== initialContent)
            }
          }),
        ],
      }),
      parent: containerRef.current,
    })
    viewRef.current = view
    setReady(true)

    // Expose content getter
    if (editorRef) {
      editorRef.current = {
        getContent: () => view.state.doc.toString(),
      }
    }

    return () => {
      view.destroy()
      viewRef.current = null
      if (editorRef) editorRef.current = null
    }
  }, [initialContent, path])

  return (
    <div
      ref={containerRef}
      className="h-full overflow-hidden"
      style={{
        background: 'var(--bg-primary)',
        opacity: ready ? 1 : 0,
        transition: 'opacity 150ms ease',
      }}
    />
  )
}
