import type { AnchorHTMLAttributes, MouseEvent } from 'react'
import { navigate } from '../../router'

interface AppLinkProps extends Omit<AnchorHTMLAttributes<HTMLAnchorElement>, 'href'> {
  href: string
  replace?: boolean
}

export function AppLink({
  href,
  onClick,
  replace = false,
  target,
  ...props
}: AppLinkProps) {
  const handleClick = (event: MouseEvent<HTMLAnchorElement>) => {
    onClick?.(event)

    if (
      event.defaultPrevented ||
      event.button !== 0 ||
      event.metaKey ||
      event.ctrlKey ||
      event.shiftKey ||
      event.altKey ||
      target === '_blank' ||
      !href.startsWith('/')
    ) {
      return
    }

    event.preventDefault()
    navigate(href, { replace })
  }

  return <a {...props} href={href} onClick={handleClick} target={target} />
}
