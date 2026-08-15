import api from "./api";

export const getAcademicStream = () => api.get('/v1/academic.stream.view')