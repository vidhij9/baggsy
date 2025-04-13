import React, { useState } from 'react';
import axios from 'axios';
import { motion } from 'framer-motion';
import { SparklesIcon } from '@heroicons/react/24/solid';
import { toast } from 'react-toastify';
const API_URL = process.env.REACT_APP_API_URL;

function Login({ setToken, setRole, setError, logout, switchView }) {
  const [isRegistering, setIsRegistering] = useState(false);
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [email, setEmail] = useState('');
  const [role, setLocalRole] = useState('employee');
  const [localError, setLocalError] = useState(null);

  const login = async (e) => {
    e.preventDefault();
    console.log("Attempting login to baggsy");
    try {
      const res = await axios.post(`${API_URL}/login`, { username, password });
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
      switchView('register'); // Default view after login
    } catch (err) {
      const errorMsg = err.response?.data?.error || 'Login failed.';
      console.error("Login error:", err.message, err.response);
      setLocalError(errorMsg);
      setError(errorMsg);
      toast.error(errorMsg, { position: 'top-center' });
      if (err.response?.status === 401) {
        logout('Invalid credentials. Please try again.');
      }
    }
  };

  const register = async (e) => {
    e.preventDefault();
    if (!username || !password || !email) {
      setLocalError('All fields are required');
      toast.error('All fields are required', { position: 'top-center' });
      return;
    }
    try {
      const res = await axios.post(`${API_URL}/register`, {
        username,
        password,
        email,
        role,
      });
      console.log("Registration response:", res.data);
      toast.success(res.data.message, { position: 'top-center' });
      setLocalError(null);
      setError(null);
      setIsRegistering(false);
      setUsername('');
      setPassword('');
      setEmail('');
      setLocalRole('employee');
    } catch (err) {
      const errorMsg = err.response?.data?.error || 'Registration failed';
      console.error("Registration error:", err.message, err.response);
      setLocalError(errorMsg);
      setError(errorMsg);
      toast.error(errorMsg, { position: 'top-center' });
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
      {isRegistering ? (
        <form onSubmit={register} className="space-y-4">
          <input
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            placeholder="Username"
            className="w-full p-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary text-accent"
            required
          />
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="Email"
            className="w-full p-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary text-accent"
            required
          />
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="Password (8+ chars, uppercase, number)"
            className="w-full p-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary text-accent"
            required
          />
          <select
            value={role}
            onChange={(e) => setLocalRole(e.target.value)}
            className="w-full p-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary text-accent"
          >
            <option value="employee">Employee</option>
            <option value="admin">Admin</option>
          </select>
          <button
            type="submit"
            className="w-full bg-primary text-white py-3 rounded-lg hover:bg-green-700 transition duration-300 flex items-center justify-center"
          >
            <SparklesIcon className="w-5 h-5 mr-2" />
            Register
          </button>
          <button
            type="button"
            onClick={() => setIsRegistering(false)}
            className="w-full bg-gray-300 text-accent py-3 rounded-lg hover:bg-gray-400 transition duration-300"
          >
            Back to Login
          </button>
        </form>
      ) : (
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
          <button
            type="button"
            onClick={() => setIsRegistering(true)}
            className="w-full bg-gray-300 text-accent py-3 rounded-lg hover:bg-gray-400 transition duration-300"
          >
            Create New Account
          </button>
        </form>
      )}
      {localError && <p className="text-red-500 text-center mt-4">{localError}</p>}
    </motion.div>
  );
}

export default Login;