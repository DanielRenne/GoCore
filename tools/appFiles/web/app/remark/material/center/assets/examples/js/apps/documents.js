/*!
 * remark (http://getbootstrapadmin.com/remark)
 * Copyright 2017 amazingsurge
 * Licensed under the Themeforest Standard Licenses
 */
(function(document, window, $) {
  'use strict';
  window.AppDocuments = App.extend({
    affixHandle: function() {
      $('#articleAffix').affix({
        offset: {
          top: 210
        }
      });
    },
    scrollHandle: function() {
      $('body').scrollspy({
        target: '#articleAffix'
      });
    },
    run: function(next) {
      this.scrollHandle();
      this.affixHandle();

      next();
    }
  });

  $(document).ready(function() {
    AppDocuments.run();
  });
})(document, window, jQuery);
