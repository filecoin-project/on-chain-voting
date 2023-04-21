const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const Dotenv = require('dotenv-webpack');
const webpack = require('webpack');
const TsconfigPathsPlugin = require('tsconfig-paths-webpack-plugin');


module.exports = {
  target:'web',
  entry: './js/src/index.tsx',
  output: {
    path: path.join(__dirname, '/dist'),
    filename: 'index_bundle.js'
  },
  resolve: {
    extensions: ['.tsx', '.ts', '.js'],
    fallback:{
      "path":require.resolve("path-browserify"),
      "os":require.resolve("os-browserify"),
      "fs":require.resolve("browserify-fs")
    },
    plugins: [new TsconfigPathsPlugin({/* options: see below */})]
  },
  devServer: {
    static: path.join(__dirname, 'public'),
    port: 3000,
    open: true,
    historyApiFallback: true
  },
  stats: {
    warningsFilter: [
      /Failed to parse source map from/,
    ],
  },
  performance: {
    hints: false
  },
  module: {
    rules: [
      {
        test: /\.tsx?$/,
        use: 'ts-loader',
        exclude: /node_modules/,
      },
      {
        test: /\.js$/,
        enforce: 'pre',
        use: ['source-map-loader'],
      },
      {
        test: /\.css$/i,
        use: [
          "style-loader",
          "css-loader",
        ]
      },
      {
        test: /\.less$/i,
        use: [
          "style-loader",
          "css-loader",
          "less-loader"
        ]
      },
    ],
  },
  plugins: [
    new Dotenv({
      path:"./.env"
    }),
    new HtmlWebpackPlugin({
      template: './js/public/index.html',
    }),
    new webpack.ProvidePlugin({
      Buffer: ['buffer', 'Buffer']
    })
    

  ],
  devtool: 'source-map',
};
