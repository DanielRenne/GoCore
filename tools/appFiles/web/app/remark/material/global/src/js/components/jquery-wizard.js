$.components.register("wizard", {
  mode: "default",
  defaults: {
    step: ".steps .step, .pearls .pearl",
    templates: {
      buttons: function() {
        var options = this.options;
        return '<div class="wizard-buttons">' +
          '<a class="btn btn-default" href="#' + this.id + '" data-wizard="back" role="button">' + options.buttonLabels.back + '</a>' +
          '<a class="btn btn-primary pull-right" href="#' + this.id + '" data-wizard="next" role="button">' + options.buttonLabels.next + '</a>' +
          '<a class="btn btn-success pull-right" href="#' + this.id + '" data-wizard="finish" role="button">' + options.buttonLabels.finish + '</a>' +
          '</div>';
      }
    }
  }
});
