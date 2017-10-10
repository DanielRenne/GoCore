(function(document, window, $) {
  'use strict';

  window.AppMailbox = App.extend({
    handleAction: function() {
      var actionBtn = $('.site-action').actionBtn().data('actionBtn');
      var $selectable = $('[data-selectable]');

      $('.site-action-toggle', '.site-action').on('click', function(e) {
        var $selected = $selectable.asSelectable('getSelected');

        if ($selected.length === 0) {
          $('#addMailForm').modal('show');
          e.stopPropagation();
        }
      });

      $('[data-action="trash"]', '.site-action').on('click', function() {
        console.log('trash');
      });

      $('[data-action="inbox"]', '.site-action').on('click', function() {
        console.log('folder');
      });

      $selectable.on('asSelectable::change', function(e, api, checked) {
        if (checked) {
          actionBtn.show();
        } else {
          actionBtn.hide();
        }
      });
    },

    handleListItem: function() {
      $('#addLabelToggle').on('click', function(e) {
        $('#addLabelForm').modal('show');
        e.stopPropagation();
      });

      $(document).on('click', '[data-tag=list-delete]', function(e) {
        bootbox.dialog({
          message: "Do you want to delete the label?",
          buttons: {
            success: {
              label: "Delete",
              className: "btn-danger",
              callback: function() {
                // $(e.target).closest('.list-group-item').remove();
              }
            }
          }
        });
      });
    },


    itemTpl: function(data) {
      return '<tr id="' + data.id + '" data-mailbox="slidePanel" ' + (data.unread === 'true' ? 'class="unread"' : '') + '>' +
        '<td class="cell-60">' +
        '<span class="checkbox-custom checkbox-primary checkbox-lg">' +
        '<input type="checkbox" class="mailbox-checkbox selectable-item" id="mail_' + data.id + '"/>' +
        '<label for="mail_' + data.id + '"></label>' +
        '</span>' +
        '</td>' +
        '<td class="cell-30 responsive-hide">' +
        '<span class="checkbox-important checkbox-default">' +
        '<input type="checkbox" class="mailbox-checkbox mailbox-important" ' + (data.starred === 'true' ? 'checked="checked"' : '') + ' id="mail_' + data.id + '_important"/>' +
        '<label for="mail_' + data.id + '_important"></label>' +
        '</span>' +
        '</td>' +
        '<td class="cell-60 responsive-hide">' +
        '<a class="avatar" href="javascript:void(0)"><img class="img-responsive" src="' + data.avatar + '" alt="..."></a>' +
        '</td>' +
        '<td>' +
        '<div class="content">' +
        '<div class="title">' + data.name + '</div>' +
        '<div class="abstract">' + data.title + '</div>' +
        '</div>' +
        '</td>' +
        '<td class="cell-30 responsive-hide">' +
        (data.attachments.length > 0 ? '<i class="icon md-attachment-alt" aria-hidden="true"></i>' : '') +
        '</td>' +
        '<td class="cell-130">' +
        '<div class="time">' + data.time + '</div>' +
        (data.group.length > 0 ? '<div class="identity"><i class="md-circle ' + data.color + '" aria-hidden="true"></i>' + data.group + '</div>' : '') +
        '</td>' +
        '</tr>';
    },

    attachmentsTpl: function(data) {
      var self = this,
        html = '';

      html += '<div class="mail-attachments">' +
        '<p><i Class="icon md-attachment-alt"></i>Attachments | <a href="javascript:void(0)">Download All</a></p>' +
        '<ul class="list-group">';

      $.each(data, function(n, item) {
        html += self.attachmentTpl(item);
      });

      html += '</ul></div>';

      return html;
    },

    attachmentTpl: function(data) {
      return '<li class="list-group-item">' +
        '<span class="name">' + data.name + '</span><span class="size">' + data.size + '</span>' +
        '<button type="button" class="btn btn-icon btn-pure btn-default"><i class="icon md-download" aria-hidden="true"></i></button>' +
        '</li>';
    },

    messagesTpl: function(data) {
      var self = this,
        html = '';

      $.each(data.messages, function(n, item) {
        html += '<section class="slidePanel-inner-section">' +
          '<div class="mail-header">' +
          '<div class="mail-header-main">' +
          '<a class="avatar" href="javascript:void(0)"><img src="' + data.avatar + '" alt="..."></a>' +
          '<div><span class="name">' + data.name + '</span></div>' +
          '<div>' +
          '<a href="javascript:void(0)" class="mailbox-panel-email">' + data.email + '</a>' +
          ' to <a href="javascript:void(0)" class="margin-right-10">me</a>' +
          '<span class="identity"><i class="md-circle red-600" aria-hidden="true"></i>' + data.group + '</span>' +
          '</div>' +
          '</div>' +
          '<div class="mail-header-right">' +
          '<span class="time">' + item.time + '</span>' +
          '<div class="btn-group btn-group-flat actions" role="group">' +
          '<button type="button" class="btn btn-icon btn-pure btn-default"><i class="icon md-star" aria-hidden="true"></i></button>' +
          '<button type="button" class="btn btn-icon btn-pure btn-default"><i class="icon md-mail-reply" aria-hidden="true"></i></button>' +
          '</div>' +
          '</div>' +
          '</div>' +
          '<div class="mail-content">' + item.content + '</div>';

        if (n === 0) {
          if (item.attachments && item.attachments.length > 0) {
            html += this.attachmentsTpl(item.attachments);
          }
        }

        html += '</section>';
      });

      return html;
    },

    initMail: function() {
      var self = this;

      $.getJSON('../../../assets/data/appsMailbox.json', function(data) {
        var $wrap = $('#mailboxTable');

        self.buildMail($wrap, data);
        self.initMailData(data);
        self.handlSlidePanelPlugin();
      });
    },

    initMailData: function(data) {
      this.mailboxData = data;
    },

    buildMail: function($wrap, data) {
      var self = this,
        $tbody = $('<tbody></tbody>');

      $.each(data, function(i, item) {
        self.buildItem($tbody, item);
      });

      $wrap.empty().append($tbody);
    },

    buildItem: function($wrap, data) {
      $wrap.append($(this.itemTpl(data)).data('mailInfo', data));
    },

    buildPanel: function() {

    },

    filter: function(flag, value) {

    },

    handlePanel: function() {
      $(document).on('click', '[data-mailbox="slidePanel"]', function(e) {
        console.log(this, $(this))
      });
    },

    handlSlidePanelPlugin: function() {
      if (typeof $.slidePanel === 'undefined') return;

      var self = this;
      var defaults = $.components.getDefaults("slidePanel");
      var options = $.extend({}, defaults, {
        template: function(options) {
          return '<div class="' + options.classes.base + ' ' + options.classes.base + '-' + options.direction + '">' +
            '<div class="' + options.classes.base + '-scrollable"><div>' +
            '<div class="' + options.classes.content + '"></div>' +
            '</div></div>' +
            '<div class="' + options.classes.base + '-handler"></div>' +
            '</div>';
        },
        afterLoad: function(object) {
          var _this = this,
            $target = $(object.target),
            info = $target.data('taskInfo');

          this.$panel.find('.' + this.options.classes.base + '-scrollable').asScrollable({
            namespace: 'scrollable',
            contentSelector: '>',
            containerSelector: '>'
          });
        },
        contentFilter: function(data, object) {
          var $target = $(object.target),
            info = $target.data('mailInfo'),
            $panel = $(data);

          $('.mailbox-panel-title', $panel).html(info.title);

          $('.slidePanel-messages', $panel).html(self.messagesTpl(info));

          return $panel;
        }
      });

      $(document).on('click', '[data-mailbox="slidePanel"]', function(e) {
        $.slidePanel.show({
          url: 'panel.tpl',
          target: $(this)
        }, options);

        e.stopPropagation();
      });
    },


    run: function(next) {
      this.handleAction();
      this.handleListItem();

      this.initMail();

      $('#addlabelForm').modal({
        show: false
      });

      $('#addMailForm').modal({
        show: false
      });

      $('.checkbox-important').on('click', function(e) {
        e.stopPropagation();
      });

      this.handleMultiSelect();
      next();
    }
  });

  $(document).ready(function() {
    AppMailbox.run();
  });
})(document, window, jQuery);
