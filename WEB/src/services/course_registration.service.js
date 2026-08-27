import api from './api'

export const createcourseregistration = (data) => api.post('/v1/course.registration.create',data)
export const getcourseregistration = (params) => api.get('/v1/course.registration.view',{params})