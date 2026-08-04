import api from "./api"

export const getAcademics = (params) => api.get('/v1/academic.view', { params })
export const createAcademic = (data) => api.post('/v1/academic.create', data)
export const updateAcademic = (id, data) => api.put(`/v1/academic.update/${id}`, data)
export const toggleAcademic = (id) => api.put(`/v1/academic.toggle/${id}`)