import api from "./api"

export const getFaculty = (params) => api.get('/v1/faculty.view',{params})

export const createFaculty = (data) => api.post('/v1/faculty.create',data)

export const getFacultyByProgrammes = (id) => api.get(`/v1/faculty.view.by.programmes/${id}`)

export const updateFaculty = (id,data) => api.put(`/v1/faculty.update/${id}`,data)

export const toggleFaculty = (id) => api.put(`/v1/faculty.toggle/${id}`)