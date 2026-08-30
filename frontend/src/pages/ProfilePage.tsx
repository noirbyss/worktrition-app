import { useEffect, useState, type ReactNode } from 'react'
import { ApiError, type Profile } from '../api'
import { useAuth } from '../auth/useAuth'
import { AppFrame } from '../components/app/AppFrame'
import { InlineMessage } from '../components/auth/InlineMessage'
import { useCurrentUserData } from '../hooks'
import { toErrorMessage } from '../utils'

const genderLabels: Record<number, string> = {
    1: 'Мужской',
    2: 'Женский',
}

const trainingLevelLabels: Record<number, string> = {
    1: 'Новичок',
    2: 'Средний',
    3: 'Продвинутый',
}

const activityLevelLabels: Record<number, string> = {
    1: 'Сидячий',
    2: 'Лёгкий',
    3: 'Средний',
    4: 'Высокий',
}

const goalLabels: Record<number, string> = {
    1: 'Снижение веса',
    2: 'Поддержание веса',
    3: 'Набор мышечной массы',
}

const trainingLocationLabels: Record<number, string> = {
    1: 'Дома',
    2: 'В зале',
}

export function ProfilePage() {
    const { getProfile } = useAuth()
    const { isLoading: isLoadingUser, loadError: userLoadError, user } = useCurrentUserData()
    const [isLoadingProfile, setIsLoadingProfile] = useState(true)
    const [profile, setProfile] = useState<Profile | null>(null)
    const [profileLoadError, setProfileLoadError] = useState<string | null>(null)
    const [profileMissing, setProfileMissing] = useState(false)

    useEffect(() => {
        let isCancelled = false

        const loadProfile = async () => {
            try {
                setIsLoadingProfile(true)
                setProfileLoadError(null)
                const nextProfile = await getProfile()

                if (!isCancelled) {
                    setProfile(nextProfile)
                    setProfileMissing(false)
                }
            } catch (error) {
                if (!isCancelled) {
                    if (error instanceof ApiError && error.status === 404) {
                        setProfile(null)
                        setProfileMissing(true)
                    } else {
                        setProfileLoadError(toErrorMessage(error, 'Не удалось загрузить данные профиля.'))
                    }
                }
            } finally {
                if (!isCancelled) {
                    setIsLoadingProfile(false)
                }
            }
        }

        void loadProfile()

        return () => {
            isCancelled = true
        }
    }, [getProfile])

    const avatarInitials = getInitials(user?.name)
    const birthDate = user ? formatShortDate(user.birthDate) : '—'
    const isLoading = isLoadingUser || isLoadingProfile

    return (
        <AppFrame
            currentUser={user}
            description="Твоя персональная информация и параметры, на основе которых формируются планы тренировок и питания."
            eyebrow="Экран 05"
            isCurrentUserLoading={isLoadingUser}
            title="Профиль"
        >
            {userLoadError ? <InlineMessage>{userLoadError}</InlineMessage> : null}
            {profileLoadError ? <InlineMessage>{profileLoadError}</InlineMessage> : null}

            <section className="card frame profile-hero">
                <div className="avatar-lg">
                    <div className="avatar-lg-inner">{avatarInitials}</div>
                </div>

                <div>
                    <h2 className="profile-name">{user?.name ?? 'Загружаем профиль...'}</h2>
                    <div className="profile-email">{user?.email ?? 'Подтягиваем данные из gateway...'}</div>
                    <div className="profile-meta">
                        <span className="meta-chip">
                            Дата рождения: <b>{birthDate}</b>
                        </span>
                    </div>
                </div>
            </section>

            {profileMissing && !isLoading ? (
                <section className="section">
                    <div className="card empty-state">
                        <div className="card-title">Профиль не найден</div>
                        <p className="panel-copy">
                            Gateway вернул `404 profile not found`. Если это произошло после завершения анкеты,
                            стоит проверить сохранение данных в user-service.
                        </p>
                    </div>
                </section>
            ) : null}

            <section className="section">
                <div className="section-head">
                    <h2 className="section-title">Физические параметры</h2>
                    <span className="section-note">основные данные</span>
                </div>
                <div className="grid g-5">
                    <InfoStatCard
                        label="Пол"
                        value={profile ? genderLabels[profile.gender] ?? 'Не указан' : getLoadingValue(isLoading)}
                    />
                    <InfoStatCard label="Возраст" value={profile ? `${profile.age} лет` : getLoadingValue(isLoading)} />
                    <InfoStatCard label="Рост" value={profile ? `${formatNumber(profile.heightCm)} см` : getLoadingValue(isLoading)} />
                    <InfoStatCard label="Текущий вес" value={profile ? `${formatNumber(profile.weightKg)} кг` : getLoadingValue(isLoading)} />
                    <InfoStatCard
                        label="Целевой вес"
                        value={
                            profile
                                ? profile.targetWeightKg
                                    ? `${formatNumber(profile.targetWeightKg)} кг`
                                    : 'Не указан'
                                : getLoadingValue(isLoading)
                        }
                    />
                </div>
            </section>

            <section className="section">
                <div className="section-head">
                    <h2 className="section-title">Тренировки и питание</h2>
                    <span className="section-note">персональные настройки</span>
                </div>
                <div className="grid g-2">
                    <div className="card">
                        <div className="card-title">Тренировочный профиль</div>
                        <div className="info-list">
                            <InfoRow
                                label="Уровень подготовки"
                                value={profile ? trainingLevelLabels[profile.trainingLevel] ?? 'Не указан' : getLoadingValue(isLoading)}
                                valueClassName="value-accent"
                            />
                            <InfoRow
                                label="Активность"
                                value={profile ? activityLevelLabels[profile.activityLevel] ?? 'Не указана' : getLoadingValue(isLoading)}
                            />
                            <InfoRow
                                label="Цель"
                                value={profile ? goalLabels[profile.goal] ?? 'Не указана' : getLoadingValue(isLoading)}
                            />
                        </div>
                    </div>

                    <div className="card">
                        <div className="card-title">Параметры тренировок</div>
                        <div className="info-list">
                            <InfoRow
                                label="Место тренировок"
                                value={profile ? trainingLocationLabels[profile.trainingLocation] ?? 'Не указано' : getLoadingValue(isLoading)}
                            />
                            <InfoRow
                                label="Дней в неделю"
                                value={profile ? formatTrainingDays(profile.trainingDaysPerWeek) : getLoadingValue(isLoading)}
                            />
                            <InfoRow
                                label="Инвентарь"
                                value={profile ? profile.equipment || 'Не указан' : getLoadingValue(isLoading)}
                            />
                        </div>
                    </div>
                </div>
            </section>

            <section className="section">
                <div className="section-head">
                    <h2 className="section-title">Пищевые предпочтения</h2>
                    <span className="section-note">учитываются при составлении рациона</span>
                </div>
                <div className="preference-grid">
                    <PreferenceCard
                        hint={buildCountLabel(profile?.allergies.length ?? 0)}
                        icon={PreferenceIconHeart}
                        title="Аллергии"
                        values={profile?.allergies ?? []}
                        variant="allergy"
                    />
                    <PreferenceCard
                        hint={buildCountLabel(profile?.excludedFoods.length ?? 0)}
                        icon={PreferenceIconBan}
                        title="Исключённые продукты"
                        values={profile?.excludedFoods ?? []}
                        variant="exclude"
                    />
                    <PreferenceCard
                        hint={buildCountLabel(profile?.foodPreferences.length ?? 0)}
                        icon={PreferenceIconLike}
                        title="Пищевые предпочтения"
                        values={profile?.foodPreferences ?? []}
                        variant="like"
                    />
                </div>
            </section>

            <footer className="foot">worktrition · персональный план тренировок и питания · профиль пользователя</footer>
        </AppFrame>
    )
}

function InfoStatCard({ label, value }: { label: string; value: string }) {
    return (
        <div className="card info-card">
            <div className="info-label">{label}</div>
            <div className="info-value">{value}</div>
        </div>
    )
}

function InfoRow({
    label,
    value,
    valueClassName,
}: {
    label: string
    value: string
    valueClassName?: string
}) {
    return (
        <div className="info-row">
            <span className="label">{label}</span>
            <span className={valueClassName ? `value ${valueClassName}` : 'value'}>{value}</span>
        </div>
    )
}

function PreferenceCard({
    hint,
    icon: Icon,
    title,
    values,
    variant,
}: {
    hint: string
    icon: () => ReactNode
    title: string
    values: string[]
    variant: 'allergy' | 'exclude' | 'like'
}) {
    return (
        <div className="card pref-card">
            <div className="pref-top">
                <div className="pref-icon">
                    <Icon />
                </div>
                <div>
                    <div className="pref-name">{title}</div>
                    <div className="pref-hint">{hint}</div>
                </div>
            </div>
            <div className="tags">
                {values.length > 0 ? (
                    values.map((value) => (
                        <span className={`tag ${variant}`} key={value}>
                            {value}
                        </span>
                    ))
                ) : (
                    <span className="tag tag--muted">Не указано</span>
                )}
            </div>
        </div>
    )
}

function PreferenceIconHeart() {
    return (
        <svg fill="none" height="18" viewBox="0 0 24 24" width="18">
            <path d="M12 3v11" stroke="currentColor" strokeLinecap="round" strokeWidth="1.6" />
            <path
                d="M8.5 6.5c1.1-2.1 2.3-3 3.5-3s2.4.9 3.5 3c1 1.9.5 4.2-1.2 5.4L12 13.5l-2.3-1.6C8 10.7 7.5 8.4 8.5 6.5Z"
                stroke="currentColor"
                strokeWidth="1.4"
            />
        </svg>
    )
}

function PreferenceIconBan() {
    return (
        <svg fill="none" height="18" viewBox="0 0 24 24" width="18">
            <circle cx="12" cy="12" r="8" stroke="currentColor" strokeWidth="1.5" />
            <path d="M8 8l8 8" stroke="currentColor" strokeLinecap="round" strokeWidth="1.5" />
        </svg>
    )
}

function PreferenceIconLike() {
    return (
        <svg fill="none" height="18" viewBox="0 0 24 24" width="18">
            <path
                d="M12 20s-7-4.4-7-10a4 4 0 0 1 7-2.6A4 4 0 0 1 19 10c0 5.6-7 10-7 10Z"
                stroke="currentColor"
                strokeLinejoin="round"
                strokeWidth="1.5"
            />
        </svg>
    )
}

function getInitials(name?: string | null) {
    const parts = (name ?? '')
        .trim()
        .split(/\s+/)
        .filter(Boolean)

    if (parts.length === 0) {
        return 'WT'
    }

    return parts
        .slice(0, 2)
        .map((part) => part[0]?.toUpperCase() ?? '')
        .join('')
}

function buildCountLabel(count: number) {
    if (count === 0) {
        return 'нет данных'
    }

    return `${count} ${pluralizeItems(count)}`
}

function pluralizeItems(count: number) {
    const mod10 = count % 10
    const mod100 = count % 100

    if (mod10 === 1 && mod100 !== 11) {
        return 'элемент'
    }

    if (mod10 >= 2 && mod10 <= 4 && (mod100 < 12 || mod100 > 14)) {
        return 'элемента'
    }

    return 'элементов'
}

function formatShortDate(value: string) {
    const parsedDate = new Date(`${value}T00:00:00`)
    if (Number.isNaN(parsedDate.getTime())) {
        return value
    }

    return new Intl.DateTimeFormat('ru-RU').format(parsedDate)
}

function formatTrainingDays(days: number) {
    if (days === 0) {
        return '0 дней'
    }

    if (days === 1) {
        return '1 день'
    }

    if (days >= 2 && days <= 4) {
        return `${days} дня`
    }

    return `${days} дней`
}

function formatNumber(value: number) {
    return Number.isInteger(value) ? String(value) : value.toFixed(1)
}

function getLoadingValue(isLoading: boolean) {
    return isLoading ? 'Загружаем...' : 'Не указано'
}
