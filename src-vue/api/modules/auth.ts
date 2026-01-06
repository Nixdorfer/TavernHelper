import { api } from '../index'
import type { UserProfile } from '@/types'
export const authApi = {
  async login(email: string, password: string, rememberMe = false) {
    const result = await api.auth.login(email, password, rememberMe)
    const token = typeof result?.data === 'string'
      ? result.data
      : (result?.data?.access_token || result?.access_token)
    if (!token) throw new Error(result?.message || '登录失败')
    return { token, raw: result }
  },
  logout(token: string) {
    return api.auth.logout(token)
  },
  getProfile(token: string): Promise<UserProfile> {
    return api.auth.getProfile(token)
  },
  getPoints(token: string, userId: string) {
    return api.auth.getPoints(token, userId)
  }
}
