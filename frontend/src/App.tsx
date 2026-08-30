import { useEffect } from 'react'
import { useAuth } from './auth/useAuth'
import { HomePage } from './pages/HomePage'
import { LoginPage } from './pages/LoginPage'
import { NotFoundPage } from './pages/NotFoundPage'
import { NutritionPage } from './pages/NutritionPage'
import { ProfilePage } from './pages/ProfilePage'
import { QuestionnairePage } from './pages/QuestionnairePage'
import { RegisterPage } from './pages/RegisterPage'
import { StatsPage } from './pages/StatsPage'
import { WorkoutsPage } from './pages/WorkoutsPage'
import { navigate, usePathname } from './router'

const questionnairePath = '/questionnaire'
const defaultAuthenticatedPath = '/app'

function App() {
    const pathname = normalizePathname(usePathname())
    const { isAuthenticated, session, status } = useAuth()
    const profileCompleted = Boolean(session?.profileCompleted)
    const authenticatedPath = profileCompleted ? defaultAuthenticatedPath : questionnairePath

    if (status === 'loading') {
        return (
            <FullscreenState
                title="Подключаем сессию"
                description="Загружаем страницу!"
            />
        )
    }

    if (pathname === '/') {
        return <Redirect replace to={isAuthenticated ? authenticatedPath : '/login'} />
    }

    if (!isAuthenticated) {
        if (pathname === '/login') {
            return <LoginPage />
        }

        if (pathname === '/register') {
            return <RegisterPage />
        }

        return <Redirect replace to="/login" />
    }

    if (!profileCompleted) {
        if (pathname === questionnairePath) {
            return <QuestionnairePage />
        }

        return <Redirect replace to={questionnairePath} />
    }

    if (pathname === '/login' || pathname === '/register') {
        return <Redirect replace to={defaultAuthenticatedPath} />
    }

    if (pathname === questionnairePath) {
        return <Redirect replace to={defaultAuthenticatedPath} />
    }

    if (pathname === '/app') {
        return <HomePage />
    }

    if (pathname === '/profile') {
        return <ProfilePage />
    }

    if (pathname === '/nutrition') {
        return <NutritionPage />
    }

    if (pathname === '/stats') {
        return <StatsPage />
    }

    if (pathname === '/workouts') {
        return <WorkoutsPage />
    }

    return (
        <NotFoundPage
            authenticatedPath={authenticatedPath}
            isAuthenticated={isAuthenticated}
        />
    )
}

function Redirect({
    replace = false,
    to,
}: {
    replace?: boolean
    to: string
}) {
    useEffect(() => {
        navigate(to, { replace })
    }, [replace, to])

    return (
        <FullscreenState
            title="Переходим дальше"
            description="Маршрут обновляется, это займет меньше секунды."
        />
    )
}

function FullscreenState({
    description,
    title,
}: {
    description: string
    title: string
}) {
    return (
        <div className="auth-page">
            <section className="auth-card auth-card--state">
                <p className="app-eyebrow">WORKTRITION</p>
                <h1 className="auth-title">{title}</h1>
                <p className="auth-subtitle">{description}</p>
            </section>
        </div>
    )
}

function normalizePathname(pathname: string) {
    if (pathname === '/') {
        return pathname
    }

    return pathname.replace(/\/+$/, '')
}

export default App
