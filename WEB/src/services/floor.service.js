import api from "./api";

export const getFloors = () => api.get('/v1/floor.view')