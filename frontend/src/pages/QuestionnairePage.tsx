import { useEffect, useState, type KeyboardEvent, type ReactNode } from 'react'
import worktritionLogo from '../assets/worktrition-logo.png'
import type { CurrentUser, SaveProfilePayload } from '../api'
import { useAuth } from '../auth/useAuth'
import { InlineMessage } from '../components/auth/InlineMessage'
import { navigate } from '../router'
import { mapErrorToFieldErrors, toErrorMessage } from '../utils'

type QuestionnaireField =
  | 'activityLevel'
  | 'allergies'
  | 'equipment'
  | 'excludedFoods'
  | 'foodPreferences'
  | 'gender'
  | 'goal'
  | 'heightCm'
  | 'targetWeightKg'
  | 'trainingDaysPerWeek'
  | 'trainingLevel'
  | 'trainingLocation'
  | 'weightKg'

type QuestionnaireDraftField = 'allergies' | 'excludedFoods' | 'foodPreferences'

interface QuestionnaireFormState {
  activityLevel: string
  allergies: string[]
  equipment: string
  excludedFoods: string[]
  foodPreferences: string[]
  gender: string
  goal: string
  heightCm: string
  targetWeightKg: string
  trainingDaysPerWeek: string
  trainingLevel: string
  trainingLocation: string
  weightKg: string
}

interface DraftState {
  allergies: string
  excludedFoods: string
  foodPreferences: string
}

const totalSteps = 5

const initialFormState: QuestionnaireFormState = {
  activityLevel: '',
  allergies: [],
  equipment: '',
  excludedFoods: [],
  foodPreferences: [],
  gender: '',
  goal: '',
  heightCm: '',
  targetWeightKg: '',
  trainingDaysPerWeek: '',
  trainingLevel: '',
  trainingLocation: '',
  weightKg: '',
}

const initialDraftState: DraftState = {
  allergies: '',
  excludedFoods: '',
  foodPreferences: '',
}

const fieldMap = {
  activity_level: 'activityLevel',
  allergies: 'allergies',
  equipment: 'equipment',
  excluded_foods: 'excludedFoods',
  food_preferences: 'foodPreferences',
  gender: 'gender',
  goal: 'goal',
  height_cm: 'heightCm',
  target_weight_kg: 'targetWeightKg',
  training_days_per_week: 'trainingDaysPerWeek',
  training_level: 'trainingLevel',
  training_location: 'trainingLocation',
  weight_kg: 'weightKg',
} satisfies Record<string, QuestionnaireField>

const fieldStepMap: Record<QuestionnaireField, number> = {
  activityLevel: 1,
  allergies: 3,
  equipment: 4,
  excludedFoods: 3,
  foodPreferences: 3,
  gender: 0,
  goal: 2,
  heightCm: 0,
  targetWeightKg: 2,
  trainingDaysPerWeek: 4,
  trainingLevel: 1,
  trainingLocation: 4,
  weightKg: 0,
}

export function QuestionnairePage() {
  const { getCurrentUser, logout, saveProfile } = useAuth()
  const [currentStep, setCurrentStep] = useState(0)
  const [currentUser, setCurrentUser] = useState<CurrentUser | null>(null)
  const [draftValues, setDraftValues] = useState(initialDraftState)
  const [fieldErrors, setFieldErrors] = useState<Partial<Record<QuestionnaireField, string>>>({})
  const [formValues, setFormValues] = useState(initialFormState)
  const [isLoadingUser, setIsLoadingUser] = useState(true)
  const [isLoggingOut, setIsLoggingOut] = useState(false)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [loadError, setLoadError] = useState<string | null>(null)
  const [submitError, setSubmitError] = useState<string | null>(null)

  useEffect(() => {
    let isCancelled = false

    const loadUser = async () => {
      try {
        setIsLoadingUser(true)
        setLoadError(null)
        const nextUser = await getCurrentUser()

        if (!isCancelled) {
          setCurrentUser(nextUser)
        }
      } catch (error) {
        if (!isCancelled) {
          setLoadError(toErrorMessage(error, 'Не удалось загрузить данные пользователя.'))
        }
      } finally {
        if (!isCancelled) {
          setIsLoadingUser(false)
        }
      }
    }

    void loadUser()

    return () => {
      isCancelled = true
    }
  }, [getCurrentUser])

  const isEquipmentEnabled = formValues.trainingLocation === '1'
  const derivedAge = currentUser ? calculateAge(currentUser.birthDate) : null

  const handleLogout = async () => {
    try {
      setIsLoggingOut(true)
      await logout()
      navigate('/login', { replace: true })
    } finally {
      setIsLoggingOut(false)
    }
  }

  const handleValueChange = (field: keyof QuestionnaireFormState, value: string) => {
    setFormValues((current) => ({
      ...current,
      [field]: value,
    }))
    setFieldErrors((current) => ({ ...current, [field]: undefined }))
    setSubmitError(null)
  }

  const handleDraftChange = (field: QuestionnaireDraftField, value: string) => {
    setDraftValues((current) => ({
      ...current,
      [field]: value,
    }))
    setFieldErrors((current) => ({ ...current, [field]: undefined }))
    setSubmitError(null)
  }

  const handleTrainingLocationChange = (value: string) => {
    setFormValues((current) => ({
      ...current,
      equipment: value === '1' ? current.equipment : '',
      trainingLocation: value,
    }))
    setFieldErrors((current) => ({
      ...current,
      equipment: undefined,
      trainingLocation: undefined,
    }))
    setSubmitError(null)
  }

  const addDynamicItem = (field: QuestionnaireDraftField) => {
    const nextValue = draftValues[field].trim()
    if (!nextValue) {
      return
    }

    setFormValues((current) => {
      const existingValues = current[field]
      const alreadyExists = existingValues.some((value) => value.toLowerCase() === nextValue.toLowerCase())
      if (alreadyExists) {
        return current
      }

      return {
        ...current,
        [field]: [...existingValues, nextValue],
      }
    })

    setDraftValues((current) => ({
      ...current,
      [field]: '',
    }))
    setFieldErrors((current) => ({ ...current, [field]: undefined }))
    setSubmitError(null)
  }

  const removeDynamicItem = (field: QuestionnaireDraftField, valueToRemove: string) => {
    setFormValues((current) => ({
      ...current,
      [field]: current[field].filter((value) => value !== valueToRemove),
    }))
    setFieldErrors((current) => ({ ...current, [field]: undefined }))
    setSubmitError(null)
  }

  const handleDynamicInputKeyDown = (
    event: KeyboardEvent<HTMLInputElement>,
    field: QuestionnaireDraftField,
  ) => {
    if (event.key !== 'Enter') {
      return
    }

    event.preventDefault()
    addDynamicItem(field)
  }

  const handleNext = async () => {
    const stepErrors = validateStep(currentStep, formValues, derivedAge)
    if (Object.keys(stepErrors).length > 0) {
      setFieldErrors((current) => ({ ...current, ...stepErrors }))
      return
    }

    if (currentStep < totalSteps - 1) {
      setCurrentStep((current) => current + 1)
      return
    }

    const validationResult = validateQuestionnaireForm(formValues, derivedAge)
    setFieldErrors(validationResult.errors)
    setSubmitError(null)

    if (Object.keys(validationResult.errors).length > 0 || !validationResult.payload) {
      setCurrentStep(findStepForErrors(validationResult.errors))
      return
    }

    try {
      setIsSubmitting(true)
      await saveProfile(validationResult.payload)
      navigate('/profile', { replace: true })
    } catch (error) {
      const mappedErrors = mapErrorToFieldErrors(error, fieldMap)
      setFieldErrors(mappedErrors.fieldErrors)
      setSubmitError(mappedErrors.formError)

      if (Object.keys(mappedErrors.fieldErrors).length > 0) {
        setCurrentStep(findStepForErrors(mappedErrors.fieldErrors))
      }
    } finally {
      setIsSubmitting(false)
    }
  }

  const handlePrevious = () => {
    if (currentStep === 0) {
      return
    }

    setCurrentStep((current) => current - 1)
  }

  const progressBars = Array.from({ length: totalSteps }, (_, index) => index <= currentStep)
  const isLastStep = currentStep === totalSteps - 1

  return (
    <div className="onboarding">
      <img alt="Worktrition" className="brand-logo" src={worktritionLogo} />

      <main className="wizard">
        <div className="wizard-toolbar">
          <div className="wizard-user">
            <strong className="wizard-user__name">{currentUser?.name ?? 'Загружаем...'}</strong>
            {currentUser?.email ? <span className="wizard-user__email">{currentUser.email}</span> : null}
          </div>
          <button
            className="btn btn-secondary wizard-exit"
            disabled={isLoggingOut}
            onClick={() => {
              void handleLogout()
            }}
            type="button"
          >
            {isLoggingOut ? 'ВЫХОД...' : 'ВЫЙТИ'}
          </button>
        </div>

        <div className="top">
          <div className="step-label">
            НАСТРОЙКА ПРОФИЛЯ · ШАГ {String(currentStep + 1).padStart(2, '0')}
          </div>
          <h1 className="wizard-title">СОЗДАЙ СВОЙ ПЛАН</h1>
          <p className="subtitle">
            Ответь на несколько вопросов. На их основе мы сможем сформировать персональные
            тренировки и план питания.
          </p>
          <div aria-label="Прогресс анкеты" className="progress">
            {progressBars.map((isActive, index) => (
              <span className={isActive ? 'active' : undefined} key={index} />
            ))}
          </div>
        </div>

        <form
          onSubmit={(event) => {
            event.preventDefault()
            void handleNext()
          }}
        >
          <section className="card">
            {loadError ? <InlineMessage>{loadError}</InlineMessage> : null}
            {submitError ? <InlineMessage>{submitError}</InlineMessage> : null}

            {renderStep(
              currentStep === 0,
              <>
                <div className="section-kicker">01 / ОСНОВА</div>
                <h2 className="wizard-step-title">РАССКАЖИ О СЕБЕ</h2>
                <p className="step-desc">Начнём с основных физических параметров.</p>

                <div className="grid-2">
                  <div className="field">
                    <label htmlFor="age">Возраст</label>
                    <input
                      className="wizard-input wizard-input--disabled"
                      disabled
                      id="age"
                      placeholder={isLoadingUser ? 'Загружаем...' : 'Недоступно'}
                      type="text"
                      value={derivedAge === null ? '' : String(derivedAge)}
                    />
                  </div>

                  <div className="field">
                    <div className="group-label">Пол</div>
                    <div className="options two">
                      <OptionCard
                        checked={formValues.gender === '1'}
                        description={null}
                        name="gender"
                        onChange={() => {
                          handleValueChange('gender', '1')
                        }}
                        title="Мужской"
                        value="1"
                      />
                      <OptionCard
                        checked={formValues.gender === '2'}
                        description={null}
                        name="gender"
                        onChange={() => {
                          handleValueChange('gender', '2')
                        }}
                        title="Женский"
                        value="2"
                      />
                    </div>
                    {fieldErrors.gender ? <p className="field-error">{fieldErrors.gender}</p> : null}
                  </div>
                </div>

                <div className="grid-2">
                  <FieldInput
                    error={fieldErrors.heightCm}
                    id="height"
                    inputMode="numeric"
                    label="Рост, см"
                    max={250}
                    min={80}
                    onChange={(value) => {
                      handleValueChange('heightCm', value)
                    }}
                    placeholder="Например, 180"
                    step="1"
                    type="number"
                    value={formValues.heightCm}
                  />
                  <FieldInput
                    error={fieldErrors.weightKg}
                    id="weight"
                    inputMode="decimal"
                    label="Вес, кг"
                    max={400}
                    min={25}
                    onChange={(value) => {
                      handleValueChange('weightKg', value)
                    }}
                    placeholder="Например, 75.5"
                    step="0.1"
                    type="number"
                    value={formValues.weightKg}
                  />
                </div>
              </>,
            )}

            {renderStep(
              currentStep === 1,
              <>
                <div className="section-kicker">02 / АКТИВНОСТЬ</div>
                <h2 className="wizard-step-title">ТЕКУЩАЯ ФОРМА</h2>
                <p className="step-desc">Это поможет подобрать подходящую нагрузку.</p>

                <div className="field">
                  <div className="group-label">Уровень подготовки</div>
                  <div className="options three">
                    <OptionCard
                      checked={formValues.trainingLevel === '1'}
                      description="Только начинаю"
                      name="training_level"
                      onChange={() => {
                        handleValueChange('trainingLevel', '1')
                      }}
                      title="Новичок"
                      value="1"
                    />
                    <OptionCard
                      checked={formValues.trainingLevel === '2'}
                      description="Есть регулярный опыт"
                      name="training_level"
                      onChange={() => {
                        handleValueChange('trainingLevel', '2')
                      }}
                      title="Средний"
                      value="2"
                    />
                    <OptionCard
                      checked={formValues.trainingLevel === '3'}
                      description="Тренируюсь давно"
                      name="training_level"
                      onChange={() => {
                        handleValueChange('trainingLevel', '3')
                      }}
                      title="Продвинутый"
                      value="3"
                    />
                  </div>
                  {fieldErrors.trainingLevel ? (
                    <p className="field-error">{fieldErrors.trainingLevel}</p>
                  ) : null}
                </div>

                <div className="field">
                  <div className="group-label">Уровень повседневной активности</div>
                  <div className="options two">
                    <OptionCard
                      checked={formValues.activityLevel === '1'}
                      description="Почти нет физической активности"
                      name="activity_level"
                      onChange={() => {
                        handleValueChange('activityLevel', '1')
                      }}
                      title="Сидячий"
                      value="1"
                    />
                    <OptionCard
                      checked={formValues.activityLevel === '2'}
                      description="Немного движения в течение дня"
                      name="activity_level"
                      onChange={() => {
                        handleValueChange('activityLevel', '2')
                      }}
                      title="Лёгкий"
                      value="2"
                    />
                    <OptionCard
                      checked={formValues.activityLevel === '3'}
                      description="Регулярная физическая активность"
                      name="activity_level"
                      onChange={() => {
                        handleValueChange('activityLevel', '3')
                      }}
                      title="Средний"
                      value="3"
                    />
                    <OptionCard
                      checked={formValues.activityLevel === '4'}
                      description="Много движения или физическая работа"
                      name="activity_level"
                      onChange={() => {
                        handleValueChange('activityLevel', '4')
                      }}
                      title="Высокий"
                      value="4"
                    />
                  </div>
                  {fieldErrors.activityLevel ? (
                    <p className="field-error">{fieldErrors.activityLevel}</p>
                  ) : null}
                </div>
              </>,
            )}

            {renderStep(
              currentStep === 2,
              <>
                <div className="section-kicker">03 / ЦЕЛЬ</div>
                <h2 className="wizard-step-title">К ЧЕМУ СТРЕМИМСЯ?</h2>
                <p className="step-desc">
                  Выбери главную цель, на которую будет ориентироваться твой план.
                </p>

                <div className="options two">
                  <OptionCard
                    checked={formValues.goal === '1'}
                    description="Дефицит калорий и поддержание активности"
                    name="goal"
                    onChange={() => {
                      handleValueChange('goal', '1')
                    }}
                    title="Снижение веса"
                    value="1"
                  />
                  <OptionCard
                    checked={formValues.goal === '2'}
                    description="Сохранить текущий вес и форму"
                    name="goal"
                    onChange={() => {
                      handleValueChange('goal', '2')
                    }}
                    title="Поддержание веса"
                    value="2"
                  />
                  <OptionCard
                    checked={formValues.goal === '3'}
                    description="Сила, прогрессия и питание"
                    name="goal"
                    onChange={() => {
                      handleValueChange('goal', '3')
                    }}
                    title="Набор мышечной массы"
                    value="3"
                  />
                </div>
                {fieldErrors.goal ? <p className="field-error">{fieldErrors.goal}</p> : null}

                <div className="field field--spaced">
                  <FieldInput
                    error={fieldErrors.targetWeightKg}
                    id="targetWeight"
                    inputMode="decimal"
                    label={
                      'Желаемый вес, кг'
                    }
                    max={400}
                    min={25}
                    onChange={(value) => {
                      handleValueChange('targetWeightKg', value)
                    }}
                    placeholder="Например, 70"
                    step="0.1"
                    type="number"
                    value={formValues.targetWeightKg}
                  />
                </div>
              </>,
            )}

            {renderStep(
              currentStep === 3,
              <>
                <div className="section-kicker">04 / ПИТАНИЕ</div>
                <h2 className="wizard-step-title">ПРЕДПОЧТЕНИЯ</h2>
                <p className="step-desc">
                  Отметь продукты и особенности питания, которые нужно учитывать при составлении
                  рациона.
                </p>

                <DynamicField
                  draftValue={draftValues.allergies}
                  error={fieldErrors.allergies}
                  field="allergies"
                  hint="Добавляй каждую аллергию отдельно."
                  inputId="allergyInput"
                  label={
                    <>
                      Аллергии <span className="optional">· если есть</span>
                    </>
                  }
                  onAdd={addDynamicItem}
                  onDraftChange={handleDraftChange}
                  onKeyDown={handleDynamicInputKeyDown}
                  onRemove={removeDynamicItem}
                  placeholder="Например: орехи"
                  values={formValues.allergies}
                />

                <DynamicField
                  draftValue={draftValues.excludedFoods}
                  error={fieldErrors.excludedFoods}
                  field="excludedFoods"
                  hint="Добавляй продукты, которые не хочешь видеть в рационе."
                  inputId="excludedFoodInput"
                  label={
                    <>
                      Исключить продукты <span className="optional">· необязательно</span>
                    </>
                  }
                  onAdd={addDynamicItem}
                  onDraftChange={handleDraftChange}
                  onKeyDown={handleDynamicInputKeyDown}
                  onRemove={removeDynamicItem}
                  placeholder="Например: печень"
                  values={formValues.excludedFoods}
                />

                <DynamicField
                  draftValue={draftValues.foodPreferences}
                  error={fieldErrors.foodPreferences}
                  field="foodPreferences"
                  hint="Например: вегетарианское питание, быстрые блюда, больше белка."
                  inputId="foodPreferenceInput"
                  label={
                    <>
                      Предпочтения в питании <span className="optional">· необязательно</span>
                    </>
                  }
                  onAdd={addDynamicItem}
                  onDraftChange={handleDraftChange}
                  onKeyDown={handleDynamicInputKeyDown}
                  onRemove={removeDynamicItem}
                  placeholder="Например: больше белка"
                  values={formValues.foodPreferences}
                />
              </>,
            )}

            {renderStep(
              currentStep === 4,
              <>
                <div className="section-kicker">05 / ТРЕНИРОВКИ</div>
                <h2 className="wizard-step-title">СОБЕРЁМ ТВОЙ ГРАФИК</h2>
                <p className="step-desc">
                  Осталось выбрать место, частоту и доступный инвентарь.
                </p>

                <div className="field">
                  <div className="group-label">Где ты тренируешься?</div>
                  <div className="options two">
                    <OptionCard
                      checked={formValues.trainingLocation === '1'}
                      description="Можно указать доступный инвентарь"
                      name="training_location"
                      onChange={() => {
                        handleTrainingLocationChange('1')
                      }}
                      title="Дома"
                      value="1"
                    />
                    <OptionCard
                      checked={formValues.trainingLocation === '2'}
                      description="Инвентарь уже доступен"
                      name="training_location"
                      onChange={() => {
                        handleTrainingLocationChange('2')
                      }}
                      title="В зале"
                      value="2"
                    />
                  </div>
                  {fieldErrors.trainingLocation ? (
                    <p className="field-error">{fieldErrors.trainingLocation}</p>
                  ) : null}
                </div>

                <div className="grid-2">
                  <div className="field">
                    <label htmlFor="trainingDays">Тренировок в неделю</label>
                    <select
                      className={fieldErrors.trainingDaysPerWeek ? 'wizard-select wizard-select--error' : 'wizard-select'}
                      id="trainingDays"
                      onChange={(event) => {
                        handleValueChange('trainingDaysPerWeek', event.target.value)
                      }}
                      value={formValues.trainingDaysPerWeek}
                    >
                      <option disabled value="">
                        Выбери количество
                      </option>
                      <option value="0">0 тренировок</option>
                      <option value="1">1 тренировка</option>
                      <option value="2">2 тренировки</option>
                      <option value="3">3 тренировки</option>
                      <option value="4">4 тренировки</option>
                      <option value="5">5 тренировок</option>
                      <option value="6">6 тренировок</option>
                      <option value="7">7 тренировок</option>
                    </select>
                    {fieldErrors.trainingDaysPerWeek ? (
                      <p className="field-error">{fieldErrors.trainingDaysPerWeek}</p>
                    ) : null}
                  </div>

                  <div className="field">
                    <label htmlFor="equipment">
                      Инвентарь{' '}
                      <span className="optional">
                        {isEquipmentEnabled
                          ? '· необязательно'
                          : '· недоступно для тренировок в зале'}
                      </span>
                    </label>
                    <input
                      className={fieldErrors.equipment ? 'wizard-input wizard-input--error' : 'wizard-input'}
                      disabled={!isEquipmentEnabled}
                      id="equipment"
                      onChange={(event) => {
                        handleValueChange('equipment', event.target.value)
                      }}
                      placeholder={
                        isEquipmentEnabled ? 'Например: гантели, резинки, турник' : ''
                      }
                      type="text"
                      value={formValues.equipment}
                    />
                    {fieldErrors.equipment ? <p className="field-error">{fieldErrors.equipment}</p> : null}
                  </div>
                </div>
              </>,
            )}

            <div className="actions">
              <button
                className="btn btn-secondary wizard-button"
                onClick={handlePrevious}
                style={currentStep === 0 ? { visibility: 'hidden' } : undefined}
                type="button"
              >
                НАЗАД
              </button>
              <div className="step-count">
                ШАГ {String(currentStep + 1).padStart(2, '0')} ИЗ {String(totalSteps).padStart(2, '0')}
              </div>
              <button className="btn btn-primary wizard-button" disabled={isSubmitting} type="submit">
                {isSubmitting ? 'СОХРАНЯЕМ...' : isLastStep ? 'СОЗДАТЬ МОЙ ПЛАН' : 'ДАЛЕЕ'}
              </button>
            </div>
          </section>
        </form>
      </main>
    </div>
  )
}

function renderStep(isActive: boolean, children: ReactNode) {
  return <div className={isActive ? 'step active' : 'step'}>{children}</div>
}

function OptionCard({
  checked,
  description,
  name,
  onChange,
  title,
  value,
}: {
  checked: boolean
  description: string | null
  name: string
  onChange: () => void
  title: string
  value: string
}) {
  return (
    <label className="option">
      <input checked={checked} name={name} onChange={onChange} type="radio" value={value} />
      <span>
        {description ? <strong>{title}</strong> : title}
        {description ? <small>{description}</small> : null}
      </span>
    </label>
  )
}

function FieldInput({
  error,
  id,
  inputMode,
  label,
  max,
  min,
  onChange,
  placeholder,
  step,
  type,
  value,
}: {
  error?: string
  id: string
  inputMode?: 'decimal' | 'numeric'
  label: ReactNode
  max?: number
  min?: number
  onChange: (value: string) => void
  placeholder?: string
  step?: string
  type: 'number' | 'text'
  value: string
}) {
  return (
    <div className="field">
      <label htmlFor={id}>{label}</label>
      <input
        className={error ? 'wizard-input wizard-input--error' : 'wizard-input'}
        id={id}
        inputMode={inputMode}
        max={max}
        min={min}
        onChange={(event) => {
          onChange(event.target.value)
        }}
        placeholder={placeholder}
        step={step}
        type={type}
        value={value}
      />
      {error ? <p className="field-error">{error}</p> : null}
    </div>
  )
}

function DynamicField({
  draftValue,
  error,
  field,
  hint,
  inputId,
  label,
  onAdd,
  onDraftChange,
  onKeyDown,
  onRemove,
  placeholder,
  values,
}: {
  draftValue: string
  error?: string
  field: QuestionnaireDraftField
  hint: string
  inputId: string
  label: ReactNode
  onAdd: (field: QuestionnaireDraftField) => void
  onDraftChange: (field: QuestionnaireDraftField, value: string) => void
  onKeyDown: (event: KeyboardEvent<HTMLInputElement>, field: QuestionnaireDraftField) => void
  onRemove: (field: QuestionnaireDraftField, value: string) => void
  placeholder: string
  values: string[]
}) {
  return (
    <div className="field">
      <label htmlFor={inputId}>{label}</label>
      <div className="dynamic-input">
        <input
          className={error ? 'wizard-input wizard-input--error' : 'wizard-input'}
          id={inputId}
          onChange={(event) => {
            onDraftChange(field, event.target.value)
          }}
          onKeyDown={(event) => {
            onKeyDown(event, field)
          }}
          placeholder={placeholder}
          type="text"
          value={draftValue}
        />
        <button
          className="add-btn"
          onClick={() => {
            onAdd(field)
          }}
          type="button"
        >
          ДОБАВИТЬ
        </button>
      </div>
      {values.length > 0 ? (
        <div className="dynamic-list">
          {values.map((value) => (
            <div className="dynamic-tag" key={value}>
              <span>{value}</span>
              <button
                aria-label={`Удалить ${value}`}
                onClick={() => {
                  onRemove(field, value)
                }}
                type="button"
              >
                ×
              </button>
            </div>
          ))}
        </div>
      ) : null}
      <div className="hint">{hint}</div>
      {error ? <p className="field-error">{error}</p> : null}
    </div>
  )
}

function validateStep(
  step: number,
  values: QuestionnaireFormState,
  derivedAge: number | null,
) {
  const errors: Partial<Record<QuestionnaireField, string>> = {}

  if (step === 0) {
    if (derivedAge === null || derivedAge < 13 || derivedAge > 120) {
      errors.heightCm = 'Не удалось корректно определить возраст по дате рождения.'
    }
    if (!values.gender) {
      errors.gender = 'Выберите пол.'
    }
    if (!isIntegerInRange(values.heightCm, 80, 250)) {
      errors.heightCm = 'Рост должен быть в диапазоне от 80 до 250 см.'
    }
    if (!isFloatInRange(values.weightKg, 25, 400)) {
      errors.weightKg = 'Вес должен быть в диапазоне от 25 до 400 кг.'
    }
  }

  if (step === 1) {
    if (!isIntegerInRange(values.trainingLevel, 1, 3)) {
      errors.trainingLevel = 'Выберите уровень тренировок.'
    }
    if (!isIntegerInRange(values.activityLevel, 1, 4)) {
      errors.activityLevel = 'Выберите уровень активности.'
    }
  }

  if (step === 2) {
    if (!isIntegerInRange(values.goal, 1, 3)) {
      errors.goal = 'Выберите цель.'
    }
    if (values.targetWeightKg.trim() && !isFloatInRange(values.targetWeightKg, 25, 400)) {
      errors.targetWeightKg = 'Целевой вес должен быть в диапазоне от 25 до 400 кг.'
    }
  }

  if (step === 3) {
    validateStringListField(values.allergies, 'allergies', errors)
    validateStringListField(values.excludedFoods, 'excludedFoods', errors)
    validateStringListField(values.foodPreferences, 'foodPreferences', errors)
  }

  if (step === 4) {
    if (!isIntegerInRange(values.trainingLocation, 1, 2)) {
      errors.trainingLocation = 'Выберите место тренировок.'
    }
    if (!isIntegerInRange(values.trainingDaysPerWeek, 0, 7)) {
      errors.trainingDaysPerWeek = 'Количество тренировок должно быть от 0 до 7.'
    }
    if (values.equipment.trim().length > 500) {
      errors.equipment = 'Описание инвентаря слишком длинное.'
    }
  }

  return errors
}

function validateQuestionnaireForm(
  values: QuestionnaireFormState,
  derivedAge: number | null,
) {
  const errors: Partial<Record<QuestionnaireField, string>> = {
    ...validateStep(0, values, derivedAge),
    ...validateStep(1, values, derivedAge),
    ...validateStep(2, values, derivedAge),
    ...validateStep(3, values, derivedAge),
    ...validateStep(4, values, derivedAge),
  }

  if (Object.keys(errors).length > 0) {
    return { errors, payload: null }
  }

  return {
    errors,
    payload: {
      activity_level: Number(values.activityLevel),
      allergies: values.allergies,
      equipment: values.equipment.trim(),
      excluded_foods: values.excludedFoods,
      food_preferences: values.foodPreferences,
      gender: Number(values.gender),
      goal: Number(values.goal),
      height_cm: Number(values.heightCm),
      target_weight_kg: values.targetWeightKg.trim() ? Number(values.targetWeightKg) : undefined,
      training_days_per_week: Number(values.trainingDaysPerWeek),
      training_level: Number(values.trainingLevel),
      training_location: Number(values.trainingLocation),
      weight_kg: Number(values.weightKg),
    } satisfies SaveProfilePayload,
  }
}

function findStepForErrors(errors: Partial<Record<QuestionnaireField, string>>) {
  const firstField = Object.keys(errors)
    .map((field) => field as QuestionnaireField)
    .sort((left, right) => fieldStepMap[left] - fieldStepMap[right])[0]

  return firstField ? fieldStepMap[firstField] : 0
}

function validateStringListField(
  values: string[],
  field: 'allergies' | 'excludedFoods' | 'foodPreferences',
  errors: Partial<Record<QuestionnaireField, string>>,
) {
  if (values.length > 50) {
    errors[field] = 'Список не должен содержать больше 50 элементов.'
    return
  }

  const hasTooLongItem = values.some((value) => value.length > 80)
  if (hasTooLongItem) {
    errors[field] = 'Каждый пункт списка должен быть короче 80 символов.'
  }
}

function isIntegerInRange(value: string, min: number, max: number) {
  const parsedValue = Number(value)
  return value.trim() !== '' && Number.isInteger(parsedValue) && parsedValue >= min && parsedValue <= max
}

function isFloatInRange(value: string, min: number, max: number) {
  const parsedValue = Number(value)
  return value.trim() !== '' && !Number.isNaN(parsedValue) && parsedValue >= min && parsedValue <= max
}

function calculateAge(birthDate: string) {
  const parsedDate = new Date(`${birthDate}T00:00:00`)
  if (Number.isNaN(parsedDate.getTime())) {
    return null
  }

  const now = new Date()
  let age = now.getFullYear() - parsedDate.getFullYear()
  const hasBirthdayPassed =
    now.getMonth() > parsedDate.getMonth() ||
    (now.getMonth() === parsedDate.getMonth() && now.getDate() >= parsedDate.getDate())

  if (!hasBirthdayPassed) {
    age -= 1
  }

  return age
}
