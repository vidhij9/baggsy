import React, { useState } from 'react';
import axios from 'axios';
import { motion } from 'framer-motion';
import { SparklesIcon } from '@heroicons/react/24/solid';
import { toast } from 'react-toastify';

function Login({ setToken, setRole, setError, logout }) {
  const [localError, setLocalError] = useState(null);
  const [username, setUsername] = useState('admin');
  const [password, setPassword] = useState('password');

  const login = async (e) => {
    e.preventDefault();
    console.log("Attempting login to http://localhost:8080/login");
    try {
      const res = await axios.post('http://localhost:8080/login', { username, password });
      console.log("Login successful:", res.data);
      setToken(res.data.token);
      localStorage.setItem('token', res.data.token);
      const tokenParts = res.data.token.split('.');
      const payload = JSON.parse(atob(tokenParts[1]));
      setRole(payload.role);
      localStorage.setItem('role', payload.role);
      axios.defaults.headers.common['Authorization'] = `Bearer ${res.data.token}`;
      toast.success('Logged into Star Agriseeds Baggsy!', { position: 'top-center' });
      setLocalError(null);
      setError(null);
    } catch (err) {
      const errorMsg = err.response?.data?.error || 'Login failed. Please check backend logs.';
      console.error("Login error:", err.message, err.response);
      setLocalError(errorMsg);
      setError(errorMsg);
      toast.error(errorMsg, { position: 'top-center' });
      if (err.response?.status === 401) {
        logout('Invalid credentials. Please try again.');
      }
    }
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="bg-white rounded-xl shadow-2xl p-8 w-full max-w-md border-t-4 border-primary"
    >
      <div className="flex items-center justify-center mb-6">
        <SparklesIcon className="w-8 h-8 text-primary mr-2" />
        <h1 className="text-3xl font-bold text-accent">Star Agriseeds</h1>
      </div>
      <form onSubmit={login} className="space-y-4">
        <input
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          placeholder="Username"
          className="w-full p-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary text-accent"
          required
        />
        <input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          placeholder="Password"
          className="w-full p-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary text-accent"
          required
        />
        <button
          type="submit"
          className="w-full bg-primary text-white py-3 rounded-lg hover:bg-green-700 transition duration-300 flex items-center justify-center"
        >
          <SparklesIcon className="w-5 h-5 mr-2" />
          Login to Baggsy
        </button>
      </form>
      {localError && <p className="text-red-500 text-center mt-4">{localError}</p>}
    </motion.div>
  );
}

export default Login;