import api from './api'

export const prestationService = {
  async list(params = {}) {
    const { data } = await api.get('/prestations', { params })
    return data
  },

  async get(id) {
    const { data } = await api.get(`/prestations/${id}`)
    return data
  },

  async create(payload) {
    const { data } = await api.post('/prestations', payload)
    return data
  },

  async update(id, payload) {
    const { data } = await api.put(`/prestations/${id}`, payload)
    return data
  },

  async updateStatus(id, status) {
    const { data } = await api.put(`/prestations/${id}/status`, { status })
    return data
  },

  async remove(id) {
    await api.delete(`/prestations/${id}`)
  },
}
