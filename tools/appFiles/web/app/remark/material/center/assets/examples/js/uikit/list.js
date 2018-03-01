/*!
 * remark (http://getbootstrapadmin.com/remark)
 * Copyright 2017 amazingsurge
 * Licensed under the Themeforest Standard Licenses
 */
(function(document, window, $) {
  'use strict';

  var Site = window.Site;

  $(document).ready(function($) {
    Site.run();

    Waves.attach('.list-group > a:not(.disabled)', ['waves-block', 'waves-classic']);
  });

})(document, window, jQuery);
