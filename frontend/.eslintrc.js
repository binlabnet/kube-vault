module.exports = {
  parserOptions: {
    ecmaVersion: 2019,
    sourceType: 'module'
  },

  env: {
    es6: true,
    browser: true
  },

  plugins: [
    "svelte3",
  ],

  overrides: [
    {
      files: ['**/*.svelte'],
      processor: 'svelte3/svelte3'
    }
  ],

  rules: {
    "quotes": [2, "double", "avoid-escape"],

    "camelcase": "off",
    "no-console": "off",
    "no-param-reassign": "off",
    "no-plusplus": "off",
    "no-underscore-dangle": "off",
    "no-bitwise": "off",
    "global-require": "off",
    "max-len": "off",
    "no-mixed-operators": "off",

    "semi": ["error", "never"],
    "space-before-function-paren": ["error", "always"],

    'import/no-extraneous-dependencies': 0,
    'import/no-unresolved': 0,
    'import/extensions': 0,
  }
}