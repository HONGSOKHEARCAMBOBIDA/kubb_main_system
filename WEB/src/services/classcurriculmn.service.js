import api from './api'

export const CreateClassCurriculum = (data) => api.post('/v1/class_curriculmn.create', data)
export const GetClassCurriculum = (params) => api.get('/v1/class_curriculmn.view', { params })
export const GetClassCurriculumWithTeacherRate = (params) => api.get('/v1/class_curriculmn.view.with.teacher.rate', { params })