import { useEffect, useCallback } from 'react'

interface ModalProps {
  isOpen: boolean
  onClose: () => void
  title: string
  children: React.ReactNode
  size?: 'sm' | 'md' | 'lg'
}

export default function Modal({ isOpen, onClose, title, children, size = 'md' }: ModalProps) {
  const handleEsc = useCallback((e: KeyboardEvent) => {
    if (e.key === 'Escape') onClose()
  }, [onClose])

  useEffect(() => {
    if (isOpen) {
      document.addEventListener('keydown', handleEsc)
      document.body.style.overflow = 'hidden'
    }
    return () => {
      document.removeEventListener('keydown', handleEsc)
      document.body.style.overflow = ''
    }
  }, [isOpen, handleEsc])

  if (!isOpen) return null

  return (
    <div className="nx-modal-backdrop" onClick={onClose}>
      <div className={`nx-modal nx-modal--${size}`} onClick={(e) => e.stopPropagation()}>
        <div className="nx-modal-header">
          <h3>{title}</h3>
          <button className="nx-modal-close" onClick={onClose}>x</button>
        </div>
        <div className="nx-modal-body">
          {children}
        </div>
      </div>
    </div>
  )
}
