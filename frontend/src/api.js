import axios from 'axios';

const API_URL = '/api';

export const movieApi = {
  getAll: () => axios.get(`${API_URL}/movies`),
  getById: (id) => axios.get(`${API_URL}/movies/${id}`),
  getMostLookedUp: () => axios.get(`${API_URL}/movies/most-looked-up`),
  create: (movie) => axios.post(`${API_URL}/movies`, movie),
  update: (id, movie) => axios.put(`${API_URL}/movies/${id}`, movie),
  delete: (id) => axios.delete(`${API_URL}/movies/${id}`),
};

export const bookApi = {
  getAll: () => axios.get(`${API_URL}/books`),
  getById: (id) => axios.get(`${API_URL}/books/${id}`),
  getMostLookedUp: () => axios.get(`${API_URL}/books/most-looked-up`),
  create: (book) => axios.post(`${API_URL}/books`, book),
  update: (id, book) => axios.put(`${API_URL}/books/${id}`, book),
  delete: (id) => axios.delete(`${API_URL}/books/${id}`),
};
