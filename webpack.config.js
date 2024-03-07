const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const webpack = require('webpack');
const TsconfigPathsPlugin = require('tsconfig-paths-webpack-plugin');
const CopyWebpackPlugin = require('copy-webpack-plugin');
const dotenv = require('dotenv');

// 根据 NODE_ENV 环境变量加载不同的配置文件
let envFile = '.env';
if (process.env.NODE_ENV === 'development') {
  envFile = '.env.development';
} else if (process.env.NODE_ENV === 'production') {
  envFile = '.env.example';
} else if (process.env.NODE_ENV === 'test') {
  envFile = '.env.test';
}

const envConfig = dotenv.config({ path: envFile }).parsed;

module.exports = {
  target:'web',
  entry: './src/index.tsx',
  output: {
    path: path.resolve(__dirname, 'dist'),
    publicPath: '/',
    filename: 'index_bundle.js'
  },
  resolve: {
    extensions: ['.tsx', '.ts', '.js'],
    fallback:{
      "path":require.resolve("path-browserify"),
      "os":require.resolve("os-browserify"),
      "fs":require.resolve("browserify-fs")
    },
    plugins: [
      new TsconfigPathsPlugin({
        configFile: './tsconfig.json',
      })
    ]
  },
  devServer: {
    static: path.join(__dirname, 'public'),
    port: 3001,
    open: true,
    historyApiFallback: true,
    proxy: {
      '/api': {
        target: 'http://192.168.11.94:9999/power_voting',
        // target: 'http://192.168.3.198:9999/power_voting',
        changeOrigin: true,
        pathRewrite: {
          '^/api': '/api'
        }
      }
    }
  },
  stats: {
    warningsFilter: [
      /Failed to parse source map from/,
      /Critical dependency: the request of a dependency is an expression/
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
          'postcss-loader'
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
    new webpack.ProvidePlugin({
      Buffer: ['buffer', 'Buffer']
    }),
    new CopyWebpackPlugin({
      patterns: [
        {
          from: 'public/images',
          to: 'images'
        }
      ]
    }),
    new HtmlWebpackPlugin({
      template: './public/index.html',
    }),
    new webpack.DefinePlugin({
      'process.env': JSON.stringify(envConfig)
    })
  ],
  devtool: 'source-map',
};