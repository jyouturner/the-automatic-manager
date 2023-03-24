module.exports = {
  env: {
    browser: true,
    es2021: true,
  },
  extends: [
    "eslint:recommended",
    "plugin:react/recommended",
    "prettier",
    "plugin:jsx-a11y/recommended",
  ],
  parserOptions: {
    parser: "@babel/eslint-parser",
    ecmaFeatures: {
      jsx: true,
    },
    ecmaVersion: "latest",
    sourceType: "module",
    allowImportExportEverywhere: true,
    requireConfigFile: false,
  },
  plugins: ["react", "prettier", "jsx-a11y"],
  rules: {
    "prettier/prettier": ["error"],
    "no-multiple-empty-lines": ["error"],
    "no-undef": "off",
    "react/react-in-jsx-scope": "off",
    "react/prop-types": "off",
  },
  settings: {
    react: {
      version: "detect",
    },
  },
};
