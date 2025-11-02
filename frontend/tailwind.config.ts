import type { Config } from 'tailwindcss';

export default {
  content: ['./index.html', './src/**/*.{vue,ts,tsx}'],
  theme: {
    extend: {
      colors: {
        brand: {
          DEFAULT: '#4A5FFF',
          50: '#EEF0FF',
          100: '#CDD2FF',
          200: '#ACB3FF',
          300: '#8B95FF',
          400: '#6A76FF',
          500: '#4A5FFF',
          600: '#3B4ACC',
          700: '#2C3799',
          800: '#1D2466',
          900: '#0E1233'
        }
      },
      boxShadow: {
        surface: '0 8px 24px rgba(15, 23, 42, 0.08)'
      }
    }
  },
  plugins: []
} satisfies Config;

