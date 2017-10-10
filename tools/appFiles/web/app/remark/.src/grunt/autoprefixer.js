module.exports = function () {
  "use strict";

  return {
    options: {
      browsers: '<%= config.autoprefixerBrowsers %>'
    },
// @ifdef processCss
    css: {
      options: {
        map: true
      },
      src: '<%= config.destination.css %>/**/*.css'
    },
// @endif
// @ifdef processSkins
    skins: {
      options: {
        map: false
      },
      src: ['<%= config.destination.skins %>/*.css', '!<%= config.destination.skins %>/*.min.css']
    },
// @endif
// @ifdef processExamples
    examples: {
      options: {
        map: false
      },
      src: ['<%= config.destination.examples %>/**/*.css', '!<%= config.destination.examples %>/**/*.min.css']
    },
// @endif
// @ifdef processFonts
    fonts: {
      options: {
        map: false
      },
      src: ['<%= config.destination.fonts %>/*/*.css', '!<%= config.destination.fonts %>/*/*.min.css']
    },
// @endif
// @ifdef processVendor
    vendor: {
      options: {
        map: false
      },
      src: ['<%= config.destination.vendor %>/*/*.css', '!<%= config.destination.vendor %>/*/*.min.css']
    }
// @endif
  };
};
