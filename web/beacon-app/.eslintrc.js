module.exports = {
  root: true,
  env: {
    node: true,
    es6: true
  },
  parserOptions: {
    ecmaVersion: 8,
    sourceType: 'module'
  },
  ignorePatterns: ['node_modules/*'],
  extends: ['eslint:recommended', 'plugin:storybook/recommended'],
  overrides: [{
    files: ['**/*.ts', '**/*.tsx'],
    parser: '@typescript-eslint/parser',
    settings: {
      react: {
        version: 'detect'
      },
      'import/resolver': {
        typescript: {}
      }
    },
    env: {
      browser: true,
      node: true,
      es6: true
    },
    plugins: ['@typescript-eslint', 'unused-imports', 'tailwindcss', 'simple-import-sort'],
    extends: ['eslint:recommended', 'plugin:import/errors', 'plugin:import/warnings', 'plugin:import/typescript', 'plugin:@typescript-eslint/recommended', 'plugin:react/recommended', 'plugin:react-hooks/recommended', 'plugin:jsx-a11y/recommended', 'plugin:prettier/recommended', 'plugin:testing-library/react'],
    rules: {
      'no-restricted-imports': ['error', {
        patterns: ['@/features/*/*']
      }],
      'linebreak-style': ['error', 'unix'],
      'react/prop-types': 'off',
      'import/default': 'off',
      'import/no-named-as-default-member': 'off',
      'import/no-named-as-default': 'off',
      'react/react-in-jsx-scope': 'off',
      'jsx-a11y/anchor-is-valid': 'off',
      '@typescript-eslint/no-unused-vars': 'off',
      '@typescript-eslint/explicit-function-return-type': ['off'],
      '@typescript-eslint/explicit-module-boundary-types': ['off'],
      '@typescript-eslint/no-empty-function': ['off'],
      '@typescript-eslint/no-explicit-any': ['off'],
      'simple-import-sort/imports': 'error',
      // Import configuration for `eslint-plugin-simple-import-sort`
      'simple-import-sort/exports': 'error',
      // Export configuration for `eslint-plugin-simple-import-sort`
      'unused-imports/no-unused-imports': 'error',
      'unused-imports/no-unused-vars': ['error', {
        argsIgnorePattern: '^_'
      }],
      'prettier/prettier': ['error', {}, {
        usePrettierrc: true
      }]
    }
  }]
};