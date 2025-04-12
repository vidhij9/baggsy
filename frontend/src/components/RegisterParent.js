import React, { useState, useCallback } from 'react';
import axios from 'axios';
import { motion } from 'framer-motion';
import { toast } from 'react-toastify';
import { ArchiveBoxIcon } from '@heroicons/react/24/solid';
import RegisterChildModal from './RegisterChildModal';
import jsQR from 'jsqr';
import debounce from 'lodash/debounce';

function RegisterParent({ setError, token }) {
  const [isLoading, setIsLoading] = useState(false);
  const [showChildModal, setShowChildModal] = useState(false);
  const [parentData, setParentData] = useState(null);
  const [photoPreview, setPhotoPreview] = useState(null);

  const validateQR = (qrCode) => {
    const parts = qrCode.split('-');
    if (parts.length !== 2 || !parts[0].startsWith('P')) return false;
    const count = parseInt(parts[1], 10);
    return !isNaN(count) && count > 0;
  };

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
      if (!validateQR(qr)) {
        throw new Error('Invalid QR format. Use P<Number>-<ChildCount> (e.g., P123-10)');
      }
      const res = await axios.post(
        'https://baggsy-backend.up.railway.app/api/register-parent',
        { qrCode: qr, type: 'parent' },
        { headers: { Authorization: `Bearer ${token}` } }
      );
      if (res.data.childCount <= 0) {
        throw new Error('Parent bag must have at least one child bag');
      }
      setParentData(res.data);
      setShowChildModal(true);
      toast.success('Parent seed bag registered! Now add child bags.', { position: 'top-center' });
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
  }, 1000, { leading: true, trailing: false }), [isLoading, token]);

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
            Take Photo to Register Parent
          </button>
          {isLoading && (
            <div className="text-accent text-center">
              Scanning... <span className="animate-pulse">5 seconds remaining</span>
            </div>
          )}
          {photoPreview && (
            <img src={photoPreview} alt="Captured QR" className="w-full rounded-lg" />
          )}
        </div>
      </motion.div>
      {showChildModal && parentData && (
        <RegisterChildModal
          parent={parentData}
          closeModal={() => {
            setShowChildModal(false);
            setParentData(null);
          }}
          setError={setError}
          token={token}
        />
      )}
    </>
  );
}

export default RegisterParent;