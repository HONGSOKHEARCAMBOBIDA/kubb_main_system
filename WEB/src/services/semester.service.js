import api from "./api";

export const getSemester = (params) => api.get('/v1/Semester.view', { params })
export const getSemesterByAcademic = (id) => api.get(`/v1/Semester.view.academic/${id}`)
export const createSemester = (data) => api.post('/v1/Semester.create', data)
export const updateSemester = (id, data) => api.put(`/v1/Semester.update/${id}`, data)
export const toggleSemester = (id) => api.put(`/v1/Semester.toggle/${id}`)