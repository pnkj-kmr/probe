import api from "./base";

// export const createItem = (item: any) => api.post('/items', item);

export const getSystemData = () => api.get("/api/dashboard/system");
export const getResourceData = () => api.get("/api/dashboard/resource");
export const getPollingData = () => api.get("/api/dashboard/poll");

export const getSearchModels = () => api.get("/api/dashboard/model");
export const getSearchResult = (payload: any) =>
  api.get("/api/dashboard/data", payload);

export const getInfo = () => api.get("/api/profile/agent/info");
export const getEnvVariables = () => api.get("/api/profile/env");
export const getProcessStatus = () => api.get("/api/profile/process/status");
