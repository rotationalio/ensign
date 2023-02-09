/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    backgroundImage: {
      hexagon: "url('/src/assets/images/tileable-hexagon.png')",
      footer: "url('/src/assets/images/footer.png')",
    },
    extend: {
      colors: {
        'icon-hover': 'rgba(217, 217, 217, 0.4)',
      },
    },
  },
  presets: [require('@rotational/beacon-foundation/lib/tailwindPreset.config')],
  safelist: [
    {
      pattern: /^(.*?)/,
    },
  ],

  plugins: [],
};
