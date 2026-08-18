import api from "./api";

export const getSchoolarshipGroup = (params) => api.get('/v1/SchoolarshipGroup.view', { params })
export const createSchoolarshipGroup = (data) => api.post('/v1/SchoolarshipGroup.create', data)
export const updateSchoolarshipGroup = (id, data) => api.put(`/v1/SchoolarshipGroup.update/${id}`, data)
export const toggleSchoolarshipGroup = (id) => api.put(`/v1/SchoolarshipGroup.toggle/${id}`)