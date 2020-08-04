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
    '@typescript-eslint/ban-ts-ignore': 0,
    'arrow-parens': 0,
    'class-methods-use-this': 0,
    'object-curly-newline': 0,
    '@typescript-eslint/member-delimiter-style': 0,
    'operator-linebreak': 0,
    'lines-between-class-members': 0,
  },
}
