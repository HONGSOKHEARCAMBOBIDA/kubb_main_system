import api from "./api";
export const getAdmission = (params) => api.get('/v1/admission.view', { params })
export const updateAdmission = (uuid,data) => api.put(`/v1/admission.update/${uuid}`,data)