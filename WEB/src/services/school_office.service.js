import api from "./api";

export const getSchoolOffices = () => api.get('/v1/school.office.view')