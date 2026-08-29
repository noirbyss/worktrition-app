import { useEffect, useState } from 'react'
import type { CurrentUser } from '../api'
import { useAuth } from '../auth/useAuth'
import { toErrorMessage } from '../utils'

export function useCurrentUserData() {
  const { getCurrentUser } = useAuth()
  const [isLoading, setIsLoading] = useState(true)
  const [loadError, setLoadError] = useState<string | null>(null)
  const [user, setUser] = useState<CurrentUser | null>(null)

  useEffect(() => {
    let isCancelled = false

    const loadUser = async () => {
      try {
        setIsLoading(true)
        setLoadError(null)
        const nextUser = await getCurrentUser()

        if (!isCancelled) {
          setUser(nextUser)
        }
      } catch (error) {
        if (!isCancelled) {
          setLoadError(toErrorMessage(error, 'Не удалось загрузить данные пользователя.'))
        }
      } finally {
        if (!isCancelled) {
          setIsLoading(false)
        }
      }
    }

    void loadUser()

    return () => {
      isCancelled = true
    }
  }, [getCurrentUser])

  return {
    isLoading,
    loadError,
    user,
  }
}
