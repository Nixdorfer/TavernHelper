import { ref, computed } from 'vue'
import { messages, localeNames } from './messages'
function getNestedValue(obj: Record<string, any>, path: string): string | undefined {
  return path.split('.').reduce((acc, part) => acc && acc[part], obj)
}
function replacePlaceholders(str: string, params?: Record<string, string | number>): string {
  if (!params) return str
  return str.replace(/\{(\w+)\}/g, (match, key) =>
    params[key] !== undefined ? String(params[key]) : match
  )
}
const currentLocale = ref('zh-CN')
export function useI18n() {
  const locale = computed({
    get: () => currentLocale.value,
    set: (value: string) => {
      if (messages[value]) {
        currentLocale.value = value
      }
    }
  })
  function t(key: string, params?: Record<string, string | number>): string {
    const message = getNestedValue(messages[currentLocale.value], key) ||
                    getNestedValue(messages['zh-CN'], key) ||
                    key
    return replacePlaceholders(message, params)
  }
  return {
    locale,
    t,
    availableLocales: Object.keys(messages),
    localeNames
  }
}
export function setLocale(locale: string) {
  if (messages[locale]) {
    currentLocale.value = locale
  }
}
export function getLocale(): string {
  return currentLocale.value
}
export { messages, localeNames }
