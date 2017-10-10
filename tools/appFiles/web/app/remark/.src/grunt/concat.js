module.exports = function (grunt) {
  "use strict";

// @ifdef processJs
  var components = grunt.file.readJSON('components.json');
  var componentsSrc = [];

  for(var component in components) {
    if(components[component]){
      componentsSrc.push('<%= config.source.js %>/components/'+component+'.js');
    }
  }
// @endif

  return {
    options: {
      banner: '<%= banner %>',
      stripBanners: false
    },
// @ifdef processJs
    js: {
      expand: true,
      cwd: '<%= config.source.js %>',
      src: ['**/*.js'],
      dest: '<%= config.destination.js %>',
    },
    components: {
      src: componentsSrc,
      dest: '<%= config.destination.js %>/components.js'
    },
// @endif
// @ifdef processExamples
    examples: {
      expand: true,
      cwd: '<%= config.source.examples %>/js',
      src: ['**/*.js'],
      dest: '<%= config.destination.examples %>/js',
    }
// @endif
  };
};
