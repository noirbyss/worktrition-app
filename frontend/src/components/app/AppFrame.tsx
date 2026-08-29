import type { ReactNode } from 'react'
import { usePathname } from '../../router'
import Logo from '../logo/logo'
import { AppLink } from '../navigation/AppLink'

interface AppFrameProps {
  actions?: ReactNode
  children: ReactNode
  description: string
  title: string
}

export function AppFrame({
  actions,
  children,
  description,
  title,
}: AppFrameProps) {
  const pathname = usePathname()

  return (
    <div className="app-page">
      <header className="app-header">
        <Logo compact />
        <nav aria-label="Основная навигация" className="app-nav">
          <AppLink
            className={pathname === '/app' ? 'app-nav__link is-active' : 'app-nav__link'}
            href="/app"
          >
            Аккаунт
          </AppLink>
          <AppLink
            className={pathname === '/profile' ? 'app-nav__link is-active' : 'app-nav__link'}
            href="/profile"
          >
            Профиль
          </AppLink>
        </nav>
      </header>

      <main className="app-main">
        <section className="app-hero">
          <div className="app-hero__copy">
            <p className="app-eyebrow">WORKTRITION AUTH</p>
            <h1 className="app-screen-title">{title}</h1>
            <p className="app-description">{description}</p>
          </div>
          {actions ? <div className="app-hero__actions">{actions}</div> : null}
        </section>
        {children}
      </main>
    </div>
  )
}
