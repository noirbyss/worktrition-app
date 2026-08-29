import { useEffect, useState } from 'react'
import type { CurrentUser } from '../api'
import { useAuth } from '../auth/useAuth'
import { AppFrame } from '../components/app/AppFrame'
import { InlineMessage } from '../components/auth/InlineMessage'
import { AppLink } from '../components/navigation/AppLink'
import { navigate } from '../router'
import { formatDate, formatDateTimeFromUnix, toErrorMessage } from '../utils'

export function HomePage() {
  const { getCurrentUser, logout, session } = useAuth()
  const [isLoading, setIsLoading] = useState(true)
  const [isLoggingOut, setIsLoggingOut] = useState(false)
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
          setLoadError(toErrorMessage(error, 'Не удалось загрузить данные аккаунта.'))
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

  const handleLogout = async () => {
    try {
      setIsLoggingOut(true)
      await logout()
      navigate('/login', { replace: true })
    } finally {
      setIsLoggingOut(false)
    }
  }

  return (
    <AppFrame
      actions={
        <div className="button-row">
          <AppLink className="btn btn--ghost btn--small" href="/profile">
            Открыть профиль
          </AppLink>
          <AppLink className="btn btn--ghost btn--small" href="/nutrition">
            Питание
          </AppLink>
          <AppLink className="btn btn--ghost btn--small" href="/workouts">
            Тренировки
          </AppLink>
          <button
            className="btn btn--small"
            disabled={isLoggingOut}
            onClick={() => {
              void handleLogout()
            }}
            type="button"
          >
            {isLoggingOut ? 'ВЫХОД...' : 'ВЫЙТИ'}
          </button>
        </div>
      }
      description="После обязательной анкеты доступ к защищенным разделам открыт. Здесь можно проверить сессию и перейти в тестовые разделы приложения."
      title="Аккаунт"
    >
      {loadError ? <InlineMessage>{loadError}</InlineMessage> : null}

      <section className="panel-grid">
        <article className="panel">
          <h2 className="panel-title">Состояние авторизации</h2>
          <p className="panel-copy">
            {session?.profileCompleted
              ? 'Анкета сохранена, доступ к профилю, питанию и тренировкам открыт.'
              : 'Сессия активна, но анкета еще не завершена.'}
          </p>
          <div className="badge-row">
            <span className="badge">{session?.profileCompleted ? 'PROFILE READY' : 'PROFILE PENDING'}</span>
            <span className="badge badge--muted">USER {session?.userId ?? 'UNKNOWN'}</span>
          </div>
        </article>

        <article className="panel">
          <h2 className="panel-title">Access token</h2>
          <dl className="detail-list">
            <div>
              <dt>Истекает</dt>
              <dd>
                {session
                  ? formatDateTimeFromUnix(session.accessTokenExpiresAt)
                  : 'Нет активной сессии'}
              </dd>
            </div>
            <div>
              <dt>Обновление</dt>
              <dd>При `401` клиент делает запрос на `/auth/refresh` автоматически.</dd>
            </div>
          </dl>
        </article>
      </section>

      <section className="panel">
        <h2 className="panel-title">Данные пользователя</h2>
        {isLoading ? (
          <p className="panel-copy">Загружаем `/users/me` через gateway...</p>
        ) : user ? (
          <dl className="detail-list">
            <div>
              <dt>Имя</dt>
              <dd>{user.name}</dd>
            </div>
            <div>
              <dt>Email</dt>
              <dd>{user.email}</dd>
            </div>
            <div>
              <dt>Дата рождения</dt>
              <dd>{formatDate(user.birthDate)}</dd>
            </div>
            <div>
              <dt>Статус анкеты</dt>
              <dd>{user.profileCompleted ? 'Заполнена' : 'Не заполнена'}</dd>
            </div>
          </dl>
        ) : (
          <p className="panel-copy">Пользователь пока не загружен.</p>
        )}
      </section>
    </AppFrame>
  )
}
