import api from "./api"

export const getMajor = (params) => api.get('/v1/Major.view',{params})

export const createMajor = (data) => api.post('/v1/Major.create',data)

export const getMajorByDepartment = (id) => api.get(`/v1/Major.view.by.department/${id}`)

export const updateMajor = (id,data) => api.put(`/v1/Major.update/${id}`,data)

export const toggleMajor = (id) => api.put(`/v1/Major.toggle/${id}`)