import { useEffect } from 'react'
import { useAuth } from './auth/useAuth'
import { HomePage } from './pages/HomePage'
import { LoginPage } from './pages/LoginPage'
import { NotFoundPage } from './pages/NotFoundPage'
import { ProfilePage } from './pages/ProfilePage'
import { RegisterPage } from './pages/RegisterPage'
import { navigate, usePathname } from './router'

function App() {
  const pathname = normalizePathname(usePathname())
  const { isAuthenticated, session, status } = useAuth()
  const authenticatedPath = session?.profileCompleted ? '/app' : '/profile'

  if (status === 'loading') {
    return (
      <FullscreenState
        title="Подключаем сессию"
        description="Проверяем access token и refresh cookie перед загрузкой страницы."
      />
    )
  }

  if (pathname === '/') {
    return <Redirect replace to={isAuthenticated ? authenticatedPath : '/login'} />
  }

  if (pathname === '/login') {
    return isAuthenticated ? <Redirect replace to={authenticatedPath} /> : <LoginPage />
  }

  if (pathname === '/register') {
    return isAuthenticated ? <Redirect replace to={authenticatedPath} /> : <RegisterPage />
  }

  if (pathname === '/app') {
    return isAuthenticated ? <HomePage /> : <Redirect replace to="/login" />
  }

  if (pathname === '/profile') {
    return isAuthenticated ? <ProfilePage /> : <Redirect replace to="/login" />
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
