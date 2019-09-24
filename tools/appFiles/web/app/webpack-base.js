const path = require('path');
const ExtractTextPlugin = require("extract-text-webpack-plugin");
const webpack = require('webpack');
const Config = require('webpack-config');
const OptimizeCssAssetsPlugin = require("optimize-css-assets-webpack-plugin");

var entries = [];

entries.push(path.join(__dirname, '/javascript/App.js'));
entries.push(path.join(__dirname, '/css/App.css'));
entries.push(path.join(__dirname, '/css/flexbox-examples/build/main.css'));
entries.push(path.join(__dirname, '/node_modules/react-flexgrid-no-bs-conflict/lib/flexgrid.css'));

module.exports = new Config.default().merge({
  entry:  entries,
  plugins: [
    // Allows error warnings but does not stop compiling.
    new webpack.NoEmitOnErrorsPlugin(),
    new OptimizeCssAssetsPlugin({
      assetNameRegExp: /\.css$/g,
      cssProcessor: require('cssnano'),
      cssProcessorOptions: { discardComments: {removeAll: true } },
      canPrint: true
    }),

    new ExtractTextPlugin("./dist/css/go-core-app.css")
  ],
  cache: true,
  module: {
    loaders: [
      {
        test: /\.json$/, loader: "json-loader"
      },
      {
        test: /\.scss$/,
        loader: "style-loader!css-loader!sass-loader"
      },
      {
        test: /\.css$/, loader: ExtractTextPlugin.extract({ fallback: 'style-loader', use: 'css-loader' })
      },
      {
          test: /\.(gif|png)$/i,
          loader: 'file-loader?name=[name].[ext]',
      },
      {
          test: /\.woff(\?v=\d+\.\d+\.\d+)?$/,
          loader: 'url?mimetype=application/font-woff',
      },
      {
          test: /\.woff2(\?v=\d+\.\d+\.\d+)?$/,
          loader: 'url?mimetype=application/font-woff2',
      },
      {
          test: /\.eot(\?v=\d+\.\d+\.\d+)?$/,
          loader: 'file-loader?name=[name].[ext]',
      },
      {
          test: /\.ttf(\?v=\d+\.\d+\.\d+)?$/,
          loader: 'file-loader?name=[name].[ext]',
      },
      {
          test: /\.svg(\?v=\d+\.\d+\.\d+)?$/,
          loader: 'file-loader?name=[name].[ext]',
      }
    ]
  }
});
