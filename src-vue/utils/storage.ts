import { logger } from '@/utils/logger'
const PREFIX = 'tavern_'

export function getItem<T = any>(key: string, defaultValue?: T): T | null {
  try {
    const item = localStorage.getItem(PREFIX + key)
    if (item === null) return defaultValue ?? null
    return JSON.parse(item)
  } catch {
    return defaultValue ?? null
  }
}

export function setItem(key: string, value: any): void {
  try {
    localStorage.setItem(PREFIX + key, JSON.stringify(value))
  } catch (e) {
    logger.error('Failed to save to localStorage:', e)
  }
}

export function removeItem(key: string): void {
  localStorage.removeItem(PREFIX + key)
}

export function clear(): void {
  const keys = Object.keys(localStorage)
  keys.forEach(key => {
    if (key.startsWith(PREFIX)) {
      localStorage.removeItem(key)
    }
  })
}

export function getSessionItem<T = any>(key: string, defaultValue?: T): T | null {
  try {
    const item = sessionStorage.getItem(PREFIX + key)
    if (item === null) return defaultValue ?? null
    return JSON.parse(item)
  } catch {
    return defaultValue ?? null
  }
}

export function setSessionItem(key: string, value: any): void {
  try {
    sessionStorage.setItem(PREFIX + key, JSON.stringify(value))
  } catch (e) {
    logger.error('Failed to save to sessionStorage:', e)
  }
}

export function removeSessionItem(key: string): void {
  sessionStorage.removeItem(PREFIX + key)
}
