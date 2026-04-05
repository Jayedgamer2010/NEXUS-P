interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string
  error?: string
  helper?: string
}

export default function Input({ label, error, helper, required, className = '', type = 'text', ...props }: InputProps) {
  return (
    <div className="nx-form-group">
      {label && (
        <label className="nx-input-label">
          {label}
          {required && <span style={{ color: '#ef4444', marginLeft: 4 }}>*</span>}
        </label>
      )}
      <input
        type={type}
        className={`nx-input ${error ? 'nx-input--error' : ''} ${className}`}
        required={required}
        {...props}
      />
      {error && <div className="nx-input-error">{error}</div>}
      {!error && helper && <div className="nx-input-helper">{helper}</div>}
    </div>
  )
}
