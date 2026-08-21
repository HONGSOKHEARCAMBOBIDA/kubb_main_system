import api from "./api";

export const EnrollmentCreate = (data) => api.post('/v1/enrollment.create',data)