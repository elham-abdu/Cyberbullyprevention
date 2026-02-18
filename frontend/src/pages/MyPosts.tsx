import React, { useState, useEffect } from 'react';
import { posts } from '../services/api';
import { Post } from '../types';
import { Link } from 'react-router-dom';
import toast from 'react-hot-toast';

const MyPosts: React.FC = () => {
  const [postsList, setPostsList] = useState<Post[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [deleting, setDeleting] = useState<number | null>(null);

  useEffect(() => {
    fetchPosts();
  }, []);

  const fetchPosts = async (): Promise<void> => {
    try {
      const response = await posts.getMyPosts();
      setPostsList(response.data);
    } catch (error) {
      console.error('Failed to fetch posts:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (postId: number): Promise<void> => {
    if (!window.confirm('Are you sure you want to delete this post?')) {
      return;
    }

    setDeleting(postId);
    try {
      await posts.delete({ post_id: postId });
      setPostsList(postsList.filter(p => p.ID !== postId));
      toast.success('Post deleted successfully');
    } catch (error) {
      console.error('Failed to delete post:', error);
    } finally {
      setDeleting(null);
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
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-900">My Posts</h1>
        <Link
          to="/posts/create"
          className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700"
        >
          Create New Post
        </Link>
      </div>

      {postsList.length === 0 ? (
        <div className="bg-white shadow sm:rounded-lg">
          <div className="px-4 py-5 sm:p-6 text-center">
            <p className="text-gray-500">You haven't created any posts yet.</p>
          </div>
        </div>
      ) : (
        <div className="bg-white shadow overflow-hidden sm:rounded-md">
          <ul className="divide-y divide-gray-200">
            {postsList.map((post) => (
              <li key={post.ID} className="px-4 py-4">
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <p className="text-sm text-gray-900 whitespace-pre-wrap">
                      {post.Content}
                    </p>
                    <div className="mt-2 flex items-center space-x-4">
                      <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                        post.IsFlagged 
                          ? 'bg-red-100 text-red-800' 
                          : 'bg-green-100 text-green-800'
                      }`}>
                        {post.IsFlagged ? 'Flagged' : 'Safe'} 
                        {post.ToxicityScore > 0 && ` (Score: ${post.ToxicityScore}%)`}
                      </span>
                      <span className="text-xs text-gray-500">
                        {new Date(post.CreatedAt).toLocaleString()}
                      </span>
                    </div>
                  </div>
                  <div className="ml-4 flex items-center space-x-2">
                    <Link
                      to={`/posts/edit/${post.ID}`}
                      className="text-sm text-blue-600 hover:text-blue-900"
                    >
                      Edit
                    </Link>
                    <button
                      onClick={() => handleDelete(post.ID)}
                      disabled={deleting === post.ID}
                      className="text-sm text-red-600 hover:text-red-900 disabled:opacity-50"
                    >
                      {deleting === post.ID ? 'Deleting...' : 'Delete'}
                    </button>
                  </div>
                </div>
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
};

export default MyPosts;
