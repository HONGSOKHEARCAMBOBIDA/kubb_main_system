import api from './api'

export const getTerm = (params) => api.get('/v1/term.view', { params })
export const createTerm = (data) => api.post('/v1/term.create', data)
export const updateTerm = (id, data) => api.put(`/v1/term.update/${id}`, data)
export const toggleTerm = (id) => api.put(`/v1/term.toggle/${id}`)
