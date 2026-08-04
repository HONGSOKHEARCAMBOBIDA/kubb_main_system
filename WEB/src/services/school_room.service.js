import api from "./api";

export const getSchoolRooms = () => api.get('/v1/school.room.view')