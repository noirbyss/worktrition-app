import type { ReactNode } from 'react'

interface InlineMessageProps {
  children: ReactNode
  tone?: 'error' | 'info' | 'success'
}

export function InlineMessage({
  children,
  tone = 'error',
}: InlineMessageProps) {
  return <p className={`inline-message inline-message--${tone}`}>{children}</p>
}
