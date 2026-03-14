import api from './api'

export const userService = {
  async list(params = {}) {
    const { data } = await api.get('/users', { params })
    return data
  },

  async get(id) {
    const { data } = await api.get(`/users/${id}`)
    return data
  },

  async update(id, payload) {
    const { data } = await api.put(`/users/${id}`, payload)
    return data
  },

  async ban(id) {
    const { data } = await api.post(`/users/${id}/ban`)
    return data
  },

  async activate(id) {
    const { data } = await api.post(`/users/${id}/activate`)
    return data
  },

  async remove(id) {
    await api.delete(`/users/${id}`)
  },
}

export const providerService = {
  async list(params = {}) {
    const { data } = await api.get('/providers', { params })
    return data
  },

  async get(id) {
    const { data } = await api.get(`/providers/${id}`)
    return data
  },

  async updateStatus(id, status) {
    const { data } = await api.put(`/providers/${id}/status`, { status })
    return data
  },
}

export const adminUserService = {
  async list() {
    const { data } = await api.get('/admins')
    return data
  },

  async getById(id) {
    const { data } = await api.get(`/admins/${id}`)
    return data
  },

  async create(payload) {
    const { data } = await api.post('/admins', payload)
    return data
  },

  async update(id, payload) {
    const { data } = await api.put(`/admins/${id}`, payload)
    return data
  },

  async remove(id) {
    await api.delete(`/admins/${id}`)
  },
}
