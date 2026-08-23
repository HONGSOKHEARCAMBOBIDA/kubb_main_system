import api from './api'

export const CreateClassOffering = (data) => api.post('/v1/class.offering.create', data)
