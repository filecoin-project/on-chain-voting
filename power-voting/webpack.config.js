const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const { DefinePlugin, ProvidePlugin, ProgressPlugin } = require('webpack');
const TsconfigPathsPlugin = require('tsconfig-paths-webpack-plugin');
const CopyWebpackPlugin = require('copy-webpack-plugin');
const dotenv = require('dotenv');

// Load different configuration files based on the NODE_ENV environment variable
const envFile = process.env.NODE_ENV === 'production' ? '.env.example' : '.env';
const envConfig = dotenv.config({ path: envFile }).parsed;

module.exports = {
  target:'web',
  entry: './src/index.tsx',
  output: {
    path: path.resolve(__dirname, 'dist'),
    publicPath: '/',
    filename: '[name].[contenthash].js'
  },
  resolve: {
    extensions: ['.tsx', '.ts', '.js'],
    plugins: [
      new TsconfigPathsPlugin({
        configFile: './tsconfig.json',
      })
    ]
  },
  devServer: {
    static: path.join(__dirname, 'public'),
    port: 4001,
    open: true,
    historyApiFallback: true,
    proxy: {
      '/api': {
        target: 'http://192.168.11.94:10000/power_voting',
        changeOrigin: true,
        pathRewrite: {
          '^/api': '/api'
        }
      },
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
      {
        test: /\.(woff2?|eot|ttf|otf)(\?.*)?$/,
        loader: 'file-loader',
        options: {
          name: '[name].[ext]',
          outputPath: 'fonts/',
        },
      },
    ],
  },
  plugins: [
    new DefinePlugin({
      'process.env': JSON.stringify(envConfig)
    }),
    new ProvidePlugin({
      Buffer: ['buffer', 'Buffer']
    }),
    new CopyWebpackPlugin({
      patterns: [
        {
          from: 'public/images',
          to: 'images'
        },
        {
          from: 'public/fonts',
          to: 'fonts'
        }
      ]
    }),
    new HtmlWebpackPlugin({
      template: './public/index.html',
    }),
    new ProgressPlugin(),
  ],
  devtool: 'source-map',
  // optimization: {
  //   minimize: true,
  // },
};