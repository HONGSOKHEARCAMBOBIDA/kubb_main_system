import api from "./api";

export const getAcademicshift = (params) => api.get('/v1/AcademicShift.view', { params })
export const getAcademicshiftByAcademic = (id) => api.get(`/v1/AcademicShift.view.academic/${id}`)
export const createAcademicshift = (data) => api.post('/v1/AcademicShift.create', data)
export const updateAcademicshift = (id, data) => api.put(`/v1/AcademicShift.update/${id}`, data)
export const toggleAcademicshift = (id) => api.put(`/v1/AcademicShift.toggle/${id}`)