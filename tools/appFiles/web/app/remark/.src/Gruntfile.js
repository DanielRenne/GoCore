module.exports = function(grunt) {
  'use strict';

  var path = require('path');

  require('load-grunt-config')(grunt, {
    // path to task.js files, defaults to grunt dir
    configPath: path.join(process.cwd(), 'grunt'),

    // auto grunt.initConfig
    init: true,

    // data passed into config.  Can use with <%= test %>
    data: {
      pkg: grunt.file.readJSON('package.json'),
      config: grunt.file.readJSON('config.json'),
      color: grunt.file.readYAML('color.yml'),
      banner: '/*!\n' +
            ' * <%= pkg.name %> (<%= pkg.homepage %>)\n' +
            ' * Copyright <%= grunt.template.today("yyyy") %> <%= pkg.author.name %>\n' +
            ' * Licensed under the <%= pkg.license %>\n' +
            ' */\n'
    },

    // can optionally pass options to load-grunt-tasks.
    // If you set to false, it will disable auto loading tasks.
    loadGruntTasks: {
      pattern: 'grunt-*',
      config: require('./package.json'),
      scope: ['devDependencies' ,'dependencies']
    }
  });

  // @ifdef processHtml
  // HTML validation task
  grunt.registerTask('validate-html', ['bootlint', 'htmllint']);

  // @endif
  // @ifdef processLint
  // lint task
  grunt.registerTask('lint', [/* @ifdef processCss */'csslint', /* @endif *//* @ifdef processJs */ 'jshint'/* @endif */]);

  // @endif
  // Clean task.
  grunt.registerTask('clean-dist', [/* @ifdef processHtml */'clean:html', /* @endif *//* @ifdef processCss */ 'clean:css', /* @endif *//* @ifdef processJs */ 'clean:js', /* @endif *//* @ifdef processSkins */ 'clean:skins',/* @endif *//* @ifdef processExamples */ 'clean:examples', /* @endif *//* @ifdef processVendor */ 'clean:vendor', /* @endif *//* @ifdef processFonts */ 'clean:fonts'/* @endif */]);

  // @ifdef processHtml
  // Html distribution task.
  grunt.registerTask('dist-html', ['clean:html', 'hb', 'prettify', 'notify:html']);

  // @endif
  // @ifdef processJs
  // JS distribution task.
  grunt.registerTask('dist-js', ['clean:js', 'concat:js', 'uglify:js', 'notify:js']);

  // @endif
  // @ifdef processCss
  // CSS distribution task.
  grunt.registerTask('less-compile', ['less:compileBootstrap', 'less:compileSite', 'less:compileExtend']);
  grunt.registerTask('dist-css', ['clean:css', 'less-compile', 'autoprefixer:css', 'csscomb:css', 'cssmin:css', 'notify:css']);

  // @endif
  // @ifdef processSkins
  // Skins distribution task.
  grunt.registerTask('dist-skins', ['clean:skins', 'less:skins', 'autoprefixer:skins', 'csscomb:skins', 'cssmin:skins', 'notify:skins']);

  // @endif
  // @ifdef processExamples
  // Examples distribution task.
  grunt.registerTask('dist-examples-js',  ['concat:examples', 'uglify:examples']);
  grunt.registerTask('dist-examples-css', ['less:examples', 'autoprefixer:examples', 'csscomb:examples', 'cssmin:examples', 'notify:examples']);
  grunt.registerTask('dist-examples',  ['clean:examples', 'dist-examples-js', 'dist-examples-css']);

  // @endif
  // @ifdef processVendor
  // Vendor distribution task.
  grunt.registerTask('dist-vendor', ['clean:vendor', 'less:vendor', 'autoprefixer:vendor', 'csscomb:vendor', 'cssmin:vendor', 'notify:vendor']);

  // @endif
  // @ifdef processFonts
  // Fonts distribution task.
  grunt.registerTask('dist-fonts', ['clean:fonts', 'less:fonts', 'autoprefixer:fonts', 'csscomb:fonts', 'cssmin:fonts', 'notify:fonts']);

  // @endif
  // Full distribution task.
  grunt.registerTask('dist', [/* @ifdef processHtml */'dist-html', /* @endif *//* @ifdef processCss */ 'dist-css', /* @endif *//* @ifdef processJs */ 'dist-js',/* @endif *//* @ifdef processSkins */ 'dist-skins',/* @endif *//* @ifdef processExamples */ 'dist-examples',/* @endif *//* @ifdef processVendor */ 'dist-vendor',/* @endif *//* @ifdef processFonts */ 'dist-fonts',/* @endif */ 'notify:all']);

  // @ifdef processBrowserSync
  // Static server.
  grunt.registerTask('serve', ['browserSync']);
  // @endif

  grunt.registerTask('default', ['dist']);
};
