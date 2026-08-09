import api from "./api"

export const getDepartment = (params) => api.get('/v1/Department.view',{params})

export const createDepartment = (data) => api.post('/v1/Department.create',data)

export const getDepartmentByFaculty = (id) => api.get(`/v1/Department.view.by.faculty/${id}`)

export const updateDepartment = (id,data) => api.put(`/v1/Department.update/${id}`,data)

export const toggleDepartment = (id) => api.put(`/v1/Department.toggle/${id}`)