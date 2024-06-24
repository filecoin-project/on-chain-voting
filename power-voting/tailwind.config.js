/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [ './src/**/*.{jsx,tsx, ts}' ],
  theme: {
    fontFamily: {
      'body': ['Calibre', '-apple-system', 'BlinkMacSystemFont', 'Helvetica', 'Arial', 'sans-serif', 'Apple Color Emoji'],
    },
    extend: {
      colors: {
        skin: {
          'block-bg': '#ffffff', // set bg color
        },
      },
    },
    listStyleType: {
      'none': 'none'
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}

