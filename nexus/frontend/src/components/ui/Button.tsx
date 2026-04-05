interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'danger' | 'warning' | 'ghost'
  size?: 'sm' | 'md' | 'lg'
  loading?: boolean
}

export default function Button({
  variant = 'primary',
  size = 'md',
  loading = false,
  children,
  className = '',
  disabled,
  ...props
}: ButtonProps) {
  return (
    <button
      className={`nx-btn nx-btn--${variant} nx-btn--${size} ${className}`}
      disabled={disabled || loading}
      {...props}
    >
      {loading ? (
        <div className="nx-spinner nx-spinner--sm" />
      ) : children}
    </button>
  )
}
