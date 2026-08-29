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
    return error.message
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
