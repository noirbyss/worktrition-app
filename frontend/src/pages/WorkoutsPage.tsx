import { useCallback, useEffect, useState } from 'react'
import { ApiError, type WorkoutDayPlan, type WorkoutStats } from '../api'
import { useAuth } from '../auth/useAuth'
import { AppFrame } from '../components/app/AppFrame'
import { InlineMessage } from '../components/auth/InlineMessage'
import { NotchBar } from '../components/app/PlaceholderUi'
import { useCurrentUserData } from '../hooks'
import {
  clearPendingPlanGeneration,
  loadPendingPlanGeneration,
  persistPendingPlanGeneration,
  subscribePendingPlanGeneration,
  type PendingPlanGeneration,
} from '../store/plan-generation'
import { toErrorMessage } from '../utils'

const dayOptions = [
  { fullLabel: 'Понедельник', shortLabel: 'Пн', value: 1 },
  { fullLabel: 'Вторник', shortLabel: 'Вт', value: 2 },
  { fullLabel: 'Среда', shortLabel: 'Ср', value: 3 },
  { fullLabel: 'Четверг', shortLabel: 'Чт', value: 4 },
  { fullLabel: 'Пятница', shortLabel: 'Пт', value: 5 },
  { fullLabel: 'Суббота', shortLabel: 'Сб', value: 6 },
  { fullLabel: 'Воскресенье', shortLabel: 'Вс', value: 7 },
] as const

const durationOptions = [20, 30, 45, 60] as const
const generationPollIntervalMs = 2500
const generationStateCheckIntervalMs = 400

export function WorkoutsPage() {
  const {
    completeWorkoutTraining,
    getGenerationStatus,
    getWorkoutDayPlan,
    getWorkoutStats,
    session,
    startGeneration,
  } = useAuth()
  const { isLoading: isLoadingUser, loadError: userLoadError, user } = useCurrentUserData()
  const [activeDay, setActiveDay] = useState<number>(() => getCurrentDayValue())
  const [dayPlan, setDayPlan] = useState<WorkoutDayPlan | null>(null)
  const [durationMinutes, setDurationMinutes] = useState<number>(durationOptions[2])
  const [generationError, setGenerationError] = useState<string | null>(null)
  const [isCompletingWorkout, setIsCompletingWorkout] = useState(false)
  const [isLoadingPlan, setIsLoadingPlan] = useState(true)
  const [isStartingGeneration, setIsStartingGeneration] = useState(false)
  const [loadError, setLoadError] = useState<string | null>(null)
  const [pendingGeneration, setPendingGeneration] = useState<PendingPlanGeneration | null>(() =>
    loadPendingPlanGeneration(),
  )
  const [stats, setStats] = useState<WorkoutStats | null>(null)
  const [updateError, setUpdateError] = useState<string | null>(null)

  const activeDayOption = dayOptions.find((option) => option.value === activeDay) ?? dayOptions[0]
  const isWorkoutGenerationPending =
    pendingGeneration?.planType === 'workout' || pendingGeneration?.planType === 'all'
  const generationActionLabel =
    pendingGeneration?.planType === 'all' ? 'Собираем стартовый план' : 'Собираем тренировки'
  const isWorkoutCompleted = dayPlan?.isCompleted ?? false
  const hasExercises = (dayPlan?.exercises.length ?? 0) > 0
  const planTypeLabel = formatTrainingTypeLabel(dayPlan?.type)
  const planStatusLabel = hasExercises
    ? isWorkoutCompleted
      ? 'выполнено'
      : 'в плане'
    : 'отдых'
  const planStatusClassName = hasExercises
    ? isWorkoutCompleted
      ? 'pill done'
      : 'pill pending'
    : 'pill'

  const readWorkoutData = useCallback(
    async (dayValue: number) => {
      const [planResult, statsResult] = await Promise.allSettled([
        getWorkoutDayPlan(dayValue),
        getWorkoutStats(),
      ])

      let nextDayPlan: WorkoutDayPlan | null = null
      let nextStats: WorkoutStats | null = null
      let nextError: string | null = null

      if (planResult.status === 'fulfilled') {
        nextDayPlan = planResult.value
      } else if (!(planResult.reason instanceof ApiError && planResult.reason.status === 404)) {
        nextError = toErrorMessage(planResult.reason, 'Не удалось загрузить план тренировок.')
      }

      if (statsResult.status === 'fulfilled') {
        nextStats = statsResult.value
      } else if (!(statsResult.reason instanceof ApiError && statsResult.reason.status === 404)) {
        nextError ??= toErrorMessage(statsResult.reason, 'Не удалось загрузить статистику тренировок.')
      }

      return {
        dayPlan: nextDayPlan,
        error: nextError,
        stats: nextStats,
      }
    },
    [getWorkoutDayPlan, getWorkoutStats],
  )

  const applyWorkoutData = useCallback(
    async (dayValue: number) => {
      const snapshot = await readWorkoutData(dayValue)
      setDayPlan(snapshot.dayPlan)
      setLoadError(snapshot.error)
      setStats(snapshot.stats)
    },
    [readWorkoutData],
  )

  useEffect(() => {
    return subscribePendingPlanGeneration(() => {
      setPendingGeneration(loadPendingPlanGeneration())
    })
  }, [])

  useEffect(() => {
    let isCancelled = false

    const loadPage = async () => {
      setIsLoadingPlan(true)

      try {
        const snapshot = await readWorkoutData(activeDay)

        if (isCancelled) {
          return
        }

        setDayPlan(snapshot.dayPlan)
        setLoadError(snapshot.error)
        setStats(snapshot.stats)
      } finally {
        if (!isCancelled) {
          setIsLoadingPlan(false)
        }
      }
    }

    void loadPage()

    return () => {
      isCancelled = true
    }
  }, [activeDay, readWorkoutData])

  useEffect(() => {
    if (!pendingGeneration) {
      return
    }

    let isCancelled = false
    let timeoutId: number | null = null

    const monitorGeneration = async () => {
      if (!pendingGeneration.generationId) {
        timeoutId = window.setTimeout(() => {
          if (!isCancelled) {
            setPendingGeneration(loadPendingPlanGeneration())
          }
        }, generationStateCheckIntervalMs)
        return
      }

      try {
        const generation = await getGenerationStatus(pendingGeneration.generationId)

        if (isCancelled) {
          return
        }

        if (generation.status === 'done') {
          clearPendingPlanGeneration()
          setGenerationError(null)
          await applyWorkoutData(activeDay)
          return
        }

        if (generation.status === 'failed') {
          clearPendingPlanGeneration()
          setGenerationError(generation.errorMessage ?? 'Не удалось сгенерировать тренировки.')
          return
        }

        timeoutId = window.setTimeout(() => {
          void monitorGeneration()
        }, generationPollIntervalMs)
      } catch (error) {
        if (isCancelled) {
          return
        }

        setGenerationError(toErrorMessage(error, 'Не удалось проверить статус генерации.'))
        timeoutId = window.setTimeout(() => {
          void monitorGeneration()
        }, generationPollIntervalMs)
      }
    }

    void monitorGeneration()

    return () => {
      isCancelled = true
      if (timeoutId !== null) {
        window.clearTimeout(timeoutId)
      }
    }
  }, [activeDay, applyWorkoutData, getGenerationStatus, pendingGeneration])

  const handleRegenerateWorkout = async () => {
    try {
      setGenerationError(null)
      setIsStartingGeneration(true)
      setUpdateError(null)

      persistPendingPlanGeneration({
        planType: 'workout',
        status: 'starting',
      })

      const generation = await startGeneration('workout')

      if (generation.status === 'failed') {
        clearPendingPlanGeneration()
        throw new Error(generation.errorMessage ?? 'Не удалось запустить генерацию тренировок.')
      }

      if (generation.status === 'done') {
        clearPendingPlanGeneration()
        await applyWorkoutData(activeDay)
      } else {
        persistPendingPlanGeneration({
          generationId: generation.generationId,
          planType: 'workout',
          status: 'pending',
        })
      }
    } catch (error) {
      clearPendingPlanGeneration()
      setUpdateError(toErrorMessage(error, 'Не удалось обновить тренировки.'))
    } finally {
      setIsStartingGeneration(false)
    }
  }

  const handleCompleteWorkout = async () => {
    if (!dayPlan || !session?.userId || !hasExercises) {
      return
    }

    try {
      setIsCompletingWorkout(true)
      setUpdateError(null)

      await completeWorkoutTraining(activeDay, durationMinutes * 60)
      await applyWorkoutData(activeDay)
    } catch (error) {
      setUpdateError(toErrorMessage(error, 'Не удалось сохранить тренировку.'))
    } finally {
      setIsCompletingWorkout(false)
    }
  }

  const isGenerationInProgress = pendingGeneration !== null || isStartingGeneration
  const isWorkoutRefreshLoading = isStartingGeneration || isWorkoutGenerationPending
  const generationButtonLabel = isGenerationInProgress
    ? 'ГЕНЕРАЦИЯ ИДЕТ...'
    : dayPlan
      ? 'ПОМЕНЯТЬ ТРЕНИРОВКИ'
      : 'СГЕНЕРИРОВАТЬ ТРЕНИРОВКИ'

  return (
    <AppFrame
      actions={
        <button
          className="btn btn-secondary header-cta"
          disabled={isGenerationInProgress}
          onClick={() => {
            void handleRegenerateWorkout()
          }}
          type="button"
        >
          {generationButtonLabel}
        </button>
      }
      currentUser={user}
      description="План по дням, отметка выполнения и статистика."
      eyebrow="Экран 02"
      isCurrentUserLoading={isLoadingUser}
      title="Тренировки"
    >
      {isWorkoutRefreshLoading ? (
        <section className="nutrition-refresh-loading" aria-busy="true" aria-live="polite">
          <div className="nutrition-refresh-loading__spinner" />
        </section>
      ) : (
        <>
          {userLoadError ? <InlineMessage>{userLoadError}</InlineMessage> : null}
          {loadError ? <InlineMessage>{loadError}</InlineMessage> : null}
          {updateError ? <InlineMessage>{updateError}</InlineMessage> : null}
          {generationError ? <InlineMessage>{generationError}</InlineMessage> : null}

          {isWorkoutGenerationPending ? (
            <section className="card frame placeholder-section nutrition-generation">
              <div className="section-kicker">AI · В ПРОЦЕССЕ</div>
              <h2 className="wizard-step-title nutrition-generation__title">{generationActionLabel}</h2>
            </section>
          ) : null}

          {stats ? (
            <section className="grid g-3 placeholder-section">
              <SummaryCard
                label="Выполнение"
                note="по активному плану"
                value={formatPercent(stats.percentagePlanFulfilled)}
              />
              <SummaryCard
                label="Серия"
                note="дней подряд"
                value={formatStreak(stats.currentStreakDays)}
              />
              <SummaryCard
                label="Время"
                note="всего"
                value={formatDuration(stats.totalTrainingTimeSeconds)}
              />
            </section>
          ) : null}

          <div className="day-tabs">
            {dayOptions.map((day) => (
              <button
                className={day.value === activeDay ? 'day-tab active' : 'day-tab'}
                key={day.value}
                onClick={() => {
                  setActiveDay(day.value)
                }}
                type="button"
              >
                {day.shortLabel}
              </button>
            ))}
          </div>

          {isLoadingPlan && !dayPlan ? (
            <section className="card empty-state">
              <div className="card-title">Загрузка</div>
              <p className="panel-copy">Подтягиваем тренировку на {activeDayOption.fullLabel.toLowerCase()}.</p>
            </section>
          ) : null}

          {!isLoadingPlan && !dayPlan ? (
            <section className="card empty-state">
              <div className="card-title">План не найден</div>
              <p className="panel-copy">Для выбранного дня тренировка пока не сохранена.</p>
            </section>
          ) : null}

          {dayPlan ? (
            <section className="grid g-2">
              <div className="card frame">
                <div className="card-title">
                  <span>План на {activeDayOption.fullLabel.toLowerCase()}</span>
                  <span className={planStatusClassName}>{planStatusLabel}</span>
                </div>

                <p className="panel-copy workout-copy">
                  {planTypeLabel} · {formatExercisesCount(dayPlan.exercises.length)}
                </p>

                {hasExercises ? (
                  dayPlan.exercises.map((exercise, index) => (
                    <div className="workout-exercise" key={`${exercise}-${index}`}>
                      <span className="workout-exercise__index">{String(index + 1).padStart(2, '0')}</span>
                      <span className="workout-exercise__name">{exercise}</span>
                    </div>
                  ))
                ) : (
                  <div className="empty-state workout-inline-empty">
                    <p className="panel-copy">На этот день тренировка не назначена.</p>
                  </div>
                )}

                {hasExercises ? (
                  <div className="workout-actions">
                    <label className="workout-duration" htmlFor="workoutDuration">
                      <span className="workout-duration__label">Длительность</span>
                      <select
                        className="wizard-select workout-duration__select"
                        id="workoutDuration"
                        onChange={(event) => {
                          setDurationMinutes(Number(event.target.value))
                        }}
                        value={durationMinutes}
                      >
                        {durationOptions.map((value) => (
                          <option key={value} value={value}>
                            {value} мин
                          </option>
                        ))}
                      </select>
                    </label>

                    <button
                      className="btn btn-secondary workout-complete-button"
                      disabled={isWorkoutCompleted || isCompletingWorkout}
                      onClick={() => {
                        void handleCompleteWorkout()
                      }}
                      type="button"
                    >
                      {isCompletingWorkout ? 'СОХРАНЯЕМ...' : isWorkoutCompleted ? 'ВЫПОЛНЕНО' : 'ЗАВЕРШИТЬ'}
                    </button>
                  </div>
                ) : null}

                <div className="reward">
                  <b>{formatDuration(durationMinutes * 60)}</b>
                  <span>сохранится в статистику после завершения</span>
                </div>
              </div>

              <div className="status-stack">
                <div className="card">
                  <div className="card-title">День</div>
                  <MetricBlock
                    active={20}
                    label="Тип"
                    value={planTypeLabel}
                    variant={getTrainingVariant(dayPlan.type)}
                  />
                  <hr className="rule" />
                  <MetricBlock
                    active={Math.min(dayPlan.exercises.length, 20)}
                    label="Упражнения"
                    marginTop
                    value={String(dayPlan.exercises.length)}
                  />
                  <MetricBlock
                    active={Math.min(Math.round((durationMinutes / durationOptions[durationOptions.length - 1]) * 20), 20)}
                    label="Сессия"
                    marginTop
                    value={`${durationMinutes} мин`}
                    variant="discipline"
                  />
                </div>

                <div className="card">
                  <div className="card-title">Прогресс</div>
                  <MetricBlock
                    active={toNotchCount(stats?.percentagePlanFulfilled ?? 0)}
                    label="Выполнение"
                    value={formatPercent(stats?.percentagePlanFulfilled ?? 0)}
                  />
                  <MetricBlock
                    active={Math.min((stats?.currentStreakDays ?? 0) * 2, 20)}
                    label="Серия"
                    marginTop
                    value={formatStreak(stats?.currentStreakDays ?? 0)}
                    variant="discipline"
                  />
                  <MetricBlock
                    active={toTrainingTimeNotches(stats?.totalTrainingTimeSeconds ?? 0)}
                    label="Время"
                    marginTop
                    value={formatDuration(stats?.totalTrainingTimeSeconds ?? 0)}
                    variant="endurance"
                  />
                </div>
              </div>
            </section>
          ) : null}

          <footer className="foot">worktrition · тренировки · backend</footer>
        </>
      )}
    </AppFrame>
  )
}

function SummaryCard({
  label,
  note,
  value,
}: {
  label: string
  note: string
  value: string
}) {
  return (
    <div className="card info-card">
      <div className="info-label">{label}</div>
      <div className="info-value">{value}</div>
      <div className="info-sub">{note}</div>
    </div>
  )
}

function MetricBlock({
  active,
  label,
  marginTop = false,
  value,
  variant,
}: {
  active: number
  label: string
  marginTop?: boolean
  value: string
  variant?: 'balance' | 'discipline' | 'endurance' | 'strength' | 'water'
}) {
  return (
    <div className={marginTop ? 'metric-block metric-block--spaced' : 'metric-block'}>
      <div className="bar-row">
        <span className="bar-label">{label}</span>
        <span className="bar-value">{value}</span>
      </div>
      <NotchBar active={active} total={20} variant={variant} />
    </div>
  )
}

function getCurrentDayValue() {
  const day = new Date().getDay()

  if (day === 0) {
    return 7
  }

  return day
}

function normalizePercentage(value: number) {
  if (!Number.isFinite(value)) {
    return 0
  }

  return Math.max(0, Math.min(100, value))
}

function toNotchCount(percentage: number) {
  return Math.max(0, Math.min(20, Math.round((normalizePercentage(percentage) / 100) * 20)))
}

function toTrainingTimeNotches(durationSeconds: number) {
  const durationMinutes = normalizeNonNegativeNumber(durationSeconds) / 60
  return Math.max(0, Math.min(20, Math.round((durationMinutes / 600) * 20)))
}

function formatPercent(value: number) {
  return `${Math.round(normalizePercentage(value))}%`
}

function formatDuration(value: number) {
  const totalMinutes = Math.round(normalizeNonNegativeNumber(value) / 60)
  const hours = Math.floor(totalMinutes / 60)
  const minutes = totalMinutes % 60

  if (hours === 0) {
    return `${minutes} мин`
  }

  if (minutes === 0) {
    return `${hours} ч`
  }

  return `${hours} ч ${minutes} мин`
}

function formatExercisesCount(count: number) {
  const safeCount = Math.max(0, Math.floor(normalizeNonNegativeNumber(count)))

  if (safeCount === 1) {
    return '1 упражнение'
  }

  if (safeCount >= 2 && safeCount <= 4) {
    return `${safeCount} упражнения`
  }

  return `${safeCount} упражнений`
}

function formatTrainingTypeLabel(value?: string) {
  const normalized = value?.trim()

  if (!normalized) {
    return 'Тренировка'
  }

  return normalized
}

function formatStreak(value: number) {
  const safeValue = Math.max(0, Math.floor(normalizeNonNegativeNumber(value)))

  if (safeValue === 1) {
    return '1 день'
  }

  if (safeValue >= 2 && safeValue <= 4) {
    return `${safeValue} дня`
  }

  return `${safeValue} дней`
}

function getTrainingVariant(value: string): 'endurance' | 'strength' {
  const normalized = value.trim().toLowerCase()

  if (
    normalized.includes('кардио') ||
    normalized.includes('бег') ||
    normalized.includes('вынослив')
  ) {
    return 'endurance'
  }

  return 'strength'
}

function normalizeNonNegativeNumber(value: number) {
  if (!Number.isFinite(value)) {
    return 0
  }

  return Math.max(0, value)
}
