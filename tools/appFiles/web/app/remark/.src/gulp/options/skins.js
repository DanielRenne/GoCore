// config
var config = require('../../config.json');
var styles = require('./styles');

styles.less.paths = [
  config.source.skins + '/less',
  config.source.less,
  config.bootstrap.less,
  config.bootstrap.mixins
];

module.exports = styles;
