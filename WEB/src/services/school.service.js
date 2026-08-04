import api from "./api";

export const getSchools = () => api.get('/v1/school.view')