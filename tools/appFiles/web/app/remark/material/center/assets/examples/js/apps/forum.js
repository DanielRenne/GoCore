/*!
 * remark (http://getbootstrapadmin.com/remark)
 * Copyright 2017 amazingsurge
 * Licensed under the Themeforest Standard Licenses
 */
(function(document, window, $) {
  'use strict';

  window.AppForum = App.extend({
    run: function(next) {
      next();
    }
  });

  $(document).ready(function() {
    AppForum.run();
  });
})(document, window, jQuery);
