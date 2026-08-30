import { useCallback, useEffect, useState } from 'react'
import { ApiError, type NutritionDayPlan, type NutritionStats } from '../api'
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
import {
  clearNutritionAccuracyPendingReset,
  loadNutritionUiState,
  markNutritionAccuracyPendingReset,
  subscribeNutritionUiState,
} from '../store/nutrition-ui'
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

const generationPollIntervalMs = 2500
const generationStateCheckIntervalMs = 400
const waterQuickActions = [250, 500] as const

export function NutritionPage() {
  const {
    completeNutritionMeal,
    completeNutritionWater,
    getGenerationStatus,
    getNutritionDayPlan,
    getNutritionStats,
    startGeneration,
  } = useAuth()
  const { isLoading: isLoadingUser, loadError: userLoadError, user } = useCurrentUserData()
  const [activeDay, setActiveDay] = useState<number>(() => getCurrentDayValue())
  const [dayPlan, setDayPlan] = useState<NutritionDayPlan | null>(null)
  const [generationError, setGenerationError] = useState<string | null>(null)
  const [isCompletingMealId, setIsCompletingMealId] = useState<number | null>(null)
  const [isCompletingWaterAmount, setIsCompletingWaterAmount] = useState<number | null>(null)
  const [isLoadingPlan, setIsLoadingPlan] = useState(true)
  const [isStartingGeneration, setIsStartingGeneration] = useState(false)
  const [loadError, setLoadError] = useState<string | null>(null)
  const [pendingGeneration, setPendingGeneration] = useState<PendingPlanGeneration | null>(() =>
    loadPendingPlanGeneration(),
  )
  const [nutritionUiState, setNutritionUiState] = useState(() => loadNutritionUiState())
  const [sessionWaterByDay, setSessionWaterByDay] = useState<Partial<Record<number, number>>>({})
  const [stats, setStats] = useState<NutritionStats | null>(null)
  const [updateError, setUpdateError] = useState<string | null>(null)

  const activeDayOption = dayOptions.find((option) => option.value === activeDay) ?? dayOptions[0]
  const todayDayValue = getCurrentDayValue()
  const isTodaySelected = activeDay === todayDayValue
  const generationActionLabel =
    pendingGeneration?.planType === 'all' ? 'Собираем ваш стартовый план' : 'Собираем новый рацион'
  const activeDaySessionWaterMl = sessionWaterByDay[activeDay] ?? 0
  const shouldMaskNutritionAccuracy = nutritionUiState.awaitingFreshNutritionAccuracy

  const readNutritionData = useCallback(
    async (dayValue: number) => {
      const [planResult, statsResult] = await Promise.allSettled([
        getNutritionDayPlan(dayValue),
        getNutritionStats(),
      ])

      let nextDayPlan: NutritionDayPlan | null = null
      let nextStats: NutritionStats | null = null
      let nextError: string | null = null

      if (planResult.status === 'fulfilled') {
        nextDayPlan = planResult.value
      } else if (!(planResult.reason instanceof ApiError && planResult.reason.status === 404)) {
        nextError = toErrorMessage(planResult.reason, 'Не удалось загрузить план питания.')
      }

      if (statsResult.status === 'fulfilled') {
        nextStats = statsResult.value
      } else if (!(statsResult.reason instanceof ApiError && statsResult.reason.status === 404)) {
        nextError ??= toErrorMessage(statsResult.reason, 'Не удалось загрузить статистику питания.')
      }

      return {
        dayPlan: nextDayPlan,
        error: nextError,
        stats: nextStats,
      }
    },
    [getNutritionDayPlan, getNutritionStats],
  )

  const applyNutritionData = useCallback(
    async (dayValue: number) => {
      const snapshot = await readNutritionData(dayValue)
      setDayPlan(snapshot.dayPlan)
      setLoadError(snapshot.error)
      setStats(snapshot.stats)
    },
    [readNutritionData],
  )

  useEffect(() => {
    return subscribePendingPlanGeneration(() => {
      setPendingGeneration(loadPendingPlanGeneration())
    })
  }, [])

  useEffect(() => {
    return subscribeNutritionUiState(() => {
      setNutritionUiState(loadNutritionUiState())
    })
  }, [])

  useEffect(() => {
    let isCancelled = false

    const loadPage = async () => {
      setIsLoadingPlan(true)

      try {
        const snapshot = await readNutritionData(activeDay)

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
  }, [activeDay, readNutritionData])

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
          if (pendingGeneration.planType === 'nutrition') {
            markNutritionAccuracyPendingReset()
          }
          clearPendingPlanGeneration()
          setGenerationError(null)
          await applyNutritionData(activeDay)
          return
        }

        if (generation.status === 'failed') {
          clearPendingPlanGeneration()
          setGenerationError(generation.errorMessage ?? 'Не удалось сгенерировать план питания.')
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
  }, [activeDay, applyNutritionData, getGenerationStatus, pendingGeneration])

  const handleRegenerateNutrition = async () => {
    try {
      setGenerationError(null)
      setIsStartingGeneration(true)
      setUpdateError(null)

      persistPendingPlanGeneration({
        planType: 'nutrition',
        status: 'starting',
      })

      const generation = await startGeneration('nutrition')

      if (generation.status === 'failed') {
        clearPendingPlanGeneration()
        throw new Error(generation.errorMessage ?? 'Не удалось запустить обновление питания.')
      }

      if (generation.status === 'done') {
        markNutritionAccuracyPendingReset()
        clearPendingPlanGeneration()
        await applyNutritionData(activeDay)
      } else {
        persistPendingPlanGeneration({
          generationId: generation.generationId,
          planType: 'nutrition',
          status: 'pending',
        })
      }
    } catch (error) {
      clearPendingPlanGeneration()
      setUpdateError(toErrorMessage(error, 'Не удалось обновить питание.'))
    } finally {
      setIsStartingGeneration(false)
    }
  }

  const handleCompleteMeal = async (mealItemId: number) => {
    try {
      setIsCompletingMealId(mealItemId)
      setUpdateError(null)
      await completeNutritionMeal(mealItemId)
      setDayPlan((current) =>
        current
          ? {
              ...current,
              meals: current.meals.map((meal) =>
                meal.id === mealItemId ? { ...meal, isCompleted: true } : meal,
              ),
            }
          : current,
      )
      clearNutritionAccuracyPendingReset()
      await applyNutritionData(activeDay)
    } catch (error) {
      setUpdateError(toErrorMessage(error, 'Не удалось отметить приём пищи.'))
    } finally {
      setIsCompletingMealId(null)
    }
  }

  const handleCompleteWater = async (amountMl: number) => {
    try {
      setIsCompletingWaterAmount(amountMl)
      setUpdateError(null)
      await completeNutritionWater(amountMl)
      setSessionWaterByDay((current) => ({
        ...current,
        [activeDay]: (current[activeDay] ?? 0) + amountMl,
      }))
      await applyNutritionData(activeDay)
    } catch (error) {
      setUpdateError(toErrorMessage(error, 'Не удалось сохранить воду.'))
    } finally {
      setIsCompletingWaterAmount(null)
    }
  }

  const isGenerationInProgress = pendingGeneration !== null || isStartingGeneration
  const generationButtonLabel = isGenerationInProgress
    ? 'ГЕНЕРАЦИЯ ИДЁТ...'
    : dayPlan
      ? 'ПОМЕНЯТЬ ПИТАНИЕ'
      : 'СГЕНЕРИРОВАТЬ ПИТАНИЕ'

  return (
    <AppFrame
      actions={
        <button
          className="btn btn-secondary header-cta"
          disabled={isGenerationInProgress}
          onClick={() => {
            void handleRegenerateNutrition()
          }}
          type="button"
        >
          {generationButtonLabel}
        </button>
      }
      currentUser={user}
      description="Рацион по дням недели, состав блюд и быстрые действия для фиксации питания и воды."
      eyebrow="Экран 03"
      isCurrentUserLoading={isLoadingUser}
      title="Питание"
    >
      {userLoadError ? <InlineMessage>{userLoadError}</InlineMessage> : null}
      {loadError ? <InlineMessage>{loadError}</InlineMessage> : null}
      {updateError ? <InlineMessage>{updateError}</InlineMessage> : null}
      {generationError ? <InlineMessage>{generationError}</InlineMessage> : null}

      {pendingGeneration ? (
        <section className="card frame placeholder-section nutrition-generation">
          <div className="section-kicker">AI · В ПРОЦЕССЕ</div>
          <h2 className="wizard-step-title nutrition-generation__title">{generationActionLabel}</h2>
        </section>
      ) : null}

      {stats ? (
        <section className="grid g-3 placeholder-section">
          <SummaryCard
            label="Точность по БЖУ"
            note={
              shouldMaskNutritionAccuracy
                ? 'после смены питания статистика начнёт собираться заново'
                : 'среднее попадание в дневные нормы'
            }
            value={formatPercent(stats.percentageComplianceNutritionFacts, shouldMaskNutritionAccuracy)}
          />
          <SummaryCard
            label="Выполнение плана"
            note="доля отмеченных приёмов пищи"
            value={formatPercent(stats.percentagePlanFulfilled)}
          />
          <SummaryCard
            label="Норма воды"
            note="среднее соблюдение водного режима"
            value={formatPercent(stats.percentageWaterStandardFulfillment)}
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
          <p className="panel-copy">
            Подтягиваем рацион на {activeDayOption.fullLabel.toLowerCase()} и статистику питания.
          </p>
        </section>
      ) : null}

      {!isLoadingPlan && !dayPlan ? (
        <section className="card empty-state">
          <div className="card-title">План пока не готов</div>
          <p className="panel-copy">
            Для выбранного дня ещё нет сохранённого рациона. Можно дождаться текущей генерации или
            запустить отдельное обновление питания.
          </p>
        </section>
      ) : null}

      {dayPlan ? (
        <section className="grid g-2">
          <div className="card frame">
            <div className="card-title">
              <span>Рацион на {activeDayOption.fullLabel.toLowerCase()}</span>
              <span className="pill pending">{formatMealsCount(dayPlan.meals.length)}</span>
            </div>

            <p className="panel-copy nutrition-meal-copy">
              Здесь собраны блюда на выбранный день. После еды можно сразу отметить приём пищи,
              чтобы обновить общую статистику по плану.
            </p>

            {dayPlan.meals.length > 0 ? (
              dayPlan.meals.map((meal) => {
                const isCompleted = meal.isCompleted
                const isSavingMeal = isCompletingMealId === meal.id

                return (
                  <div className="meal" key={meal.id}>
                    <div className="meal-top">
                      <span className="meal-name">{meal.name}</span>
                      <span className={isCompleted ? 'pill done' : 'pill pending'}>
                        {isCompleted ? 'отмечено' : 'в плане'}
                      </span>
                    </div>

                    <p className="meal-items nutrition-meal-recipe">
                      {meal.recipe.trim() || 'Рецепт для этого блюда не был передан backend.'}
                    </p>

                    <div className="macro-row">
                      <span className="macro">
                        <b>{formatNumber(meal.nutritionFacts.calories)}</b> ккал
                      </span>
                      <span className="macro">
                        Б <b>{formatNumber(meal.nutritionFacts.protein)}</b> г
                      </span>
                      <span className="macro">
                        Ж <b>{formatNumber(meal.nutritionFacts.fat)}</b> г
                      </span>
                      <span className="macro">
                        У <b>{formatNumber(meal.nutritionFacts.carb)}</b> г
                      </span>
                    </div>

                    <div className="nutrition-meal-actions">
                      <button
                        className="btn btn-secondary nutrition-meal-button"
                        disabled={isCompleted || isSavingMeal}
                        onClick={() => {
                          void handleCompleteMeal(meal.id)
                        }}
                        type="button"
                      >
                        {isSavingMeal ? 'СОХРАНЯЕМ...' : isCompleted ? 'ОТМЕЧЕНО' : 'ОТМЕТИТЬ КАК СЪЕДЕНОЕ'}
                      </button>
                    </div>
                  </div>
                )
              })
            ) : (
              <div className="empty-state nutrition-inline-empty">
                <p className="panel-copy">Для выбранного дня backend пока не вернул ни одного блюда.</p>
              </div>
            )}
          </div>

          <div className="status-stack">
            <div className="card">
              <div className="card-title">Норма на день</div>

              <MacroBlock active={20} label="Калории" value={`${formatNumber(dayPlan.nutritionFacts.calories)} ккал`} />

              <hr className="rule" />

              <MacroBlock
                active={20}
                label="Белки"
                value={`${formatNumber(dayPlan.nutritionFacts.protein)} г`}
                variant="balance"
              />
              <MacroBlock
                active={20}
                label="Жиры"
                marginTop
                value={`${formatNumber(dayPlan.nutritionFacts.fat)} г`}
                variant="balance"
              />
              <MacroBlock
                active={20}
                label="Углеводы"
                marginTop
                value={`${formatNumber(dayPlan.nutritionFacts.carb)} г`}
                variant="balance"
              />
            </div>

            <div className="card">
              <div className="card-title">Прогресс по плану</div>

              <MacroBlock
                active={toNotchCount(stats?.percentageComplianceNutritionFacts ?? 0, shouldMaskNutritionAccuracy)}
                label="Точность по БЖУ"
                value={formatPercent(stats?.percentageComplianceNutritionFacts ?? 0, shouldMaskNutritionAccuracy)}
                variant="balance"
              />
              <MacroBlock
                active={toNotchCount(stats?.percentagePlanFulfilled ?? 0)}
                label="Выполнение плана"
                marginTop
                value={formatPercent(stats?.percentagePlanFulfilled ?? 0)}
              />
              <MacroBlock
                active={toNotchCount(stats?.percentageWaterStandardFulfillment ?? 0)}
                label="Вода"
                marginTop
                value={formatPercent(stats?.percentageWaterStandardFulfillment ?? 0)}
                variant="water"
              />

              <div className="reward">
                <b>{formatMealsCount(dayPlan.meals.length)}</b>
                <span>в дневном рационе. Статистика справа обновляется после отметок о еде и воде.</span>
              </div>
            </div>

            <div className="card">
              <div className="card-title">Вода</div>
              <div className="nutrition-water-day">
                <div>
                  <div className="nutrition-water-day__label">Открыт день</div>
                  <div className="nutrition-water-day__value">{activeDayOption.fullLabel}</div>
                </div>
                <span className={isTodaySelected ? 'pill done' : 'pill pending'}>
                  {isTodaySelected ? 'сегодня' : 'просмотр'}
                </span>
              </div>

              <div className="nutrition-water-goal">{formatWaterGoal(dayPlan.waterGoalMl)} в день</div>

              <div className="nutrition-water-stats">
                <div className="nutrition-water-stat">
                  <span className="nutrition-water-stat__label">Цель для дня</span>
                  <strong className="nutrition-water-stat__value">{formatWaterGoal(dayPlan.waterGoalMl)}</strong>
                </div>
                <div className="nutrition-water-stat">
                  <span className="nutrition-water-stat__label">Добавлено в этой сессии</span>
                  <strong className="nutrition-water-stat__value">{formatWaterGoal(activeDaySessionWaterMl)}</strong>
                </div>
              </div>

              <p className="panel-copy nutrition-water-note">
                {isTodaySelected
                  ? 'Добавляйте воду быстрыми кнопками для выбранного дня.'
                  : 'Для этого дня показана отдельная цель по воде.'}
              </p>

              <div className="button-row nutrition-water-actions">
                {waterQuickActions.map((amountMl) => (
                  <button
                    className="btn btn-secondary nutrition-water-button"
                    disabled={!isTodaySelected || isCompletingWaterAmount !== null}
                    key={amountMl}
                    onClick={() => {
                      void handleCompleteWater(amountMl)
                    }}
                    type="button"
                  >
                    {isCompletingWaterAmount === amountMl ? 'СОХРАНЯЕМ...' : `+ ${amountMl} МЛ`}
                  </button>
                ))}
              </div>

            </div>
          </div>
        </section>
      ) : null}

      <footer className="foot">worktrition · питание по дням недели · актуальные данные из backend</footer>
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

function MacroBlock({
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
  variant?: 'balance' | 'water'
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

function toNotchCount(percentage: number, maskValue = false) {
  if (maskValue) {
    return 0
  }

  return Math.max(0, Math.min(20, Math.round((normalizePercentage(percentage) / 100) * 20)))
}

function normalizePercentage(value: number) {
  if (!Number.isFinite(value)) {
    return 0
  }

  return Math.max(0, Math.min(100, value))
}

function formatPercent(value: number, maskValue = false) {
  if (maskValue) {
    return 'Нет данных'
  }

  return `${Math.round(normalizePercentage(value))}%`
}

function formatNumber(value: number) {
  if (Number.isInteger(value)) {
    return String(value)
  }

  return value.toFixed(1)
}

function formatMealsCount(count: number) {
  if (count === 1) {
    return '1 приём пищи'
  }

  if (count >= 2 && count <= 4) {
    return `${count} приёма пищи`
  }

  return `${count} приёмов пищи`
}

function formatWaterGoal(value: number) {
  if (value >= 1000) {
    return `${(value / 1000).toFixed(value % 1000 === 0 ? 0 : 1)} л`
  }

  return `${value} мл`
}
