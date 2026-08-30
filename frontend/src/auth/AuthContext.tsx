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
  type GenerationResult,
  type LoginPayload,
  type NutritionDayPlan,
  type NutritionStats,
  type PlanType,
  type Profile,
  type RegisterPayload,
  type SaveProfilePayload,
  type SaveProfileResult,
  type WorkoutDayPlan,
  type WorkoutStats,
  completeNutritionMeal as completeNutritionMealRequest,
  completeNutritionWater as completeNutritionWaterRequest,
  completeWorkoutTraining as completeWorkoutTrainingRequest,
  getGenerationStatus as getGenerationStatusRequest,
  getCurrentUser as fetchCurrentUser,
  getNutritionDayPlan as getNutritionDayPlanRequest,
  getNutritionStats as getNutritionStatsRequest,
  getProfile as fetchProfile,
  getWorkoutDayPlan as getWorkoutDayPlanRequest,
  getWorkoutStats as getWorkoutStatsRequest,
  login as loginRequest,
  logout as logoutRequest,
  refresh as refreshRequest,
  register as registerRequest,
  saveProfile as saveProfileRequest,
  startGeneration as startGenerationRequest,
} from '../api'
import { loadAuthSession, persistAuthSession } from './auth-storage'

type AuthStatus = 'anonymous' | 'authenticated' | 'loading'

interface AuthContextValue {
  completeNutritionMeal: (mealItemId: number) => Promise<void>
  completeNutritionWater: (amountMl: number) => Promise<void>
  completeWorkoutTraining: (dayOfWeek: number | string, durationSeconds: number) => Promise<void>
  getGenerationStatus: (generationId: string) => Promise<GenerationResult>
  getCurrentUser: () => Promise<CurrentUser>
  getNutritionDayPlan: (dayOfWeek: number | string) => Promise<NutritionDayPlan>
  getNutritionStats: () => Promise<NutritionStats>
  getProfile: () => Promise<Profile>
  getWorkoutDayPlan: (dayOfWeek: number | string) => Promise<WorkoutDayPlan>
  getWorkoutStats: () => Promise<WorkoutStats>
  isAuthenticated: boolean
  login: (payload: LoginPayload) => Promise<AuthSession>
  logout: () => Promise<void>
  refreshSession: () => Promise<AuthSession>
  register: (payload: RegisterPayload) => Promise<AuthSession>
  saveProfile: (payload: SaveProfilePayload) => Promise<SaveProfileResult>
  session: AuthSession | null
  startGeneration: (planType: PlanType) => Promise<GenerationResult>
  status: AuthStatus
}

const initialSession = loadAuthSession()

const AuthContext = createContext<AuthContextValue | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [session, setSession] = useState<AuthSession | null>(initialSession)
  const [status, setStatus] = useState<AuthStatus>(() => {
    if (!initialSession) {
      return 'anonymous'
    }

    return isSessionFresh(initialSession) ? 'authenticated' : 'loading'
  })
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
    const currentSession = sessionRef.current

    if (!currentSession) {
      setStatus('anonymous')
      void refreshSession().catch(() => {
        if (!sessionRef.current) {
          setStatus('anonymous')
        }
      })
      return
    }

    if (isSessionFresh(currentSession)) {
      setStatus('authenticated')
      return
    }

    void refreshSession().catch(() => {
      clearSession()
    })
  }, [clearSession, refreshSession])

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

  const startGeneration = useCallback(
    (planType: PlanType) =>
      withAuthorizedSession((activeSession) =>
        startGenerationRequest(activeSession.accessToken, planType),
      ),
    [withAuthorizedSession],
  )

  const getGenerationStatus = useCallback(
    (generationId: string) =>
      withAuthorizedSession((activeSession) =>
        getGenerationStatusRequest(activeSession.accessToken, generationId),
      ),
    [withAuthorizedSession],
  )

  const getNutritionDayPlan = useCallback(
    (dayOfWeek: number | string) =>
      withAuthorizedSession((activeSession) =>
        getNutritionDayPlanRequest(activeSession.accessToken, dayOfWeek),
      ),
    [withAuthorizedSession],
  )

  const getNutritionStats = useCallback(
    () => withAuthorizedSession((activeSession) => getNutritionStatsRequest(activeSession.accessToken)),
    [withAuthorizedSession],
  )

  const completeNutritionMeal = useCallback(
    async (mealItemId: number) => {
      await withAuthorizedSession((activeSession) =>
        completeNutritionMealRequest(activeSession.accessToken, mealItemId),
      )
    },
    [withAuthorizedSession],
  )

  const completeNutritionWater = useCallback(
    async (amountMl: number) => {
      await withAuthorizedSession((activeSession) =>
        completeNutritionWaterRequest(activeSession.accessToken, amountMl),
      )
    },
    [withAuthorizedSession],
  )

  const getWorkoutDayPlan = useCallback(
    (dayOfWeek: number | string) =>
      withAuthorizedSession((activeSession) =>
        getWorkoutDayPlanRequest(activeSession.accessToken, dayOfWeek),
      ),
    [withAuthorizedSession],
  )

  const getWorkoutStats = useCallback(
    () => withAuthorizedSession((activeSession) => getWorkoutStatsRequest(activeSession.accessToken)),
    [withAuthorizedSession],
  )

  const completeWorkoutTraining = useCallback(
    async (dayOfWeek: number | string, durationSeconds: number) => {
      await withAuthorizedSession((activeSession) =>
        completeWorkoutTrainingRequest(activeSession.accessToken, dayOfWeek, durationSeconds),
      )
    },
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
      completeNutritionMeal,
      completeNutritionWater,
      completeWorkoutTraining,
      getGenerationStatus,
      getCurrentUser,
      getNutritionDayPlan,
      getNutritionStats,
      getProfile,
      getWorkoutDayPlan,
      getWorkoutStats,
      isAuthenticated: status === 'authenticated',
      login,
      logout,
      refreshSession,
      register,
      saveProfile,
      session,
      startGeneration,
      status,
    }),
    [
      completeNutritionMeal,
      completeNutritionWater,
      completeWorkoutTraining,
      getCurrentUser,
      getGenerationStatus,
      getNutritionDayPlan,
      getNutritionStats,
      getProfile,
      getWorkoutDayPlan,
      getWorkoutStats,
      login,
      logout,
      refreshSession,
      register,
      saveProfile,
      session,
      startGeneration,
      status,
    ],
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
