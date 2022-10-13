/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./src/**/*.{js,jsx,ts,tsx}"],
  theme: {
    backgroundImage: {
      'hero': "url('/src/components/layout/img/hero.png')",
      'footer': "url('/src/components/layout/img/foot.png')",
      'wave': "url('/src/components/layout/img/footer-wave.png')",
    },
    extend: {}
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}
