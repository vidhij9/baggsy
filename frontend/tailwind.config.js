/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./src/**/*.{js,jsx,ts,tsx}"],
  theme: {
    extend: {
      colors: {
        green: "#4CAF50",
        gold: "#FFD700",
        white: "#FFFFFF",
        lightGreen: "#E8F5E9",
        darkGreen: "#2E7D32",
      },
      fontFamily: {
        sans: ["Roboto", "Poppins", "Arial", "sans-serif"],
      },
    },
  },
  plugins: [],
};
