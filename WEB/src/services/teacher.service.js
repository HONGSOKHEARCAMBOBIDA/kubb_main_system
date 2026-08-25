import api from "./api"
  
export const getTeacher = (params) => api.get('/v1/teacher.view',{params})
export const getTeacherFilter = (params) => api.get('/v1/teacher.view.filter',{params})
export const createTeacher = (data) => api.post('/v1/teacher.create',data)
export const updateTeacher = (uuid,data) => api.put(`/v1/teacher.update/${uuid}`,data)
export const toggleTeacher = (uuid) => api.put(`/v1/teacher.toggle/${uuid}`)