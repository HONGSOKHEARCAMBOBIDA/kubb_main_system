import api from './api'

export const CreateClassCurriculum = (data) => api.post('/v1/class_curriculmn.create', data)
