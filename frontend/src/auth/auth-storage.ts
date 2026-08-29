import type { AuthSession } from '../api'

const STORAGE_KEY = 'worktrition.auth.session'

export function loadAuthSession() {
  try {
    const rawValue = window.localStorage.getItem(STORAGE_KEY)
    if (!rawValue) {
      return null
    }

    const parsedValue = JSON.parse(rawValue)
    if (!isAuthSession(parsedValue)) {
      window.localStorage.removeItem(STORAGE_KEY)
      return null
    }

    return parsedValue
  } catch {
    return null
  }
}

export function persistAuthSession(session: AuthSession | null) {
  try {
    if (!session) {
      window.localStorage.removeItem(STORAGE_KEY)
      return
    }

    window.localStorage.setItem(STORAGE_KEY, JSON.stringify(session))
  } catch {
    return
  }
}

function isAuthSession(value: unknown): value is AuthSession {
  if (!value || typeof value !== 'object') {
    return false
  }

  const session = value as Partial<AuthSession>

  return (
    typeof session.accessToken === 'string' &&
    typeof session.accessTokenExpiresAt === 'number' &&
    typeof session.profileCompleted === 'boolean' &&
    typeof session.userId === 'string'
  )
}
