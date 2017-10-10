// config
var config = require('../../config.json');

module.exports = {
  browserSync: {
    server: true,
    startPath: 'html/index.html'
  },
  files: [
    config.destination.css + '/*.css',
    config.destination.js + '/**/*.js',
  ],
};
