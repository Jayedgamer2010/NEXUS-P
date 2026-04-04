import './Input.css';

interface InputProps extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'children'> {
  label?: string;
  error?: string;
  helperText?: string;
}

export default function Input({
  label,
  error,
  helperText,
  id,
  className = '',
  ...props
}: InputProps) {
  const inputId = id || `input-${Math.random().toString(36).slice(2)}`;

  return (
    <div className={`input-wrapper ${className}`}>
      {label && (
        <label htmlFor={inputId} className="input-label">
          {label}
        </label>
      )}
      <input
        id={inputId}
        className={`input-field ${error ? 'input-field--error' : ''}`}
        aria-invalid={!!error}
        aria-describedby={error ? `${inputId}-error` : helperText ? `${inputId}-helper` : undefined}
        {...props}
      />
      {error && (
        <p id={`${inputId}-error`} className="input-error">
          {error}
        </p>
      )}
      {helperText && !error && (
        <p id={`${inputId}-helper`} className="input-helper">
          {helperText}
        </p>
      )}
    </div>
  );
}
