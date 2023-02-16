/** @type {import('tailwindcss').Config} */

module.exports = {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    backgroundImage: {
      footer: "url('/src/assets/images/footer.png')",
    },
    extend: {
      colors: {
        'icon-hover': 'rgba(217, 217, 217, 0.4)',
      },
      keyframes: {
      // Toast
      'toast-hide': {
        '0%': { opacity: 1 },
        '100%': { opacity: 0 },
      },
      'toast-slide-in-right': {
        '0%': { transform: `translateX(calc(100% + 1rem))` },
        '100%': { transform: 'translateX(0)' },
      },
      'toast-slide-in-bottom': {
        '0%': { transform: `translateY(calc(100% + 1rem))` },
        '100%': { transform: 'translateY(0)' },
      },
      'toast-swipe-out-x': {
        '0%': { transform: 'translateX(var(--radix-toast-swipe-end-x))' },
        '100%': {
          transform: `translateX(calc(100% + 1rem))`,
        },
      },
      'toast-swipe-out-y': {
        '0%': { transform: 'translateY(var(--radix-toast-swipe-end-y))' },
        '100%': {
          transform: `translateY(calc(100% + 1rem))`,
        },
      },
    },
    animation: {
      // Toast
      'toast-hide': 'toast-hide 100ms ease-in forwards',
      'toast-slide-in-right':
        'toast-slide-in-right 150ms cubic-bezier(0.16, 1, 0.3, 1)',
      'toast-slide-in-bottom':
        'toast-slide-in-bottom 150ms cubic-bezier(0.16, 1, 0.3, 1)',
      'toast-swipe-out-x': 'toast-swipe-out-x 100ms ease-out forwards',
      'toast-swipe-out-y': 'toast-swipe-out-y 100ms ease-out forwards',
    },
    },
  },
  presets: [
    require('@rotational/beacon-foundation/lib/tailwindPreset.config'),
  ],
  safelist: [
    {
      pattern: /^(.*?)/,
    },
  ],

   plugins: [require('@tailwindcss/forms'), require('tailwindcss-radix')()],
};
