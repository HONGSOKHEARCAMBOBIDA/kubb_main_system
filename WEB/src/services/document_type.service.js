import api from "./api";

export const getDocumentType = () => api.get('v1/document.type.view')