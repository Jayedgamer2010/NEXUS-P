import React from 'react';
import './Button.css';

type ButtonVariant = 'primary' | 'secondary' | 'danger' | 'warning' | 'ghost';
type ButtonSize = 'sm' | 'md' | 'lg';

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: ButtonVariant;
  size?: ButtonSize;
  loading?: boolean;
  children: React.ReactNode;
}

export default function Button({
  variant = 'primary',
  size = 'md',
  loading = false,
  className = '',
  disabled,
  children,
  ...rest
}: ButtonProps) {
  const classes = [
    'nexus-btn',
    `btn--${variant}`,
    `btn--${size}`,
    loading && 'btn--loading',
    (disabled || loading) && 'btn--disabled',
    className,
  ]
    .filter(Boolean)
    .join(' ');

  return (
    <button className={classes} disabled={disabled || loading} {...rest}>
      {loading && <span className="btn-spinner" />}
      {children}
    </button>
  );
}
