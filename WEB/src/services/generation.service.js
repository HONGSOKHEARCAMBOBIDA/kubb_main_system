import api from "./api";

export const getGenerations = () => api.get('/v1/generation.view')