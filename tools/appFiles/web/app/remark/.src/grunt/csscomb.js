module.exports = function () {
  "use strict";

  return {
    options: {
      config: '<%= config.source.less %>/.csscomb.json'
    },
// @ifdef processCss
    css: {
      expand: true,
      cwd: '<%= config.destination.css %>',
      src: ['**/*.css', '!**/*.min.css'],
      dest: '<%= config.destination.css %>/'
    },
// @endif
// @ifdef processSkins
    skins: {
      expand: true,
      cwd: '<%= config.destination.skins %>',
      src: ['*.css', '!*.min.css'],
      dest: '<%= config.destination.skins %>'
    },
// @endif
// @ifdef processExamples
    examples: {
      expand: true,
      cwd: '<%= config.destination.examples %>/css',
      src: ['**/*.css', '!**/*.min.css'],
      dest: '<%= config.destination.examples %>/css'
    },
// @endif
// @ifdef processFonts
    fonts: {
      expand: true,
      cwd: '<%= config.destination.fonts %>',
      src: ['*/*.css', '!*/*.min.css'],
      dest: '<%= config.destination.fonts %>'
    },
// @endif
// @ifdef processVendor
    vendor: {
      expand: true,
      cwd: '<%= config.destination.vendor %>',
      src: ['*/*.css', '!*/*.min.css'],
      dest: '<%= config.destination.vendor %>'
    }
// @endif
  };
};
