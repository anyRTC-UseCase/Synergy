module.exports = {
  purge: ["./src/**/*.html", "./src/**/*.vue", "./src/**/*.jsx"],
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {
      colors: {
        synergy: {
          blue_300: "#0087FF",
          blue_400: "#3F72FF",
          gray_200: "#999999",
          gray_400: "#4D4D4D",
          black_400: "#1A1A1A",
          green_400: "#00D067",
          red_400: "#FF432B",
          yellow_400: "#FF9A2B",
        },
      },
      minHeight: {
        96: "24rem",
      },
    },
  },
  variants: {
    extend: {},
  },
  plugins: [],
};
