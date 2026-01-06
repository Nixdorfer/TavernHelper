import { api } from '@/api'

let enabled = false

export const initLogger = (isEnabled: boolean) => {
  enabled = isEnabled
}

const formatArgs = (args: any[]): string => {
  return args.map(arg => {
    if (typeof arg === 'object') {
      try {
        return JSON.stringify(arg)
      } catch {
        return String(arg)
      }
    }
    return String(arg)
  }).join(' ')
}

const writeLog = (level: string, args: any[]) => {
  if (!enabled) return
  api.system.writeVueLog(level, formatArgs(args)).catch(() => {})
}

export const logger = {
  log: (...args: any[]) => {
    if (!enabled) return
    console.log(...args)
    writeLog('LOG', args)
  },
  info: (...args: any[]) => {
    if (!enabled) return
    console.info(...args)
    writeLog('INFO', args)
  },
  warn: (...args: any[]) => {
    if (!enabled) return
    console.warn(...args)
    writeLog('WARN', args)
  },
  error: (...args: any[]) => {
    if (!enabled) return
    console.error(...args)
    writeLog('ERROR', args)
  },
  debug: (...args: any[]) => {
    if (!enabled) return
    console.debug(...args)
    writeLog('DEBUG', args)
  },
}

export default logger
