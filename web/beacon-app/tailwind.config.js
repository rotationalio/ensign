/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    backgroundImage: {
      "hexagon": "url('/src/assets/images/tileable-hexagon.png')"
    },
  },
  presets: [require('@rotational/beacon-foundation/lib/tailwindPreset.config')],
  safelist: [
    {
      pattern: /^(.*?)/,
    },
  ],
};
