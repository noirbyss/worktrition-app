const STORAGE_KEY = 'worktrition.nutrition-ui'
const STORAGE_EVENT = 'worktrition:nutrition-ui'

interface NutritionUiState {
  awaitingFreshNutritionAccuracy: boolean
}

const defaultState: NutritionUiState = {
  awaitingFreshNutritionAccuracy: false,
}

export function loadNutritionUiState() {
  try {
    const rawValue = window.localStorage.getItem(STORAGE_KEY)
    if (!rawValue) {
      return defaultState
    }

    const parsedValue = JSON.parse(rawValue)
    if (!isNutritionUiState(parsedValue)) {
      window.localStorage.removeItem(STORAGE_KEY)
      return defaultState
    }

    return parsedValue
  } catch {
    return defaultState
  }
}

export function markNutritionAccuracyPendingReset() {
  persistNutritionUiState({
    awaitingFreshNutritionAccuracy: true,
  })
}

export function clearNutritionAccuracyPendingReset() {
  persistNutritionUiState({
    awaitingFreshNutritionAccuracy: false,
  })
}

export function subscribeNutritionUiState(onChange: () => void) {
  window.addEventListener(STORAGE_EVENT, onChange)

  return () => {
    window.removeEventListener(STORAGE_EVENT, onChange)
  }
}

function persistNutritionUiState(value: NutritionUiState) {
  try {
    window.localStorage.setItem(STORAGE_KEY, JSON.stringify(value))
  } catch {
    return
  } finally {
    window.dispatchEvent(new Event(STORAGE_EVENT))
  }
}

function isNutritionUiState(value: unknown): value is NutritionUiState {
  if (!value || typeof value !== 'object') {
    return false
  }

  const state = value as Partial<NutritionUiState>

  return typeof state.awaitingFreshNutritionAccuracy === 'boolean'
}
