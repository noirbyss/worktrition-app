import { useEffect, useState } from 'react'
import { ApiError, type CurrentUser, type Profile } from '../api'
import { useAuth } from '../auth/useAuth'
import { AppFrame } from '../components/app/AppFrame'
import { InlineMessage } from '../components/auth/InlineMessage'
import { AppLink } from '../components/navigation/AppLink'
import { formatDate, formatList, toErrorMessage } from '../utils'

const genderLabels: Record<number, string> = {
  1: 'Мужской',
  2: 'Женский',
}

const trainingLevelLabels: Record<number, string> = {
  1: 'Начальный',
  2: 'Средний',
  3: 'Продвинутый',
}

const activityLevelLabels: Record<number, string> = {
  1: 'Низкая активность',
  2: 'Легкая активность',
  3: 'Умеренная активность',
  4: 'Высокая активность',
}

const goalLabels: Record<number, string> = {
  1: 'Снижение веса',
  2: 'Поддержание веса',
  3: 'Набор мышц',
}

const trainingLocationLabels: Record<number, string> = {
  1: 'Дома',
  2: 'В зале',
}

export function ProfilePage() {
  const { getCurrentUser, getProfile, session } = useAuth()
  const [isLoading, setIsLoading] = useState(true)
  const [loadError, setLoadError] = useState<string | null>(null)
  const [profile, setProfile] = useState<Profile | null>(null)
  const [profileMissing, setProfileMissing] = useState(false)
  const [user, setUser] = useState<CurrentUser | null>(null)

  useEffect(() => {
    let isCancelled = false

    const loadProfile = async () => {
      try {
        setIsLoading(true)
        setLoadError(null)
        const nextUser = await getCurrentUser()
        let nextProfile: Profile | null = null
        let isCurrentProfileMissing = false

        try {
          nextProfile = await getProfile()
        } catch (error) {
          if (error instanceof ApiError && error.status === 404) {
            isCurrentProfileMissing = true
          } else {
            throw error
          }
        }

        if (!isCancelled) {
          setUser(nextUser)
          setProfile(nextProfile)
          setProfileMissing(isCurrentProfileMissing)
        }
      } catch (error) {
        if (!isCancelled) {
          setLoadError(toErrorMessage(error, 'Не удалось загрузить профиль пользователя.'))
        }
      } finally {
        if (!isCancelled) {
          setIsLoading(false)
        }
      }
    }

    void loadProfile()

    return () => {
      isCancelled = true
    }
  }, [getCurrentUser, getProfile])

  return (
    <AppFrame
      actions={
        <div className="button-row">
          <AppLink className="btn btn--ghost btn--small" href="/app">
            Назад в аккаунт
          </AppLink>
        </div>
      }
      description="Страница читает `/profile` и показывает понятное состояние, если анкета еще не создана."
      title="Профиль"
    >
      {loadError ? <InlineMessage>{loadError}</InlineMessage> : null}

      <section className="panel">
        <h2 className="panel-title">Сводка</h2>
        {isLoading ? (
          <p className="panel-copy">Загружаем данные профиля...</p>
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
              <dt>Флаг профиля</dt>
              <dd>{session?.profileCompleted ? 'Заполнен' : 'Не заполнен'}</dd>
            </div>
          </dl>
        ) : (
          <p className="panel-copy">Нет данных пользователя.</p>
        )}
      </section>

      {profile ? (
        <section className="panel-grid">
          <article className="panel">
            <h2 className="panel-title">Физические параметры</h2>
            <dl className="detail-list">
              <div>
                <dt>Возраст</dt>
                <dd>{profile.age} лет</dd>
              </div>
              <div>
                <dt>Пол</dt>
                <dd>{genderLabels[profile.gender] ?? 'Не указан'}</dd>
              </div>
              <div>
                <dt>Рост</dt>
                <dd>{profile.heightCm} см</dd>
              </div>
              <div>
                <dt>Вес</dt>
                <dd>{profile.weightKg} кг</dd>
              </div>
              <div>
                <dt>Целевой вес</dt>
                <dd>{profile.targetWeightKg ? `${profile.targetWeightKg} кг` : 'Не указан'}</dd>
              </div>
            </dl>
          </article>

          <article className="panel">
            <h2 className="panel-title">Тренировки и питание</h2>
            <dl className="detail-list">
              <div>
                <dt>Уровень тренировки</dt>
                <dd>{trainingLevelLabels[profile.trainingLevel] ?? 'Не указан'}</dd>
              </div>
              <div>
                <dt>Активность</dt>
                <dd>{activityLevelLabels[profile.activityLevel] ?? 'Не указана'}</dd>
              </div>
              <div>
                <dt>Цель</dt>
                <dd>{goalLabels[profile.goal] ?? 'Не указана'}</dd>
              </div>
              <div>
                <dt>Место тренировок</dt>
                <dd>{trainingLocationLabels[profile.trainingLocation] ?? 'Не указано'}</dd>
              </div>
              <div>
                <dt>Дней в неделю</dt>
                <dd>{profile.trainingDaysPerWeek}</dd>
              </div>
              <div>
                <dt>Инвентарь</dt>
                <dd>{profile.equipment || 'Не указан'}</dd>
              </div>
            </dl>
          </article>

          <article className="panel panel--wide">
            <h2 className="panel-title">Списки предпочтений</h2>
            <dl className="detail-list">
              <div>
                <dt>Аллергии</dt>
                <dd>{formatList(profile.allergies)}</dd>
              </div>
              <div>
                <dt>Исключенные продукты</dt>
                <dd>{formatList(profile.excludedFoods)}</dd>
              </div>
              <div>
                <dt>Пищевые предпочтения</dt>
                <dd>{formatList(profile.foodPreferences)}</dd>
              </div>
            </dl>
          </article>
        </section>
      ) : profileMissing && !isLoading ? (
        <section className="panel">
          <h2 className="panel-title">Профиль еще не создан</h2>
          <p className="panel-copy">
            Пользователь уже авторизован, но gateway вернул `404 profile not found`.
            Это ожидаемое состояние после регистрации, пока анкета не заполнена.
          </p>
        </section>
      ) : null}
    </AppFrame>
  )
}
