import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App'; // Import the main App component
import './index.css'; // Import global CSS (optional)

// Create a root DOM node where the React app will render
const root = ReactDOM.createRoot(document.getElementById('root'));

// Render the App component into the root DOM node
root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
