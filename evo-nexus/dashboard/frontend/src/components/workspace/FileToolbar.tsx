import {
  FilePlus, FolderPlus, Upload, RefreshCw, Edit3, Download,
  Pencil, Trash2, Save, X, Share2,
} from 'lucide-react'

export type EditorMode = 'preview' | 'edit'

interface FileToolbarProps {
  selectedPath: string | null
  isDir: boolean
  mode: EditorMode
  isDirty: boolean
  onNewFile: () => void
  onNewFolder: () => void
  onUpload: () => void
  onRefresh: () => void
  onEdit: () => void
  onSave: () => void
  onCancel: () => void
  onRename: () => void
  onDelete: () => void
  onDownload: () => void
  onShare: () => void
}

interface ToolbarButtonProps {
  icon: React.ReactNode
  label: string
  onClick: () => void
  variant?: 'default' | 'accent' | 'warning' | 'danger'
  disabled?: boolean
}

function ToolbarButton({ icon, label, onClick, variant = 'default', disabled = false }: ToolbarButtonProps) {
  const baseStyle: React.CSSProperties = {
    display: 'flex',
    alignItems: 'center',
    gap: '6px',
    padding: '6px 12px',
    borderRadius: '6px',
    fontSize: '13px',
    fontWeight: 500,
    cursor: disabled ? 'not-allowed' : 'pointer',
    opacity: disabled ? 0.4 : 1,
    transition: 'all 150ms ease',
    border: '1px solid transparent',
    outline: 'none',
  }

  const variantStyle = (): React.CSSProperties => {
    switch (variant) {
      case 'accent':
        return {
          background: 'rgba(0,255,167,0.12)',
          color: 'var(--evo-green)',
          borderColor: 'rgba(0,255,167,0.3)',
        }
      case 'warning':
        return {
          background: 'rgba(245,158,11,0.12)',
          color: 'var(--warning)',
          borderColor: 'rgba(245,158,11,0.4)',
        }
      case 'danger':
        return {
          color: '#f87171',
          background: 'transparent',
          borderColor: 'transparent',
        }
      default:
        return {
          color: 'var(--text-secondary)',
          background: 'transparent',
          borderColor: 'transparent',
        }
    }
  }

  const handleMouseEnter = (e: React.MouseEvent<HTMLButtonElement>) => {
    if (disabled) return
    switch (variant) {
      case 'accent':
        e.currentTarget.style.background = 'rgba(0,255,167,0.2)'
        break
      case 'warning':
        e.currentTarget.style.background = 'rgba(245,158,11,0.2)'
        break
      case 'danger':
        e.currentTarget.style.background = 'rgba(239,68,68,0.1)'
        e.currentTarget.style.color = '#ef4444'
        break
      default:
        e.currentTarget.style.background = 'var(--surface-hover)'
        e.currentTarget.style.color = 'var(--text-primary)'
    }
    e.currentTarget.style.transform = 'scale(1)'
  }

  const handleMouseLeave = (e: React.MouseEvent<HTMLButtonElement>) => {
    const vs = variantStyle()
    e.currentTarget.style.background = vs.background as string
    e.currentTarget.style.color = vs.color as string
    e.currentTarget.style.transform = 'scale(1)'
  }

  const handleMouseDown = (e: React.MouseEvent<HTMLButtonElement>) => {
    if (!disabled) e.currentTarget.style.transform = 'scale(0.97)'
  }

  const handleMouseUp = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.currentTarget.style.transform = 'scale(1)'
  }

  return (
    <button
      onClick={disabled ? undefined : onClick}
      disabled={disabled}
      style={{ ...baseStyle, ...variantStyle() }}
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
      onMouseDown={handleMouseDown}
      onMouseUp={handleMouseUp}
    >
      {icon}
      <span className="hidden sm:inline">{label}</span>
    </button>
  )
}

function Separator() {
  return (
    <div
      style={{
        width: '1px',
        height: '20px',
        background: 'var(--border)',
        margin: '0 4px',
        flexShrink: 0,
      }}
    />
  )
}

export default function FileToolbar({
  selectedPath,
  isDir,
  mode,
  isDirty,
  onNewFile,
  onNewFolder,
  onUpload,
  onRefresh,
  onEdit,
  onSave,
  onCancel,
  onRename,
  onDelete,
  onDownload,
  onShare,
}: FileToolbarProps) {
  const hasSelection = selectedPath !== null
  const isFile = hasSelection && !isDir
  const isEditMode = mode === 'edit'

  return (
    <div
      className="flex items-center gap-0.5 px-3 py-2 overflow-x-auto"
      style={{
        minHeight: '44px',
        flexShrink: 0,
      }}
    >
      {/* Group: create + upload + refresh — always visible */}
      <ToolbarButton
        icon={<FilePlus size={14} />}
        label="Novo arquivo"
        onClick={onNewFile}
      />
      <ToolbarButton
        icon={<FolderPlus size={14} />}
        label="Nova pasta"
        onClick={onNewFolder}
      />
      <ToolbarButton
        icon={<Upload size={14} />}
        label="Upload"
        onClick={onUpload}
      />

      {/* Separator */}
      {(hasSelection || isEditMode) && <Separator />}

      {/* Edit mode buttons */}
      {isFile && isEditMode && (
        <>
          <ToolbarButton
            icon={
              isDirty
                ? <Save size={14} className={isDirty ? 'animate-pulse' : ''} />
                : <Save size={14} />
            }
            label="Salvar"
            onClick={onSave}
            variant={isDirty ? 'warning' : 'accent'}
          />
          <ToolbarButton
            icon={<X size={14} />}
            label="Cancelar"
            onClick={onCancel}
          />
          <Separator />
        </>
      )}

      {/* Preview mode edit button */}
      {isFile && !isEditMode && (
        <>
          <ToolbarButton
            icon={<Edit3 size={14} />}
            label="Editar"
            onClick={onEdit}
            variant="accent"
          />
          <Separator />
        </>
      )}

      {/* Download — visible when file selected */}
      {isFile && (
        <ToolbarButton
          icon={<Download size={14} />}
          label="Download"
          onClick={onDownload}
        />
      )}

      {/* Share — visible when file selected and not in edit mode */}
      {isFile && !isEditMode && (
        <ToolbarButton
          icon={<Share2 size={14} />}
          label="Compartilhar"
          onClick={onShare}
        />
      )}

      {/* Rename — visible when anything selected */}
      {hasSelection && (
        <ToolbarButton
          icon={<Pencil size={14} />}
          label="Renomear"
          onClick={onRename}
        />
      )}

      {/* Delete — visible when anything selected */}
      {hasSelection && (
        <ToolbarButton
          icon={<Trash2 size={14} />}
          label="Excluir"
          onClick={onDelete}
          variant="danger"
        />
      )}

      {/* Separator + Refresh — always last */}
      <Separator />
      <ToolbarButton
        icon={<RefreshCw size={14} />}
        label="Atualizar"
        onClick={onRefresh}
      />
    </div>
  )
}
