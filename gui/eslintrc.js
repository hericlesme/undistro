module.exports = {
  root: true,
  parserOptions: { ecmaVersion: 8 },
  parser: '@typescript-eslint/parser',
  settings: {
    react: { version: 'detect' },
  },
  ignorePatterns: ['node_modules/*', '.next/*', '.out/*', '!.prettierrc.js'],
  env: {
    browser: true,
    es6: true,
    node: true,
    jest: true
  },
  plugins: ['jsx-a11y', '@typescript-eslint', 'react-hooks'],
  extends: [
    'next/core-web-vitals',
    'eslint:recommended',
    'plugin:react/recommended',
    'plugin:jsx-a11y/recommended',
    'plugin:prettier/recommended',
    'plugin:@typescript-eslint/recommended',
    'prettier'
  ],
  rules: {
    'prettier/prettier': ['error', {}, { usePrettierrc: true }],
    'react-hooks/rules-of-hooks': 'error',
    'react-hooks/exhaustive-deps': 'warn',
    'react/react-in-jsx-scope': 'off',
    'react/prop-types': 'off', // TS types for component props
    'jsx-a11y/anchor-is-valid': [
      'error',
      {
        components: ['Link'],
        specialLink: ['hrefLeft', 'hrefRight'],
        aspects: ['invalidHref', 'preferButton']
      }
    ],
    '@typescript-eslint/explicit-module-boundary-types': 'off'
  }
}
