export function mapErrorToFieldErrors<T extends string>(
  error: unknown,
  fieldMap: Record<string, T>,
) {
  const message = toErrorMessage(error, 'Не удалось выполнить запрос. Попробуйте еще раз.')
  const separatorIndex = message.indexOf(':')

  if (separatorIndex === -1) {
    return {
      fieldErrors: {} as Partial<Record<T, string>>,
      formError: message,
    }
  }

  const rawField = message.slice(0, separatorIndex).trim()
  const mappedField = fieldMap[rawField]
  if (!mappedField) {
    return {
      fieldErrors: {} as Partial<Record<T, string>>,
      formError: message,
    }
  }

  return {
    fieldErrors: {
      [mappedField]: message.slice(separatorIndex + 1).trim(),
    } as Partial<Record<T, string>>,
    formError: null,
  }
}

export function toErrorMessage(error: unknown, fallback: string) {
  if (error instanceof Error && error.message.trim()) {
    return translateBackendError(error.message)
  }

  return fallback
}

export function formatDate(value: string) {
  const parsedDate = new Date(`${value}T00:00:00`)
  if (Number.isNaN(parsedDate.getTime())) {
    return value
  }

  return new Intl.DateTimeFormat('ru-RU', {
    day: '2-digit',
    month: 'long',
    year: 'numeric',
  }).format(parsedDate)
}

export function formatDateTimeFromUnix(value: number) {
  return new Intl.DateTimeFormat('ru-RU', {
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    month: 'long',
    year: 'numeric',
  }).format(new Date(value * 1000))
}

export function formatList(values: string[]) {
  return values.length > 0 ? values.join(', ') : 'Не указано'
}

function translateBackendError(message: string) {
  const normalizedMessage = message.trim()

  const fieldTranslations: Record<string, string> = {
    'birth_date: age must be between 13 and 120': 'birth_date: Возраст должен быть от 13 до 120 лет.',
    'birth_date: cannot be in the future': 'birth_date: Дата рождения не может быть в будущем.',
    'birth_date: is required': 'birth_date: Укажите дату рождения.',
    'birth_date: must use YYYY-MM-DD format': 'birth_date: Дата должна быть в формате YYYY-MM-DD.',
    'email: is not a valid email address': 'email: Укажите корректный email.',
    'email: is required': 'email: Введите email.',
    'email: is too long': 'email: Email слишком длинный.',
    'name: is required': 'name: Введите имя.',
    'name: is too long': 'name: Имя слишком длинное.',
    'name: must be at least 2 characters long': 'name: Имя должно содержать минимум 2 символа.',
    'password: is required': 'password: Введите пароль.',
    'password: is too long': 'password: Пароль слишком длинный.',
    'password: must be at least 8 characters long': 'password: Пароль должен содержать минимум 8 символов.',
    'password: must contain at least one letter and one digit':
      'password: Пароль должен содержать хотя бы одну букву и одну цифру.',
  }

  const commonTranslations: Record<string, string> = {
    'access token expired': 'Сессия истекла. Войдите снова.',
    'authorization bearer token is required': 'Требуется авторизация.',
    'invalid credentials': 'Неправильный логин или пароль.',
    'invalid JSON body': 'Некорректные данные формы.',
    'internal error': 'Внутренняя ошибка сервера. Попробуйте еще раз.',
    'profile not found': 'Профиль пока не заполнен.',
    'refresh token cookie is required': 'Сессия истекла. Войдите снова.',
    'token expired': 'Сессия истекла. Войдите снова.',
    'token revoked': 'Сессия завершена. Войдите снова.',
    'user already exists': 'Пользователь с таким email уже существует.',
    'user not found': 'Пользователь не найден.',
  }

  return fieldTranslations[normalizedMessage] ?? commonTranslations[normalizedMessage] ?? normalizedMessage
}
