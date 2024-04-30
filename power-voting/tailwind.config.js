/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [ './src/**/*.{jsx,tsx}' ],
  theme: {
    fontFamily: {
      'body': ['Calibre', '-apple-system', 'BlinkMacSystemFont', 'Helvetica', 'Arial', 'sans-serif', 'Apple Color Emoji'],
    },
    extend: {},
    listStyleType: {
      'none': 'none'
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}

