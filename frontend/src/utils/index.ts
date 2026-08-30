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
    'activity_level: must be specified': 'activity_level: Выберите уровень активности.',
    'allergies: items cannot be empty': 'allergies: Удалите пустые элементы из списка аллергий.',
    'allergies: items must contain no more than 80 characters':
      'allergies: Каждый пункт списка аллергий должен быть короче 80 символов.',
    'allergies: must contain no more than 50 items':
      'allergies: Список аллергий не должен содержать больше 50 элементов.',
    'age: must be between 13 and 120': 'age: Возраст должен быть от 13 до 120 лет.',
    'birth_date: age must be between 13 and 120': 'birth_date: Возраст должен быть от 13 до 120 лет.',
    'birth_date: cannot be in the future': 'birth_date: Дата рождения не может быть в будущем.',
    'birth_date: is required': 'birth_date: Укажите дату рождения.',
    'birth_date: must use YYYY-MM-DD format': 'birth_date: Дата должна быть в формате YYYY-MM-DD.',
    'equipment: is too long': 'equipment: Описание инвентаря слишком длинное.',
    'email: is not a valid email address': 'email: Укажите корректный email.',
    'email: is required': 'email: Введите email.',
    'email: is too long': 'email: Email слишком длинный.',
    'excluded_foods: items cannot be empty':
      'excluded_foods: Удалите пустые элементы из списка исключенных продуктов.',
    'excluded_foods: items must contain no more than 80 characters':
      'excluded_foods: Каждый исключенный продукт должен быть короче 80 символов.',
    'excluded_foods: must contain no more than 50 items':
      'excluded_foods: Список исключенных продуктов не должен содержать больше 50 элементов.',
    'food_preferences: items cannot be empty':
      'food_preferences: Удалите пустые элементы из списка пищевых предпочтений.',
    'food_preferences: items must contain no more than 80 characters':
      'food_preferences: Каждый пункт списка предпочтений должен быть короче 80 символов.',
    'food_preferences: must contain no more than 50 items':
      'food_preferences: Список предпочтений не должен содержать больше 50 элементов.',
    'gender: must be specified': 'gender: Выберите пол.',
    'goal: must be specified': 'goal: Выберите цель.',
    'height_cm: must be between 80 and 250': 'height_cm: Рост должен быть в диапазоне от 80 до 250 см.',
    'name: is required': 'name: Введите имя.',
    'name: is too long': 'name: Имя слишком длинное.',
    'name: must be at least 2 characters long': 'name: Имя должно содержать минимум 2 символа.',
    'password: is required': 'password: Введите пароль.',
    'password: is too long': 'password: Пароль слишком длинный.',
    'password: must be at least 8 characters long': 'password: Пароль должен содержать минимум 8 символов.',
    'password: must contain at least one letter and one digit':
      'password: Пароль должен содержать хотя бы одну букву и одну цифру.',
    'target_weight_kg: must be between 25 and 400':
      'target_weight_kg: Целевой вес должен быть в диапазоне от 25 до 400 кг.',
    'training_days_per_week: must be between 0 and 7':
      'training_days_per_week: Количество тренировок должно быть от 0 до 7.',
    'training_level: must be specified': 'training_level: Выберите уровень тренировок.',
    'training_location: must be specified': 'training_location: Выберите место тренировок.',
    'weight_kg: must be between 25 and 400': 'weight_kg: Вес должен быть в диапазоне от 25 до 400 кг.',
  }

  const commonTranslations: Record<string, string> = {
    'access token expired': 'Сессия истекла. Войдите снова.',
    'authorization bearer token is required': 'Требуется авторизация.',
    'day_of_week is required': 'Не удалось определить день недели для плана питания.',
    'day_of_week must be between 1 and 7': 'День недели должен быть в диапазоне от 1 до 7.',
    'generation not found': 'Генерация плана не найдена.',
    'invalid credentials': 'Неправильный логин или пароль.',
    'invalid JSON body': 'Некорректные данные формы.',
    'internal error': 'Внутренняя ошибка сервера. Попробуйте еще раз.',
    'meal item not found': 'Приём пищи не найден в активном плане.',
    'plan for user not found': 'План питания для пользователя пока не найден.',
    'profile not found': 'Профиль пока не заполнен.',
    'refresh token cookie is required': 'Сессия истекла. Войдите снова.',
    'token expired': 'Сессия истекла. Войдите снова.',
    'token revoked': 'Сессия завершена. Войдите снова.',
    'user already exists': 'Пользователь с таким email уже существует.',
    'user not found': 'Пользователь не найден.',
  }

  return fieldTranslations[normalizedMessage] ?? commonTranslations[normalizedMessage] ?? normalizedMessage
}
