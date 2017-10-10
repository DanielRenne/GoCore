/**
 * Js
 * configuration
 * object
 *
 */
// config
var config = require('../../config.json');
var pkg = require('../../package.json');
var banner = '/*!\n' +
          ' * ' + pkg.name + ' (' + pkg.homepage + ')\n' +
          ' * Copyright ' + new Date().getFullYear() + ' ' + pkg.author.name + '\n' +
          ' * Licensed under the ' + pkg.license + '\n' +
          ' */\n';

module.exports = {
  banner: banner,
  less: {
    strictMath: false,
    paths: [
      config.source.less,
      config.bootstrap.less,
      config.bootstrap.mixins
    ]
  },
  autoprefixer: {
    browsers: config.autoprefixerBrowsers
  },
  csscomb: {
    configPath: '.csscomb.json'
  },
  cleanCSS: {
    compatibility: 'ie8',
    keepSpecialComments: '*',
    advanced: false
  }
};
