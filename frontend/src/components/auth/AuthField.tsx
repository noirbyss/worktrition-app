import type { InputHTMLAttributes } from 'react'

interface AuthFieldProps extends InputHTMLAttributes<HTMLInputElement> {
  error?: string
  label: string
}

export function AuthField({ error, id, label, ...props }: AuthFieldProps) {
  return (
    <div className="form-group">
      <label htmlFor={id}>{label}</label>
      <input
        {...props}
        className={error ? 'input input--error' : 'input'}
        id={id}
      />
      {error ? <p className="field-error">{error}</p> : null}
    </div>
  )
}
