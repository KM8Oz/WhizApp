/* eslint-env node */

module.exports = {
  env: { browser: true, es2020: true },
  extends: [
    'eslint:recommended',
    'plugin:react/recommended',
    'plugin:react/jsx-runtime',
    'plugin:react-hooks/recommended',
  ],
  parserOptions: { ecmaVersion: 'latest', sourceType: 'module' },
  settings: { react: { version: '18.2' } },
  plugins: ['react-refresh'],
  rules: {
    'react-refresh/only-export-components': [
      'warn',
      { allowConstantExport: true },
    ],
  },
  // loaders: [
  //   {
  //     test: /\.(jpe?g|png)$/i,
  //     loaders: [
  //       'file-loader',
  //       'webp-loader?{quality: 13}'
  //     ],
  //     loader: multi(
  //       'file-loader?name=[name].[ext].webp!webp-loader?{quality: 95}'
  //       'file-loader?name=[name].[ext]',
  //     )
  //   }
  // ],
  // "rules": {
  //   "import/no-webpack-loader-syntax": "off"
  // },
}
