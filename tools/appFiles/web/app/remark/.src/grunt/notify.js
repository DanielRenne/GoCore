module.exports = function () {
  "use strict";

  return {
    options: {
      enabled: true,
      duration: 2
    },
// @ifdef processHtml
    html: {
      options: {
        message: 'Html Generated!'
      }
    },
// @endif
// @ifdef processJs
    js: {
      options: {
        message: 'JS Generated!'
      }
    },
// @endif
// @ifdef processCss
    css: {
      options: {
        message: 'CSS Generated!'
      }
    },
// @endif
// @ifdef processSkins
    skins: {
      options: {
        message: 'Skins Generated!'
      }
    },
// @endif
// @ifdef processExamples
    examples: {
      options: {
        message: 'Examples Generated!'
      }
    },
// @endif
// @ifdef processFonts
    fonts: {
      options: {
        message: 'Fonts Generated!'
      }
    },
// @endif
// @ifdef processVendor
    vendor: {
      options: {
        message: 'Vendor Generated!'
      }
    },
// @endif
    all: {
      options: {
        message: 'All Generated!'
      }
    }
  };
};
