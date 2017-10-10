module.exports = function () {
  "use strict";

  return {
    options: {
      strictMath: false,
      paths: [
        '<%= config.source.less %>',
        '<%= config.bootstrap.less %>',
        '<%= config.bootstrap.mixins %>'
      ]
    },
// @ifdef processCss
    compileBootstrap: {
      options: {
        strictMath: true
      },
      src: '<%= config.source.less %>/bootstrap.less',
      dest: '<%= config.destination.css %>/bootstrap.css'
    },
    compileExtend: {
      options: {
        strictMath: true
      },
      src: '<%= config.source.less %>/bootstrap-extend.less',
      dest: '<%= config.destination.css %>/bootstrap-extend.css'
    },
    compileSite: {
      options: {
        strictMath: true
      },
      src: '<%= config.source.less %>/site.less',
      dest: '<%= config.destination.css %>/site.css'
    },
// @endif
// @ifdef processSkins
    skins: {
      options: {
        strictMath: true,
        paths: [
          '<%= config.source.skins %>/less',
          '<%= config.source.less %>',
          '<%= config.bootstrap.less %>',
          '<%= config.bootstrap.mixins %>'
        ]
      },
      expand: true,
      cwd: '<%= config.source.skins %>',
      src: ['*.less'],
      dest: '<%= config.destination.skins %>',
      ext: '.css',
      extDot: 'last'
    },
// @endif
// @ifdef processExamples
    examples: {
      expand: true,
      cwd: '<%= config.source.examples %>/less',
      src: ['**/*.less'],
      dest: '<%= config.destination.examples %>/css',
      ext: '.css',
      extDot: 'last'
    },
// @endif
// @ifdef processFonts
    fonts: {
      expand: true,
      cwd: '<%= config.source.fonts %>',
      src: ['*/*.less', '!*/_*.less'],
      dest: '<%= config.destination.fonts %>',
      ext: '.css',
      extDot: 'last'
    },
// @endif
// @ifdef processVendor
    vendor: {
      expand: true,
      cwd: '<%= config.source.vendor %>',
      src: ['*/*.less', '!*/settings.less'],
      dest: '<%= config.destination.vendor %>',
      ext: '.css',
      extDot: 'last'
    },
// @endif
  };
};
