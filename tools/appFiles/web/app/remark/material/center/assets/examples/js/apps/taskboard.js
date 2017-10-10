/*!
 * remark (http://getbootstrapadmin.com/remark)
 * Copyright 2017 amazingsurge
 * Licensed under the Themeforest Standard Licenses
 */
(function(document, window, $) {
  'use strict';

  window.AppTaskboard = App.extend({
    //TPL
    stageTpl: function(title) {
      return '<li class="taskboard-stage">' +
        '<header class="taskboard-stage-header">' +
        '<div class="taskboard-stage-actions pull-right">' +
        '<div class="dropdown">' +
        '<a class="dropdown-toggle" data-toggle="dropdown" href="#" aria-expanded="false"><i class="icon md-chevron-down" aria-hidden="true"></i></a>' +
        '<ul class="dropdown-menu bullet" role="menu">' +
        '<li role="presentation" class="taskboard-stage-rename"><a href="javascript:void(0)" role="menuitem"><i class="icon md-edit" aria-hidden="true"></i>Rename</a></li>' +
        '<li role="presentation" class="taskboard-stage-delete" ><a href="javascript:void(0)" role="menuitem"><i class="icon md-delete" aria-hidden="true"></i>Delete</a></li>' +
        '<li class="taskboard-stage-rename-wrap">' +
        '<div class="form-group">' +
        '<input class="form-control taskboard-stage-rename-input" type="text" value="' + title + '" name="name">' +
        '</div>' +
        '<button class="btn btn-primary btn-block taskboard-stage-rename-save" type="button">Save</button>' +
        '</li>' +
        '</ul>' +
        '</div>' +
        '</div>' +
        '<h5 class="taskboard-stage-title">' + title + '</h5>' +
        '</header>' +
        '<div class="taskboard-stage-content">' +
        '<ul class="list-group taskboard-list">' +
        '</ul>' +
        '<div class="action-wrap">' +
        '<a class="add-item-toggle" href="#"><i class="icon md-plus" aria-hidden="true"></i>Add Task</a>' +
        '<div class="add-item-wrap">' +
        '<form class="add-item" role="form" method="post" action="#">' +
        '<div class="form-group">' +
        '<label class="control-label margin-bottom-15" for="name">Task name:</label>' +
        '<input class="form-control" type="text" placeholder="Task name" name="name">' +
        '</div>' +
        '<div class="form-group text-right">' +
        '<a class="btn btn-sm btn-white btn-flat add-item-cancel">Cancel</a>' +
        '<button type="button" class="btn btn-primary add-item-add">Add</button>' +
        '</div>' +
        '</form>' +
        '</div>' +
        '</div>' +
        '</div>' +
        '</li>';
    },

    taskTpl: function(data) {
      return '<li class="list-group-item priority-' + data.priority + '" data-taskboard="slidePanel" data-url="panel.tpl">' +
        '<div class="checkbox-custom checkbox-primary">' +
        '<input type="checkbox" ' + (data.complete ? 'checked="checked"' : '') + ' name="checkbox">' +
        '<label class="task-title">' + data.title + '</label>' +
        '</div>' +
        '<div class="task-badges"></div>' +
        '<ul class="task-members">' +
        '<li><img class="avatar avatar-sm" src="../../../../global/portraits/5.jpg"></li>' +
        '</div>' +
        '</li>';
    },

    badgesTpl: function(type, content) {
      var html = '';
      switch (type) {
        case 'duedate':
          html = '<span class="task-badge task-badge-subtask icon md-calendar">' + content + '</span>';
          break;
        case 'subtasks':
          html = '<span class="task-badge task-badge-subtask icon md-format-list-bulleted">' + content + '</span>';
          break;
        case 'attachments':
          html = '<span class="task-badge task-badge-attachments icon md-attachment-alt">' + content + '</span>';
          break;
        case 'comments':
          html = '<span class="task-badge task-badge-comments icon md-comment">' + content + '</span>';
          break;
      }
      return html;
    },

    membersTpl: function(src) {
      return '<li><img class="avatar avatar-sm" src="' + src + '"></li>';
    },

    subtaskTpl: function(data) {
      return '<li class="list-group-item subtask">' +
        '<div class="checkbox-custom checkbox-primary">' +
        '<input type="checkbox" ' + (data.complete ? 'checked="checked"' : '') + ' name="checkbox">' +
        '<label class="title">' + data.title + '</label>' +
        '</div>' +
        '<div class="subtask-editor">' +
        '<form>' +
        '<div class="form-group">' +
        '<input class="form-control subtask-title" type="text" name="title">' +
        '</div>' +
        '<div class="form-group">' +
        '<button class="btn btn-primary subtask-editor-save" type="button">Save</button>' +
        '<a class="btn btn-sm btn-white subtask-editor-delete" href="javascript:void(0)">Delete</a>' +
        '</div>' +
        '</form>' +
        '</div>' +
        '</li>';
    },

    attachmentTpl: function(data) {
      return '<li class="list-group-item">' +
        '<div class="meida">' +
        '<div class="media-left">' +
        '<div class="attachments-image">' +
        '<img src="' + data.src + '">' +
        '</div>' +
        '</div>' +
        '<div class="media-body">' +
        '<p><span class="name">' + data.title + '</span><span</p>' +
        '<p>' +
        '<span class="size">' + data.size + '</span>' +
        '<span class="attachments-actions">' +
        '<button class="btn btn-icon btn-pure" type="button">' +
        '<i class="icon md-download" aria-hidden="true"></i>' +
        '</button>' +
        '<button class="btn btn-icon btn-pure" type="button">' +
        '<i class="icon md-delete" aria-hidden="true"></i>' +
        '</button>' +
        '</span>' +
        '</p>' +
        '</div>' +
        '</div>' +
        '</li>';
    },

    commentTpl: function(src, user, time, content) {
      return '<div class="comment media">' +
        '<div class="media-left">' +
        '<a class="avatar avatar-lg" href="javascript:void(0)">' +
        '<img src="' + src + '" alt="...">' +
        '</a>' +
        '</div>' +
        '<div class="media-body">' +
        '<div class="comment-body">' +
        '<a class="comment-author" href="javascript:void(0)">' + user + '</a>' +
        '<div class="comment-meta">' +
        '<span class="date">' + time + '</span>' +
        '</div>' +
        '<div class="comment-content"><p>' + content + '</p></div>' +
        '</div>' +
        '</div>' +
        '</div>';
    },

    dataTpl: function() {
      var data = {
        "status": false,
        "title": "",
        "description": "",
        "priority": "normal",
        "duedate": "",
        "members": [],
        "subtasks": [],
        "attachments": [],
        "comments": []
      };
      return data;
    },

    //Init Page
    init: function() {
      var self = this;
      $.getJSON('../../../assets/data/taskboard.json', function(data) {
        var $wrap = $('#taskboard-stages');
        self.buildStage($wrap, data);
        self.initSortable();
      });
    },

    buildStage: function($wrap, data) {
      if (data.length === 0) return;

      var self = this;
      $.each(data, function(n, info) {
        var $stage = $(self.stageTpl(info.title, info.type));
        self.buildTask($stage, info.tasks);
        $wrap.append($stage);
      });
    },

    buildTask: function($wrap, data, once) {
      if (data.length === 0) return;

      var self = this,
        $container = $('.taskboard-list', $wrap);
      if (once) {
        var $task = $(self.taskTpl(data));
        self.buildBadges($task, data);
        $task.data('taskInfo', data);
        $wrap.append($task);
      } else {
        $.each(data, function(n, info) {
          var $task = $(self.taskTpl(info));
          self.buildBadges($task, info);
          self.buildMembers($task, info.members);
          $task.data('taskInfo', info);
          $container.append($task);
        });
      }
    },

    buildBadges: function($wrap, data) {
      var self = this,
        html = '',
        duedate = data.duedate,
        subtasks = data.subtasks,
        attachments = data.attachments,
        comments = data.comments;

      if (duedate.length > 0) {
        html += self.badgesTpl('duedate', duedate.split(/\//, 2).join("/"));
      }

      if (subtasks.length > 0) {
        var num = 0;
        $.each(subtasks, function(n, i) {
          if (i.complete) num++;
        });

        html += self.badgesTpl('subtasks', num + '/' + subtasks.length);
      }

      if (attachments.length > 0) {
        html += self.badgesTpl('attachments', attachments.length);
      }

      if (comments.length > 0) {
        html += self.badgesTpl('comments', comments.length);
      }

      $wrap.find('.task-badges').html(html);
    },

    buildMembers: function($wrap, data) {
      var self = this,
        html = '';
      if (data.length === 0) return;
      $.each(data, function(i, n) {
        html += self.membersTpl(n.img);
      });
      $wrap.find('.task-members').html(html);
    },

    //Sortable
    initSortable: function() {
      $('.taskboard-stages').sortable({
        handle: ".taskboard-stage-header"
      });
      $('.taskboard-stage .list-group').sortable({
        connectWith: ".taskboard-stage .list-group"
      });
    },

    //Stage
    handleAddStage: function() {
      var self = this;

      $(document).on('click', '.site-floataction', function() {
        var $model = $('#addStageFrom');

        $('input', $model).val('');
        $('option:first', $('select', $model)).prop("selected", 'selected');
      });

      $(document).on('click', '#taskboard-stage-creat', function() {
        var $this = $(this),
          $model = $('#addStageFrom'),
          $name = $('[name="name"]', $model);

        $('.taskboard-stages').append(self.stageTpl($name.val()));
        self.initSortable();
      });
    },

    handleDeleteStage: function() {
      $(document).on('click', '.taskboard-stage-delete', function() {
        var $this = $(this);
        bootbox.dialog({
          message: "Do you want to delete the stage?",
          buttons: {
            success: {
              label: "Delete",
              className: "btn-danger",
              callback: function() {
                $this.closest('.taskboard-stage').remove();
              }
            }
          }
        });
      });
    },

    getStage: function($task) {
      return $task.closest('.taskboard-stage');
    },

    //Stage Dropdown
    initStageDropdown: function() {
      $(document).on('click', '.taskboard-stage-actions .dropdown-toggle', function() {
        $(this).next('.dropdown-menu').removeClass('is-edit');

        //judge dropdown side
      });
    },

    handleStageRename: function() {
      $(document).on('click', '.taskboard-stage-rename', function(e) {
        var $header = $(this).closest('.taskboard-stage-header'),
          $menu = $(this).closest('.dropdown-menu'),
          $input = $('.taskboard-stage-rename-input', $menu),
          $title = $('.taskboard-stage-title', $header);

        $menu.toggleClass('is-edit');
        $input.val('').focus().val($title.html());
        e.stopPropagation();
      });

      $(document).on('click', '.taskboard-stage-rename-save', function() {
        var $header = $(this).closest('.taskboard-stage-header'),
          $input = $('.taskboard-stage-rename-input', $header),
          $title = $('.taskboard-stage-title', $header),
          value = $input.val();

        if (value.length === 0) return;

        $title.html(value);
      });
    },

    //Task
    handleAddTask: function() {
      var self = this;

      $(document).on('click', '.add-item-toggle, .add-item-add, .add-item-cancel', function() {
        var $this = $(this),
          $wrap = $this.closest('.action-wrap'),
          $input = $('[name="name"]', $wrap);

        $wrap.toggleClass('action-open');
        if ($this.hasClass('add-item-toggle')) {
          $input.val('');
        }

        if ($this.hasClass('add-item-toggle')) {
          $(document).on('click.add-item', function(e) {
            var $target = $(e.target);
            if ($target.closest('.add-item-wrap').length === 0) {
              $wrap.removeClass('action-open');
              $(document).off('click.add-item');
            }
          });
        } else {
          $(document).off('click.add-item');
        }
      });

      $(document).on('click', '.add-item-add', function() {
        var $this = $(this),
          $wrap = $this.closest('.action-wrap'),
          $input = $('[name="name"]', $wrap),
          $list = $('.taskboard-list', $this.closest('.taskboard-stage-content')),
          data = self.dataTpl();

        if ($input.val().length === 0) {
          return;
        }

        data.title = $input.val();
        self.buildTask($list, data, true);
      });
    },

    handleDeleteTask: function() {
      $(document).on('click', '.taskboard-task-delete', function() {
        var $this = $(this);
        bootbox.dialog({
          message: "Do you want to delete the task?",
          buttons: {
            success: {
              label: "Delete",
              className: "btn-danger",
              callback: function() {
                $this.closest('.slidePanel').data('slidePanel').target.remove();
                $('.slidePanel-close').trigger('click');
              }
            }
          }
        });
      });
    },

    handleTaskInput: function() {
      var self = this;
      $(document).on('click', '.taskboard-list .checkbox-custom input', function(e) {
        var $this = $(this),
          $target = $this.closest('.list-group-item');

        self.dataChange($target, 'complete', $this.prop("checked"));
        e.stopPropagation();
      });
    },

    //Init SlidePanel
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

          this.$panel.find('#task-description').markdown();
          if (info.duedate.length > 0) {
            this.$panel.find('#taskDatepicker').data('date', info.duedate);
          }
          this.$panel.find('#taskDatepicker').datepicker({
            autoclose: false,
            todayHighlight: true,
          }).on('changeDate', function() {
            $('#taskDatepickerInput').val(
              _this.$panel.find('#taskDatepicker').datepicker('getFormattedDate')
            );
          });;


          this.$panel.data('slidePanel', object);

          $(document).off('click.slidePanelDatepicker');
          $(document).on('click.slidePanelDatepicker', 'span, td, th', function(e) {
            e.stopPropagation();
          });
        },
        afterShow: function() {
          var self = this;
          $(document).on('click.slidePanelShow', function(e) {
            if ($(e.target).closest('.slidePanel').length === 0 && $(e.target).closest('body').length === 1) {
              self.hide();
            }
          });
        },
        afterHide: function() {
          $(document).off('click.slidePanelShow');
          $(document).off('click.slidePanelDatepicker');
        },
        contentFilter: function(data, object) {
          var $target = $(object.target),
            info = $target.data('taskInfo'),
            $panel = $(data),
            $checked;

          $('.stage-name', $panel).html($('.taskboard-stage-title', self.getStage($target)).html());

          $('.task-title', $panel).html(info.title);

          switch (info.priority) {
            case 'high':
              $checked = $('#priorityHigh', $panel);
              break;
            case 'urgent':
              $checked = $('#priorityUrgent', $panel);
              break;
            default:
              $checked = $('#priorityNormal', $panel);
              break;
          }
          $checked.prop("checked", true);

          self.handleSelective($('[data-plugin="jquery-selective"]', $panel), info.members);

          if (info.description.length === 0) {
            $('.description', $panel).addClass('is-empty');
          } else {
            $('.description-content', $panel).html(info.description);
          }

          if (info.subtasks.length !== 0) {
            $.each(info.subtasks, function(n, subtask) {
              var $subtask = $(self.subtaskTpl(subtask));
              $('.subtasks-list', $panel).append($subtask);
            });
            $('.subtasks', $panel).toggleClass('is-show');
          }

          if (info.attachments.length !== 0) {
            $.each(info.attachments, function(n, attachment) {
              var $attachment = $(self.attachmentTpl(attachment));
              $('.attachments-list', $panel).append($attachment);
            });
            $('.attachments', $panel).toggleClass('is-show');
          }

          if (info.comments.length !== 0) {
            $.each(info.comments, function(n, comment) {
              var $comment = $(self.commentTpl(comment.src, comment.user, comment.time, comment.content));
              $('.comments-history', $panel).append($comment);
            });
          }

          return $panel;
        }
      });

      $(document).on('click', '[data-taskboard="slidePanel"]', function(e) {
        var $target = $(e.target).closest('.list-group-item');
        $.slidePanel.show({
          url: $(this).data('url'),
          target: $target
        }, options);

        e.stopPropagation();
      });

      $(document).on('click', '#fileuploadToggle', function() {
        $('#fileupload').trigger('click');
      });
    },

    //SlidePanel Section Handle
    handleSelective: function($target, selected) {
      var self = this;
      var getSelected = function() {
        var _this = this,
          arr = [];
        $.each(this._options.getOptions(this), function(n, option) {
          $.each(_this.options.local, function(i, user) {
            if (user.id === $(option).val()) {
              arr.push(user);
            }
          });
        });
        return arr;
      };
      var members = [{
        id: 'uid_1',
        name: 'Herman Beck',
        img: '../../../../global/portraits/1.jpg'
      }, {
        id: 'uid_2',
        name: 'Mary Adams',
        img: '../../../../global/portraits/2.jpg'
      }, {
        id: 'uid_3',
        name: 'Caleb Richards',
        img: '../../../../global/portraits/3.jpg'
      }, {
        id: 'uid_4',
        name: 'June Lane',
        img: '../../../../global/portraits/4.jpg'
      }, {
        id: 'uid_5',
        name: 'Edward Fletcher',
        img: '../../../../global/portraits/5.jpg'
      }, {
        id: 'uid_6',
        name: 'Crystal Bates',
        img: '../../../../global/portraits/6.jpg'
      }];

      $target.selective({
        namespace: 'addMember',
        local: members,
        selected: selected,
        buildFromHtml: false,
        tpl: {
          optionValue: function(data) {
            return data.id;
          },
          frame: function() {
            return '<div class="' + this.namespace + '">' +
              this.options.tpl.items.call(this) +
              '<div class="' + this.namespace + '-trigger">' +
              this.options.tpl.triggerButton.call(this) +
              '<div class="' + this.namespace + '-trigger-dropdown">' +
              this.options.tpl.list.call(this) +
              '</div>' +
              '</div>' +
              '</div>'
          },
          triggerButton: function() {
            return '<div class="' + this.namespace + '-trigger-button"><i class="md-plus"></i></div>';
          },
          listItem: function(data) {
            return '<li class="' + this.namespace + '-list-item"><img class="avatar" src="' + data.img + '">' + data.name + '</li>';
          },
          item: function(data) {
            return '<li class="' + this.namespace + '-item"><img class="avatar" src="' + data.img + '">' +
              this.options.tpl.itemRemove.call(this) +
              '</li>';
          },
          itemRemove: function() {
            return '<span class="' + this.namespace + '-remove"><i class="md-minus-circle"></i></span>';
          },
          option: function(data) {
            return '<option value="' + this.options.tpl.optionValue.call(this, data) + '">' + data.name + '</option>';
          }
        },
        onAfterItemAdd: function() {
          var $target = this.$el.closest('.slidePanel').data('slidePanel').target,
            arr = getSelected.call(this);
          self.dataChange($target, 'members', arr);
        },
        onAfterItemRemove: function() {
          var $target = this.$el.closest('.slidePanel').data('slidePanel').target,
            arr = getSelected.call(this);
          self.dataChange($target, 'members', arr);
        }
      });
    },

    handlePriority: function() {
      var self = this;
      $(document).on('click', '[name="priorities"]', function() {
        var $this = $(this),
          $target = $this.closest('.slidePanel').data('slidePanel').target;

        self.dataChange($target, 'priority', $this.data('priority'));
      });
    },

    handleEditor: function() {
      var self = this;
      $(document).on('click', '.slidePanel .task-title, .taskboard-task-edit, .description-toggle', function() {
        var $this = $(this),
          $target = $this.closest('.slidePanel').data('slidePanel').target,
          data = $target.data('taskInfo');

        $('#task-title').val(data.title);
        $('#task-description').val(data.description);
        $this.closest('.slidePanel').find('.task-main').addClass('is-edit');
      });

      $(document).on('click', '.task-main-editor-save', function() {
        var $this = $(this),
          $target = $this.closest('.slidePanel').data('slidePanel').target,
          data = $target.data('taskInfo');

        self.dataChange($target, 'title', $('#task-title').val());
        self.dataChange($target, 'description', $('#task-description').val());

        $this.closest('.slidePanel').find('.task-main').removeClass('is-edit');
        if ($('#task-description').val().length === 0) {
          $('.description').addClass('is-empty');
        } else {
          $('.description').removeClass('is-empty');
        }
      });

      $(document).on('click', '.task-main-editor-cancel', function() {
        $(this).closest('.slidePanel').find('.task-main').removeClass('is-edit');
      });
    },

    handleSubtasks: function() {
      var self = this;
      $(document).on('click', '.subtask-toggle', function() {
        var length = $('.subtask').length,
          $input = $('.subtasks-add .subtask-title'),
          $subtasks = $('.subtasks');

        $input.val('');
        if (length === 0) {
          $subtasks.addClass('is-show');
        }
        $subtasks.addClass('is-edit');

        $input.focus();

        $(document).on('click.subtask-add', function(e) {
          var $target = $(e.target);
          if ($target.closest($('.subtasks-add')).length === 0) {
            $subtasks.removeClass('is-edit');
            $(document).off('click.subtask-add');
          }
        });
      });

      $(document).on('click', '.subtask-add-save', function() {
        var length = $('.subtask').length,
          $subtasks = $('.subtasks'),
          $input = $('.subtasks-add .subtask-title'),
          value = $input.val(),
          $target = $(this).closest('.slidePanel').data('slidePanel').target;

        if (value.length === 0) {
          if (length === 0) {
            $subtasks.removeClass('is-show');
          }
        } else {
          var data = {
              'title': value,
              'complete': false
            },
            $subtask = $(self.subtaskTpl(data));

          $('.subtasks-list').append($subtask);
          self.dataChange($target, 'subtasks', data, length);
        }
        $input.val('').focus();
      });

      $(document).on('click', '.subtask-add-cancel', function() {
        $('.subtasks').removeClass('is-edit');
        $(document).off('click.subtask-add');
      });

      $(document).on('click', '.subtask input', function() {
        var $this = $(this),
          $target = $this.closest('.slidePanel').data('slidePanel').target,
          $subtask = $this.closest('.subtask'),
          index = $subtask.index();

        self.dataChange($target, 'subtasks', $this.prop("checked"), index, 'complete');
      });

      $(document).on('click', '.subtask .title', function() {
        var $this = $(this),
          $target = $this.closest('.slidePanel').data('slidePanel').target,
          data = $target.data('taskInfo'),
          $subtask = $this.closest('.subtask'),
          index = $subtask.index(),
          $input = $('.subtask-title', $subtask);

        $subtask.addClass('is-edit');
        $input.val('').focus().val(data.subtasks[index].title);

        $(document).on('click.subtask', function(e) {
          var $target = $(e.target);
          if ($target.closest($subtask).length === 0) {
            $subtask.removeClass('is-edit');
            $(document).off('click.subtask');
          }
        });
      });

      $(document).on('click', '.subtask-editor-save', function() {
        var $this = $(this),
          $target = $this.closest('.slidePanel').data('slidePanel').target,
          data = $target.data('taskInfo'),
          $subtask = $this.closest('.subtask'),
          index = $subtask.index();

        self.dataChange($target, 'subtasks', $('.subtask-title', $subtask).val(), index, 'title');
        $subtask.removeClass('is-edit');
        $(document).off('click.subtask');
      });

      $(document).on('click', '.subtask-editor-delete', function(e) {
        var $this = $(this);

        bootbox.dialog({
          message: "Do you want to delete the subtask?",
          buttons: {
            success: {
              label: "Delete",
              className: "btn-danger",
              callback: function() {
                var $target = $this.closest('.slidePanel').data('slidePanel').target,
                  data = $target.data('taskInfo'),
                  $subtask = $this.closest('.subtask'),
                  index = $subtask.index();
                self.dataDelete($target, 'subtasks', index);
                $subtask.remove();
                $(document).off('click.subtask');
                if ($('.subtask').length === 0) {
                  $('.subtasks').removeClass('is-show');
                }
              }
            }
          }
        });

      });
    },

    handleDatepicker: function() {
      var self = this;
      $(document).on('click', '.due-date-save', function() {
        var $this = $(this),
          $target = $this.closest('.slidePanel').data('slidePanel').target,
          value = $('#taskDatepickerInput').val();
        if (value.length > 0) {
          self.dataChange($target, 'duedate', value);
        }
      });
      $(document).on('click', '.due-date-delete', function() {
        var $this = $(this),
          $target = $this.closest('.slidePanel').data('slidePanel').target,
          data = $target.data('taskInfo');
        if (data.duedate.length === 0) return;
        self.dataDelete($target, 'duedate');
        $('#taskDatepicker').datepicker('clearDates');
      });
    },

    //Data
    dataDelete: function($target, name, index) {
      if (index) {
        $target.data('taskInfo')[name].splice(index, 1);
        this.dataDeleteResponse($target, name, index);
      } else {
        $target.data('taskInfo')[name] = '';
        this.dataChangeResponse($target, name);
      }
    },

    dataDeleteResponse: function($target, name, index) {
      switch (name) {
        case 'duedate':
          this.buildBadges($target, $target.data('taskInfo'));
          break;
        case 'subtasks':
          this.buildBadges($target, $target.data('taskInfo'));
          break;
      }
    },

    dataChange: function($target, name, content, index, subname) {
      if (content.length === 0) return;
      if (index !== undefined) {
        if (subname) {
          $target.data('taskInfo')[name][index][subname] = content;
        } else {
          $target.data('taskInfo')[name][index] = content;
        }
      } else {
        $target.data('taskInfo')[name] = content;
      }
      this.dataChangeResponse($target, name, content, index, subname);
    },

    dataChangeResponse: function($target, name, content, index, subname) {
      switch (name) {
        case 'title':
          $('.task-title', $target).html(content);
          $('.slidePanel .task-title').html(content);
          break;
        case 'description':
          $('.slidePanel .description-content').html(content);
          break;
        case 'priority':
          $target.removeClass('priority-normal priority-high priority-urgent').addClass('priority-' + content);
          break;
        case 'duedate':
          this.buildBadges($target, $target.data('taskInfo'));
          break;
        case 'members':
          this.buildMembers($target, $target.data('taskInfo').members);
          break;
        case 'subtasks':
          if (subname === 'title') {
            $('.title', $('.subtasks-list .subtask')[index]).html(content);
          } else {
            this.buildBadges($target, $target.data('taskInfo'));
          }
          break;
        case 'attachments':
          break;
        case 'comments':
          break;
      }
    },

    run: function(next) {
      var self = this;

      this.init();

      this.handleAddStage();
      this.handleDeleteStage();

      this.handleAddTask();
      this.handleDeleteTask();

      this.handleTaskInput();

      this.initStageDropdown();
      this.handleStageRename();

      this.handleDatepicker();
      this.handlSlidePanelPlugin();

      this.handleEditor();
      this.handleSubtasks();
      this.handlePriority();


      next();
    }
  });

  $(document).ready(function() {
    AppTaskboard.run();
  });
})(document, window, jQuery);
