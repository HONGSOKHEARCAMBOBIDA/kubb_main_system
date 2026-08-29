import api from "./api"

export const getSubject = (params) => api.get('/v1/Subject.view',{params})

export const createSubject = (data) => api.post('/v1/Subject.create',data)

export const getSubjectByMajor = (id) => api.get(`/v1/Subject.view.by.major/${id}`)

export const updateSubject = (id,data) => api.put(`/v1/Subject.update/${id}`,data)

export const toggleSubject = (id) => api.put(`/v1/Subject.toggle/${id}`)

export const creategradecomponent = (data) => api.post('/v1/grade.component.create',data)