/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./src/**/*.{js,jsx,ts,tsx}"], // Ensure this scans all src files
  theme: {
    extend: {
      colors: {
        primary: "#4CAF50", // Green
        secondary: "#FFD700", // Gold
        background: "#F9F9F9", // Light background
        dark: "#2E7D32", // Darker green
        white: "#FFFFFF",
        accent: "#2E7D32",
        lightGray: "#F9F9F9",
      },
      fontFamily: {
        sans: ["Roboto", "Poppins", "Arial", "sans-serif"],
      },
    },
  },
  plugins: [],
};

