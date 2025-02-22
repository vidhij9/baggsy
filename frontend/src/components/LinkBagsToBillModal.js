import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { motion } from 'framer-motion';
import { toast } from 'react-toastify';
import { DocumentTextIcon } from '@heroicons/react/24/solid';

function LinkBagsToBillModal({ setError, closeModal, token }) {
  const [billID, setBillID] = useState('');
  const [capacity, setCapacity] = useState(1);
  const [parentIDs, setParentIDs] = useState([]);
  const [unlinkedParents, setUnlinkedParents] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [fetchLoading, setFetchLoading] = useState(true);

  useEffect(() => {
    if (token) {
      fetchUnlinkedParents();
    }
  }, [token]);

  const fetchUnlinkedParents = async () => {
    setFetchLoading(true);
    try {
      const res = await axios.get('http://localhost:8080/api/unlinked-parents', {
        headers: { Authorization: `Bearer ${token}` },
      });
      console.log("Unlinked parents response:", res.data);
      setUnlinkedParents(Array.isArray(res.data) ? res.data : []);
      setError(null);
    } catch (err) {
      const errorMsg = err.response?.data?.error || 'Failed to fetch unlinked parents';
      console.error("Fetch unlinked parents error:", err);
      setUnlinkedParents([]);
      setError(errorMsg);
      toast.error(errorMsg, { position: 'top-center' });
    } finally {
      setFetchLoading(false);
    }
  };

  const handleAddParent = (id) => {
    if (parentIDs.length < capacity && !parentIDs.includes(id)) {
      setParentIDs([...parentIDs, id]);
    }
  };

  const handleRemoveParent = (id) => {
    setParentIDs(parentIDs.filter((pid) => pid !== id));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!billID.trim()) {
      setError('Bill ID is required');
      toast.error('Bill ID is required', { position: 'top-center' });
      return;
    }
    setIsLoading(true);
    try {
      await axios.post(
        'http://localhost:8080/api/link-bags-to-bill',
        { billID, parentIDs, capacity },
        { headers: { Authorization: `Bearer ${token}` } }
      );
      setBillID('');
      setCapacity(1);
      setParentIDs([]);
      fetchUnlinkedParents(); // Refresh list
      closeModal();
      toast.success('Bags linked to bill successfully!', { position: 'top-center' });
      setError(null);
    } catch (err) {
      if (err.response?.status === 401) {
        setError('Session expired. Please log in again.');
        toast.error('Session expired. Logging out...', { position: 'top-center' });
      } else {
        setError(err.response?.data?.error || 'Linking failed');
        toast.error(err.response?.data?.error || 'Linking failed', { position: 'top-center' });
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
      <div className="bg-white p-6 rounded-xl shadow-lg border-t-4 border-primary w-full max-w-lg">
        <h2 className="text-2xl font-semibold text-accent mb-4 flex items-center">
          <DocumentTextIcon className="w-6 h-6 text-primary mr-2" />
          Link Bags to Bill
        </h2>
        <form onSubmit={handleSubmit} className="space-y-4">
          <input
            value={billID}
            onChange={(e) => setBillID(e.target.value)}
            placeholder="Bill ID"
            className="w-full p-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary text-accent"
            required
          />
          <input
            type="number"
            value={capacity}
            onChange={(e) => setCapacity(Math.max(1, parseInt(e.target.value) || 1))}
            min="1"
            placeholder="Number of Parent Bags"
            className="w-full p-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary text-accent"
            required
          />
          <div className="max-h-40 overflow-y-auto">
            {fetchLoading ? (
              <p className="text-accent text-center">Loading unlinked parents...</p>
            ) : unlinkedParents.length === 0 ? (
              <p className="text-accent text-center">No unlinked parents found.</p>
            ) : (
              unlinkedParents.map((parent) => (
                <div key={parent.bag.id} className="flex justify-between items-center py-2 border-b">
                  <span>{parent.bag.qrCode}</span>
                  {parentIDs.includes(parent.bag.id) ? (
                    <button
                      type="button"
                      onClick={() => handleRemoveParent(parent.bag.id)}
                      className="text-red-500 hover:text-red-700"
                    >
                      Remove
                    </button>
                  ) : (
                    <button
                      type="button"
                      onClick={() => handleAddParent(parent.bag.id)}
                      disabled={parentIDs.length >= capacity}
                      className={`text-primary hover:text-green-700 ${
                        parentIDs.length >= capacity ? 'opacity-50 cursor-not-allowed' : ''
                      }`}
                    >
                      Add
                    </button>
                  )}
                </div>
              ))
            )}
          </div>
          <p className="text-accent">Selected: {parentIDs.length}/{capacity}</p>
          <div className="flex space-x-2">
            <button
              type="submit"
              disabled={isLoading || parentIDs.length !== capacity}
              className={`flex-1 py-3 rounded-lg text-white flex items-center justify-center ${
                isLoading || parentIDs.length !== capacity ? 'bg-gray-400' : 'bg-secondary hover:bg-yellow-600'
              } transition duration-300`}
            >
              {isLoading ? (
                <svg className="animate-spin h-5 w-5 mr-2" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z" />
                </svg>
              ) : (
                <DocumentTextIcon className="w-5 h-5 mr-2" />
              )}
              Link
            </button>
            <button
              type="button"
              onClick={closeModal}
              className="flex-1 py-3 rounded-lg bg-gray-300 text-accent hover:bg-gray-400 transition duration-300"
            >
              Cancel
            </button>
          </div>
        </form>
      </div>
    </motion.div>
  );
}

export default LinkBagsToBillModal;