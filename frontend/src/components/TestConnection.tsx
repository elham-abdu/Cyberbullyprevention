import React, { useEffect, useState } from 'react';

const TestConnection: React.FC = () => {
  const [status, setStatus] = useState<string>('Testing...');

  useEffect(() => {
    const testConnection = async () => {
      try {
        const response = await fetch('http://localhost:8080/login', {
          method: 'OPTIONS'
        });
        setStatus('✅ Backend is reachable!');
      } catch (error) {
        setStatus('❌ Cannot reach backend. Make sure it\'s running on port 8080');
      }
    };
    testConnection();
  }, []);

  return (
    <div className="p-4 bg-gray-100 rounded">
      <h2>Connection Test: {status}</h2>
    </div>
  );
};

export default TestConnection;
