import api from "./api";

export const getProvince = () => api.get('v1/province.view')
export const getDistrict = (id) => api.get(`v1/district.view/${id}`)
export const getCommunce = (id) => api.get(`v1/communce.view/${id}`)
export const getVillage = (id) => api.get(`v1/village.view/${id}`)