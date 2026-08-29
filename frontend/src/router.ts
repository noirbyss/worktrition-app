import { useSyncExternalStore } from 'react'

const NAVIGATION_EVENT = 'worktrition:navigate'

export function usePathname() {
  return useSyncExternalStore(subscribe, getSnapshot, () => '/')
}

export function navigate(to: string, options?: { replace?: boolean }) {
  const pathname = to.startsWith('/') ? to : `/${to}`

  if (window.location.pathname === pathname) {
    return
  }

  if (options?.replace) {
    window.history.replaceState(null, '', pathname)
  } else {
    window.history.pushState(null, '', pathname)
  }

  window.dispatchEvent(new Event(NAVIGATION_EVENT))
}

function subscribe(onStoreChange: () => void) {
  window.addEventListener('popstate', onStoreChange)
  window.addEventListener(NAVIGATION_EVENT, onStoreChange)

  return () => {
    window.removeEventListener('popstate', onStoreChange)
    window.removeEventListener(NAVIGATION_EVENT, onStoreChange)
  }
}

function getSnapshot() {
  return window.location.pathname || '/'
}
