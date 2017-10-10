module.exports = function () {
  "use strict";

  return {
    options: {
      // TODO: disable `zeroUnits` optimization once clean-css 3.2 is released
      //    and then simplify the fix for https://github.com/twbs/bootstrap/issues/14837 accordingly
      compatibility: 'ie8',
      keepSpecialComments: '*',
      advanced: false
    },
// @ifdef processCss
    css: {
      files: [
        {
          expand: true,
          cwd: '<%= config.destination.css %>',
          src: ['**/*.css', '!**/*.min.css'],
          dest: '<%= config.destination.css %>',
          ext: '.min.css',
          extDot: 'last'
        }
      ]
    },
// @endif
// @ifdef processSkins
    skins: {
      files: [
        {
          expand: true,
          cwd: '<%= config.destination.skins %>',
          src: ['**/*.css', '!**/*.min.css'],
          dest: '<%= config.destination.skins %>',
          ext: '.min.css',
          extDot: 'last'
        }
      ]
    },
// @endif
// @ifdef processExamples
    examples: {
      files: [
        {
          expand: true,
          cwd: '<%= config.destination.examples %>/css',
          src: ['**/*.css', '!**/*.min.css'],
          dest: '<%= config.destination.examples %>/css',
          ext: '.min.css',
          extDot: 'last'
        }
      ]
    },
// @endif
// @ifdef processFonts
    fonts: {
      files: [
        {
          expand: true,
          cwd: '<%= config.destination.fonts %>',
          src: ['*/*.css', '!*/*.min.css'],
          dest: '<%= config.destination.fonts %>',
          ext: '.min.css',
          extDot: 'last'
        }
      ]
    },
// @endif
// @ifdef processVendor
    vendor: {
      files: [
        {
          expand: true,
          cwd: '<%= config.destination.vendor %>',
          src: ['*/*.css', '!*/*.min.css'],
          dest: '<%= config.destination.vendor %>',
          ext: '.min.css',
          extDot: 'last'
        }
      ]
    }
// @endif
  };
};
