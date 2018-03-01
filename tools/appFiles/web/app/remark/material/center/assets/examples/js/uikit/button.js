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

    Waves.attach('.page-content .btn-flat');
    Waves.attach('.page-content .btn-round', ['waves-round', 'waves-light']);
    Waves.attach('.page-content .btn-pure', ['waves-circle', 'waves-classic']);
    Waves.attach('.page-content .btn-floating', ['waves-float', 'waves-light']);
  });

})(document, window, jQuery);
