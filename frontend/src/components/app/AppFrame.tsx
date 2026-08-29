import { useState, type ReactNode } from 'react'
import type { CurrentUser } from '../../api'
import { useAuth } from '../../auth/useAuth'
import { navigate, usePathname } from '../../router'
import worktritionLogo from '../../assets/worktrition-logo.png'
import { AppLink } from '../navigation/AppLink'

interface AppNavigationItem {
  href: string
  label: string
}

interface AppSidebarProfile {
  badge: string
  meta: string
  name: string
}

interface AppFrameProps {
  actions?: ReactNode
  children: ReactNode
  currentUser?: CurrentUser | null
  description: string
  eyebrow?: string
  isCurrentUserLoading?: boolean
  navigationItems?: AppNavigationItem[]
  sidebarProfile?: AppSidebarProfile
  title: string
}

const defaultNavigationItems: AppNavigationItem[] = [
  { href: '/app', label: 'Персонаж' },
  { href: '/workouts', label: 'Тренировки' },
  { href: '/nutrition', label: 'Питание' },
  { href: '/stats', label: 'Статистика' },
  { href: '/profile', label: 'Профиль' },
]

export function AppFrame({
  actions,
  children,
  currentUser,
  description,
  eyebrow,
  isCurrentUserLoading = false,
  navigationItems = defaultNavigationItems,
  sidebarProfile,
  title,
}: AppFrameProps) {
  const pathname = normalizePathname(usePathname())
  const { logout, session } = useAuth()
  const [isLoggingOut, setIsLoggingOut] = useState(false)
  const sidebarName =
    sidebarProfile?.name ?? currentUser?.name ?? (isCurrentUserLoading ? 'Загружаем...' : 'Пользователь')
  const sidebarEmail =
    sidebarProfile?.meta ??
    currentUser?.email ??
    (session?.userId ? `ID ${session.userId.slice(0, 8)}` : 'Авторизованный аккаунт')
  const avatarBadge = sidebarProfile?.badge ?? getInitials(currentUser?.name)

  const handleLogout = async () => {
    try {
      setIsLoggingOut(true)
      await logout()
      navigate('/login', { replace: true })
    } finally {
      setIsLoggingOut(false)
    }
  }

  return (
    <div className="shell">
      <aside className="nav">
        <AppLink className="logo" href="/app">
          <img alt="Worktrition" src={worktritionLogo} />
        </AppLink>

        {navigationItems.length > 0 ? (
          <nav aria-label="Основная навигация">
            <ul className="nav-list">
              {navigationItems.map((item, index) => (
                <li key={item.href}>
                  <AppLink
                    className={pathname === item.href ? 'active' : undefined}
                    href={item.href}
                  >
                    <span className="ic">{renderNavIcon(item.href)}</span>
                    <span className="num">{String(index + 1).padStart(2, '0')}</span>
                    {item.label}
                  </AppLink>
                </li>
              ))}
            </ul>
          </nav>
        ) : null}

        <div className="nav-foot">
          <div className="avatar-mini">
            <svg className="ring" viewBox="0 0 40 40">
              <circle cx="20" cy="20" fill="none" r="17" stroke="rgba(255,255,255,.1)" strokeWidth="3" />
              <circle
                cx="20"
                cy="20"
                fill="none"
                r="17"
                stroke="url(#sidebarGradient)"
                strokeDasharray="106.8"
                strokeDashoffset="19.2"
                strokeLinecap="round"
                strokeWidth="3"
                transform="rotate(-90 20 20)"
              />
              <defs>
                <linearGradient id="sidebarGradient" x1="0" x2="1" y1="0" y2="1">
                  <stop offset="0" stopColor="#FFD23F" />
                  <stop offset="1" stopColor="#FF7A1A" />
                </linearGradient>
              </defs>
            </svg>
            <div className="core">{avatarBadge}</div>
          </div>

          <div className="nav-foot__copy">
            <div className="nav-foot-name">{sidebarName}</div>
            <div className="nav-foot-level">{sidebarEmail}</div>
          </div>
        </div>

        <button
          className="sidebar-logout"
          disabled={isLoggingOut}
          onClick={() => {
            void handleLogout()
          }}
          type="button"
        >
          {isLoggingOut ? 'ВЫХОД...' : 'ВЫЙТИ'}
        </button>
      </aside>

      <main className="shell-main">
        <section className="page-head">
          <div className="page-head__copy">
            <div className="page-eyebrow">{eyebrow ?? 'WORKTRITION'}</div>
            <h1 className="page-title">{title}</h1>
            <p className="page-desc">{description}</p>
          </div>
          {actions ? <div className="page-head__actions">{actions}</div> : null}
        </section>
        {children}
      </main>
    </div>
  )
}

function normalizePathname(pathname: string) {
  if (pathname === '/') {
    return pathname
  }

  return pathname.replace(/\/+$/, '')
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

function renderNavIcon(href: string) {
  switch (href) {
    case '/app':
      return (
        <svg fill="none" viewBox="0 0 24 24">
          <path d="M4 12.5 12 4l8 8.5" stroke="currentColor" strokeLinejoin="round" strokeWidth="1.5" />
          <path d="M6.5 10.5V20h11V10.5" stroke="currentColor" strokeWidth="1.5" />
        </svg>
      )
    case '/workouts':
      return (
        <svg fill="none" viewBox="0 0 24 24">
          <rect height="8" rx="1" stroke="currentColor" strokeWidth="1.4" width="4" x="1.5" y="8" />
          <rect height="8" rx="1" stroke="currentColor" strokeWidth="1.4" width="4" x="18.5" y="8" />
          <line stroke="currentColor" strokeWidth="1.6" x1="6.5" x2="17.5" y1="12" y2="12" />
          <line stroke="currentColor" strokeWidth="1.6" x1="4.5" x2="4.5" y1="9.5" y2="14.5" />
          <line stroke="currentColor" strokeWidth="1.6" x1="19.5" x2="19.5" y1="9.5" y2="14.5" />
        </svg>
      )
    case '/nutrition':
      return (
        <svg fill="none" viewBox="0 0 24 24">
          <path d="M4 12h16a8 8 0 0 1-16 0Z" stroke="currentColor" strokeLinejoin="round" strokeWidth="1.4" />
          <line stroke="currentColor" strokeLinecap="round" strokeWidth="1.4" x1="12" x2="12" y1="3" y2="6.5" />
        </svg>
      )
    case '/stats':
      return (
        <svg fill="none" viewBox="0 0 24 24">
          <rect height="7" stroke="currentColor" strokeWidth="1.4" width="4" x="4" y="14" />
          <rect height="12" stroke="currentColor" strokeWidth="1.4" width="4" x="10" y="9" />
          <rect height="17" stroke="currentColor" strokeWidth="1.4" width="4" x="16" y="4" />
        </svg>
      )
    case '/profile':
      return (
        <svg fill="none" viewBox="0 0 24 24">
          <circle cx="12" cy="8" r="3.5" stroke="currentColor" strokeWidth="1.5" />
          <path d="M5 21c.8-4.1 3.2-6 7-6s6.2 1.9 7 6" stroke="currentColor" strokeLinecap="round" strokeWidth="1.5" />
        </svg>
      )
    default:
      return (
        <svg fill="none" viewBox="0 0 24 24">
          <circle cx="12" cy="12" r="8" stroke="currentColor" strokeWidth="1.5" />
        </svg>
      )
  }
}
