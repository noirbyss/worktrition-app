import {
  createContext,
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
  type ReactNode,
} from 'react'
import {
  ApiError,
  type AuthSession,
  type CurrentUser,
  type LoginPayload,
  type Profile,
  type RegisterPayload,
  type SaveProfilePayload,
  type SaveProfileResult,
  getCurrentUser as fetchCurrentUser,
  getProfile as fetchProfile,
  login as loginRequest,
  logout as logoutRequest,
  refresh as refreshRequest,
  register as registerRequest,
  saveProfile as saveProfileRequest,
} from '../api'
import { loadAuthSession, persistAuthSession } from './auth-storage'

type AuthStatus = 'anonymous' | 'authenticated' | 'loading'

interface AuthContextValue {
  getCurrentUser: () => Promise<CurrentUser>
  getProfile: () => Promise<Profile>
  isAuthenticated: boolean
  login: (payload: LoginPayload) => Promise<AuthSession>
  logout: () => Promise<void>
  refreshSession: () => Promise<AuthSession>
  register: (payload: RegisterPayload) => Promise<AuthSession>
  saveProfile: (payload: SaveProfilePayload) => Promise<SaveProfileResult>
  session: AuthSession | null
  status: AuthStatus
}

const initialSession = loadAuthSession()

const AuthContext = createContext<AuthContextValue | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [session, setSession] = useState<AuthSession | null>(initialSession)
  const [status, setStatus] = useState<AuthStatus>(() =>
    isSessionFresh(initialSession) ? 'authenticated' : 'loading',
  )
  const refreshPromiseRef = useRef<Promise<AuthSession> | null>(null)
  const sessionRef = useRef<AuthSession | null>(initialSession)

  const clearSession = useCallback(() => {
    sessionRef.current = null
    persistAuthSession(null)
    setSession(null)
    setStatus('anonymous')
  }, [])

  const applySession = useCallback((nextSession: AuthSession) => {
    sessionRef.current = nextSession
    persistAuthSession(nextSession)
    setSession(nextSession)
    setStatus('authenticated')
  }, [])

  const markProfileCompleted = useCallback(() => {
    const currentSession = sessionRef.current
    if (!currentSession || currentSession.profileCompleted) {
      return
    }

    applySession({
      ...currentSession,
      profileCompleted: true,
    })
  }, [applySession])

  const refreshSession = useCallback(async () => {
    if (refreshPromiseRef.current) {
      return refreshPromiseRef.current
    }

    const refreshPromise = refreshRequest()
      .then((nextSession) => {
        applySession(nextSession)
        return nextSession
      })
      .catch((error: unknown) => {
        if (error instanceof ApiError && error.status === 401) {
          clearSession()
        } else if (!sessionRef.current) {
          setStatus('anonymous')
        }

        throw error
      })
      .finally(() => {
        refreshPromiseRef.current = null
      })

    refreshPromiseRef.current = refreshPromise

    return refreshPromise
  }, [applySession, clearSession])

  useEffect(() => {
    if (isSessionFresh(sessionRef.current)) {
      setStatus('authenticated')
      return
    }

    void refreshSession().catch(() => {
      if (!sessionRef.current) {
        setStatus('anonymous')
      }
    })
  }, [refreshSession])

  const getValidSession = useCallback(async () => {
    const currentSession = sessionRef.current
    if (currentSession && isSessionFresh(currentSession)) {
      return currentSession
    }

    return refreshSession()
  }, [refreshSession])

  const withAuthorizedSession = useCallback(
    async <T,>(request: (activeSession: AuthSession) => Promise<T>) => {
      let activeSession = await getValidSession()

      try {
        return await request(activeSession)
      } catch (error) {
        if (!(error instanceof ApiError) || error.status !== 401) {
          throw error
        }

        activeSession = await refreshSession()

        return request(activeSession)
      }
    },
    [getValidSession, refreshSession],
  )

  const login = useCallback(
    async (payload: LoginPayload) => {
      const nextSession = await loginRequest(payload)
      applySession(nextSession)
      return nextSession
    },
    [applySession],
  )

  const register = useCallback(
    async (payload: RegisterPayload) => {
      const nextSession = await registerRequest(payload)
      applySession(nextSession)
      return nextSession
    },
    [applySession],
  )

  const logout = useCallback(async () => {
    try {
      await logoutRequest()
    } finally {
      clearSession()
    }
  }, [clearSession])

  const getCurrentUser = useCallback(
    () => withAuthorizedSession((activeSession) => fetchCurrentUser(activeSession.accessToken)),
    [withAuthorizedSession],
  )

  const getProfile = useCallback(
    () => withAuthorizedSession((activeSession) => fetchProfile(activeSession.accessToken)),
    [withAuthorizedSession],
  )

  const saveProfile = useCallback(
    async (payload: SaveProfilePayload) => {
      const response = await withAuthorizedSession((activeSession) =>
        saveProfileRequest(activeSession.accessToken, payload),
      )

      if (response.profileCompleted) {
        markProfileCompleted()
      }

      return response
    },
    [markProfileCompleted, withAuthorizedSession],
  )

  const value = useMemo<AuthContextValue>(
    () => ({
      getCurrentUser,
      getProfile,
      isAuthenticated: status === 'authenticated',
      login,
      logout,
      refreshSession,
      register,
      saveProfile,
      session,
      status,
    }),
    [getCurrentUser, getProfile, login, logout, refreshSession, register, saveProfile, session, status],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export { AuthContext }

function isSessionFresh(activeSession: AuthSession | null) {
  if (!activeSession) {
    return false
  }

  return activeSession.accessTokenExpiresAt * 1000 - Date.now() > 30_000
}
