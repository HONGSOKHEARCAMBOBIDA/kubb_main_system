import api from "./api";

export const getTerm = (param) => api.get('/v1/term.view',{params})