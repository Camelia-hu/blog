/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './index.html',
    './src/**/*.{vue,js,ts,jsx,tsx}',
  ],
  theme: {
    extend: {
      colors: {
        // Catppuccin Mocha 调色板（MkDocs Material 风格）
        ctp: {
          base:     '#1e1e2e',
          mantle:   '#181825',
          crust:    '#11111b',
          surface0: '#313244',
          surface1: '#45475a',
          surface2: '#585b70',
          overlay0: '#6c7086',
          overlay1: '#7f849c',
          text:     '#cdd6f4',
          subtext:  '#a6adc8',
          teal:     '#94e2d5',
          sky:      '#89dceb',
          blue:     '#89b4fa',
          sapphire: '#74c7ec',
          mauve:    '#cba6f7',
          green:    '#a6e3a1',
          yellow:   '#f9e2af',
          peach:    '#fab387',
          red:      '#f38ba8',
          pink:     '#f5c2e7',
          lavender: '#b4befe',
        }
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
        mono: ['JetBrains Mono', 'Fira Code', 'monospace'],
      },
      typography: {
        DEFAULT: {
          css: {
            color: '#cdd6f4',
            maxWidth: '100%',
            a: { color: '#89b4fa', '&:hover': { color: '#89dceb' } },
            strong: { color: '#cdd6f4' },
            h1: { color: '#cdd6f4' },
            h2: { color: '#cdd6f4' },
            h3: { color: '#cdd6f4' },
            h4: { color: '#cdd6f4' },
            code: {
              color: '#89dceb',
              backgroundColor: '#313244',
              borderRadius: '4px',
              padding: '2px 6px',
              fontWeight: '400',
              '&::before': { content: 'none' },
              '&::after': { content: 'none' },
            },
            pre: {
              backgroundColor: '#181825',
              color: '#cdd6f4',
              borderRadius: '8px',
              border: '1px solid #313244',
            },
            'pre code': {
              backgroundColor: 'transparent',
              padding: '0',
              color: 'inherit',
            },
            blockquote: {
              color: '#a6adc8',
              borderLeftColor: '#89dceb',
            },
            hr: { borderColor: '#313244' },
            table: { color: '#cdd6f4' },
            thead: { borderBottomColor: '#45475a' },
            'tbody tr': { borderBottomColor: '#313244' },
            th: { color: '#cdd6f4' },
          }
        }
      }
    },
  },
  plugins: [
    require('@tailwindcss/typography'),
  ],
}
