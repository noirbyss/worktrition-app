const API_BASE_URL = (import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080').replace(/\/+$/, '')

type HttpMethod = 'GET' | 'POST'

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

function getErrorMessage(payload: unknown, status: number) {
  if (payload && typeof payload === 'object' && 'error' in payload) {
    const { error } = payload
    if (typeof error === 'string' && error.trim() !== '') {
      return error
    }
  }

  return `Request failed with status ${status}`
}
