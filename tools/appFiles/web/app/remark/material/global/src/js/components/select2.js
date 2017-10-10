$.components.register("select2", {
  mode: "init",
  defaults: {
    width: "style",
    theme: "material"
  },
  init: function(context) {
    if (!$.fn.select2) return;
    var defaults = $.components.getDefaults("select2");

    $('[data-plugin="select2"]', context).each(function() {
      var options = $.extend({}, defaults, $(this).data());

      $(this).select2(options);
    });
  }
});
