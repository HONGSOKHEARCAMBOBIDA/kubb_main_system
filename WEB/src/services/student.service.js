import api from "./api";

export const createStudent = (data) => api.post('v1/student.create',data)
export const getStudent = (params) => api.get('v1/student.view', { params })