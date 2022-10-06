/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./src/**/*.{js,jsx,ts,tsx}"],
  theme: { backgroundImage: {
    'hero': ("url('/src/assets/hero.png')")
  },
    extend: {},
  },
  plugins: [],
}
