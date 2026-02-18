import React, { useEffect, useState } from 'react';
import { useAuth } from '../context/AuthContext';
import { posts } from '../services/api';
import { Post } from '../types';
import { Link } from 'react-router-dom';

const Dashboard: React.FC = () => {
  const { user } = useAuth();
  const [recentPosts, setRecentPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    fetchRecentPosts();
  }, []);

  const fetchRecentPosts = async (): Promise<void> => {
    try {
      const response = await posts.getMyPosts();
      setRecentPosts(response.data.slice(0, 5));
    } catch (error) {
      console.error('Failed to fetch posts:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="bg-white shadow overflow-hidden sm:rounded-lg mb-6">
        <div className="px-4 py-5 sm:px-6">
          <h3 className="text-lg leading-6 font-medium text-gray-900">
            Welcome, {user?.Email}!
          </h3>
          <p className="mt-1 max-w-2xl text-sm text-gray-500">
            Your role: {user?.Role}
          </p>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 mb-8">
        <div className="bg-white overflow-hidden shadow rounded-lg">
          <div className="px-4 py-5 sm:p-6">
            <dt className="text-sm font-medium text-gray-500 truncate">Total Posts</dt>
            <dd className="mt-1 text-3xl font-semibold text-gray-900">{recentPosts.length}</dd>
          </div>
        </div>

        <div className="bg-white overflow-hidden shadow rounded-lg">
          <div className="px-4 py-5 sm:p-6">
            <dt className="text-sm font-medium text-gray-500 truncate">Flagged Posts</dt>
            <dd className="mt-1 text-3xl font-semibold text-red-600">
              {recentPosts.filter(p => p.IsFlagged).length}
            </dd>
          </div>
        </div>

        <div className="bg-white overflow-hidden shadow rounded-lg">
          <div className="px-4 py-5 sm:p-6">
            <dt className="text-sm font-medium text-gray-500 truncate">Safe Posts</dt>
            <dd className="mt-1 text-3xl font-semibold text-green-600">
              {recentPosts.filter(p => !p.IsFlagged).length}
            </dd>
          </div>
        </div>
      </div>

      <div className="flex space-x-4 mb-6">
        <Link
          to="/posts/create"
          className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700"
        >
          Create New Post
        </Link>
        <Link
          to="/posts/my-posts"
          className="inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
        >
          View All Posts
        </Link>
      </div>

      <div className="bg-white shadow overflow-hidden sm:rounded-md">
        <div className="px-4 py-5 sm:px-6">
          <h3 className="text-lg leading-6 font-medium text-gray-900">Recent Posts</h3>
        </div>
        <ul className="divide-y divide-gray-200">
          {loading ? (
            <li className="px-4 py-4 text-center text-gray-500">Loading...</li>
          ) : recentPosts.length === 0 ? (
            <li className="px-4 py-4 text-center text-gray-500">No posts yet</li>
          ) : (
            recentPosts.map((post) => (
              <li key={post.ID} className="px-4 py-4">
                <div className="flex items-center justify-between">
                  <div className="flex-1">
                    <p className="text-sm text-gray-900">{post.Content}</p>
                    <p className="text-xs text-gray-500 mt-1">
                      Created: {new Date(post.CreatedAt).toLocaleDateString()}
                    </p>
                    <div className="mt-2">
                      <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                        post.IsFlagged 
                          ? 'bg-red-100 text-red-800' 
                          : 'bg-green-100 text-green-800'
                      }`}>
                        {post.IsFlagged ? 'Flagged' : 'Safe'} 
                        {post.ToxicityScore > 0 && ` (Score: ${post.ToxicityScore}%)`}
                      </span>
                    </div>
                  </div>
                  <Link
                    to={`/posts/edit/${post.ID}`}
                    className="ml-4 text-sm text-blue-600 hover:text-blue-900"
                  >
                    Edit
                  </Link>
                </div>
              </li>
            ))
          )}
        </ul>
      </div>
    </div>
  );
};

export default Dashboard;
