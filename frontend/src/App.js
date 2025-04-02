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
  const [view, setView] = useState('login');
  const [showLinkModal, setShowLinkModal] = useState(false);
  const [refreshBills, setRefreshBills] = useState(false);

  useEffect(() => {
    if (token) {
      axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
      validateToken();
    } else {
      setView('login');
    }
  }, [token]);

  const validateToken = async () => {
    try {
      const res = await axios.get('https://baggsy.app/api/bags', { headers: { Authorization: `Bearer ${token}` } });
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
    if (!window.confirm('Are you sure you want to log out?')) return;
    setToken(null);
    setRole(null);
    localStorage.removeItem('token');
    localStorage.removeItem('role');
    delete axios.defaults.headers.common['Authorization'];
    toast.info(message, { position: 'top-center' });
    setView('login');
    setError(null);
  };

  const handleLinkSuccess = () => {
    setRefreshBills((prev) => !prev);
  };

  const switchView = (newView) => {
    setView(newView);
    setError(null);
  };

  const renderView = () => {
    switch (view) {
      case 'login':
        return <Login setToken={setToken} setRole={setRole} setError={setError} logout={logout} switchView={switchView} />;
      case 'register':
        return <RegisterParent setError={setError} token={token} />;
      case 'listBags':
        return role === 'admin' ? <ListBags setError={setError} token={token} /> : <p className="text-red-500">Unauthorized</p>;
      case 'listBills':
        return role === 'admin' ? (
          <ListBills setError={setError} token={token} refresh={refreshBills} />
        ) : (
          <p className="text-red-500">Unauthorized</p>
        );
      case 'search':
        return role === 'admin' ? <Search setError={setError} token={token} /> : <p className="text-red-500">Unauthorized</p>;
      default:
        return <Login setToken={setToken} setRole={setRole} setError={setError} logout={logout} switchView={switchView} />;
    }
  };

  return (
    <div className="min-h-screen flex flex-col items-center justify-center p-4">
      <ToastContainer />
      <div className="w-full max-w-5xl">
        {token && (
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
                onClick={() => switchView('register')}
                className={`py-2 px-4 rounded-lg ${view === 'register' ? 'bg-primary text-white' : 'bg-gray-200'}`}
              >
                Register
              </button>
              {role === 'admin' && (
                <>
                  <button
                    onClick={() => switchView('listBags')}
                    className={`py-2 px-4 rounded-lg ${view === 'listBags' ? 'bg-primary text-white' : 'bg-gray-200'}`}
                  >
                    List Bags
                  </button>
                  <button
                    onClick={() => switchView('listBills')}
                    className={`py-2 px-4 rounded-lg ${view === 'listBills' ? 'bg-primary text-white' : 'bg-gray-200'}`}
                  >
                    List Bills
                  </button>
                  <button
                    onClick={() => switchView('search')}
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
        )}
        {error && <p className="text-red-500 text-center mb-6">{error}</p>}
        {renderView()}
        {showLinkModal && (
          <LinkBagsToBillModal
            setError={setError}
            closeModal={() => setShowLinkModal(false)}
            token={token}
            onSuccess={handleLinkSuccess}
          />
        )}
      </div>
    </div>
  );
}

export default App;