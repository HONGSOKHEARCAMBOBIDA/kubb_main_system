import api from "./api";

export const getGenerations = (params) => api.get('/v1/generation.view', { params })
export const getGenerationByAcademic = (id) => api.get(`/v1/generation.view.academic/${id}`)
export const createGeneration = (data) => api.post('/v1/generation.create', data)
export const updateGeneration = (id, data) => api.put(`/v1/generation.update/${id}`, data)
export const toggleGeneration = (id) => api.put(`/v1/generation.toggle/${id}`)