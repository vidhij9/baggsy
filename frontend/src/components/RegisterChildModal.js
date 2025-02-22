import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { motion } from 'framer-motion';
import { toast } from 'react-toastify';
import { ArchiveBoxIcon } from '@heroicons/react/24/solid';

function RegisterChildModal({ parent, closeModal, setError }) {
  const [qr, setQr] = useState('');
  const [currentCount, setCurrentCount] = useState(0);
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsLoading(true);
    try {
      const res = await axios.post('http://localhost:8080/api/register-child', {
        qrCode: qr,
        type: 'child',
        parentId: parent.id,
      });
      setCurrentCount(res.data.currentCount);
      setQr('');
      if (res.data.currentCount === res.data.capacity) {
        closeModal();
        toast.success('All child bags registered! Parent bag completed.', { position: 'top-center' });
      } else {
        toast.success('Child seed bag registered!', { position: 'top-center' });
      }
      setError(null);
    } catch (err) {
      if (err.response?.status === 401) {
        setError('Session expired. Please log in again.');
        toast.error('Session expired. Logging out...', { position: 'top-center' });
      } else {
        setError(err.response?.data?.error || 'Registration failed');
        toast.error(err.response?.data?.error || 'Registration failed', { position: 'top-center' });
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      className="fixed inset-0 bg-gray-600 bg-opacity-50 flex items-center justify-center z-50"
    >
      <div className="bg-white p-6 rounded-xl shadow-lg border-t-4 border-primary w-full max-w-md">
        <h2 className="text-2xl font-semibold text-accent mb-4 flex items-center">
          <ArchiveBoxIcon className="w-6 h-6 text-primary mr-2" />
          Register Child Bags for {parent.qrCode}
        </h2>
        <p className="text-accent mb-4">
          Current: {currentCount}/{parent.childCount}
        </p>
        <form onSubmit={handleSubmit} className="space-y-4">
          <input
            value={qr}
            onChange={(e) => setQr(e.target.value)}
            placeholder="Scan Child QR"
            className="w-full p-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary text-accent"
            required
          />
          <div className="w-full bg-gray-200 rounded-full h-3">
            <motion.div
              className="bg-primary h-3 rounded-full"
              initial={{ width: 0 }}
              animate={{ width: `${(currentCount / parent.childCount) * 100}%` }}
              transition={{ duration: 0.5 }}
            />
          </div>
          <button
            type="submit"
            disabled={isLoading}
            className={`w-full py-3 rounded-lg text-white flex items-center justify-center ${
              isLoading ? 'bg-gray-400' : 'bg-primary hover:bg-green-700'
            } transition duration-300`}
          >
            {isLoading ? (
              <svg className="animate-spin h-5 w-5 mr-2" viewBox="0 0 24 24">
                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z" />
              </svg>
            ) : (
              <ArchiveBoxIcon className="w-5 h-5 mr-2" />
            )}
            Register Child
          </button>
        </form>
      </div>
    </motion.div>
  );
}

export default RegisterChildModal;