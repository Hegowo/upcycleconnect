/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        brand: {
          navy:          '#103652',
          'navy-light':  '#174d72',
          blue:          '#1B5B88',
          'blue-light':  '#2470a8',
          'dark-green':  '#10522b',
          green:         '#1B8848',
          'green-light': '#22a357',
        },
      },
    },
  },
  plugins: [],
}
