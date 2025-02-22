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
  const [view, setView] = useState('register');
  const [showLinkModal, setShowLinkModal] = useState(false);

  useEffect(() => {
    if (token) {
      axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
      validateToken();
    }
  }, [token]);

  const validateToken = async () => {
    try {
      const res = await axios.get('http://localhost:8080/api/bags', { headers: { Authorization: `Bearer ${token}` } });
      if (res.status === 200 && role !== localStorage.getItem('role')) {
        setRole(null);
        logout();
      }
    } catch (err) {
      if (err.response?.status === 401) {
        logout('Session expired. Please log in again.');
      }
    }
  };

  const logout = (message = 'Logged out successfully.') => {
    setToken(null);
    setRole(null);
    localStorage.removeItem('token');
    localStorage.removeItem('role');
    delete axios.defaults.headers.common['Authorization'];
    toast.info(message, { position: 'top-center' });
  };

  const renderView = () => {
    switch (view) {
      case 'register':
        return <RegisterParent setError={setError} token={token} />;
      case 'listBags':
        return role === 'admin' ? <ListBags setError={setError} token={token} /> : <p className="text-red-500">Unauthorized</p>;
      case 'listBills':
        return role === 'admin' ? <ListBills setError={setError} token={token} /> : <p className="text-red-500">Unauthorized</p>;
      case 'search':
        return role === 'admin' ? <Search setError={setError} token={token} /> : <p className="text-red-500">Unauthorized</p>;
      default:
        return <RegisterParent setError={setError} token={token} />;
    }
  };

  return (
    <div className="min-h-screen flex flex-col items-center justify-center p-4">
      <ToastContainer />
      {!token ? (
        <Login setToken={setToken} setRole={setRole} setError={setError} logout={logout} />
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
            <div className="flex space-x-4">
              <button
                onClick={() => setView('register')}
                className={`py-2 px-4 rounded-lg ${view === 'register' ? 'bg-primary text-white' : 'bg-gray-200'}`}
              >
                Register
              </button>
              {role === 'admin' && (
                <>
                  <button
                    onClick={() => setView('listBags')}
                    className={`py-2 px-4 rounded-lg ${view === 'listBags' ? 'bg-primary text-white' : 'bg-gray-200'}`}
                  >
                    List Bags
                  </button>
                  <button
                    onClick={() => setView('listBills')}
                    className={`py-2 px-4 rounded-lg ${view === 'listBills' ? 'bg-primary text-white' : 'bg-gray-200'}`}
                  >
                    List Bills
                  </button>
                  <button
                    onClick={() => setView('search')}
                    className={`py-2 px-4 rounded-lg ${view === 'search' ? 'bg-primary text-white' : 'bg-gray-200'}`}
                  >
                    Search
                  </button>
                </>
              )}
              {(role === 'employee' || role === 'admin') && (
                <button
                  onClick={() => setShowLinkModal(true)}
                  className="py-2 px-4 rounded-lg bg-secondary text-white hover:bg-yellow-600"
                >
                  Link Bags to Bill
                </button>
              )}
              <button
                onClick={() => logout()}
                className="bg-red-500 text-white py-2 px-4 rounded-lg hover:bg-red-600"
              >
                Logout
              </button>
            </div>
          </motion.header>
          {error && <p className="text-red-500 text-center mb-6">{error}</p>}
          {renderView()}
          {showLinkModal && (
            <LinkBagsToBillModal setError={setError} closeModal={() => setShowLinkModal(false)} token={token} />
          )}
        </div>
      )}
    </div>
  );
}

export default App;