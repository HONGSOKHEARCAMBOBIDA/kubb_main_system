import api from "./api";

export const getAcademicSection = (params) => api.get('/v1/AcademicSection.view', { params })
export const getAcademicSectionByAcademic = (id) => api.get(`/v1/AcademicSection.view.academic/${id}`)
export const createAcademicSection = (data) => api.post('/v1/AcademicSection.create', data)
export const updateAcademicSection = (id, data) => api.put(`/v1/AcademicSection.update/${id}`, data)
export const toggleAcademicSection = (id) => api.put(`/v1/AcademicSection.toggle/${id}`)