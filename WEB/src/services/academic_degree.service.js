import api from "./api";

export const getAcademicDegree = (params) => api.get('/v1/AcademicDegree.view', { params })
export const getAcademicDegreeByAcademic = (id) => api.get(`/v1/AcademicDegree.view.by.academic/${id}`)
export const createAcademicDegree = (data) => api.post('/v1/AcademicDegree.create', data)
export const updateAcademicDegree = (id, data) => api.put(`/v1/AcademicDegree.update/${id}`, data)
export const toggleAcademicDegree = (id) => api.put(`/v1/AcademicDegree.toggle/${id}`)