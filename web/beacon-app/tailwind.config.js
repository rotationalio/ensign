/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  presets: [require('@rotational/beacon-foundation/lib/tailwindPreset.config')],
  safelist: [
    {
      pattern: /^(.*?)/,
    },
  ],
};
