import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authService } from '@/services/authService'
import { useToast } from '@/utils/useToast'

export const useAuthStore = defineStore('auth', () => {
  const user  = ref(null)
  const token = ref(localStorage.getItem('admin_token') || null)

  const isAuthenticated = computed(() => !!token.value)
  const role            = computed(() => user.value?.role || null)
  const isSuperAdmin    = computed(() => role.value === 'super_admin')
  const isAdmin         = computed(() => ['admin', 'super_admin'].includes(role.value))
  const fullName        = computed(() => user.value ? `${user.value.first_name} ${user.value.last_name}` : '')

  async function login(email, password) {
    const data = await authService.login(email, password)
    token.value = data.token
    user.value  = data.user
    localStorage.setItem('admin_token', data.token)
    localStorage.setItem('admin_user', JSON.stringify(data.user))
  }

  async function logout() {
    try {
      await authService.logout()
    } catch {
    } finally {
      token.value = null
      user.value  = null
      localStorage.removeItem('admin_token')
      localStorage.removeItem('admin_user')
    }
  }

  async function fetchMe() {
    const data = await authService.me()
    user.value = data
    localStorage.setItem('admin_user', JSON.stringify(data))
  }

  async function initAuth() {
    const savedToken = localStorage.getItem('admin_token')
    const savedUser  = localStorage.getItem('admin_user')

    if (!savedToken) return

    token.value = savedToken

    if (savedUser) {
      try {
        user.value = JSON.parse(savedUser)
      } catch {}
    }

    try {
      await fetchMe()
    } catch {
      token.value = null
      user.value  = null
      localStorage.removeItem('admin_token')
      localStorage.removeItem('admin_user')
    }
  }

  return {
    user,
    token,
    isAuthenticated,
    role,
    isSuperAdmin,
    isAdmin,
    fullName,
    login,
    logout,
    fetchMe,
    initAuth,
  }
})
