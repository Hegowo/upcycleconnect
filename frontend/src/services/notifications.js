import api from './api'
import { userApi } from './publicApi'

export const adminNotifications = {
  async list(limit = 30) {
    const { data } = await api.get('/notifications', { params: { limit } })
    return data
  },
  async unreadCount() {
    const { data } = await api.get('/notifications/unread-count')
    return data.unread_count || 0
  },
  async markRead(id) {
    await api.post(`/notifications/${id}/read`)
  },
  async markAllRead() {
    await api.post('/notifications/read-all')
  },
}

export const userNotifications = {
  async list(limit = 30) {
    return userApi(`/notifications?limit=${limit}`)
  },
  async unreadCount() {
    const data = await userApi('/notifications/unread-count')
    return data.unread_count || 0
  },
  async markRead(id) {
    await userApi(`/notifications/${id}/read`, { method: 'POST' })
  },
  async markAllRead() {
    await userApi('/notifications/read-all', { method: 'POST' })
  },
  async registerPushToken(playerId) {
    await userApi('/notifications/push-token', {
      method: 'POST',
      body: JSON.stringify({ player_id: playerId }),
    })
  },
}
