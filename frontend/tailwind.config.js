module.exports = {
  content: ["./src/**/*.{js,jsx,ts,tsx}"],
  theme: {
    extend: {
      colors: {
        primary: "#4CAF50", // Green
        secondary: "#FFD700", // Gold
        background: "#F9F9F9", // Light background
        dark: "#2E7D32", // Darker green
        white: "#FFFFFF",
      },
      fontFamily: {
        sans: ["Roboto", "Poppins", "Arial", "sans-serif"],
      },
    },
  },
  plugins: [],
};