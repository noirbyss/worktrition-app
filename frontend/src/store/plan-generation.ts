import type { PlanType } from '../api'

const STORAGE_KEY = 'worktrition.plan-generation'
const STORAGE_EVENT = 'worktrition:plan-generation'

export interface PendingPlanGeneration {
  generationId?: string
  planType: PlanType
  status: 'pending' | 'starting'
  updatedAt: number
}

export function loadPendingPlanGeneration() {
  try {
    const rawValue = window.sessionStorage.getItem(STORAGE_KEY)
    if (!rawValue) {
      return null
    }

    const parsedValue = JSON.parse(rawValue)
    if (!isPendingPlanGeneration(parsedValue)) {
      window.sessionStorage.removeItem(STORAGE_KEY)
      return null
    }

    return parsedValue
  } catch {
    return null
  }
}

export function persistPendingPlanGeneration(
  value: Omit<PendingPlanGeneration, 'updatedAt'>,
) {
  const nextValue: PendingPlanGeneration = {
    ...value,
    updatedAt: Date.now(),
  }

  try {
    window.sessionStorage.setItem(STORAGE_KEY, JSON.stringify(nextValue))
  } catch {
    return nextValue
  } finally {
    window.dispatchEvent(new Event(STORAGE_EVENT))
  }

  return nextValue
}

export function clearPendingPlanGeneration() {
  try {
    window.sessionStorage.removeItem(STORAGE_KEY)
  } catch {
    return
  } finally {
    window.dispatchEvent(new Event(STORAGE_EVENT))
  }
}

export function subscribePendingPlanGeneration(onChange: () => void) {
  window.addEventListener(STORAGE_EVENT, onChange)

  return () => {
    window.removeEventListener(STORAGE_EVENT, onChange)
  }
}

function isPendingPlanGeneration(value: unknown): value is PendingPlanGeneration {
  if (!value || typeof value !== 'object') {
    return false
  }

  const pendingGeneration = value as Partial<PendingPlanGeneration>

  return (
    (pendingGeneration.planType === 'all' ||
      pendingGeneration.planType === 'nutrition' ||
      pendingGeneration.planType === 'workout') &&
    (pendingGeneration.status === 'pending' || pendingGeneration.status === 'starting') &&
    typeof pendingGeneration.updatedAt === 'number' &&
    (pendingGeneration.generationId === undefined || typeof pendingGeneration.generationId === 'string')
  )
}
