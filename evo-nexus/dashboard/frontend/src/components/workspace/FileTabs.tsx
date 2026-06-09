import { useRef, useState, useEffect, useCallback } from 'react'
import { X } from 'lucide-react'

export interface TabEntry {
  path: string
}

interface FileTabsProps {
  tabs: TabEntry[]
  activePath: string | null
  dirtyPaths: Set<string>
  onSwitch: (path: string) => void
  onClose: (path: string) => void
  onCloseAll?: () => void
  onCloseOthers?: (path: string) => void
  onCloseToLeft?: (path: string) => void
  onCloseToRight?: (path: string) => void
}

function fileName(path: string): string {
  return path.split('/').pop() ?? path
}

interface ContextMenu {
  x: number
  y: number
  tabPath: string
  tabIndex: number
}

export default function FileTabs({ tabs, activePath, dirtyPaths, onSwitch, onClose, onCloseAll, onCloseOthers, onCloseToLeft, onCloseToRight }: FileTabsProps) {
  const scrollRef = useRef<HTMLDivElement>(null)
  const [ctxMenu, setCtxMenu] = useState<ContextMenu | null>(null)

  // Close context menu on click outside or Escape
  useEffect(() => {
    if (!ctxMenu) return
    const handleClose = () => setCtxMenu(null)
    const handleKey = (e: KeyboardEvent) => { if (e.key === 'Escape') setCtxMenu(null) }
    document.addEventListener('click', handleClose)
    document.addEventListener('keydown', handleKey)
    return () => {
      document.removeEventListener('click', handleClose)
      document.removeEventListener('keydown', handleKey)
    }
  }, [ctxMenu])

  const handleContextMenu = useCallback((e: React.MouseEvent, tabPath: string, tabIndex: number) => {
    e.preventDefault()
    setCtxMenu({ x: e.clientX, y: e.clientY, tabPath, tabIndex })
  }, [])

  if (tabs.length === 0) return null

  const hasLeft = ctxMenu ? ctxMenu.tabIndex > 0 : false
  const hasRight = ctxMenu ? ctxMenu.tabIndex < tabs.length - 1 : false
  const hasOthers = tabs.length > 1

  return (
    <>
      <div
        ref={scrollRef}
        className="flex items-end overflow-x-auto flex-shrink-0 file-tabs-scroll"
        style={{
          background: 'var(--bg-card)',
          borderBottom: '1px solid var(--border)',
          scrollbarWidth: 'none',
          msOverflowStyle: 'none',
          minHeight: '36px',
        }}
      >
        <style>{`
          .file-tabs-scroll::-webkit-scrollbar { display: none; }
          .file-tab-close { opacity: 0; }
          .file-tab:hover .file-tab-close,
          .file-tab.active .file-tab-close { opacity: 1; }
        `}</style>
        {tabs.map((tab, index) => {
          const isActive = tab.path === activePath
          const isDirty = dirtyPaths.has(tab.path)

          return (
            <button
              key={tab.path}
              className={`file-tab flex items-center gap-1.5 px-3 py-2 text-xs flex-shrink-0 transition-colors relative select-none${isActive ? ' active' : ''}`}
              style={{
                maxWidth: '200px',
                minWidth: '80px',
                borderBottom: isActive ? '2px solid var(--evo-green)' : '2px solid transparent',
                color: isActive ? 'var(--text-primary)' : 'var(--text-muted)',
                background: isActive ? 'var(--surface-active)' : 'transparent',
                cursor: 'pointer',
                whiteSpace: 'nowrap',
              }}
              onClick={() => onSwitch(tab.path)}
              onContextMenu={e => handleContextMenu(e, tab.path, index)}
              onMouseDown={e => {
                if (e.button === 1) {
                  e.preventDefault()
                  onClose(tab.path)
                }
              }}
              onMouseEnter={e => {
                if (!isActive) e.currentTarget.style.background = 'var(--surface-hover)'
              }}
              onMouseLeave={e => {
                if (!isActive) e.currentTarget.style.background = 'transparent'
              }}
              title={tab.path}
            >
              {isDirty && (
                <span
                  style={{
                    width: '6px',
                    height: '6px',
                    borderRadius: '50%',
                    background: 'var(--warning)',
                    flexShrink: 0,
                    display: 'inline-block',
                  }}
                />
              )}

              <span
                className="truncate flex-1 min-w-0"
                style={{ maxWidth: isDirty ? '130px' : '150px' }}
              >
                {fileName(tab.path)}
              </span>

              <span
                className="file-tab-close flex items-center justify-center flex-shrink-0 rounded transition-colors"
                style={{ width: '16px', height: '16px' }}
                onClick={e => {
                  e.stopPropagation()
                  onClose(tab.path)
                }}
                onMouseEnter={e => { e.currentTarget.style.background = 'var(--border)' }}
                onMouseLeave={e => { e.currentTarget.style.background = 'transparent' }}
              >
                <X size={10} />
              </span>
            </button>
          )
        })}
      </div>

      {/* Context menu */}
      {ctxMenu && (
        <div
          className="fixed z-[100] rounded-lg border shadow-xl py-1"
          style={{
            left: ctxMenu.x,
            top: ctxMenu.y,
            background: 'var(--bg-card)',
            borderColor: 'var(--border)',
            minWidth: '200px',
            boxShadow: '0 8px 32px rgba(0,0,0,0.4)',
          }}
        >
          <CtxMenuItem
            label="Fechar"
            shortcut="Ctrl+W"
            onClick={() => { onClose(ctxMenu.tabPath); setCtxMenu(null) }}
          />
          {hasOthers && onCloseOthers && (
            <CtxMenuItem
              label="Fechar outras"
              onClick={() => { onCloseOthers(ctxMenu.tabPath); setCtxMenu(null) }}
            />
          )}
          {hasLeft && onCloseToLeft && (
            <CtxMenuItem
              label="Fechar todas à esquerda"
              onClick={() => { onCloseToLeft(ctxMenu.tabPath); setCtxMenu(null) }}
            />
          )}
          {hasRight && onCloseToRight && (
            <CtxMenuItem
              label="Fechar todas à direita"
              onClick={() => { onCloseToRight(ctxMenu.tabPath); setCtxMenu(null) }}
            />
          )}
          <div style={{ borderTop: '1px solid var(--border)', margin: '4px 0' }} />
          {onCloseAll && (
            <CtxMenuItem
              label="Fechar todas"
              onClick={() => { onCloseAll(); setCtxMenu(null) }}
            />
          )}
        </div>
      )}
    </>
  )
}

function CtxMenuItem({ label, shortcut, onClick }: { label: string; shortcut?: string; onClick: () => void }) {
  return (
    <button
      className="w-full flex items-center justify-between px-3 py-1.5 text-xs transition-colors"
      style={{ color: 'var(--text-secondary)' }}
      onClick={onClick}
      onMouseEnter={e => { e.currentTarget.style.background = 'var(--surface-hover)'; e.currentTarget.style.color = 'var(--text-primary)' }}
      onMouseLeave={e => { e.currentTarget.style.background = 'transparent'; e.currentTarget.style.color = 'var(--text-secondary)' }}
    >
      <span>{label}</span>
      {shortcut && <span style={{ color: 'var(--text-muted)', fontSize: '10px' }}>{shortcut}</span>}
    </button>
  )
}
