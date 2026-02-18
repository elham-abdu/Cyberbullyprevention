import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { Toaster } from 'react-hot-toast';
import { AuthProvider } from './context/AuthContext';
import Navbar from './components/Navbar';
import ProtectedRoute from './components/ProtectedRoute';
import AdminRoute from './components/AdminRoute';
import Login from './pages/Login';
import Register from './pages/Register';
import Dashboard from './pages/Dashboard';
import CreatePost from './pages/CreatePost';
import MyPosts from './pages/MyPosts';
import EditPost from './pages/EditPost';
import AdminDashboard from './pages/AdminDashboard';

const App: React.FC = () => {
  return (
    <Router>
      <AuthProvider>
        <div className="min-h-screen bg-gray-50">
          <Navbar />
          <Toaster position="top-right" />
          <Routes>
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />
            <Route path="/" element={<Navigate to="/dashboard" />} />
            
            <Route path="/dashboard" element={
              <ProtectedRoute>
                <Dashboard />
              </ProtectedRoute>
            } />
            
            <Route path="/posts/create" element={
              <ProtectedRoute>
                <CreatePost />
              </ProtectedRoute>
            } />
            
            <Route path="/posts/my-posts" element={
              <ProtectedRoute>
                <MyPosts />
              </ProtectedRoute>
            } />
            
            <Route path="/posts/edit/:id" element={
              <ProtectedRoute>
                <EditPost />
              </ProtectedRoute>
            } />
            
            <Route path="/admin" element={
              <AdminRoute>
                <AdminDashboard />
              </AdminRoute>
            } />
          </Routes>
        </div>
      </AuthProvider>
    </Router>
  );
};

export default App;
