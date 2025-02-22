import React, { useState } from 'react';
import axios from 'axios';
import { motion } from 'framer-motion';
import { toast } from 'react-toastify';
import { ArchiveBoxIcon } from '@heroicons/react/24/solid';
import RegisterChildModal from './RegisterChildModal';

function RegisterParent({ setError }) {
  const [qr, setQr] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [showChildModal, setShowChildModal] = useState(false);
  const [parentData, setParentData] = useState(null);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsLoading(true);
    try {
      const res = await axios.post('http://localhost:8080/api/register-parent', {
        qrCode: qr,
        type: 'parent',
      });
      setParentData(res.data);
      setShowChildModal(true);
      setQr('');
      toast.success('Parent seed bag registered! Now add child bags.', { position: 'top-center' });
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
    <>
      <motion.div
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        className="bg-white p-6 rounded-xl shadow-lg border-t-4 border-primary"
      >
        <h2 className="text-2xl font-semibold text-accent mb-4 flex items-center">
          <ArchiveBoxIcon className="w-6 h-6 text-primary mr-2" />
          Register Parent Bag
        </h2>
        <form onSubmit={handleSubmit} className="space-y-4">
          <input
            value={qr}
            onChange={(e) => setQr(e.target.value)}
            placeholder="Scan Parent QR (e.g., P123-10)"
            className="w-full p-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary text-accent"
            required
          />
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
            Register Parent
          </button>
        </form>
      </motion.div>
      {showChildModal && parentData && (
        <RegisterChildModal
          parent={parentData}
          closeModal={() => {
            setShowChildModal(false);
            setParentData(null);
          }}
          setError={setError}
        />
      )}
    </>
  );
}

export default RegisterParent;