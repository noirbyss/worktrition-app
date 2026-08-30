const API_BASE_URL = (import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080').replace(/\/+$/, '')

type HttpMethod = 'GET' | 'POST' | 'PUT'

interface RequestOptions {
  accessToken?: string
  body?: unknown
  method?: HttpMethod
}

interface AuthResponseDto {
  access_token: string
  access_token_expires_at: number
  profile_completed: boolean
  user_id: string
}

interface CurrentUserDto {
  birth_date: string
  email: string
  name: string
  profile_completed: boolean
  user_id: string
}

interface ProfileDto {
  activity_level: number
  age: number
  allergies: string[]
  equipment: string
  excluded_foods: string[]
  food_preferences: string[]
  gender: number
  goal: number
  height_cm: number
  target_weight_kg?: number
  training_days_per_week: number
  training_level: number
  training_location: number
  user_id: string
  weight_kg: number
}

interface SaveProfileResponseDto {
  bmi: number
  profile_completed: boolean
}

interface GenerationResponseDto {
  error_message?: string
  generation_id: string
  status: string
}

interface NutritionFactsDto {
  calories: number
  protein: number
  fat: number
  carb: number
}

interface NutritionMealDto {
  id: number
  is_completed?: boolean
  name: string
  recipe: string
  nutrition_facts?: NutritionFactsDto | null
}

interface NutritionDayPlanDto {
  meal_items?: NutritionMealDto[]
  nutrition_facts?: NutritionFactsDto | null
  water_goal_ml: number
}

interface NutritionStatsDto {
  percentage_compliance_nutrition_facts: number
  percentage_plan_fulfilled: number
  percentage_water_standard_fulfillment: number
}

interface WorkoutDayPlanDto {
  day_of_week: number
  exercises: string[]
  is_completed: boolean
  type: string
}

interface WorkoutStatsDto {
  current_streak_days: number
  percentage_plan_fulfilled: number
  total_training_time_seconds: number
}

export interface LoginPayload {
  email: string
  password: string
}

export interface RegisterPayload extends LoginPayload {
  birth_date: string
  name: string
}

export interface AuthSession {
  accessToken: string
  accessTokenExpiresAt: number
  profileCompleted: boolean
  userId: string
}

export interface SaveProfilePayload {
  activity_level: number
  allergies: string[]
  equipment: string
  excluded_foods: string[]
  food_preferences: string[]
  gender: number
  goal: number
  height_cm: number
  target_weight_kg?: number
  training_days_per_week: number
  training_level: number
  training_location: number
  weight_kg: number
}

export interface CurrentUser {
  birthDate: string
  email: string
  name: string
  profileCompleted: boolean
  userId: string
}

export interface Profile {
  activityLevel: number
  age: number
  allergies: string[]
  equipment: string
  excludedFoods: string[]
  foodPreferences: string[]
  gender: number
  goal: number
  heightCm: number
  targetWeightKg?: number
  trainingDaysPerWeek: number
  trainingLevel: number
  trainingLocation: number
  userId: string
  weightKg: number
}

export interface SaveProfileResult {
  bmi: number
  profileCompleted: boolean
}

export type PlanType = 'all' | 'nutrition' | 'workout'

export type GenerationStatus = 'done' | 'failed' | 'pending' | 'unspecified'

export interface GenerationResult {
  errorMessage?: string
  generationId: string
  status: GenerationStatus
}

export interface NutritionFacts {
  calories: number
  protein: number
  fat: number
  carb: number
}

export interface NutritionMeal {
  id: number
  isCompleted: boolean
  name: string
  nutritionFacts: NutritionFacts
  recipe: string
}

export interface NutritionDayPlan {
  meals: NutritionMeal[]
  nutritionFacts: NutritionFacts
  waterGoalMl: number
}

export interface NutritionStats {
  percentageComplianceNutritionFacts: number
  percentagePlanFulfilled: number
  percentageWaterStandardFulfillment: number
}

export interface WorkoutDayPlan {
  dayOfWeek: number
  exercises: string[]
  isCompleted: boolean
  type: string
}

export interface WorkoutStats {
  currentStreakDays: number
  percentagePlanFulfilled: number
  totalTrainingTimeSeconds: number
}

export class ApiError extends Error {
  public readonly payload: unknown
  public readonly status: number

  public constructor(message: string, status: number, payload?: unknown) {
    super(message)
    this.name = 'ApiError'
    this.payload = payload
    this.status = status
  }
}

export async function register(payload: RegisterPayload) {
  const response = await request<AuthResponseDto>('/auth/register', {
    body: payload,
    method: 'POST',
  })

  return mapAuthResponse(response)
}

export async function login(payload: LoginPayload) {
  const response = await request<AuthResponseDto>('/auth/login', {
    body: payload,
    method: 'POST',
  })

  return mapAuthResponse(response)
}

export async function refresh() {
  const response = await request<AuthResponseDto>('/auth/refresh', {
    method: 'POST',
  })

  return mapAuthResponse(response)
}

export async function logout() {
  return request<{ success: boolean }>('/auth/logout', {
    method: 'POST',
  })
}

export async function getCurrentUser(accessToken: string) {
  const response = await request<CurrentUserDto>('/users/me', {
    accessToken,
    method: 'GET',
  })

  return {
    birthDate: response.birth_date,
    email: response.email,
    name: response.name,
    profileCompleted: response.profile_completed,
    userId: response.user_id,
  } satisfies CurrentUser
}

export async function getProfile(accessToken: string) {
  const response = await request<ProfileDto>('/profile', {
    accessToken,
    method: 'GET',
  })

  return {
    activityLevel: response.activity_level,
    age: response.age,
    allergies: response.allergies ?? [],
    equipment: response.equipment,
    excludedFoods: response.excluded_foods ?? [],
    foodPreferences: response.food_preferences ?? [],
    gender: response.gender,
    goal: response.goal,
    heightCm: response.height_cm,
    targetWeightKg: response.target_weight_kg,
    trainingDaysPerWeek: response.training_days_per_week,
    trainingLevel: response.training_level,
    trainingLocation: response.training_location,
    userId: response.user_id,
    weightKg: response.weight_kg,
  } satisfies Profile
}

export async function saveProfile(accessToken: string, payload: SaveProfilePayload) {
  const response = await request<SaveProfileResponseDto>('/profile', {
    accessToken,
    body: payload,
    method: 'PUT',
  })

  return {
    bmi: response.bmi,
    profileCompleted: response.profile_completed,
  } satisfies SaveProfileResult
}

export async function startGeneration(accessToken: string, planType: PlanType) {
  const response = await request<GenerationResponseDto>('/ai/generations', {
    accessToken,
    body: {
      plan_type: planType,
    },
    method: 'POST',
  })

  return mapGenerationResponse(response)
}

export async function getGenerationStatus(accessToken: string, generationId: string) {
  const response = await request<GenerationResponseDto>(
    `/ai/generations/${encodeURIComponent(generationId)}`,
    {
      accessToken,
      method: 'GET',
    },
  )

  return mapGenerationResponse(response)
}

export async function getNutritionDayPlan(accessToken: string, dayOfWeek: number | string) {
  const response = await request<NutritionDayPlanDto>(
    `/nutrition/plan?day_of_week=${encodeURIComponent(String(dayOfWeek))}`,
    {
      accessToken,
      method: 'GET',
    },
  )

  return {
    meals: (response.meal_items ?? []).map((meal) => ({
      id: meal.id,
      isCompleted: meal.is_completed ?? false,
      name: meal.name,
      nutritionFacts: mapNutritionFacts(meal.nutrition_facts),
      recipe: meal.recipe,
    })),
    nutritionFacts: mapNutritionFacts(response.nutrition_facts),
    waterGoalMl: response.water_goal_ml,
  } satisfies NutritionDayPlan
}

export async function getNutritionStats(accessToken: string) {
  const response = await request<NutritionStatsDto>('/nutrition/stats', {
    accessToken,
    method: 'GET',
  })

  return {
    percentageComplianceNutritionFacts: response.percentage_compliance_nutrition_facts,
    percentagePlanFulfilled: response.percentage_plan_fulfilled,
    percentageWaterStandardFulfillment: response.percentage_water_standard_fulfillment,
  } satisfies NutritionStats
}

export async function completeNutritionMeal(accessToken: string, mealItemId: number) {
  return request<{ success: boolean }>('/nutrition/meals/complete', {
    accessToken,
    body: {
      meal_item_id: mealItemId,
    },
    method: 'POST',
  })
}

export async function completeNutritionWater(accessToken: string, amountMl: number) {
  return request<{ success: boolean }>('/nutrition/water/complete', {
    accessToken,
    body: {
      amount_ml: amountMl,
    },
    method: 'POST',
  })
}

export async function getWorkoutDayPlan(accessToken: string, dayOfWeek: number | string) {
  const response = await request<WorkoutDayPlanDto>(
    `/workout/plan?day_of_week=${encodeURIComponent(String(dayOfWeek))}`,
    {
      accessToken,
      method: 'GET',
    },
  )

  return {
    dayOfWeek: response.day_of_week,
    exercises: response.exercises ?? [],
    isCompleted: response.is_completed ?? false,
    type: response.type,
  } satisfies WorkoutDayPlan
}

export async function getWorkoutStats(accessToken: string) {
  const response = await request<WorkoutStatsDto>('/workout/stats', {
    accessToken,
    method: 'GET',
  })

  return {
    currentStreakDays: response.current_streak_days,
    percentagePlanFulfilled: response.percentage_plan_fulfilled,
    totalTrainingTimeSeconds: response.total_training_time_seconds,
  } satisfies WorkoutStats
}

export async function completeWorkoutTraining(
  accessToken: string,
  dayOfWeek: number | string,
  durationSeconds: number,
) {
  return request<{ success: boolean }>('/workout/training/complete', {
    accessToken,
    body: {
      day_of_week: dayOfWeek,
      duration_seconds: durationSeconds,
    },
    method: 'POST',
  })
}

async function request<T>(path: string, options: RequestOptions) {
  const headers = new Headers()
  const method = options.method ?? 'GET'

  if (options.body !== undefined) {
    headers.set('Content-Type', 'application/json')
  }

  if (options.accessToken) {
    headers.set('Authorization', `Bearer ${options.accessToken}`)
  }

  const response = await fetch(`${API_BASE_URL}${path}`, {
    body: options.body === undefined ? undefined : JSON.stringify(options.body),
    credentials: 'include',
    headers,
    method,
  })

  const isJson = response.headers.get('content-type')?.includes('application/json')
  const payload = isJson ? await response.json() : null

  if (!response.ok) {
    const message = getErrorMessage(payload, response.status)
    throw new ApiError(message, response.status, payload)
  }

  return payload as T
}

function mapAuthResponse(response: AuthResponseDto) {
  return {
    accessToken: response.access_token,
    accessTokenExpiresAt: response.access_token_expires_at,
    profileCompleted: response.profile_completed,
    userId: response.user_id,
  } satisfies AuthSession
}

function mapGenerationResponse(response: GenerationResponseDto) {
  return {
    errorMessage: response.error_message,
    generationId: response.generation_id,
    status: normalizeGenerationStatus(response.status),
  } satisfies GenerationResult
}

function mapNutritionFacts(response?: NutritionFactsDto | null) {
  return {
    calories: response?.calories ?? 0,
    protein: response?.protein ?? 0,
    fat: response?.fat ?? 0,
    carb: response?.carb ?? 0,
  } satisfies NutritionFacts
}

function normalizeGenerationStatus(status: string): GenerationStatus {
  if (status === 'done' || status === 'failed' || status === 'pending') {
    return status
  }

  return 'unspecified'
}

function getErrorMessage(payload: unknown, status: number) {
  if (payload && typeof payload === 'object' && 'error' in payload) {
    const { error } = payload
    if (typeof error === 'string' && error.trim() !== '') {
      return error
    }
  }

  return `Request failed with status ${status}`
}
