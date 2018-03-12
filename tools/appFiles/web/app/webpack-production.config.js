const path = require('path');
const buildPath = path.resolve(__dirname, 'javascript');
const nodeModulesPath = path.resolve(__dirname, 'node_modules');
const ExtractTextPlugin = require("extract-text-webpack-plugin");
const webpack = require('webpack');
const Config = require('webpack-config');


module.exports = new Config.default().extend('webpack-base.js').merge({
    devtool: 'source-map',

    output: {
      path: buildPath, // Path of output file
      filename: 'go-core-app.js' // Name of output file
    },
    module: {
      loaders: [
        {
          test: /\.js$/, // All .js files
          loaders: ['babel-loader'], // react-hot is like browser sync and babel loads jsx and es6-7
          exclude: [nodeModulesPath]
        }
      ]
    },
    plugins: [
      new webpack.DefinePlugin({
        "process.env": {
           NODE_ENV: JSON.stringify("production")
         }
      }),
      new webpack.optimize.UglifyJsPlugin({
            compress: {
                warnings: false
            }
      })
    ]
});