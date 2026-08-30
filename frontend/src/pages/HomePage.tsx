import { useEffect, useState } from 'react'
import {
  ApiError,
  type GamificationCharacter,
  type NutritionDayPlan,
  type NutritionStats,
  type WorkoutDayPlan,
  type WorkoutStats,
} from '../api'
import { useAuth } from '../auth/useAuth'
import { AppFrame } from '../components/app/AppFrame'
import { NotchBar } from '../components/app/PlaceholderUi'
import { InlineMessage } from '../components/auth/InlineMessage'
import { useCurrentUserData } from '../hooks'
import { toErrorMessage } from '../utils'

const statDefinitions = [
  {
    description: 'Растет за силовые тренировки и плотные сессии.',
    key: 'strength',
    label: 'Сила',
    variant: 'strength' as const,
  },
  {
    description: 'Повышается через кардио, шаги и длинные тренировки.',
    key: 'endurance',
    label: 'Выносливость',
    variant: 'endurance' as const,
  },
  {
    description: 'Набирается за регулярность и отсутствие пропусков.',
    key: 'discipline',
    label: 'Дисциплина',
    variant: 'discipline' as const,
  },
  {
    description: 'Растет через питание, воду и ровный режим.',
    key: 'balance',
    label: 'Баланс',
    variant: 'balance' as const,
  },
] as const

type DashboardState = {
  character: GamificationCharacter | null
  nutritionPlan: NutritionDayPlan | null
  nutritionStats: NutritionStats | null
  workoutPlan: WorkoutDayPlan | null
  workoutStats: WorkoutStats | null
}

const initialDashboardState: DashboardState = {
  character: null,
  nutritionPlan: null,
  nutritionStats: null,
  workoutPlan: null,
  workoutStats: null,
}

export function HomePage() {
  const { getGamificationCharacter, getNutritionDayPlan, getNutritionStats, getWorkoutDayPlan, getWorkoutStats } =
    useAuth()
  const { isLoading: isLoadingUser, loadError: userLoadError, user } = useCurrentUserData()
  const [dashboard, setDashboard] = useState<DashboardState>(initialDashboardState)
  const [isLoading, setIsLoading] = useState(true)
  const [loadError, setLoadError] = useState<string | null>(null)

  useEffect(() => {
    let isCancelled = false

    const loadDashboard = async () => {
      const today = getCurrentDayValue()

      try {
        setIsLoading(true)
        setLoadError(null)

        const [characterResult, workoutStatsResult, nutritionStatsResult, workoutPlanResult, nutritionPlanResult] =
          await Promise.allSettled([
            getGamificationCharacter(),
            getWorkoutStats(),
            getNutritionStats(),
            getWorkoutDayPlan(today),
            getNutritionDayPlan(today),
          ])

        if (isCancelled) {
          return
        }

        const errors: string[] = []

        const nextDashboard: DashboardState = {
          character: getSettledValue(characterResult, errors, 'Не удалось загрузить персонажа.'),
          nutritionPlan: getSettledValue(nutritionPlanResult, errors, 'Не удалось загрузить питание на сегодня.', {
            allowNotFound: true,
          }),
          nutritionStats: getSettledValue(nutritionStatsResult, errors, 'Не удалось загрузить статистику питания.', {
            allowNotFound: true,
          }),
          workoutPlan: getSettledValue(workoutPlanResult, errors, 'Не удалось загрузить тренировку на сегодня.', {
            allowNotFound: true,
          }),
          workoutStats: getSettledValue(workoutStatsResult, errors, 'Не удалось загрузить статистику тренировок.', {
            allowNotFound: true,
          }),
        }

        setDashboard(nextDashboard)
        setLoadError(errors.length > 0 ? errors.join(' ') : null)
      } finally {
        if (!isCancelled) {
          setIsLoading(false)
        }
      }
    }

    void loadDashboard()

    return () => {
      isCancelled = true
    }
  }, [getGamificationCharacter, getNutritionDayPlan, getNutritionStats, getWorkoutDayPlan, getWorkoutStats])

  const { character, nutritionPlan, nutritionStats, workoutPlan, workoutStats } = dashboard
  const sidebarProfile = character
    ? {
        badge: String(character.level),
        meta: `${character.currentXp} / ${character.nextLevelXp} XP`,
        name: user?.name ?? 'Персонаж Worktrition',
      }
    : undefined

  const dominantStat = character ? getDominantStat(character) : null
  const xpProgressPercent = getProgressPercent(character?.currentXp ?? 0, character?.nextLevelXp ?? 0)
  const allMealsCompleted =
    nutritionPlan?.meals.length ? nutritionPlan.meals.every((meal) => meal.isCompleted) : false
  const hasWorkoutAssigned = Boolean(workoutPlan && workoutPlan.exercises.length > 0)
  const isWorkoutCompleted = Boolean(workoutPlan?.isCompleted)

  return (
    <AppFrame
      currentUser={user}
      description=""
      eyebrow="Главная"
      isCurrentUserLoading={isLoadingUser}
      sidebarProfile={sidebarProfile}
      title="Персонаж"
    >
      {userLoadError ? <InlineMessage>{userLoadError}</InlineMessage> : null}
      {loadError ? <InlineMessage>{loadError}</InlineMessage> : null}

      {isLoading && !character ? (
        <section className="card empty-state">
          <div className="card-title">Подключаем персонажа</div>
          <p className="panel-copy">Загружаем уровень, характеристики и квесты на сегодня.</p>
        </section>
      ) : null}

      {!isLoading && !character ? (
        <section className="card empty-state">
          <div className="card-title">Персонаж пока недоступен</div>
          <p className="panel-copy">Данные персонажа пока не загрузились. Попробуйте обновить страницу чуть позже.</p>
        </section>
      ) : null}

      {character ? (
        <>
          <section className="card frame placeholder-section gamification-hero">
            <div className="gamification-hero__main">
              <div className="ring-wrap ring-wrap--xl">
                <svg viewBox="0 0 148 148">
                  <circle cx="74" cy="74" fill="none" r="63" stroke="rgba(255,255,255,.08)" strokeWidth="10" />
                  <circle
                    cx="74"
                    cy="74"
                    fill="none"
                    r="63"
                    stroke="url(#gamificationHeroGradient)"
                    strokeDasharray="395.8"
                    strokeDashoffset={395.8 - (395.8 * xpProgressPercent) / 100}
                    strokeLinecap="round"
                    strokeWidth="10"
                  />
                  <defs>
                    <linearGradient id="gamificationHeroGradient" x1="0" x2="1" y1="0" y2="1">
                      <stop offset="0" stopColor="#FFD23F" />
                      <stop offset="1" stopColor="#FF7A1A" />
                    </linearGradient>
                  </defs>
                </svg>
                <div className="ring-center">
                  <div className="ring-lv">LEVEL</div>
                  <div className="ring-num ring-num--xl">{character.level}</div>
                </div>
              </div>

              <div className="gamification-hero__copy">
                <h2 className="hero-name">{user?.name ?? 'Ваш персонаж'}</h2>
                <div className="hero-title">{getLevelTitle(character.level, dominantStat?.label)}</div>

                <div className="gamification-chip-row">
                  <StatusChip label="До уровня" value={`${Math.max(character.nextLevelXp - character.currentXp, 0)} XP`} />
                  <StatusChip label="HP" value={formatHp(character.hp)} />
                  <StatusChip label="Серия" value={formatStreak(character.currentStreak)} />
                </div>

                <div className="bar-row hero-progress-copy">
                  <span className="bar-label">Опыт персонажа</span>
                  <span className="bar-value">
                    {character.currentXp} / {character.nextLevelXp} XP
                  </span>
                </div>
                <NotchBar active={toNotchCount(xpProgressPercent, 34)} large total={34} />
              </div>
            </div>
          </section>

          <section className="grid g-4 placeholder-section gamification-summary-grid">
            <InfoCard
              label="Уровень"
              note="текущий этап персонажа"
              value={`Lv. ${character.level}`}
            />
            <InfoCard
              label="Серия"
              note="активных дней подряд"
              value={formatStreak(character.currentStreak)}
            />
            <InfoCard
              label="Тонус"
              note="ресурс за регулярность"
              value={formatHp(character.hp)}
            />
            <InfoCard
              label="Архетип"
              note="по доминирующей характеристике"
              value={dominantStat?.label ?? 'Не определен'}
            />
          </section>

          <section className="grid g-2 placeholder-section">
            <div className="card">
              <div className="card-title">Характеристики персонажа</div>
              <p className="panel-copy gamification-section-copy">
                Каждая характеристика растет от реальных действий в приложении. Чем выше показатель, тем ярче профиль персонажа.
              </p>

              {statDefinitions.map((item) => (
                <AttributeRow
                  description={item.description}
                  key={item.key}
                  label={item.label}
                  value={character[item.key]}
                  variant={item.variant}
                />
              ))}
            </div>

            <div className="card frame">
              <div className="card-title">Ежедневные квесты</div>
              <p className="panel-copy gamification-section-copy">
                Здесь собран сегодняшний цикл: движение, питание и вода. Данные подтягиваются из активных планов и статистики.
              </p>

              <QuestRow
                label={hasWorkoutAssigned ? formatWorkoutQuestLabel(workoutPlan) : 'Сегодня без тренировки'}
                note={
                  hasWorkoutAssigned
                    ? isWorkoutCompleted
                      ? 'Награда уже учтена в персонаже.'
                      : 'Завершите тренировку на странице тренировок.'
                    : 'День отдыха тоже помогает держать режим.'
                }
                status={hasWorkoutAssigned ? (isWorkoutCompleted ? 'done' : 'pending') : 'neutral'}
                value={hasWorkoutAssigned ? (isWorkoutCompleted ? 'выполнено' : 'в плане') : 'отдых'}
              />
              <QuestRow
                label={nutritionPlan ? formatMealsQuestLabel(nutritionPlan) : 'Рацион на сегодня не найден'}
                note={
                  nutritionPlan
                    ? allMealsCompleted
                      ? 'Баланс уже обновлен после отметок о еде.'
                      : 'Отмечайте блюда на странице питания.'
                    : 'Сначала нужен сохраненный план питания.'
                }
                status={nutritionPlan ? (allMealsCompleted ? 'done' : 'pending') : 'neutral'}
                value={nutritionPlan ? (allMealsCompleted ? 'готово' : 'в процессе') : 'нет плана'}
              />
              <QuestRow
                label={nutritionPlan ? `Вода ${formatWaterGoal(nutritionPlan.waterGoalMl)}` : 'Водный режим не задан'}
                note={
                  nutritionPlan
                    ? `Среднее соблюдение сейчас ${formatPercent(nutritionStats?.percentageWaterStandardFulfillment ?? 0)}.`
                    : 'Появится после генерации дневного рациона.'
                }
                status={
                  nutritionPlan
                    ? (nutritionStats?.percentageWaterStandardFulfillment ?? 0) >= 80
                      ? 'done'
                      : 'pending'
                    : 'neutral'
                }
                value={
                  nutritionPlan
                    ? (nutritionStats?.percentageWaterStandardFulfillment ?? 0) >= 80
                      ? 'ритм держится'
                      : 'нужен фокус'
                    : 'ожидание'
                }
              />

              <div className="reward gamification-reward">
                <b>Ежедневный цикл</b>
                <span>Отмечайте тренировки, питание и воду, чтобы персонаж стабильно рос без провалов по ритму.</span>
              </div>
            </div>
          </section>

          <section className="placeholder-section">
            <div className="card">
              <div className="card-title">Боевые метрики</div>
              <p className="panel-copy gamification-section-copy">
                Это агрегированная картина по тренировкам и питанию, которая уже влияет на развитие персонажа.
              </p>

              <MetricBlock
                active={toNotchCount(workoutStats?.percentagePlanFulfilled ?? 0)}
                label="Тренировки"
                value={formatPercent(workoutStats?.percentagePlanFulfilled ?? 0)}
                variant="strength"
              />
              <MetricBlock
                active={toNotchCount(nutritionStats?.percentagePlanFulfilled ?? 0)}
                label="Питание"
                marginTop
                value={formatPercent(nutritionStats?.percentagePlanFulfilled ?? 0)}
                variant="balance"
              />
              <MetricBlock
                active={toNotchCount(nutritionStats?.percentageComplianceNutritionFacts ?? 0)}
                label="Точность по БЖУ"
                marginTop
                value={formatPercent(nutritionStats?.percentageComplianceNutritionFacts ?? 0)}
                variant="discipline"
              />
              <MetricBlock
                active={toNotchCount(nutritionStats?.percentageWaterStandardFulfillment ?? 0)}
                label="Вода"
                marginTop
                value={formatPercent(nutritionStats?.percentageWaterStandardFulfillment ?? 0)}
                variant="water"
              />
            </div>
          </section>
        </>
      ) : null}
    </AppFrame>
  )
}

function InfoCard({
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

function AttributeRow({
  description,
  label,
  value,
  variant,
}: {
  description: string
  label: string
  value: number
  variant: 'balance' | 'discipline' | 'endurance' | 'strength'
}) {
  return (
    <div className="gamification-attribute">
      <div className="stat-row">
        <span className="stat-label">
          <span className={`stat-dot stat-dot--${variant}`} />
          {label}
        </span>
        <NotchBar active={toAttributeNotchCount(value)} total={20} variant={variant} />
        <span className="stat-value">{value}</span>
      </div>
      <p className="gamification-attribute__description">{description}</p>
    </div>
  )
}

function QuestRow({
  label,
  note,
  status,
  value,
}: {
  label: string
  note: string
  status: 'done' | 'neutral' | 'pending'
  value: string
}) {
  return (
    <div className="gamification-quest">
      <div className="gamification-quest__top">
        <div className="gamification-quest__label">{label}</div>
        <span className={status === 'done' ? 'pill done' : status === 'pending' ? 'pill pending' : 'pill'}>
          {value}
        </span>
      </div>
      <p className="gamification-quest__note">{note}</p>
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

function StatusChip({ label, value }: { label: string; value: string }) {
  return (
    <div className="gamification-chip">
      <span className="gamification-chip__label">{label}</span>
      <strong className="gamification-chip__value">{value}</strong>
    </div>
  )
}

function getSettledValue<T>(
  result: PromiseSettledResult<T>,
  errors: string[],
  fallbackMessage: string,
  options?: { allowNotFound?: boolean },
) {
  if (result.status === 'fulfilled') {
    return result.value
  }

  if (options?.allowNotFound && result.reason instanceof ApiError && result.reason.status === 404) {
    return null
  }

  errors.push(toErrorMessage(result.reason, fallbackMessage))
  return null
}

function getCurrentDayValue() {
  const day = new Date().getDay()

  if (day === 0) {
    return 7
  }

  return day
}

function getLevelTitle(level: number, dominantStatLabel?: string) {
  if (level >= 12) {
    return `Мастер режима · ведущая черта: ${dominantStatLabel ?? 'универсальность'}`
  }

  if (level >= 8) {
    return `Стратег формы · ведущая черта: ${dominantStatLabel ?? 'баланс'}`
  }

  if (level >= 4) {
    return `Уверенный прогресс · ведущая черта: ${dominantStatLabel ?? 'дисциплина'}`
  }

  return `Старт пути · ведущая черта: ${dominantStatLabel ?? 'рост'}`
}

function getDominantStat(character: GamificationCharacter) {
  return statDefinitions.reduce((best, current) =>
    character[current.key] > character[best.key] ? current : best,
  )
}

function getProgressPercent(current: number, total: number) {
  if (total <= 0) {
    return 0
  }

  return normalizePercentage((current / total) * 100)
}

function toAttributeNotchCount(value: number) {
  return Math.max(0, Math.min(20, Math.round((normalizeNonNegativeNumber(value) / 40) * 20)))
}

function toNotchCount(percentage: number, total = 20) {
  return Math.max(0, Math.min(total, Math.round((normalizePercentage(percentage) / 100) * total)))
}

function normalizePercentage(value: number) {
  if (!Number.isFinite(value)) {
    return 0
  }

  return Math.max(0, Math.min(100, value))
}

function normalizeNonNegativeNumber(value: number) {
  if (!Number.isFinite(value)) {
    return 0
  }

  return Math.max(0, value)
}

function formatPercent(value: number) {
  return `${Math.round(normalizePercentage(value))}%`
}

function formatHp(value: number) {
  return `${normalizeNonNegativeNumber(value).toFixed(1)} / 6.0`
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

function formatWorkoutQuestLabel(workoutPlan: WorkoutDayPlan | null) {
  if (!workoutPlan) {
    return 'Тренировка пока не назначена'
  }

  const type = workoutPlan.type.trim() || 'Тренировка'
  return `${type} · ${workoutPlan.exercises.length} упражнений`
}

function formatMealsQuestLabel(nutritionPlan: NutritionDayPlan) {
  const totalMeals = nutritionPlan.meals.length
  const completedMeals = nutritionPlan.meals.filter((meal) => meal.isCompleted).length

  return `Рацион на день · ${completedMeals}/${totalMeals} приемов отмечено`
}

function formatWaterGoal(value: number) {
  if (value >= 1000) {
    return `${(value / 1000).toFixed(value % 1000 === 0 ? 0 : 1)} л`
  }

  return `${value} мл`
}
