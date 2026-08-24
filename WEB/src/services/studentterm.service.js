import api from "./api";

export const studentTermCreate = (data) => api.post('/v1/student.term.create',data)
export const studentTermGet = (params) => api.get('/v1/student.term.view',{params})