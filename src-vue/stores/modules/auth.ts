import { defineStore } from 'pinia'
import { api } from '@/api'
import { logger } from '@/utils/logger'
export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: '',
    userId: '',
    userName: '',
    points: 0
  }),
  getters: {
    isLoggedIn: (state) => !!state.token
  },
  actions: {
    setAuth(token: string, userId: string, userName: string) {
      this.token = token
      this.userId = userId
      this.userName = userName
    },
    setPoints(points: number) {
      this.points = points
    },
    async login(email: string, password: string, rememberMe: boolean) {
      const result = await api.auth.login(email, password, rememberMe)
      const token = typeof result?.data === 'string' ? result.data : (result?.data?.access_token || result?.access_token)
      if (!token) {
        throw new Error(result?.message || result?.data?.message || '登录失败')
      }
      this.token = token
      const profile = await api.auth.getProfile(token)
      if (profile?.data) {
        this.userName = profile.data.name || profile.data.email
        this.userId = profile.data.id
      }
      const points = await api.auth.getPoints(token, this.userId)
      if (points?.data) {
        this.points = points.data.points || 0
      }
      localStorage.setItem('auth', JSON.stringify({ token: this.token, userId: this.userId, userName: this.userName, points: this.points }))
    },
    async logout() {
      try {
        if (this.token) {
          await api.auth.logout(this.token)
        }
      } catch (e) {
        logger.error('登出接口调用失败:', e)
      }
      this.token = ''
      this.userId = ''
      this.userName = ''
      this.points = 0
      localStorage.removeItem('auth')
    },
    restoreFromStorage() {
      const saved = localStorage.getItem('auth')
      if (saved) {
        try {
          const data = JSON.parse(saved)
          this.token = data.token || ''
          this.userId = data.userId || ''
          this.userName = data.userName || ''
          this.points = data.points || 0
        } catch (e) {
          logger.error('恢复auth失败:', e)
        }
      }
    }
  }
})
