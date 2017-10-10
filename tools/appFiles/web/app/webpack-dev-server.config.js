const path = require('path');
const buildPath = path.resolve(__dirname, 'javascript');
const nodeModulesPath = path.resolve(__dirname, 'node_modules');
const webpack = require('webpack');
const Config = require('webpack-config');

module.exports = new Config.default().extend('webpack-base.js').merge({
    debug: true,
    // Server Configuration options
    devServer: {
      contentBase: 'markup', // Relative directory for base of server
      devtool: 'eval',
      hot: true, // Live-reload
      inline: true,
      port: 3000, // Port Number
      host: '0.0.0.0' // Change to '0.0.0.0' for external facing server
    },

    devtool: 'eval',
    entry: [
      'webpack/hot/dev-server',
      'webpack/hot/only-dev-server'
    ],
    output: {
      pathinfo: true,
      publicPath: 'http://localhost:3000/',
      path: buildPath, // Path of output file
      filename: './dist/javascript/go-core-app.js'
    },
    plugins: [
      // Enables Hot Modules Replacement
      new webpack.HotModuleReplacementPlugin()
    ],
    module: {
      loaders: [
        {
          // React-hot loader and
          test: /\.js$/, // All .js files
          loaders: ['react-hot', 'babel-loader'], // react-hot is like browser sync and babel loads jsx and es6-7
          exclude: [nodeModulesPath]
        }
      ]
    }
});
