import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { posts } from '../services/api';
import { Post } from '../types';
import toast from 'react-hot-toast';

const CreatePost: React.FC = () => {
  const [content, setContent] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);
  const [toxicityResult, setToxicityResult] = useState<{ score: number; flagged: boolean } | null>(null);
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>): Promise<void> => {
    e.preventDefault();
    
    if (!content.trim()) {
      toast.error('Content is required');
      return;
    }

    setLoading(true);
    try {
      const response = await posts.create({ content });
      setToxicityResult({
        score: response.data.ToxicityScore,
        flagged: response.data.IsFlagged
      });
      toast.success('Post created successfully!');
      setTimeout(() => navigate('/posts/my-posts'), 2000);
    } catch (error) {
      console.error('Failed to create post:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="bg-white shadow sm:rounded-lg">
        <div className="px-4 py-5 sm:p-6">
          <h3 className="text-lg leading-6 font-medium text-gray-900">
            Create New Post
          </h3>
          <div className="mt-2 max-w-xl text-sm text-gray-500">
            <p>Write your post below. It will be automatically analyzed for toxicity.</p>
          </div>

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
                  disabled={loading}
                />
              </div>
            </div>

            {toxicityResult && (
              <div className={`mb-4 p-4 rounded-md ${
                toxicityResult.flagged ? 'bg-red-50' : 'bg-green-50'
              }`}>
                <p className={`text-sm ${
                  toxicityResult.flagged ? 'text-red-700' : 'text-green-700'
                }`}>
                  Toxicity Score: {toxicityResult.score}%
                </p>
                <p className={`text-sm font-medium ${
                  toxicityResult.flagged ? 'text-red-700' : 'text-green-700'
                }`}>
                  Status: {toxicityResult.flagged ? '⚠️ Flagged for review' : '✅ Safe'}
                </p>
              </div>
            )}

            <div className="flex justify-end space-x-3">
              <button
                type="button"
                onClick={() => navigate('/dashboard')}
                className="bg-white py-2 px-4 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={loading}
                className="inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
              >
                {loading ? 'Creating...' : 'Create Post'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default CreatePost;
