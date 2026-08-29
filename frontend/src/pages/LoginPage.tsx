import { useState, type FormEvent } from 'react'
import { useAuth } from '../auth/useAuth'
import { AuthField } from '../components/auth/AuthField'
import { AuthLayout } from '../components/auth/AuthLayout'
import { InlineMessage } from '../components/auth/InlineMessage'
import { navigate } from '../router'
import { mapErrorToFieldErrors } from '../utils'

type LoginField = 'email' | 'password'

const fieldMap = {
  email: 'email',
  password: 'password',
} satisfies Record<string, LoginField>

export function LoginPage() {
  const { login } = useAuth()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [fieldErrors, setFieldErrors] = useState<Partial<Record<LoginField, string>>>({})
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [submitError, setSubmitError] = useState<string | null>(null)

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()

    const nextFieldErrors = validateLoginForm({ email, password })
    setFieldErrors(nextFieldErrors)
    setSubmitError(null)

    if (Object.keys(nextFieldErrors).length > 0) {
      return
    }

    try {
      setIsSubmitting(true)
      const session = await login({
        email: email.trim(),
        password,
      })

      navigate(session.profileCompleted ? '/app' : '/profile', { replace: true })
    } catch (error) {
      const mappedErrors = mapErrorToFieldErrors(error, fieldMap)
      setFieldErrors(mappedErrors.fieldErrors)
      setSubmitError(mappedErrors.formError)
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <AuthLayout
      subtitle="Форма отправляет данные в API gateway и покажет ошибку сервера прямо на экране."
      switchHref="/register"
      switchLabel="Зарегистрироваться"
      switchText="Нет аккаунта?"
      title="ВХОД В АККАУНТ"
    >
      <form onSubmit={handleSubmit}>
        <AuthField
          autoComplete="email"
          disabled={isSubmitting}
          error={fieldErrors.email}
          id="login-email"
          label="Email"
          name="email"
          onChange={(event) => {
            setEmail(event.target.value)
            setFieldErrors((current) => ({ ...current, email: undefined }))
            setSubmitError(null)
          }}
          required
          type="email"
          value={email}
        />

        <AuthField
          autoComplete="current-password"
          disabled={isSubmitting}
          error={fieldErrors.password}
          id="login-password"
          label="Пароль"
          name="password"
          onChange={(event) => {
            setPassword(event.target.value)
            setFieldErrors((current) => ({ ...current, password: undefined }))
            setSubmitError(null)
          }}
          required
          type="password"
          value={password}
        />

        {submitError ? <InlineMessage>{submitError}</InlineMessage> : null}

        <button className="btn" disabled={isSubmitting} type="submit">
          {isSubmitting ? 'ВХОД...' : 'ВОЙТИ'}
        </button>
      </form>
    </AuthLayout>
  )
}

function validateLoginForm(values: { email: string; password: string }) {
  const errors: Partial<Record<LoginField, string>> = {}

  if (!values.email.trim()) {
    errors.email = 'Введите email.'
  }

  if (!values.password) {
    errors.password = 'Введите пароль.'
  }

  return errors
}
