import React, { useState, useEffect } from 'react';
import { admin } from '../services/api';
import { Post } from '../types';
import toast from 'react-hot-toast';

const AdminDashboard: React.FC = () => {
  const [flaggedPosts, setFlaggedPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [processing, setProcessing] = useState<number | null>(null);

  useEffect(() => {
    fetchFlaggedPosts();
  }, []);

  const fetchFlaggedPosts = async (): Promise<void> => {
    try {
      const response = await admin.getFlaggedPosts();
      setFlaggedPosts(response.data);
    } catch (error) {
      console.error('Failed to fetch flagged posts:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleMarkSafe = async (postId: number): Promise<void> => {
    setProcessing(postId);
    try {
      await admin.markPostSafe({ post_id: postId });
      setFlaggedPosts(flaggedPosts.filter(p => p.ID !== postId));
      toast.success('Post marked as safe');
    } catch (error) {
      console.error('Failed to mark post as safe:', error);
    } finally {
      setProcessing(null);
    }
  };

  const handleDelete = async (postId: number): Promise<void> => {
    if (!window.confirm('Are you sure you want to delete this flagged post?')) {
      return;
    }

    setProcessing(postId);
    try {
      await admin.deletePost({ post_id: postId });
      setFlaggedPosts(flaggedPosts.filter(p => p.ID !== postId));
      toast.success('Post deleted successfully');
    } catch (error) {
      console.error('Failed to delete post:', error);
    } finally {
      setProcessing(null);
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Admin Dashboard</h1>

      <div className="bg-white shadow sm:rounded-lg mb-6">
        <div className="px-4 py-5 sm:p-6">
          <h2 className="text-lg font-medium text-gray-900 mb-4">Flagged Posts</h2>
          
          {flaggedPosts.length === 0 ? (
            <p className="text-gray-500 text-center py-4">No flagged posts to review</p>
          ) : (
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Post ID
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      User ID
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Content
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Toxicity Score
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Created At
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {flaggedPosts.map((post) => (
                    <tr key={post.ID}>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                        {post.ID}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                        {post.UserID}
                      </td>
                      <td className="px-6 py-4 text-sm text-gray-900 max-w-md">
                        <p className="truncate">{post.Content}</p>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                        <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-red-100 text-red-800">
                          {post.ToxicityScore}%
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {new Date(post.CreatedAt).toLocaleDateString()}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                        <button
                          onClick={() => handleMarkSafe(post.ID)}
                          disabled={processing === post.ID}
                          className="text-green-600 hover:text-green-900 mr-4 disabled:opacity-50"
                        >
                          Mark Safe
                        </button>
                        <button
                          onClick={() => handleDelete(post.ID)}
                          disabled={processing === post.ID}
                          className="text-red-600 hover:text-red-900 disabled:opacity-50"
                        >
                          Delete
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default AdminDashboard;
