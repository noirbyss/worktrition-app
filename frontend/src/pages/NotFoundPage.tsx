import { AppFrame } from '../components/app/AppFrame'
import { AuthLayout } from '../components/auth/AuthLayout'
import { AppLink } from '../components/navigation/AppLink'

interface NotFoundPageProps {
  authenticatedPath: string
  isAuthenticated: boolean
}

export function NotFoundPage({
  authenticatedPath,
  isAuthenticated,
}: NotFoundPageProps) {
  if (!isAuthenticated) {
    return (
      <AuthLayout
        switchHref="/login"
        switchLabel="Вернуться ко входу"
        switchText="Маршрут не найден."
        title="404"
      >
        <p className="auth-subtitle">
          Такой страницы нет. Вернитесь на экран авторизации и продолжим оттуда.
        </p>
      </AuthLayout>
    )
  }

  return (
    <AppFrame
      actions={
        <div className="button-row">
          <AppLink className="btn btn--small" href={authenticatedPath}>
            На главную страницу
          </AppLink>
        </div>
      }
      description="Маршрут не найден, но ваша сессия сохранена."
      title="404"
    >
      <section className="panel">
        <h2 className="panel-title">Страница не существует</h2>
        <p className="panel-copy">
          Можно безопасно вернуться на основной экран без повторного входа в аккаунт.
        </p>
      </section>
    </AppFrame>
  )
}
