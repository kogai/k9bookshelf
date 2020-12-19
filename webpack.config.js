const fs = require('fs');
const path = require('path');
const webpack = require('webpack');
const TerserJSPlugin = require('terser-webpack-plugin');
const { CleanWebpackPlugin } = require('clean-webpack-plugin');

module.exports = {
  entry: './storefront/index.ts',
  module: {
    rules: [
      {
        test: /\.tsx?$/,
        use: 'ts-loader',
        exclude: /node_modules/,
      },
    ],
  },
  resolve: {
    extensions: ['.ts'],
  },
  devtool: 'hidden-source-map',
  optimization: {
    minimizer: [new TerserJSPlugin({})],
  },
  mode: process.env.NODE_ENV,
  output: {
    path: path.resolve(__dirname, 'theme','assets'),
    filename: 'app.min.js',
  },
  plugins: [
    new CleanWebpackPlugin(),
    new webpack.DefinePlugin({
      'process.env.BACKEND_ROOT': JSON.stringify(process.env.LOGROCKET_APP_NAME)
    }),
  ],
};
