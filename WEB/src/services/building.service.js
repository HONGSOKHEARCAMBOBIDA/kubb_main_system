import api from "./api";

export const getBuildings = () => api.get('/v1/building.view')