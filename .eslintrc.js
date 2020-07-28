module.exports = {
  root: true,
  env: {
    node: true,
  },
  extends: ['plugin:vue/essential', '@vue/airbnb', '@vue/typescript/recommended'],
  parserOptions: {
    ecmaVersion: 2020,
  },
  rules: {
    'no-console': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
    'no-debugger': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
    'comma-dangle': 0,
    semi: 0,
    'max-classes-per-file': 0,
    'no-plusplus': 0,
    'import/order': 0,
  },
}
