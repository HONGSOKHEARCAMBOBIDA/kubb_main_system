import api from "./api";

export const getFeediscountGroup = (params) => api.get('/v1/FeediscountGroup.view', { params })
export const createFeediscountGroup = (data) => api.post('/v1/FeediscountGroup.create', data)
export const updateFeediscountGroup = (id, data) => api.put(`/v1/FeediscountGroup.update/${id}`, data)
export const toggleFeediscountGroup = (id) => api.put(`/v1/FeediscountGroup.toggle/${id}`)