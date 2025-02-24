import React, { useState, useCallback } from 'react';
import axios from 'axios';
import { motion } from 'framer-motion';
import { toast } from 'react-toastify';
import { ArchiveBoxIcon } from '@heroicons/react/24/solid';
import jsQR from 'jsqr';
import debounce from 'lodash/debounce';

function RegisterChildModal({ parent, closeModal, setError, token }) {
  const [currentCount, setCurrentCount] = useState(0);
  const [isLoading, setIsLoading] = useState(false);
  const [photoPreview, setPhotoPreview] = useState(null);

  const handlePhotoCapture = useCallback(debounce(async () => {
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
      if (!qr) {
        throw new Error('Child QR code is required');
      }
      const res = await axios.post(
        'http://localhost:8080/api/register-child',
        { qrCode: qr, type: 'child', parentId: parent.id },
        { headers: { Authorization: `Bearer ${token}` } }
      );
      setCurrentCount(res.data.currentCount);
      if (res.data.currentCount === res.data.capacity) {
        closeModal();
        toast.success('All child bags registered! Parent bag completed.', { position: 'top-center' });
      } else {
        toast.success('Child seed bag registered!', { position: 'top-center' });
      }
      setError(null);
    } catch (err) {
      if (err.name === 'NotAllowedError' || err.name === 'PermissionDeniedError') {
        setError('Camera access denied. Please enable camera permissions in your browser settings.');
        toast.error('Camera access denied. Please enable permissions.', { position: 'top-center' });
      } else if (err.response?.status === 401) {
        setError('Session expired. Please log in again.');
        toast.error('Session expired. Logging out...', { position: 'top-center' });
      } else {
        setError(err.message || err.response?.data?.error || 'Registration failed');
        toast.error(err.message || err.response?.data?.error || 'Registration failed', { position: 'top-center' });
      }
    } finally {
      setIsLoading(false);
      setPhotoPreview(null);
    }
  }, 1000, { leading: true, trailing: false }), [isLoading, parent.id, token]);

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
        <div className="space-y-4">
          <button
            onClick={handlePhotoCapture}
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
            Take Photo to Register Child
          </button>
          {isLoading && (
            <div className="text-accent text-center">
              Scanning... <span className="animate-pulse">5 seconds remaining</span>
            </div>
          )}
          {photoPreview && (
            <img src={photoPreview} alt="Captured QR" className="w-full rounded-lg" />
          )}
          <div className="w-full bg-gray-200 rounded-full h-3">
            <motion.div
              className="bg-primary h-3 rounded-full"
              initial={{ width: 0 }}
              animate={{ width: `${(currentCount / parent.childCount) * 100}%` }}
              transition={{ duration: 0.5 }}
            />
          </div>
        </div>
      </div>
    </motion.div>
  );
}

export default RegisterChildModal;