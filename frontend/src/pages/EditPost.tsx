import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { posts } from '../services/api';
import { Post } from '../types';
import toast from 'react-hot-toast';

const EditPost: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [content, setContent] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(true);
  const [saving, setSaving] = useState<boolean>(false);
  const navigate = useNavigate();

  useEffect(() => {
    fetchPost();
  }, [id]);

  const fetchPost = async (): Promise<void> => {
    try {
      const response = await posts.getMyPosts();
      const post = response.data.find(p => p.ID === parseInt(id || '0'));
      if (post) {
        setContent(post.Content);
      } else {
        toast.error('Post not found');
        navigate('/posts/my-posts');
      }
    } catch (error) {
      console.error('Failed to fetch post:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>): Promise<void> => {
    e.preventDefault();
    
    if (!content.trim()) {
      toast.error('Content is required');
      return;
    }

    setSaving(true);
    try {
      await posts.edit({ 
        post_id: parseInt(id || '0'), 
        content 
      });
      toast.success('Post updated successfully');
      navigate('/posts/my-posts');
    } catch (error) {
      console.error('Failed to update post:', error);
    } finally {
      setSaving(false);
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
    <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="bg-white shadow sm:rounded-lg">
        <div className="px-4 py-5 sm:p-6">
          <h3 className="text-lg leading-6 font-medium text-gray-900">
            Edit Post
          </h3>

          <form onSubmit={handleSubmit} className="mt-5">
            <div className="mb-4">
              <label htmlFor="content" className="block text-sm font-medium text-gray-700">
                Post Content
              </label>
              <div className="mt-1">
                <textarea
                  id="content"
                  rows={6}
                  value={content}
                  onChange={(e) => setContent(e.target.value)}
                  className="shadow-sm focus:ring-blue-500 focus:border-blue-500 block w-full sm:text-sm border border-gray-300 rounded-md"
                  placeholder="Write your post here..."
                  disabled={saving}
                />
              </div>
            </div>

            <div className="flex justify-end space-x-3">
              <button
                type="button"
                onClick={() => navigate('/posts/my-posts')}
                className="bg-white py-2 px-4 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={saving}
                className="inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
              >
                {saving ? 'Saving...' : 'Save Changes'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default EditPost;
