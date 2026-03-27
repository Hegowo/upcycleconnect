import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

const BASE = '/api/v1'

async function apiFetch(path, options = {}, token = null) {
  const headers = { 'Content-Type': 'application/json', Accept: 'application/json' }
  if (token) headers['Authorization'] = `Bearer ${token}`
  const res = await fetch(`${BASE}${path}`, { ...options, headers })
  const json = await res.json().catch(() => ({}))
  if (!res.ok) throw Object.assign(new Error(json.message || 'Erreur réseau'), { status: res.status, data: json })
  return json
}

export const useUserAuthStore = defineStore('userAuth', () => {
  const user  = ref(null)
  const token = ref(localStorage.getItem('user_token') || null)

  const isLoggedIn = computed(() => !!token.value)
  const fullName   = computed(() => user.value ? `${user.value.first_name} ${user.value.last_name}` : '')
  const initials   = computed(() => {
    if (!user.value) return ''
    return ((user.value.first_name?.[0] || '') + (user.value.last_name?.[0] || '')).toUpperCase()
  })

  async function register(payload) {
    await apiFetch('/auth/register', { method: 'POST', body: JSON.stringify(payload) })
  }

  async function login(email, password) {
    const data = await apiFetch('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    })
    token.value = data.token
    user.value  = data.user
    localStorage.setItem('user_token', data.token)
    localStorage.setItem('user_data', JSON.stringify(data.user))
  }

  async function logout() {
    try {
      if (token.value) {
        await apiFetch('/auth/logout', { method: 'POST' }, token.value)
      }
    } catch {}
    token.value = null
    user.value  = null
    localStorage.removeItem('user_token')
    localStorage.removeItem('user_data')
  }

  async function fetchMe() {
    const data = await apiFetch('/auth/me', {}, token.value)
    user.value = data
    localStorage.setItem('user_data', JSON.stringify(data))
  }

  async function init() {
    const savedToken = localStorage.getItem('user_token')
    const savedUser  = localStorage.getItem('user_data')
    if (!savedToken) return
    token.value = savedToken
    if (savedUser) {
      try { user.value = JSON.parse(savedUser) } catch {}
    }
    try {
      await fetchMe()
    } catch {
      token.value = null
      user.value  = null
      localStorage.removeItem('user_token')
      localStorage.removeItem('user_data')
    }
  }

  return { user, token, isLoggedIn, fullName, initials, register, login, logout, fetchMe, init }
})
