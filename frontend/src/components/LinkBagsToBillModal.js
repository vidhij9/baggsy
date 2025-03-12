import React, { useState, useCallback } from 'react';
import axios from 'axios';
import { motion } from 'framer-motion';
import { toast } from 'react-toastify';
import { DocumentTextIcon } from '@heroicons/react/24/solid';
import jsQR from 'jsqr';
import debounce from 'lodash/debounce';

function LinkBagsToBillModal({ setError, closeModal, token, onSuccess }) {
  const [billID, setBillID] = useState('');
  const [capacity, setCapacity] = useState(1);
  const [parentIDs, setParentIDs] = useState([]);
  const [linkingParents, setLinkingParents] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [photoPreview, setPhotoPreview] = useState(null);

  const validateQR = async (qr) => {
    try {
      const res = await axios.get(`https://baggsy-env.eba-ppg7bx4x.ap-south-1.elasticbeanstalk.com/api/bag/${encodeURIComponent(qr)}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      const bag = res.data.bag;
      if (bag.type !== 'parent') {
        throw new Error('Only parent bags can be linked');
      }
      if (bag.linked) {
        throw new Error('This parent bag is already linked to a bill');
      }
      return bag.id;
    } catch (err) {
      throw new Error(err.response?.data?.error || err.message || 'Invalid or unavailable parent bag');
    }
  };

  const handlePhotoCapture = useCallback(debounce(async () => {
    if (isLoading || parentIDs.length >= capacity) return;
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
      const qrCode = code.data;
      const parentId = await validateQR(qrCode);
      if (parentIDs.includes(parentId)) {
        throw new Error('This parent bag is already added');
      }
      setParentIDs([...parentIDs, parentId]);
      setLinkingParents([...linkingParents, qrCode]);
      toast.success(`Added ${qrCode} (${parentIDs.length + 1}/${capacity})`, { position: 'top-center' });
      setError(null);
    } catch (err) {
      if (err.name === 'NotAllowedError' || err.name === 'PermissionDeniedError') {
        setError('Camera access denied. Please enable camera permissions in your browser settings.');
        toast.error('Camera access denied. Please enable permissions.', { position: 'top-center' });
      } else {
        setError(err.message);
        toast.error(err.message, { position: 'top-center' });
      }
    } finally {
      setIsLoading(false);
      setPhotoPreview(null);
    }
  }, 1000, { leading: true, trailing: false }), [isLoading, capacity, parentIDs, token, linkingParents]);

  const handleRemoveParent = (index) => {
    const newParentIDs = [...parentIDs];
    const newLinkingParents = [...linkingParents];
    newParentIDs.splice(index, 1);
    newLinkingParents.splice(index, 1);
    setParentIDs(newParentIDs);
    setLinkingParents(newLinkingParents);
    toast.info(`Removed ${linkingParents[index]}`, { position: 'top-center' });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!billID.trim()) {
      setError('Bill ID is required');
      toast.error('Bill ID is required', { position: 'top-center' });
      return;
    }
    if (parentIDs.length !== capacity) {
      setError(`Please add exactly ${capacity} parent bag${capacity > 1 ? 's' : ''}`);
      toast.error(`Please add exactly ${capacity} parent bag${capacity > 1 ? 's' : ''}`, { position: 'top-center' });
      return;
    }
    setIsLoading(true);
    try {
      const response = await axios.post(
        'https://baggsy-env.eba-ppg7bx4x.ap-south-1.elasticbeanstalk.com/api/link-bags-to-bill',
        { billID, parentIDs, capacity },
        { headers: { Authorization: `Bearer ${token}` } }
      );
      console.log("Link response:", response.data);
      setBillID('');
      setCapacity(1);
      setParentIDs([]);
      setLinkingParents([]);
      setError(null);
      onSuccess();
      toast.success('Bags linked to bill successfully!', { position: 'top-center' });
      closeModal();
    } catch (err) {
      const errorMsg = err.response?.data?.error || 'Linking failed';
      console.error("Link error:", err.response || err);
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
            disabled={isLoading}
          />
          <input
            type="number"
            value={capacity}
            onChange={(e) => setCapacity(Math.max(1, parseInt(e.target.value) || 1))}
            min="1"
            placeholder="Number of Parent Bags"
            className="w-full p-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary text-accent"
            required
            disabled={isLoading || parentIDs.length > 0}
          />
          <div className="space-y-2">
            <button
              type="button"
              onClick={handlePhotoCapture}
              disabled={isLoading || parentIDs.length >= capacity}
              className={`w-full py-3 rounded-lg text-white flex items-center justify-center ${
                isLoading || parentIDs.length >= capacity ? 'bg-gray-400' : 'bg-primary hover:bg-green-700'
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
              Take Photo to Add Parent Bag
            </button>
            {isLoading && (
              <div className="text-accent text-center">
                Scanning... <span className="animate-pulse">5 seconds remaining</span>
              </div>
            )}
            {photoPreview && (
              <img src={photoPreview} alt="Captured QR" className="w-full rounded-lg" />
            )}
            {linkingParents.length > 0 && (
              <div className="max-h-40 overflow-y-auto">
                {linkingParents.map((qr, index) => (
                  <div key={qr} className="flex justify-between items-center py-2 border-b">
                    <span>{qr}</span>
                    <button
                      type="button"
                      onClick={() => handleRemoveParent(index)}
                      className="text-red-500 hover:text-red-700"
                      disabled={isLoading}
                    >
                      Remove
                    </button>
                  </div>
                ))}
              </div>
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
              disabled={isLoading}
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