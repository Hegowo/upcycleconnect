import api from './api'

export const eventService = {
  async list(params = {}) {
    const { data } = await api.get('/events', { params })
    return data
  },

  async get(id) {
    const { data } = await api.get(`/events/${id}`)
    return data
  },

  async create(payload) {
    const { data } = await api.post('/events', payload)
    return data
  },

  async update(id, payload) {
    const { data } = await api.put(`/events/${id}`, payload)
    return data
  },

  async updateStatus(id, status) {
    const { data } = await api.put(`/events/${id}/status`, { status })
    return data
  },

  async remove(id) {
    await api.delete(`/events/${id}`)
  },
}
