import React, { useState, useCallback } from 'react';
import axios from 'axios';
import { motion } from 'framer-motion';
import { toast } from 'react-toastify';
import { MagnifyingGlassIcon } from '@heroicons/react/24/solid';
import jsQR from 'jsqr';
import debounce from 'lodash/debounce';

function Search({ setError, token }) {
  const [bill, setBill] = useState('');
  const [result, setResult] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const [expandedParents, setExpandedParents] = useState({});
  const [photoPreview, setPhotoPreview] = useState(null);

  const validateBill = (billID) => billID.trim().length > 0;

  const handlePhotoCaptureForQR = useCallback(debounce(async () => {
    if (isLoading) return;
    setIsLoading(true);
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ video: { facingMode: 'environment' } });
      const video = document.createElement('video');
      video.srcObject = stream;
      video.play();
      const canvas = document.createElement('canvas');
      canvas.width = 640;
      canvas.height = 480;
      const context = canvas.getContext('2d');
      await new Promise(resolve => setTimeout(resolve, 5000)); // Extended to 5 seconds
      context.drawImage(video, 0, 0, canvas.width, canvas.height);
      const imageData = context.getImageData(0, 0, canvas.width, canvas.height);
      const code = jsQR(imageData.data, imageData.width, imageData.height);
      stream.getTracks().forEach(track => track.stop());
      setPhotoPreview(canvas.toDataURL('image/png'));
      if (!code) {
        throw new Error('No QR code detected in photo. Please try again with a clear image.');
      }
      const qr = code.data;
      const res = await axios.get(`https://baggsy-backend.up.railway.app/api/bag/${encodeURIComponent(qr)}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      setResult({ type: 'bag', data: res.data });
      setError(null);
      toast.success('Bag found!', { position: 'top-center' });
    } catch (err) {
      if (err.name === 'NotAllowedError' || err.name === 'PermissionDeniedError') {
        setError('Camera access denied. Please enable camera permissions in your browser settings.');
        toast.error('Camera access denied. Please enable permissions.', { position: 'top-center' });
      } else if (err.response?.status === 401) {
        setError('Session expired. Please log in again.');
        toast.error('Session expired. Logging out...', { position: 'top-center' });
      } else {
        setError(err.message || err.response?.data?.error || 'Search failed');
        toast.error(err.message || err.response?.data?.error || 'Search failed', { position: 'top-center' });
      }
    } finally {
      setIsLoading(false);
      setPhotoPreview(null);
    }
  }, 1000, { leading: true, trailing: false }), [isLoading, token]);

  const searchByBill = async () => {
    if (!validateBill(bill)) {
      setError('Invalid Bill ID');
      toast.error('Invalid Bill ID', { position: 'top-center' });
      return;
    }
    setIsLoading(true);
    try {
      const res = await axios.get(`https://baggsy-backend.up.railway.app/api/bill/${encodeURIComponent(bill)}`, {
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
        <div className="space-y-2">
          <button
            onClick={handlePhotoCaptureForQR}
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
              <MagnifyingGlassIcon className="w-5 h-5 mr-2" />
            )}
            Take Photo to Search Bag
          </button>
          {isLoading && (
            <div className="text-accent text-center">
              Scanning... <span className="animate-pulse">5 seconds remaining</span>
            </div>
          )}
          {photoPreview && (
            <img src={photoPreview} alt="Captured QR" className="w-full rounded-lg" />
          )}
          <div className="flex space-x-2">
            <input
              value={bill}
              onChange={(e) => setBill(e.target.value)}
              placeholder="Bill ID"
              className="w-full p-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary text-accent"
              disabled={isLoading}
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