import { useEffect, useRef } from 'react'
import { AlertTriangle, Trash2, X } from 'lucide-react'

export type ConfirmVariant = 'delete' | 'overwrite' | 'discard'

interface ConfirmDialogProps {
  variant: ConfirmVariant
  filename: string
  currentPath?: string
  onConfirm: () => void
  onCancel: () => void
}

interface DialogContent {
  title: string
  body: React.ReactNode
  confirmLabel: string
  cancelLabel: string
  destructive: boolean
}

function getContent(variant: ConfirmVariant, filename: string, currentPath?: string): DialogContent {
  switch (variant) {
    case 'delete':
      return {
        title: 'Mover para lixeira?',
        body: (
          <p style={{ color: 'var(--text-secondary)' }} className="text-sm leading-relaxed">
            <strong className="text-[--text-primary]">{filename}</strong> será movido para{' '}
            <code className="text-xs bg-[#1e293b] px-1.5 py-0.5 rounded">.trash/</code>. Esta ação pode ser desfeita manualmente.
          </p>
        ),
        confirmLabel: 'Mover para lixeira',
        cancelLabel: 'Cancelar',
        destructive: true,
      }
    case 'overwrite':
      return {
        title: 'Substituir arquivo?',
        body: (
          <p style={{ color: 'var(--text-secondary)' }} className="text-sm leading-relaxed">
            <strong className="text-[--text-primary]">{filename}</strong> já existe. Deseja substituir?
          </p>
        ),
        confirmLabel: 'Substituir',
        cancelLabel: 'Cancelar',
        destructive: true,
      }
    case 'discard':
      return {
        title: 'Descartar alterações?',
        body: (
          <p style={{ color: 'var(--text-secondary)' }} className="text-sm leading-relaxed">
            Você tem alterações não salvas em{' '}
            <strong className="text-[--text-primary]">{filename || currentPath || 'este arquivo'}</strong>.
          </p>
        ),
        confirmLabel: 'Descartar',
        cancelLabel: 'Continuar editando',
        destructive: true,
      }
  }
}

export default function ConfirmDialog({ variant, filename, currentPath, onConfirm, onCancel }: ConfirmDialogProps) {
  const content = getContent(variant, filename, currentPath)
  const cancelRef = useRef<HTMLButtonElement>(null)

  useEffect(() => {
    cancelRef.current?.focus()
    const handleKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onCancel()
    }
    window.addEventListener('keydown', handleKey)
    return () => window.removeEventListener('keydown', handleKey)
  }, [onCancel])

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center"
      style={{ background: 'rgba(0,0,0,0.6)' }}
      onClick={(e) => { if (e.target === e.currentTarget) onCancel() }}
    >
      <div
        className="relative w-full max-w-md mx-4"
        style={{
          background: 'var(--bg-card)',
          border: '1px solid var(--border)',
          borderRadius: '12px',
          boxShadow: '0 25px 50px rgba(0,0,0,0.5)',
        }}
        role="dialog"
        aria-modal="true"
        aria-labelledby="confirm-dialog-title"
      >
        {/* Header */}
        <div className="flex items-start gap-3 p-5 pb-3">
          <div className="flex-shrink-0 mt-0.5">
            <AlertTriangle size={20} className="text-[--warning]" style={{ color: 'var(--warning)' }} />
          </div>
          <div className="flex-1 min-w-0">
            <h2
              id="confirm-dialog-title"
              className="text-base font-semibold"
              style={{ color: 'var(--text-primary)' }}
            >
              {content.title}
            </h2>
          </div>
          <button
            onClick={onCancel}
            className="flex-shrink-0 rounded-lg p-1 transition-colors"
            style={{ color: 'var(--text-muted)' }}
            onMouseEnter={e => (e.currentTarget.style.color = 'var(--text-secondary)')}
            onMouseLeave={e => (e.currentTarget.style.color = 'var(--text-muted)')}
            aria-label="Fechar"
          >
            <X size={16} />
          </button>
        </div>

        {/* Body */}
        <div className="px-5 pb-5">
          {content.body}

          {/* Actions */}
          <div className="flex items-center justify-end gap-2 mt-5">
            <button
              ref={cancelRef}
              onClick={onCancel}
              className="px-4 py-2 text-sm rounded-lg transition-colors"
              style={{
                color: 'var(--text-secondary)',
                border: '1px solid var(--border)',
                background: 'transparent',
              }}
              onMouseEnter={e => {
                e.currentTarget.style.background = 'var(--surface-hover)'
                e.currentTarget.style.color = 'var(--text-primary)'
              }}
              onMouseLeave={e => {
                e.currentTarget.style.background = 'transparent'
                e.currentTarget.style.color = 'var(--text-secondary)'
              }}
            >
              {content.cancelLabel}
            </button>
            <button
              onClick={onConfirm}
              className="px-4 py-2 text-sm font-medium rounded-lg transition-colors flex items-center gap-1.5"
              style={{
                background: content.destructive ? 'var(--danger)' : 'var(--evo-green)',
                color: content.destructive ? '#fff' : '#0C111D',
                border: 'none',
              }}
              onMouseEnter={e => {
                e.currentTarget.style.background = content.destructive ? '#dc2626' : '#00e699'
              }}
              onMouseLeave={e => {
                e.currentTarget.style.background = content.destructive ? 'var(--danger)' : 'var(--evo-green)'
              }}
            >
              {variant === 'delete' && <Trash2 size={14} />}
              {content.confirmLabel}
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
