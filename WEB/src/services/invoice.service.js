import api from "./api";

export const invoicecreate = (data) => api.post('v1/invoice.create',data)