/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}"
  ],
  theme: {
    backgroundImage: {
      'footer': "url('/src/assets/images/footer.png')",
    },
    extend: {},
  },
  plugins: [],
}
