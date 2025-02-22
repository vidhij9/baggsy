import React, { useState, useEffect } from 'react';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import Login from './components/Login';
import RegisterParent from './components/RegisterParent';
import ListBags from './components/ListBags';
import LinkBagsToBillModal from './components/LinkBagsToBillModal';
import ListBills from './components/ListBills';
import Search from './components/Search';
import axios from 'axios';
import { motion } from 'framer-motion';
import { SparklesIcon } from '@heroicons/react/24/solid';

function App() {
  const [token, setToken] = useState(localStorage.getItem('token'));
  const [role, setRole] = useState(localStorage.getItem('role'));
  const [error, setError] = useState(null);
  const [showLinkModal, setShowLinkModal] = useState(false);

  useEffect(() => {
    if (token) {
      axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
    }
  }, [token]);

  const logout = () => {
    setToken(null);
    setRole(null);
    localStorage.removeItem('token');
    localStorage.removeItem('role');
    delete axios.defaults.headers.common['Authorization'];
    toast.info('Logged out successfully.', { position: 'top-center' });
  };

  return (
    <div className="min-h-screen flex flex-col items-center justify-center p-4">
      <ToastContainer />
      {!token ? (
        <Login setToken={setToken} setRole={setRole} setError={setError} />
      ) : (
        <div className="w-full max-w-5xl">
          <motion.header
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="flex justify-between items-center mb-8 bg-white p-4 rounded-xl shadow-md"
          >
            <div className="flex items-center">
              <SparklesIcon className="w-8 h-8 text-primary mr-2" />
              <h1 className="text-3xl font-bold text-accent">Star Agriseeds Baggsy</h1>
            </div>
            <button
              onClick={logout}
              className="bg-red-500 text-white py-2 px-4 rounded-lg hover:bg-red-600 transition duration-300"
            >
              Logout
            </button>
          </motion.header>
          {error && <p className="text-red-500 text-center mb-6">{error}</p>}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <RegisterParent setError={setError} />
            {role === "admin" && <ListBags setError={setError} />}
            {(role === "employee" || role === "admin") && (
              <button
                onClick={() => setShowLinkModal(true)}
                className="bg-secondary text-white py-3 rounded-lg hover:bg-yellow-600 transition duration-300"
              >
                Link Bags to Bill
              </button>
            )}
            {role === "admin" && <ListBills setError={setError} />}
            {role === "admin" && <Search setError={setError} />}
          </div>
          {showLinkModal && (
            <LinkBagsToBillModal setError={setError} closeModal={() => setShowLinkModal(false)} />
          )}
        </div>
      )}
    </div>
  );
}

export default App;
