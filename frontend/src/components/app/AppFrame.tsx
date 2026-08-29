import type { ReactNode } from 'react'
import { usePathname } from '../../router'
import Logo from '../logo/logo'
import { AppLink } from '../navigation/AppLink'

interface AppNavigationItem {
  href: string
  label: string
}

interface AppFrameProps {
  actions?: ReactNode
  children: ReactNode
  description: string
  navigationItems?: AppNavigationItem[]
  title: string
}

const defaultNavigationItems: AppNavigationItem[] = [
  { href: '/app', label: 'Аккаунт' },
  { href: '/profile', label: 'Профиль' },
  { href: '/nutrition', label: 'Питание' },
  { href: '/workouts', label: 'Тренировки' },
]

export function AppFrame({
  actions,
  children,
  description,
  navigationItems = defaultNavigationItems,
  title,
}: AppFrameProps) {
  const pathname = usePathname()

  return (
    <div className="app-page">
      <header className="app-header">
        <Logo compact />
        {navigationItems.length > 0 ? (
          <nav aria-label="Основная навигация" className="app-nav">
            {navigationItems.map((item) => (
              <AppLink
                className={pathname === item.href ? 'app-nav__link is-active' : 'app-nav__link'}
                href={item.href}
                key={item.href}
              >
                {item.label}
              </AppLink>
            ))}
          </nav>
        ) : null}
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
