/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./src/**/*.{jsx,tsx, ts}'],
  theme: {
    fontFamily: {
      'body': ['SuisseIntl','Calibre', '-apple-system', 'BlinkMacSystemFont', 'Helvetica', 'Arial', 'sans-serif', 'Apple Color Emoji'],
    },
    fontSize: {
      'xs': '0.75rem',  // 12px
      'sm': '0.875rem', // 14px
      'base': '1rem',   // 16px (default)
      'lg': '1.125rem', // 18px
      'xl': '1.25rem',  // 20px
      '2xl': '1.5rem',  // 24px
      '3xl': '1.875rem', // 30px
      '4xl': '2.25rem',  // 36px
      '5xl': '3rem',     // 48px
    },
    fontWeight: {
      thin: 100,
      extralight: 200,
      light: 300,
      normal: 400,
      medium: 500,
      semibold: 600,
      bold: 700,
      extrabold: 800,
      black: 900,
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

