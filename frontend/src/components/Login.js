import React from 'react';
import axios from 'axios';
import { motion } from 'framer-motion';
import { SparklesIcon } from '@heroicons/react/24/solid';
import { toast } from 'react-toastify';

function Login() {
  const [error, setError] = useState(null); // Initialize error as null

  const login = async () => {
    setError(null); // Clear any previous errors when the login is attempted

    try {
      // Simulate an API call or some login logic
      // Replace this with your actual login logic
      const response = await simulateLogin(); 

      if (response.success) {
        // Handle successful login (e.g., redirect)
        console.log("Login successful!");
      } else {
          setError(response.message); //set error from response
      }
    } catch (err) {
      setError("An unexpected error occurred during login.");
      console.error(err);
    }
  };
  //Function to test
    async function simulateLogin() {
      // Simulate network delay
      await new Promise(resolve => setTimeout(resolve, 1000));
  
      // Simulate a 50% chance of login success
      const success = Math.random() > 0.5;
      if (success) {
        return { success: true };
      } else {
        return { success: false, message: "Invalid username or password." };
      }
    }

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gray-100">
      <div className="bg-white p-8 rounded-lg shadow-md w-96">
        <h1 className="text-3xl font-bold mb-6 text-center">Login</h1>
        {/* Your form inputs here */}

        {/* Login Button */}
        <button
          onClick={login} // Ensure this references the defined login function
          className="w-full bg-primary text-white py-3 rounded-lg hover:bg-green-700 transition duration-300 flex items-center justify-center"
        >
          <SparklesIcon className="w-5 h-5 mr-2" />
          Login to Baggsy
        </button>

        {/* Error Display */}
        {error && <p className="text-red-500 text-center mt-4">{error}</p>}
      </div>
    </div>
  );
}

export default Login;