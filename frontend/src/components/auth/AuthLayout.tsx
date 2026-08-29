import type { ReactNode } from 'react'
import Logo from '../logo/logo'
import { AppLink } from '../navigation/AppLink'

interface AuthLayoutProps {
  children: ReactNode
  subtitle?: string
  switchHref: string
  switchLabel: string
  switchText: string
  title: string
}

export function AuthLayout({
  children,
  subtitle,
  switchHref,
  switchLabel,
  switchText,
  title,
}: AuthLayoutProps) {
  return (
    <div className="auth-page">
      <div className="auth-shell">
        <div className="auth-brand">
          <Logo />
        </div>
        <section className="auth-card">
          <h1 className="auth-title">{title}</h1>
          {subtitle ? <p className="auth-subtitle">{subtitle}</p> : null}
          {children}
          <p className="auth-switch">
            {switchText} <AppLink href={switchHref}>{switchLabel}</AppLink>
          </p>
        </section>
      </div>
    </div>
  )
}
