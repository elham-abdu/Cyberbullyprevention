import axios from 'axios';
import { Post, LoginResponse } from '../types';
import toast from 'react-hot-toast';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    console.log('Making request to:', config.url, config.data); // Debug log
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => {
    console.log('Response received:', response.data); // Debug log
    return response;
  },
  (error) => {
    console.error('API Error:', error.response?.data || error.message); // Debug log
    
    if (error.code === 'ERR_NETWORK') {
      toast.error('Cannot connect to server. Make sure backend is running on port 8080');
    } else if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
      toast.error('Session expired. Please login again.');
    } else if (error.response?.status === 403) {
      toast.error('You do not have permission to perform this action');
    } else if (error.response?.data) {
      // Show the actual error message from backend
      toast.error(error.response.data);
    } else {
      toast.error('An error occurred. Please try again.');
    }
    return Promise.reject(error);
  }
);

interface LoginData {
  email: string;
  password: string;
}

interface RegisterData {
  email: string;
  password: string;
}

interface CreatePostData {
  content: string;
}

interface EditPostData {
  post_id: number;
  content: string;
}

interface DeletePostData {
  post_id: number;
}

// Auth endpoints
export const auth = {
  register: (data: RegisterData) => api.post('/register', data),
  login: (data: LoginData) => api.post<LoginResponse>('/login', data),
  me: () => api.get('/me'),
};

// Post endpoints
export const posts = {
  create: (data: CreatePostData) => api.post<Post>('/me/posts/create', data),
  getMyPosts: () => api.get<Post[]>('/me/posts'),
  edit: (data: EditPostData) => api.put<Post>('/me/posts/edit', data),
  delete: (data: DeletePostData) => api.delete('/me/posts/delete', { data }),
};

// Admin endpoints
export const admin = {
  getDashboard: () => api.get('/admin/dashboard'),
  getFlaggedPosts: () => api.get<Post[]>('/admin/flagged-posts'),
  markPostSafe: (data: { post_id: number }) => api.post('/admin/posts/mark-safe', data),
  deletePost: (data: { post_id: number }) => api.delete('/admin/posts/delete-flagged', { data }),
};

export default api;
