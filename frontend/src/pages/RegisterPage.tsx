import { useState, type FormEvent } from 'react'
import { useAuth } from '../auth/useAuth'
import { AuthField } from '../components/auth/AuthField'
import { AuthLayout } from '../components/auth/AuthLayout'
import { InlineMessage } from '../components/auth/InlineMessage'
import { navigate } from '../router'
import { mapErrorToFieldErrors } from '../utils'

type RegisterField = 'birthDate' | 'email' | 'name' | 'password'

const fieldMap = {
  birth_date: 'birthDate',
  email: 'email',
  name: 'name',
  password: 'password',
} satisfies Record<string, RegisterField>

export function RegisterPage() {
  const { register } = useAuth()
  const [birthDate, setBirthDate] = useState('')
  const [email, setEmail] = useState('')
  const [name, setName] = useState('')
  const [password, setPassword] = useState('')
  const [fieldErrors, setFieldErrors] = useState<Partial<Record<RegisterField, string>>>({})
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [submitError, setSubmitError] = useState<string | null>(null)

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()

    const nextFieldErrors = validateRegisterForm({
      birthDate,
      email,
      name,
      password,
    })
    setFieldErrors(nextFieldErrors)
    setSubmitError(null)

    if (Object.keys(nextFieldErrors).length > 0) {
      return
    }

    try {
      setIsSubmitting(true)
      const session = await register({
        birth_date: birthDate,
        email: email.trim(),
        name: name.trim(),
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
      subtitle="Регистрация отправляется в API gateway, а refresh token сохраняется в HttpOnly cookie."
      switchHref="/login"
      switchLabel="Войти"
      switchText="Уже есть аккаунт?"
      title="РЕГИСТРАЦИЯ"
    >
      <form onSubmit={handleSubmit}>
        <AuthField
          autoComplete="name"
          disabled={isSubmitting}
          error={fieldErrors.name}
          id="register-name"
          label="Имя"
          name="name"
          onChange={(event) => {
            setName(event.target.value)
            setFieldErrors((current) => ({ ...current, name: undefined }))
            setSubmitError(null)
          }}
          required
          type="text"
          value={name}
        />

        <AuthField
          autoComplete="email"
          disabled={isSubmitting}
          error={fieldErrors.email}
          id="register-email"
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
          autoComplete="new-password"
          disabled={isSubmitting}
          error={fieldErrors.password}
          id="register-password"
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

        <AuthField
          disabled={isSubmitting}
          error={fieldErrors.birthDate}
          id="register-birth-date"
          label="Дата рождения"
          name="birth_date"
          onChange={(event) => {
            setBirthDate(event.target.value)
            setFieldErrors((current) => ({ ...current, birthDate: undefined }))
            setSubmitError(null)
          }}
          required
          type="date"
          value={birthDate}
        />

        {submitError ? <InlineMessage>{submitError}</InlineMessage> : null}

        <button className="btn" disabled={isSubmitting} type="submit">
          {isSubmitting ? 'РЕГИСТРАЦИЯ...' : 'ЗАРЕГИСТРИРОВАТЬСЯ'}
        </button>
      </form>
    </AuthLayout>
  )
}

function validateRegisterForm(values: {
  birthDate: string
  email: string
  name: string
  password: string
}) {
  const errors: Partial<Record<RegisterField, string>> = {}
  const trimmedName = values.name.trim()
  const trimmedEmail = values.email.trim()

  if (!trimmedName) {
    errors.name = 'Введите имя.'
  } else if (trimmedName.length < 2) {
    errors.name = 'Имя должно содержать минимум 2 символа.'
  }

  if (!trimmedEmail) {
    errors.email = 'Введите email.'
  }

  if (!values.password) {
    errors.password = 'Введите пароль.'
  } else if (values.password.length < 8) {
    errors.password = 'Пароль должен содержать минимум 8 символов.'
  } else if (!/(?=.*\p{L})(?=.*\d)/u.test(values.password)) {
    errors.password = 'Пароль должен содержать хотя бы одну букву и одну цифру.'
  }

  if (!values.birthDate) {
    errors.birthDate = 'Укажите дату рождения.'
  } else {
    const age = calculateAge(values.birthDate)

    if (age === null) {
      errors.birthDate = 'Дата должна быть в формате YYYY-MM-DD.'
    } else if (age < 13 || age > 120) {
      errors.birthDate = 'Возраст должен быть от 13 до 120 лет.'
    }
  }

  return errors
}

function calculateAge(birthDate: string) {
  const parsedDate = new Date(`${birthDate}T00:00:00`)
  if (Number.isNaN(parsedDate.getTime())) {
    return null
  }

  const now = new Date()
  if (parsedDate.getTime() > now.getTime()) {
    return -1
  }

  let age = now.getFullYear() - parsedDate.getFullYear()
  const hasBirthdayPassed =
    now.getMonth() > parsedDate.getMonth() ||
    (now.getMonth() === parsedDate.getMonth() && now.getDate() >= parsedDate.getDate())

  if (!hasBirthdayPassed) {
    age -= 1
  }

  return age
}
