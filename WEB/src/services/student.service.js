import api from "./api";

export const createStudent = (data) => api.post('v1/student.create',data)
export const getStudent = (params) => api.get('v1/student.view', { params })
export const updateStudent = (id,data) => api.put(`v1/student.update/${id}`,data)
export const getStudentCategory = () => api.get('/v1/student.category.view')