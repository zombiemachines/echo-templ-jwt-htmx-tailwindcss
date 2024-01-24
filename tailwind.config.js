/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./views/**/*.{html,js,templ,go,css}",
    "./static/**/*.{html,js,templ,go,css}",
    "./static/css/**/*.{html,js,css,templ,go}"],
  theme: {
    extend: {
      colors: {
        clifford: '#67e8f9',
        border: '#06b6d4',
      },
      fontFamily: {
        'sans': ["Quicksand"],
      },
      dropShadow: {
        '2xl': '0 8px 8px rgba(255, 255, 255, 0.10)',
        '3xl': [
          '0 15px 15px rgba(255, 255, 255, 0.10)'
        ]
      },
    }
  },
  // plugins: [require("@tailwindcss/forms"), require("@tailwindcss/typography")],
}


