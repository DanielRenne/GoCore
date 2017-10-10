var gulp         = require('gulp');
var browserSync = require('browser-sync');
var reload = browserSync.reload;

// config
var config = require('../../config.json');

// options
var options = require('../options/serve');


module.exports = function () {
  browserSync(options.browserSync);

  gulp.watch(options.files).on('change', reload);
};
