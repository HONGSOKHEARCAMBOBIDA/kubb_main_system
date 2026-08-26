import api from './api'

export const ScheduleCreate = (data) => api.post('/v1/schedule.create', data)
