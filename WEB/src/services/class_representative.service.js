import api from "./api";

export const CreateClassRepresentative = (data) => api.post('/v1/class.representative.create',data)