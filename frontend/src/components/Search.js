import React, { useState } from 'react';
import axios from 'axios';
import { motion } from 'framer-motion';
import { toast } from 'react-toastify';
import { MagnifyingGlassIcon } from '@heroicons/react/24/solid';

function Search({ setError, token }) {
  const [qr, setQr] = useState('');
  const [bill, setBill] = useState('');
  const [result, setResult] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const [expandedParents, setExpandedParents] = useState({});

  const validateQR = (qrCode) => qrCode.trim().length > 0;
  const validateBill = (billID) => billID.trim().length > 0;

  const searchByQr = async () => {
    if (!validateQR(qr)) {
      setError('Invalid QR code');
      toast.error('Invalid QR code', { position: 'top-center' });
      return;
    }
    setIsLoading(true);
    try {
      const res = await axios.get(`http://localhost:8080/api/bag/${encodeURIComponent(qr)}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      setResult({ type: 'bag', data: res.data });
      setError(null);
      toast.success('Bag found!', { position: 'top-center' });
    } catch (err) {
      const errorMsg = err.response?.data?.error || 'Search failed';
      if (err.response?.status === 401) {
        setError('Session expired. Please log in again.');
        toast.error('Session expired. Logging out...', { position: 'top-center' });
      } else {
        setError(errorMsg);
        toast.error(errorMsg, { position: 'top-center' });
      }
    } finally {
      setIsLoading(false);
    }
  };

  const searchByBill = async () => {
    if (!validateBill(bill)) {
      setError('Invalid Bill ID');
      toast.error('Invalid Bill ID', { position: 'top-center' });
      return;
    }
    setIsLoading(true);
    try {
      const res = await axios.get(`http://localhost:8080/api/bill/${encodeURIComponent(bill)}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      setResult({ type: 'bill', data: res.data });
      setError(null);
      toast.success('Bill found!', { position: 'top-center' });
    } catch (err) {
      const errorMsg = err.response?.data?.error || 'Search failed';
      if (err.response?.status === 401) {
        setError('Session expired. Please log in again.');
        toast.error('Session expired. Logging out...', { position: 'top-center' });
      } else {
        setError(errorMsg);
        toast.error(errorMsg, { position: 'top-center' });
      }
    } finally {
      setIsLoading(false);
    }
  };

  const toggleExpandParent = (parentId) => {
    setExpandedParents((prev) => ({ ...prev, [parentId]: !prev[parentId] }));
  };

  return (
    <motion.div
      initial={{ opacity: 0, scale: 0.95 }}
      animate={{ opacity: 1, scale: 1 }}
      className="bg-white p-6 rounded-xl shadow-lg border-t-4 border-primary"
    >
      <h2 className="text-2xl font-semibold text-accent mb-4 flex items-center">
        <MagnifyingGlassIcon className="w-6 h-6 text-primary mr-2" />
        Search
      </h2>
      <div className="space-y-4">
        <div className="flex space-x-2">
          <input
            value={qr}
            onChange={(e) => setQr(e.target.value)}
            placeholder="Bag QR"
            className="w-full p-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary text-accent"
          />
          <button
            onClick={searchByQr}
            disabled={isLoading || !qr}
            className={`py-2 px-4 rounded-lg text-white ${
              isLoading || !qr ? 'bg-gray-400' : 'bg-primary hover:bg-green-700'
            } transition duration-300`}
          >
            {isLoading && qr ? '...' : 'QR'}
          </button>
        </div>
        <div className="flex space-x-2">
          <input
            value={bill}
            onChange={(e) => setBill(e.target.value)}
            placeholder="Bill ID"
            className="w-full p-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary text-accent"
          />
          <button
            onClick={searchByBill}
            disabled={isLoading || !bill}
            className={`py-2 px-4 rounded-lg text-white ${
              isLoading || !bill ? 'bg-gray-400' : 'bg-secondary hover:bg-yellow-600'
            } transition duration-300`}
          >
            {isLoading && bill ? '...' : 'Bill'}
          </button>
        </div>
        {isLoading && <p className="text-accent text-center">Searching...</p>}
        {result && !isLoading && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="mt-4 p-4 bg-background rounded-lg border border-primary"
          >
            {result.type === 'bag' && (
              <>
                <p><strong>Bag:</strong> {result.data.bag.qrCode} ({result.data.bag.type})</p>
                {result.data.bag.type === 'parent' && result.data.children && result.data.children.length > 0 ? (
                  <>
                    <p><strong>Children:</strong></p>
                    <ul>
                      {result.data.children.map((child) => (
                        <li key={child.id}>{child.qrCode}</li>
                      ))}
                    </ul>
                  </>
                ) : result.data.bag.type === 'parent' ? (
                  <p>No children linked.</p>
                ) : null}
                {result.data.parentQR && <p><strong>Parent:</strong> {result.data.parentQR}</p>}
                {result.data.billID && <p><strong>Bill ID:</strong> {result.data.billID}</p>}
              </>
            )}
            {result.type === 'bill' && (
              <>
                <p><strong>Bill ID:</strong> {result.data.billID}</p>
                <div>
                  {result.data.bags && result.data.bags.length > 0 ? (
                    result.data.bags.map((bag) => (
                      <div key={bag.id} className="border-b py-1">
                        <div
                          className="flex justify-between items-center cursor-pointer"
                          onClick={() => toggleExpandParent(bag.id)}
                        >
                          <span>{bag.qrCode}</span>
                          <span>{expandedParents[bag.id] ? '▲' : '▼'}</span>
                        </div>
                        {expandedParents[bag.id] && bag.children && (
                          <ul className="ml-4">
                            {bag.children.map((child) => (
                              <li key={child.id}>{child.qrCode}</li>
                            ))}
                          </ul>
                        )}
                      </div>
                    ))
                  ) : (
                    <p>No bags linked to this bill.</p>
                  )}
                </div>
              </>
            )}
          </motion.div>
        )}
      </div>
    </motion.div>
  );
}

export default Search;