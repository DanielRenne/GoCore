/*! jQuery asBreadcrumbs - v0.1.0 - 2016-04-05
* https://github.com/amazingSurge/jquery-asBreadcrumbs
* Copyright (c) 2016 amazingSurge; Licensed GPL */
(function($, document, window, undefined) {
    "use strict";

    var pluginName = 'asBreadcrumbs';

    var Plugin = $[pluginName] = function(element, options) {
        this.element = element;
        this.$element = $(element);

        this.options = $.extend({}, Plugin.defaults, options, this.$element.data());

        // this._plugin = NAME;
        this.namespace = this.options.namespace;

        this.$element.addClass(this.namespace);
        // flag
        this.disabled = false;
        this.initialized = false;
        this.isCreated = false;

        this.$children = this.options.getItem(this.$element);
        this.$firstChild = this.$children.eq(0);

        this.$dropdownWrap = null;
        this.$dropdownMenu = null;

        this.gap = 6;
        this.childrenInfo = [];

        this._trigger('init');
        this.init();
    };

    Plugin.prototype = {
        constructor: Plugin,
        init: function() {
            var self = this;

            this.$element.addClass(this.namespace + '-' + this.options.overflow);

            this.generateChildrenInfo();
            this.createDropdown();

            this.render();

            if (this.options.responsive) {
                $(window).on('resize', this._throttle(function() {
                    self.resize.call(self);
                }, 250));
            }

            this.initialized = true;
            this._trigger('ready');
        },
        _trigger: function(eventType) {
            var method_arguments = Array.prototype.slice.call(arguments, 1),
                data = [this].concat(method_arguments);

            // event
            this.$element.trigger(pluginName + '::' + eventType, data);

            // callback
            eventType = eventType.replace(/\b\w+\b/g, function(word) {
                return word.substring(0, 1).toUpperCase() + word.substring(1);
            });
            var onFunction = 'on' + eventType;
            if (typeof this.options[onFunction] === 'function') {
                this.options[onFunction].apply(this, method_arguments);
            }
        },
        generateChildrenInfo: function() {
            var self = this;

            this.$children.each(function() {
                var $this = $(this);
                self.childrenInfo.push({
                    $this: $this,
                    outerWidth: $this.outerWidth(),
                    $content: $(self.options.dropdownContent($this.text())).attr("href", self.options.getItem($this).attr("href"))
                });
            });
            if (this.options.overflow === "left") {
                this.childrenInfo.reverse();
            }

            this.childrenLength = this.childrenInfo.length;
        },
        createDropdown: function() {
            if (this.isCreated === true) {
                return;
            }

            var dropdown = this.options.dropdown();
            this.$dropdownWrap = this.$firstChild.clone().removeClass().addClass(this.namespace + '-dropdown dropdown ' + this.options.itemClass).html(dropdown).hide();
            this.$dropdownMenu = this.$dropdownWrap.find('.dropdown-menu');

            this._createDropdownItem();

            if (this.options.overflow === 'right') {
                this.$dropdownWrap.appendTo(this.$element);
            } else {
                this.$dropdownWrap.prependTo(this.$element);
            }

            this._createEllipsis();

            this.isCreated = true;
        },
        render: function() {
            var dropdownWidth = this.getDropdownWidth(),
                childrenWidthTotal = 0,
                childWidth = 0,
                width = 0;

            for (var i = 0, l = this.childrenLength; i < l; i++) {

                width = this.getWidth();
                childWidth = this.childrenInfo[i].outerWidth;

                childrenWidthTotal += childWidth;

                if (childrenWidthTotal + dropdownWidth > width) {
                    this._showDropdown(i);
                } else {
                    this._hideDropdown(i);
                }
            }
        },
        resize: function() {
            this._trigger('resize');

            this.render();
        },
        getDropdownWidth: function() {
            return this.$dropdownWrap.outerWidth() + (this.options.ellipsis ? this.$ellipsis.outerWidth() : 0);
        },
        getWidth: function() {
            var width = 0,
                self = this;

            this.$element.children().each(function() {
                if ($(this).css('display') === 'inline-block' && $(this).css('float') === 'none') {
                    width += self.gap;
                }
            });
            return this.$element.width() - width;
        },
        _createEllipsis: function() {
            if (!this.options.ellipsis) {
                return;
            }

            this.$ellipsis = this.$firstChild.clone().removeClass().addClass(this.namespace + '-ellipsis ' + this.options.itemClass).html(this.options.ellipsis);

            if (this.options.overflow === 'right') {
                this.$ellipsis.insertBefore(this.$dropdownWrap).hide();
            } else {
                this.$ellipsis.insertAfter(this.$dropdownWrap).hide();
            }
        },
        _createDropdownItem: function() {
            for (var i = 0, l = this.childrenLength; i < l; i++) {
                this.childrenInfo[i].$content.appendTo(this.$dropdownMenu).hide();
            }
        },
        _showDropdown: function(i) {
            this.childrenInfo[i].$content.show();
            this.childrenInfo[i].$this.hide();
            this.$dropdownWrap.show();
            this.$ellipsis.css("display", "inline-block");
        },
        _hideDropdown: function(i) {
            this.childrenInfo[i].$this.css("display", "inline-block");
            this.childrenInfo[i].$content.hide();
            this.$dropdownWrap.hide();
            this.$ellipsis.hide();
        },
        _throttle: function(func, wait) {
            var _now = Date.now || function() {
                return new Date().getTime();
            };
            var context, args, result;
            var timeout = null;
            var previous = 0;
            var later = function() {
                previous = _now();
                timeout = null;
                result = func.apply(context, args);
                context = args = null;
            };
            return function() {
                var now = _now();
                var remaining = wait - (now - previous);
                context = this;
                args = arguments;
                if (remaining <= 0) {
                    clearTimeout(timeout);
                    timeout = null;
                    previous = now;
                    result = func.apply(context, args);
                    context = args = null;
                } else if (!timeout) {
                    timeout = setTimeout(later, remaining);
                }
                return result;
            };
        },
        destroy: function() {
            // detached events first
            // then remove all js generated html
            this.$element.children().css("display", "");
            this.$dropdownWrap.remove();
            if (this.options.ellipsis) {
                this.$ellipsis.remove();
            }
            this.isCreated = false;
            this.$element.data(pluginName, null);
            $(window).off("resize");
            $(window).off(".asBreadcrumbs");
            this._trigger('destroy');
        }
    };

    Plugin.defaults = {
        namespace: pluginName,
        overflow: "left",
        ellipsis: "&#8230;",
        dropicon: "caret",
        responsive: true,
        itemClass: "",

        dropdown: function() {
            return '<div class=\"dropdown\">' +
                '<a href=\"javascript:void(0);\" class=\"' + this.namespace + '-toggle\" data-toggle=\"dropdown\"><i class=\"' + this.dropicon + '\"></i></a>' +
                '<div class=\"' + this.namespace + '-menu dropdown-menu\"></div>' +
                '</div>';
        },
        dropdownContent: function(value) {
            return '<a class=\"dropdown-item\">' + value + '</a>';
        },
        getItem: function($parent) {
            return $parent.children();
        },

        // callback
        onInit: null,
        onReady: null
    };

    $.fn[pluginName] = function(options) {
        if (typeof options === 'string') {
            var method = options;
            var method_arguments = Array.prototype.slice.call(arguments, 1);

            if (/^\_/.test(method)) {
                return false;
            } else if ((/^(get)/.test(method))) {
                var api = this.first().data(pluginName);
                if (api && typeof api[method] === 'function') {
                    return api[method].apply(api, method_arguments);
                }
            } else {
                return this.each(function() {
                    var api = $.data(this, pluginName);
                    if (api && typeof api[method] === 'function') {
                        api[method].apply(api, method_arguments);
                    }
                });
            }
        } else {
            return this.each(function() {
                if (!$.data(this, pluginName)) {
                    $.data(this, pluginName, new Plugin(this, options));
                }
            });
        }
    };
})(jQuery, document, window);
