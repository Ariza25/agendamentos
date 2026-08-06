import { AlertTriangle } from 'lucide-react'
import type { ReactNode } from 'react'

import { Button } from './Button'
import { Modal } from './Modal'

type ConfirmModalProps = {
  open: boolean
  title: string
  message: ReactNode
  confirmLabel?: string
  cancelLabel?: string
  loading?: boolean
  onCancel: () => void
  onConfirm: () => void | Promise<void>
}

export function ConfirmModal({
  open,
  title,
  message,
  confirmLabel = 'Confirmar',
  cancelLabel = 'Cancelar',
  loading = false,
  onCancel,
  onConfirm,
}: ConfirmModalProps) {
  return (
    <Modal
      footer={
        <>
          <Button disabled={loading} onClick={onCancel} variant="secondary">
            {cancelLabel}
          </Button>
          <Button disabled={loading} loading={loading} onClick={() => void onConfirm()} variant="danger">
            {confirmLabel}
          </Button>
        </>
      }
      onClose={loading ? () => undefined : onCancel}
      open={open}
      size="sm"
      title={title}
    >
      <div className="flex items-start gap-3">
        <span className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-red-50 text-red-700">
          <AlertTriangle className="h-5 w-5" />
        </span>
        <div className="text-sm text-[var(--color-ink-soft)]">{message}</div>
      </div>
    </Modal>
  )
}
