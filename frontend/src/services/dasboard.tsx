import api from './api';

export const getItems = () => api.get('/items');
export const createItem = (item: any) => api.post('/items', item);

