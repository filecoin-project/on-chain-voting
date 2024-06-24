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
          'block-bg': '#ffffff', // 设置自定义背景颜色
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

