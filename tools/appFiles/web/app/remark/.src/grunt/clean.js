module.exports = function () {
  "use strict";

  return {
// @ifdef processHtml
    html: '<%= config.html %>',
// @endif
// @ifdef processCss
    css: '<%= config.destination.css %>',
// @endif
// @ifdef processJs
    js: '<%= config.destination.js %>',
// @endif
// @ifdef processSkins
    skins: '<%= config.destination.skins %>/**/*.css',
// @endif
// @ifdef processExamples
    examples: ['<%= config.destination.examples %>/css/**/*.css', '<%= config.destination.examples %>/js/**/*.js'],
// @endif
// @ifdef processFonts
    fonts: '<%= config.destination.fonts %>/*/*.css',
// @endif
// @ifdef processVendor
    vendor: '<%= config.destination.vendor %>/*/*.css'
// @endif
  };
};
