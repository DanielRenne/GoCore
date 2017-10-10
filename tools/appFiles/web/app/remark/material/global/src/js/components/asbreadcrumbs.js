$.components.register("breadcrumb", {
  mode: "init",
  defaults: {
    namespace: "breadcrumb",
    dropdown: function() {
      return '<div class=\"dropdown\">' +
        '<a href=\"javascript:void(0);\" data-toggle="dropdown"><i class=\"' + this.dropicon + '\"></i></a>' +
        '<ul class=\"' + this.namespace + '-menu dropdown-menu\" role="menu"></ul>' +
        '</div>';
    },
    dropdownContent: function(value) {
      return '<li><a href=\"javascript:void(0);\">' + value + '</a></li>';
    }
  },
  init: function(context) {
    if (!$.fn.asBreadcrumbs) return;
    var defaults = $.components.getDefaults("breadcrumb");

    $('[data-plugin="breadcrumb"]', context).each(function() {
      var options = $.extend({}, defaults, $(this).data());

      $(this).asBreadcrumbs(options);
    });
  }
});
