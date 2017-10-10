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
  });

  // Example editable Table
  // ----------------------
  $('#editableTable').editableTableWidget().numericInputExample().find('td:first').focus();

})(document, window, jQuery);
