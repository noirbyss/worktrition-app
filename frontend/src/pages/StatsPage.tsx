import { useEffect, useState } from 'react'
import {
  ApiError,
  type GamificationCharacter,
  type NutritionStats,
  type Profile,
  type WorkoutStats,
} from '../api'
import { useAuth } from '../auth/useAuth'
import { AppFrame } from '../components/app/AppFrame'
import { NotchBar } from '../components/app/PlaceholderUi'
import { InlineMessage } from '../components/auth/InlineMessage'
import { useCurrentUserData } from '../hooks'
import { toErrorMessage } from '../utils'

type StatsState = {
  character: GamificationCharacter | null
  nutritionStats: NutritionStats | null
  profile: Profile | null
  workoutStats: WorkoutStats | null
}

const initialStatsState: StatsState = {
  character: null,
  nutritionStats: null,
  profile: null,
  workoutStats: null,
}

const goalLabels: Record<number, string> = {
  1: 'Снижение веса',
  2: 'Поддержание веса',
  3: 'Набор мышечной массы',
}

const trainingLevelLabels: Record<number, string> = {
  1: 'Новичок',
  2: 'Средний',
  3: 'Продвинутый',
}

export function StatsPage() {
  const { getGamificationCharacter, getNutritionStats, getProfile, getWorkoutStats } = useAuth()
  const { isLoading: isLoadingUser, loadError: userLoadError, user } = useCurrentUserData()
  const [isLoading, setIsLoading] = useState(true)
  const [loadError, setLoadError] = useState<string | null>(null)
  const [stats, setStats] = useState<StatsState>(initialStatsState)

  useEffect(() => {
    let isCancelled = false

    const loadStats = async () => {
      try {
        setIsLoading(true)
        setLoadError(null)

        const [profileResult, workoutStatsResult, nutritionStatsResult, characterResult] =
          await Promise.allSettled([
            getProfile(),
            getWorkoutStats(),
            getNutritionStats(),
            getGamificationCharacter(),
          ])

        if (isCancelled) {
          return
        }

        const errors: string[] = []

        setStats({
          character: getSettledValue(characterResult, errors, 'Не удалось загрузить прогресс персонажа.', {
            allowNotFound: true,
          }),
          nutritionStats: getSettledValue(
            nutritionStatsResult,
            errors,
            'Не удалось загрузить статистику питания.',
            { allowNotFound: true },
          ),
          profile: getSettledValue(profileResult, errors, 'Не удалось загрузить данные профиля.', {
            allowNotFound: true,
          }),
          workoutStats: getSettledValue(
            workoutStatsResult,
            errors,
            'Не удалось загрузить статистику тренировок.',
            { allowNotFound: true },
          ),
        })
        setLoadError(errors.length > 0 ? errors.join(' ') : null)
      } finally {
        if (!isCancelled) {
          setIsLoading(false)
        }
      }
    }

    void loadStats()

    return () => {
      isCancelled = true
    }
  }, [getGamificationCharacter, getNutritionStats, getProfile, getWorkoutStats])

  const { character, nutritionStats, profile, workoutStats } = stats
  const sidebarProfile = character
    ? {
        badge: String(character.level),
        meta: `${character.currentXp} / ${character.nextLevelXp} XP`,
        name: user?.name ?? 'Персонаж Worktrition',
      }
    : undefined
  const bmi = profile ? calculateBmi(profile.heightCm, profile.weightKg) : null
  const weightGoalChart = buildWeightGoalChart(profile)
  const hasAnyStats = Boolean(profile || workoutStats || nutritionStats)

  return (
    <AppFrame
      currentUser={user}
      description="Реальный прогресс: вес, тренировки и питание на основе данных из backend."
      eyebrow="Экран 04"
      isCurrentUserLoading={isLoadingUser}
      sidebarProfile={sidebarProfile}
      title="Статистика"
    >
      {userLoadError ? <InlineMessage>{userLoadError}</InlineMessage> : null}
      {loadError ? <InlineMessage>{loadError}</InlineMessage> : null}

      {isLoading && !hasAnyStats ? (
        <section className="card empty-state">
          <div className="card-title">Загрузка статистики</div>
          <p className="panel-copy">Подтягиваем профиль, тренировки и питание из backend.</p>
        </section>
      ) : null}

      {!isLoading && !hasAnyStats ? (
        <section className="card empty-state">
          <div className="card-title">Статистика пока недоступна</div>
          <p className="panel-copy">
            Для этой страницы пока не нашлось данных. Обычно они появляются после заполнения
            профиля и генерации планов.
          </p>
        </section>
      ) : null}

      {hasAnyStats ? (
        <>
          <section className="grid g-2 placeholder-grid">
            <div className="card">
              <div className="card-title">Вес и цель</div>
              <WeightGoalChart chart={weightGoalChart} />
              <p className="panel-copy">
                {weightGoalChart.note}
              </p>
            </div>

            <div className="card">
              <div className="card-title">Физические показатели</div>
              <SimpleStatRow
                label="Текущий вес"
                value={profile ? `${formatWeight(profile.weightKg)} кг` : 'Нет данных'}
              />
              <SimpleStatRow
                label="Целевой вес"
                value={profile ? formatTargetWeight(profile.targetWeightKg) : 'Нет данных'}
              />
              <SimpleStatRow
                label="До цели"
                value={profile ? formatWeightDelta(profile.weightKg, profile.targetWeightKg) : 'Нет данных'}
              />
              <SimpleStatRow label="ИМТ" value={bmi === null ? 'Нет данных' : bmi.toFixed(1)} />
              <SimpleStatRow
                label="Цель"
                value={profile ? goalLabels[profile.goal] ?? 'Не указана' : 'Нет данных'}
              />
            </div>
          </section>

          <section className="grid g-2">
            <div className="card">
              <div className="card-title">Тренировки</div>
              <ProgressStatRow
                active={toNotchCount(workoutStats?.percentagePlanFulfilled ?? 0)}
                label="Выполнение плана"
                value={formatPercent(workoutStats?.percentagePlanFulfilled ?? 0, workoutStats !== null)}
              />
              <ProgressStatRow
                active={Math.min((workoutStats?.currentStreakDays ?? 0) * 2, 18)}
                label="Серия дней"
                value={formatStreak(workoutStats?.currentStreakDays ?? 0, workoutStats !== null)}
                variant="discipline"
              />
              <ProgressStatRow
                active={toTrainingTimeNotches(workoutStats?.totalTrainingTimeSeconds ?? 0)}
                label="Общее время"
                value={formatDuration(workoutStats?.totalTrainingTimeSeconds ?? 0, workoutStats !== null)}
                variant="endurance"
              />
              <ProgressStatRow
                active={Math.min(profile?.trainingDaysPerWeek ?? 0, 7) * 2}
                label="Дней в неделю"
                value={formatTrainingDays(profile?.trainingDaysPerWeek, profile !== null)}
                variant="strength"
              />
              <ProgressStatRow
                active={Math.min((profile?.trainingLevel ?? 0) * 6, 18)}
                label="Уровень"
                value={formatTrainingLevel(profile?.trainingLevel, profile !== null)}
                variant="strength"
              />
            </div>

            <div className="card">
              <div className="card-title">Питание</div>
              <ProgressStatRow
                active={toNotchCount(nutritionStats?.percentageComplianceNutritionFacts ?? 0)}
                label="Соблюдение КБЖУ"
                value={formatPercent(
                  nutritionStats?.percentageComplianceNutritionFacts ?? 0,
                  nutritionStats !== null,
                )}
                variant="balance"
              />
              <ProgressStatRow
                active={toNotchCount(nutritionStats?.percentagePlanFulfilled ?? 0)}
                label="Выполнение плана"
                value={formatPercent(nutritionStats?.percentagePlanFulfilled ?? 0, nutritionStats !== null)}
                variant="balance"
              />
              <ProgressStatRow
                active={toNotchCount(nutritionStats?.percentageWaterStandardFulfillment ?? 0)}
                label="Норма воды"
                value={formatPercent(
                  nutritionStats?.percentageWaterStandardFulfillment ?? 0,
                  nutritionStats !== null,
                )}
                variant="water"
              />
              <ProgressStatRow
                active={Math.min(character?.balance ?? 0, 18)}
                label="Баланс персонажа"
                value={formatCharacterBalance(character?.balance, character !== null)}
                variant="balance"
              />
              <ProgressStatRow
                active={Math.min(character?.currentStreak ?? 0, 18)}
                label="Общая серия"
                value={formatStreak(character?.currentStreak ?? 0, character !== null)}
                variant="discipline"
              />
            </div>
          </section>
        </>
      ) : null}

      <footer className="foot">worktrition · статистика · backend</footer>
    </AppFrame>
  )
}

function WeightGoalChart({
  chart,
}: {
  chart: {
    labelPoints: Array<{ label: string; value?: number; x: number; y: number }>
    note: string
    polylinePoints: string
    tickLabels: Array<{ label: string; x: number }>
  }
}) {
  return (
    <svg className="weight-chart" height="160" viewBox="0 0 460 160" width="100%">
      <defs>
        <linearGradient id="weightChartGradient" x1="0" x2="1" y1="0" y2="0">
          <stop offset="0" stopColor="#FFD23F" />
          <stop offset="1" stopColor="#FF7A1A" />
        </linearGradient>
      </defs>

      <polyline
        fill="none"
        points={chart.polylinePoints}
        stroke="url(#weightChartGradient)"
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth="3"
      />

      <g fill="#FF7A1A">
        {chart.labelPoints.map((point) => (
          <circle cx={point.x} cy={point.y} key={`${point.label}-${point.value}`} r="4" />
        ))}
      </g>

      <g fill="#9195A8" fontFamily="IBM Plex Mono, monospace" fontSize="11">
        {chart.labelPoints.map((point) => (
          <text key={`${point.label}-value`} x={point.x - 16} y={point.y - 10}>
            {typeof point.value === 'number' ? `${formatWeight(point.value)} кг` : '—'}
          </text>
        ))}
      </g>

      <line stroke="rgba(255,255,255,.1)" strokeWidth="1" x1="0" x2="460" y1="145" y2="145" />

      <g fill="#565A70" fontFamily="IBM Plex Mono, monospace" fontSize="10">
        {chart.tickLabels.map((item) => (
          <text key={item.label} x={item.x} y="158">
            {item.label}
          </text>
        ))}
      </g>
    </svg>
  )
}

function ProgressStatRow({
  active,
  label,
  value,
  variant,
}: {
  active: number
  label: string
  value: string
  variant?: 'balance' | 'discipline' | 'endurance' | 'strength' | 'water'
}) {
  return (
    <div className="stat-row">
      <span className="stat-label">{label}</span>
      <NotchBar active={active} total={18} variant={variant} />
      <span className="stat-value">{value}</span>
    </div>
  )
}

function SimpleStatRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="stat-row stat-row--simple">
      <span className="stat-label">{label}</span>
      <span />
      <span className="stat-value">{value}</span>
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

function buildWeightGoalChart(profile: Profile | null) {
  if (!profile) {
    return {
      labelPoints: [
        { label: 'сейчас', x: 40, y: 84 },
        { label: 'цель', x: 420, y: 84 },
      ],
      note: 'Профиль ещё не загружен, поэтому график веса пока недоступен.',
      polylinePoints: '40,84 420,84',
      tickLabels: [
        { label: 'сейчас', x: 24 },
        { label: 'цель', x: 404 },
      ],
    }
  }

  const current = profile.weightKg
  if (typeof profile.targetWeightKg !== 'number') {
    return {
      labelPoints: [{ label: 'сейчас', value: current, x: 40, y: 84 }],
      note: 'В backend пока есть только текущий вес. После сохранения целевого веса график покажет направление к цели.',
      polylinePoints: '40,84 420,84',
      tickLabels: [
        { label: 'сейчас', x: 24 },
        { label: 'цель не задана', x: 338 },
      ],
    }
  }

  const target = profile.targetWeightKg
  const midpoint = (current + target) / 2
  const values = [current, midpoint, target]
  const min = Math.min(...values)
  const max = Math.max(...values)
  const range = Math.max(max - min, 1)
  const xPositions = [40, 230, 420]
  const yPositions = values.map((value) => {
    const normalized = (value - min) / range
    return 118 - normalized * 68
  })

  return {
    labelPoints: [
      { label: 'сейчас', value: current, x: xPositions[0], y: yPositions[0] },
      { label: 'цель', value: target, x: xPositions[2], y: yPositions[2] },
    ],
    note: 'Показываем текущий и целевой вес из профиля. История замеров пока не приходит отдельным endpoint.',
    polylinePoints: xPositions.map((x, index) => `${x},${yPositions[index]}`).join(' '),
    tickLabels: [
      { label: 'сейчас', x: 24 },
      { label: 'середина', x: 206 },
      { label: 'цель', x: 394 },
    ],
  }
}

function calculateBmi(heightCm: number, weightKg: number) {
  if (!Number.isFinite(heightCm) || !Number.isFinite(weightKg) || heightCm <= 0) {
    return null
  }

  const heightM = heightCm / 100
  return Math.round((weightKg / (heightM * heightM)) * 10) / 10
}

function toNotchCount(percentage: number) {
  return Math.max(0, Math.min(18, Math.round((normalizePercentage(percentage) / 100) * 18)))
}

function toTrainingTimeNotches(durationSeconds: number) {
  const durationMinutes = normalizeNonNegativeNumber(durationSeconds) / 60
  return Math.max(0, Math.min(18, Math.round((durationMinutes / 600) * 18)))
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

function formatPercent(value: number, hasData: boolean) {
  if (!hasData) {
    return 'Нет данных'
  }

  return `${Math.round(normalizePercentage(value))}%`
}

function formatDuration(value: number, hasData: boolean) {
  if (!hasData) {
    return 'Нет данных'
  }

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

function formatStreak(value: number, hasData: boolean) {
  if (!hasData) {
    return 'Нет данных'
  }

  const safeValue = Math.max(0, Math.floor(normalizeNonNegativeNumber(value)))

  if (safeValue === 1) {
    return '1 день'
  }

  if (safeValue >= 2 && safeValue <= 4) {
    return `${safeValue} дня`
  }

  return `${safeValue} дней`
}

function formatTrainingDays(days: number | undefined, hasData: boolean) {
  if (!hasData || typeof days !== 'number') {
    return 'Нет данных'
  }

  if (days === 1) {
    return '1 день'
  }

  if (days >= 2 && days <= 4) {
    return `${days} дня`
  }

  return `${days} дней`
}

function formatTrainingLevel(level: number | undefined, hasData: boolean) {
  if (!hasData || typeof level !== 'number') {
    return 'Нет данных'
  }

  return trainingLevelLabels[level] ?? 'Не указан'
}

function formatCharacterBalance(balance: number | undefined, hasData: boolean) {
  if (!hasData || typeof balance !== 'number') {
    return 'Нет данных'
  }

  return `${Math.round(normalizeNonNegativeNumber(balance))}`
}

function formatTargetWeight(targetWeightKg?: number) {
  if (typeof targetWeightKg !== 'number') {
    return 'Не указан'
  }

  return `${formatWeight(targetWeightKg)} кг`
}

function formatWeightDelta(currentWeightKg: number, targetWeightKg?: number) {
  if (typeof targetWeightKg !== 'number') {
    return 'Не указана'
  }

  const delta = Math.abs(currentWeightKg - targetWeightKg)
  if (delta < 0.05) {
    return 'Цель достигнута'
  }

  return `${formatWeight(delta)} кг`
}

function formatWeight(value: number) {
  return Number.isInteger(value) ? String(value) : value.toFixed(1)
}
