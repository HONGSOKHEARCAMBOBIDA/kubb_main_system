import api from "./api";

export const EnrollmentCreate = (data) => api.post('/v1/enrollment.create',data)
export const EnrollmentUpdate = (uuid,data) => api.put(`/v1/enrollment.update/${uuid}`,data)
export const EnrollmentDelete = (uuid) => api.put(`/v1/enrollment.delete/${uuid}`)