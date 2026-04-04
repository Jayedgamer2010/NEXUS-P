import React from 'react';
import './Alert.css';

type AlertVariant = 'success' | 'error' | 'warning' | 'info';

interface AlertProps {
  variant?: AlertVariant;
  title?: string;
  children: React.ReactNode;
  onClose?: () => void;
}

export default function Alert({
  variant = 'info',
  title,
  children,
  onClose,
}: AlertProps) {
  return (
    <div className={`alert alert--${variant}`} role="alert">
      <div className="alert-content">
        {title && <strong className="alert-title">{title}</strong>}
        <div className="alert-message">{children}</div>
      </div>
      {onClose && (
        <button className="alert-close" onClick={onClose} aria-label="Close alert">
          &times;
        </button>
      )}
    </div>
  );
}
