import api from "./api";
export const getAdmission = (params) => api.get('/v1/admission.view', { params })
