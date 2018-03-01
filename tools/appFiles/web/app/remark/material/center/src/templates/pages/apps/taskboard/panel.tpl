<header class="slidePanel-header">
  <div class="overlay-top overlay-panel overlay-background bg-light-green-600">
    <div class="slidePanel-actions btn-group btn-group-flat" aria-label="actions" role="group">
      <div class="dropdown pull-left">
        <button type="button" class="btn btn-pure dropdown-toggle icon md-calendar" data-toggle="dropdown"
        aria-hidden="true"></button>
        <div class="dropdown-menu dropdown-menu-right bullet taskboard-task-datepicker">
          <div id="taskDatepicker"></div>
          <input type="hidden" id="taskDatepickerInput" />
          <div class="padding-10">
            <button class="btn btn-primary due-date-save">Save</button>
            <a class="btn btn-sm btn-white due-date-delete" href="javascript:void(0)">Delete</a>
          </div>
        </div>
      </div>
      <button type="button" class="btn btn-pure icon md-format-list-bulleted subtask-toggle"
      aria-hidden="true"></button>
      <div class="fileupload pull-left">
        <button id="fileuploadToggle" type="button" class="btn btn-pure icon md-attachment-alt"
        aria-hidden="true"></button>
        <input id="fileupload" type="file" name="upload">
      </div>
      <div class="dropdown pull-left">
        <button type="button" class="btn btn-pure dropdown-toggle icon md-chevron-down"
        data-toggle="dropdown" aria-hidden="true"></button>
        <ul class="dropdown-menu dropdown-menu-right bullet" role="menu">
          <li role="presentation" class="taskboard-task-edit">
            <a href="javascript:void(0)" role="menuitem"><i class="icon md-edit" aria-hidden="true"></i>Edit Task</a>
          </li>
          <li role="presentation" class="taskboard-task-delete">
            <a href="javascript:void(0)" role="menuitem"><i class="icon md-delete" aria-hidden="true"></i>Delete Task</a>
          </li>
        </ul>
      </div>
      <button type="button" class="btn btn-pure slidePanel-close icon md-close" aria-hidden="true"></button>
    </div>
    <h4 class="stage-name"></h4>
  </div>
</header>
<div class="slidePanel-inner">
  <section class="slidePanel-inner-section">
    <div class="task-main">
      <div class="task-main-surface">
        <div class="priority pull-right">
          <label class="margin-right-20">Priority:</label>
          <ul class="list-unstyled list-inline inline-block">
            <li class="radio-custom radio-normal">
              <input type="radio" class="icheckbox-grey" id="priorityNormal" data-priority="normal"
              name="priorities">
              <label for="priorityNormal">Normal</label>
            </li>
            <li class="radio-custom radio-high">
              <input type="radio" class="icheckbox-grey" id="priorityHigh" data-priority="high"
              name="priorities">
              <label for="priorityHigh">High</label>
            </li>
            <li class="radio-custom radio-urgent">
              <input type="radio" class="icheckbox-grey" id="priorityUrgent" data-priority="urgent"
              name="priorities">
              <label for="priorityUrgent">Urgent</label>
            </li>
          </ul>
        </div>
        <div class="caption task-title"></div>
        <div class="add-members">
          <select multiple="multiple" data-plugin="jquery-selective"></select>
        </div>

        <div class="description">
          <p class="description-content"></p>
          <a href="#" class="description-toggle">Add description.</a>
        </div>
      </div>
      <div class="task-main-editor">
        <form>
          <div class="form-group">
            <input id="task-title" class="form-control" type="text" name="title">
          </div>
          <div class="form-group">
            <textarea id="task-description" name="content" data-provide="markdown" data-iconlibrary="fa"
            rows="6"></textarea>
          </div>
          <div class="form-group">
            <button class="btn btn-primary task-main-editor-save" type="button">Save</button>
            <a class="btn btn-sm btn-white task-main-editor-cancel" href="javascript:void(0)">Cancel</a>
          </div>
        </form>
      </div>
    </div>
    <div class="task-section subtasks">
      <label><i class="icon md-format-list-bulleted margin-right-10"></i>Subtasks</label>
      <ul class="list-group list-group-full subtasks-list"></ul>
      <div class="subtasks-add">
        <form>
          <div class="form-group">
            <input class="form-control subtask-title" type="text" name="title">
          </div>
          <div class="form-group">
            <button class="btn btn-primary subtask-add-save" type="button">Save</button>
            <a class="btn btn-sm btn-white subtask-add-cancel" href="javascript:void(0)">Cancel</a>
          </div>
        </form>
      </div>
    </div>
    <div class="task-section attachments">
      <label><i class="icon md-attachment-alt margin-right-10"></i>Attachments | <a href="#">Download All</a></label>
      <ul class="list-group attachments-list"></ul>
    </div>
  </section>
  <section class="slidePanel-inner-section">
    <div class="comments">
      <label><i class="icon md-comments margin-right-10"></i>Comments</label>
      <div class="comments-history"></div>
      <div class="comment media">
        <div class="media-left">
          <a class="avatar avatar-lg" href="javascript:void(0)">
            <img src="{{global}}/portraits/5.jpg" alt="...">
          </a>
        </div>
        <div class="media-body">
          <div class="comment-body">
            <form class="comment-reply margin-top-5" method="post" action="#">
              <textarea class="form-control margin-bottom-15" placeholder="Reply" rows="4"></textarea>
              <button class="btn btn-primary" type="button" data-dismiss="modal">Reply</button>
            </form>
          </div>
        </div>
      </div>
    </div>
  </section>
</div>