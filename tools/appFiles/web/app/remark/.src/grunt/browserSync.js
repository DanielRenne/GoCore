module.exports = function () {
  "use strict";

  return {
    bsFiles: {
      src : [
        '<%= config.destination.css %>/*.css',
        '<%= config.destination.js %>/**/*.js',
      ]
    },
    options: {
      server: true,
      startPath: 'html/index.html'
    }
  };
};
