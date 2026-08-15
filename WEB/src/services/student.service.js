import api from "./api";

export const createStudent = (data) => api.post('v1/student.create',data)