import api from "./api"

export const getMajorTerm = (params) => api.get('/v1/major.term.view',{params})

export const createMajorTerm = (data) => api.post('/v1/major.term.create',data)

export const updateMajorTerm = (id,data) => api.put(`/v1/major.term.update/${id}`,data)