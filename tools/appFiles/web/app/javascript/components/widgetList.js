import React, {Component} from 'react';
import {findDOMNode} from 'react-dom';
import BaseComponent from './base';
import {Table, TableBody, TableFooter, TableHeader, TableHeaderColumn, TableRow, TableRowColumn}
  from 'material-ui/Table';
import ReactDOM from 'react-dom';
import IconButton from 'material-ui/IconButton';
import RaisedButton from 'material-ui/RaisedButton';
import FlatButton from 'material-ui/FlatButton';
import ContentLowPriority from 'material-ui/svg-icons/content/low-priority';
import HardwareKeyboardArrowLeft from 'material-ui/svg-icons/hardware/keyboard-arrow-left';
import HardwareKeyboardArrowRight from 'material-ui/svg-icons/hardware/keyboard-arrow-right';
import HardwareKeyboardArrowDown from 'material-ui/svg-icons/hardware/keyboard-arrow-down';
import HardwareKeyboardArrowUp from 'material-ui/svg-icons/hardware/keyboard-arrow-up';
import ActionSearch from 'material-ui/svg-icons/action/search';
import TextField from 'material-ui/TextField';
import SelectField from 'material-ui/SelectField';
import MenuItem from 'material-ui/MenuItem';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import getMuiTheme from 'material-ui/styles/getMuiTheme';
import {blueGrey100, blueGrey400, deepOrange600, blueGrey900,  blueGrey600, indigo600, grey500, grey100, green500} from 'material-ui/styles/colors';
import {fade} from 'material-ui/utils/colorManipulator';
import ActionAdd from 'material-ui/svg-icons/content/add';

const noFieldValue = <span style={{color: blueGrey100, fontSize: '13px'}}>-VAL-</span>;
const paginationRowIcon = "IconButton";
const paginationRowRaisedButton =  "RaisedButton";
const paginationRowFlatButton = "FlatButton";
const totalHeaderCount = 2;
const totalFooterHeight = 82;

class WidgetList extends BaseComponent {
  constructor(props, context) {
    super(props, context);
    this.consts = {
      paginationButtons: "paginationButtons",
      paginationRowIcon: paginationRowIcon,
      paginationRowRaisedButton: paginationRowRaisedButton,
      paginationRowFlatButton: paginationRowFlatButton,
    };
    this.headerButtonRef;
    this.headerLeftRef;
    this.textFieldSearchRef;
    this.headerCenterRef;

    this.rowButtons = [];

    props.rowButtons.map((row) => {
      if (!row.hasOwnProperty('func')) {
        row.func = (row) => {
          return true
        };
      }
      this.rowButtons.push(row);
    });
    this.checkboxButtons = [];

    var checkboxButtons = [];
    if (props.checkboxButtonsShowAlways) {
      checkboxButtons = props.checkboxButtons.map((v) => v.button );
    }

    this.rowButtonClick = (func, row) => {
      func(row);
    };

    this.multiCheckboxActionButtonClick = (func) => {
      func(this.selectedData);
    };

    this.handlePerPageChange = (event, index, value) => {
      this.buildSearchRequest({PerPage: value});
    };

    this.handleSearchChange = (event) => {
      this.buildSearchRequest({Criteria: event.target.value, Page: 1});
    };

    this.handleSearchChangeStateOnly = (event) => {
      this.setComponentState({SearchVal: event.target.value});
    };

    this.selectedData = [];

    this.rowSelected = (rows) => {
      if (rows == 'all') {
        rows = [];
        var i = 0;
        this.globs.times(this.state.PerPage)(() => {
          rows.push(i);
          i++;
        });
      } else if (rows == 'none') {
        rows = [];
      }

      this.selectedData = [];

      Object.keys(rows).forEach((k) => {
        this.selectedData.push(this.state.data[rows[k]]);
      });

      this.selectedRows = rows;

      if (this.checkboxButtonsShowAlways) {
        this.setComponentState({
          checkboxButtons: this.props.checkboxButtons.map((v) => v.button)
        }, () => {
          this.setComponentState({addButtonOffset: this.getLeftForButtonOffset()})
        });
      } else {
        this.setComponentState({
          checkboxButtons: this.props.checkboxButtons.map((v) => {
            if (this.selectedData.length > 0) {
              if (v.hasOwnProperty('func') && v.func(this.selectedData)) {
                return v.button;
              } else if (!v.hasOwnProperty('func')) {
                return v.button;
              }
            }
            return ""
          })
        }, () => {
          this.setComponentState({addButtonOffset: this.getLeftForButtonOffset()})
        })
      }
      return true;
    };

    if (props.controller)  {
      var controller = props.controller;
    } else {
      var controller = props.name;
    }

    this.hasSorted = (this.props.listViewModel.SortBy);
    this.controller = controller;
    this.warns = {};
    this.selectedRows = [];

    if (!window.pageState.hasOwnProperty(props.dataKey)) {
      console.error("Bad dataKey [" + props.dataKey + "].  Ensure your pageState has this key with an array of objects for the WidgetList to paginate");
    }
    if (props.data) {
      var data = props.data;
    } else {
      var data = (!window.pageState[props.dataKey]) ? [] : window.pageState[props.dataKey];
    }

    this.state = {
      // User changable state
      SearchFields: this.props.listViewModel.SearchFields,
      PerPage: this.props.listViewModel.PerPage,
      Page: this.props.listViewModel.Page,
      SortBy: this.props.listViewModel.SortBy,
      SortDirection: this.props.listViewModel.SortDirection,
      Criteria: this.props.listViewModel.Criteria,
      CustomCriteria: this.props.listViewModel.CustomCriteria,
      IsDefaultFilter: this.props.listViewModel.IsDefaultFilter,
      ListTitle: this.globs.translate(this.props.listViewModel.ListTitle),
      data: data,
      tableHeight: '300px',
      widthTable: $(window).width(),
      widthRight: 0,
      addButtonOffset: 362 - 49,
      widthWindow: $(window).width(),
      heightTable: $(window).height(),
      heightWindow: $(window).height(),
      checkboxButtons: checkboxButtons
    };

    $.each(this.getResizeChanges(false), (k, v) => {
      this.state[k] = v;
    });

    this.resizeEvent = (e) => this.handleResize(e);

  }

  getLeftForButtonOffset() {
    return $(findDOMNode(this.headerButtonRef)).width() - 49;
  }

  getFields(width) {
    var newFields = [];
    this.props.fields.map((row) => {
      if (!row.hasOwnProperty('sortable')) {
        row.sortable = true;
      }
      if (!row.hasOwnProperty('func')) {
        row.func = (row, currentValue) => currentValue;
      }
      if (!row.hasOwnProperty('tooltip')) {
        row.tooltip = row.headerDisplay;
      }
      if (!row.hasOwnProperty('tooltipKey')) {
        row.tooltipKey = row.stateKey;
      }
      if (!row.hasOwnProperty('responsiveKeep')) {
        row.responsiveKeep = false;
      }
      if (!row.hasOwnProperty('headerDisplay')) {
        console.error("Missing field `headerDisplay` for the following fields prop: ", row);
      }
      if (!row.hasOwnProperty('stateKey')) {
        console.error("The field " + row.headerDisplay + " needs a configuration called stateKey which maps to the key in which we would show the row value of the data");
      }
      if (!row.hasOwnProperty('sortOn')) {
        row.sortOn = row.stateKey;
      }

      if (width < 1025 && !row.responsiveKeep) {
        return true;
      }

      if (width < 400 && row.responsiveKeep && newFields.length == 1) {
        return true;
      }

      newFields.push(row);
    });

    if (newFields.length == 0) {
      newFields.push(this.props.fields[0]);
      if (width > 550 && width < 1024) {
        try {
          this.props.fields[1]
          newFields.push(this.props.fields[1]);
        } catch(e) { }
      }
    }

    if (this.rowButtons.length > 0) {
      newFields.push({
        tooltip: "",
        headerDisplay: "",
        sortable: false,
        sortOn: this.consts.paginationButtons,
        responsiveKeep: true,
        stateKey: this.consts.paginationButtons,

      });
    }
    return newFields;
  }

  handleResize(e) {
    this.setComponentState(this.getResizeChanges(true), () => {
      this.setComponentState({fieldsWithButton: this.getFields(this.state.widthWindow), addButtonOffset: this.getLeftForButtonOffset()});
    });
  }

  getResizeChanges(initBlank=true) {
    var bodySelector = $(this.props.bodySelector);
    var updates = {};
    if (!initBlank) {
      updates = {
        fieldsWithButton: this.getFields($(window).width()),
      };
    }
    // dynamically fit
    let heightOffsetCenter = 0;
    let heightOffsetLeft = 0;
    let heightOffset = 0;
    if ($(window).height() < 500) {
      updates.tableHeight = ((this.getRowHeight() * 5)) + 'px';
    } else {
      if (initBlank) {
        heightOffset = $(findDOMNode(this)).height();
      }
      var sparePixels = bodySelector.height() - (heightOffset + this.props.offsetHeightToList);
      var possibleTotalHeightToWorkWith = bodySelector.height() - this.props.offsetHeightToList - ((this.getRowHeight() * totalHeaderCount) + totalFooterHeight);
      var maxHeight = this.getRowHeight() * this.state.PerPage;
      if (possibleTotalHeightToWorkWith > maxHeight) {
        possibleTotalHeightToWorkWith = maxHeight;
      }
      updates.tableHeight = possibleTotalHeightToWorkWith + 'px';
    }
    updates.widthWindow = $(window).width();
    updates.widthTable = bodySelector.width() - (($(window).width() < 1025) ? 0: this.props.marginRight);

    if (initBlank) {
      heightOffsetCenter = $(findDOMNode(this.headerCenterRef)).width();
      heightOffsetLeft = $(findDOMNode(this.headerLeftRef)).width();
    }
    updates.widthRight = updates.widthTable -  heightOffsetCenter -  heightOffsetLeft;
    updates.heightWindow = $(window).height();
    return updates;
  }

  componentDidMount() {
    window.addEventListener('resize', this.resizeEvent);
  }

  getSearchTotalRows(rowCount=0) {
    if (this.state.Criteria == "") {
      return "";
    }
    if (this.state.Criteria != "" && rowCount > 0) {
      return " " + this.globs.__(window.appContent.WidgetListCount, {TotalCnt: rowCount})
    }
    return "";
  }

  componentWillUnmount() {
    window.removeEventListener('resize', this.resizeEvent);
  }


  showCheckBoxes() {
    return (this.state.data && this.state.data.length > 0 && this.props.showCheckboxes);
  }

  getRowHeight() {
    return this.props.rowHeight;
  }

  buildSearchRequest(params={}) {
    this.selectedRows = [];
    var newParams = {};
    window.SearchChanged = false;
    newParams.SearchFields = this.state.SearchFields;
    newParams.PerPage = this.state.PerPage;
    newParams.Page = this.state.Page;
    newParams.Criteria = this.state.Criteria;
    newParams.SortBy = this.state.SortBy;
    newParams.IsDefaultFilter = this.state.IsDefaultFilter;
    newParams.CustomCriteria = this.state.CustomCriteria;
    newParams.ListTitle = this.globs.translate(this.state.ListTitle);
    newParams.SortDirection = this.state.SortDirection;
    var uriParams = Object.assign(newParams, params);
    uriParams.PerPage = parseInt(uriParams.PerPage);
    uriParams.Page = parseInt(uriParams.Page);
    window.pageState.WidgetList = {};
    window.pageState.WidgetList.SortBy = uriParams.SortBy;
    window.pageState.WidgetList.SearchFields = uriParams.SearchFields;
    window.pageState.WidgetList.SortDirection = uriParams.SortDirection;
    window.pageState.WidgetList.Page = uriParams.Page;
    window.pageState.WidgetList.PerPage = uriParams.PerPage;
    window.pageState.WidgetList.CustomCriteria = uriParams.CustomCriteria;
    window.pageState.WidgetList.Criteria = uriParams.Criteria;
    window.pageState.WidgetList.IsDefaultFilter = uriParams.IsDefaultFilter;
    window.pageState.WidgetList.ListTitle = uriParams.ListTitle;
    window.goCore.setLoaderFromExternal({loading: true});
      if (window.pageState.WidgetList.Criteria != "" && window.IsSearching) {
        window.ReRunSearch = true;
      }
      window.IsSearching = true;
    window.api.get({action:"Search", leaveStateAlone: true, uriParams: newParams, controller: this.controller, callback: (vm) => {
      this.setComponentState({
        data: vm[this.props.dataKey],
        ListTitle: this.globs.translate(vm.WidgetList.ListTitle),
        SortBy: window.pageState.WidgetList.SortBy,
        SearchFields: window.pageState.WidgetList.SearchFields,
        SortDirection: window.pageState.WidgetList.SortDirection,
        Page: window.pageState.WidgetList.Page,
        PerPage: window.pageState.WidgetList.PerPage,
        CustomCriteria: window.pageState.WidgetList.CustomCriteria,
        Criteria: window.pageState.WidgetList.Criteria,
        IsDefaultFilter: window.pageState.WidgetList.IsDefaultFilter,
      }, () => {
          window.goCore.setLoaderFromExternal({loading: false});
          if (window.ReRunSearch) {
            window.ReRunSearch = false;
            this.buildSearchRequest({Criteria: window.pageState.WidgetList.Criteria, Page: 1});
          }
          window.IsSearching = false;
        });
      }});
  }

  getOpposite(direction) {
    return (direction === '-') ? '+': '-';
  }

  buildHeaderColumnSortLink(headerRow) {
    var style = {fontSize: '18px'};
    if (headerRow.sortable) {
      var isSortingThisCol = (this.hasSorted && this.state.SortBy == headerRow.sortOn);

      return (<a href="javascript:" style={style} onClick={() => {
        var newSort = this.getOpposite(this.state.SortDirection);
        this.hasSorted = true;
        this.buildSearchRequest({SortBy: headerRow.sortOn, SortDirection: newSort, Page: 1});
      }}>
        {headerRow.headerDisplay}
        <span style={{position: 'absolute', top: 9}}>
          {(isSortingThisCol && this.state.SortDirection == '-') ?
              <IconButton>
                <HardwareKeyboardArrowDown/>
              </IconButton>
              : null}
          {(isSortingThisCol && this.state.SortDirection == '+') ?
              <IconButton>
                <HardwareKeyboardArrowUp/>
              </IconButton>
              : null}
        </span>
      </a>)
    } else {
      return (<span style={style}>{headerRow.headerDisplay}</span>)
    }
  }

  nextPage() {
    var page = this.state.Page + 1;
    return page;
  }


  previousPage() {
    if (this.state.Page == 1) {
      var page = 1;
    } else {
      var page = this.state.Page - 1;
    }
    return page;
  }

  render() {
    try {
      this.logRender();

      var footerButtonCss = {
        mediumIcon: {
          width: 48,
          color: blueGrey100,
          height: 48,
          marginTop: 15
        },
        medium: {
          width: 96,
          height: 70,
          color: blueGrey100
        },
      };
      var maxSpacer = (this.state.Page < 10) ? 10: 20;

      var perPageOptions = [];

      let pages = {
        "10": window.appContent.Numbers10,
        "25": window.appContent.Numbers25,
        "50": window.appContent.Numbers50,
        "100": window.appContent.Numbers100
      };

      Object.keys(pages).forEach((k) => {
        perPageOptions.push(<MenuItem
          key={k}
          value={k}
          primaryText={pages[k]}
          />);
      });
      var header = (
        <TableRow key="1search">
          {(this.props.showCheckboxes) ? <TableHeaderColumn ref={(c) => this.headerLeftRef = c}  colSpan={1} className="widgetListColumn" style={{width: 24, backgroundColor: blueGrey400, height: 50}}></TableHeaderColumn> : null}

          <TableHeaderColumn ref={(c) => this.headerCenterRef = c} colSpan={3} tooltip={""} style={{textAlign: 'left', fontSize: 22, color: blueGrey100, backgroundColor: blueGrey400}}>
            {(!this.state.IsDefaultFilter) ? <IconButton iconStyle={{color: blueGrey100}} onTouchTap={() => this.buildSearchRequest({CustomCriteria: "", Page: 1, Criteria: "", IsDefaultFilter: true, ListTitle: ""})}>
              <ContentLowPriority/>
            </IconButton>: null}
            {(this.selectedData.length > 0) ? this.__(window.appContent.WidgetListSelectedItems, {count: this.selectedData.length}) : this.state.ListTitle + this.getSearchTotalRows(this.props.listViewModel.TotalResults)}
          </TableHeaderColumn>

          <TableHeaderColumn ref={(c) => this.headerButtonRef = c} colSpan={this.state.fieldsWithButton.length - 3} style={{ textAlign: 'right', backgroundColor: blueGrey400}}>
            <div style={{height: 70, paddingTop: 16, color: blueGrey100, marginRight: (this.props.addRecordOnClick) ? 34: 0 }}>
              {
                this.state.checkboxButtons.map((button, buttonIndex) => {
                  if (button == "") {
                    return;
                  }
                  var btnTypeName = this.props.checkboxButtonTypes;
                  if (btnTypeName == this.consts.paginationRowIcon || btnTypeName == this.consts.paginationRowRaisedButton || btnTypeName == this.consts.paginationRowFlatButton) {
                    var newProps = {};

                    newProps.style = {marginRight: 15};
                    Object.keys(button.props).forEach((k) => {
                      if (k == 'onTouchTap') {
                        newProps['onTouchTap'] = () => {
                          this.multiCheckboxActionButtonClick((rows) => (button.props[k])(rows));
                        }
                      } else if (k == 'label' && this.state.widthWindow < 1024) {
                        //skip
                        newProps.style.minWidth = 40;
                      } else {
                        newProps[k] = button.props[k];
                      }
                    });
                    newProps.key = buttonIndex + '-footerButton';
                    core.Debug.Dump(newProps);

                    if (btnTypeName == this.consts.paginationRowIcon) {
                      return (
                          <IconButton {...newProps}>
                            {button.props.children}
                          </IconButton>)
                    } else if (btnTypeName == this.consts.paginationRowRaisedButton) {
                      return (
                          <RaisedButton {...newProps}>
                            {button.props.children}
                          </RaisedButton>)

                    } else if (btnTypeName == this.consts.paginationRowFlatButton) {
                      return (
                          <FlatButton {...newProps}>
                            {button.props.children}
                          </FlatButton>)
                    }
                  } else {
                    console.log(btnTypeName, 'Pagination Icon button type Not implemented');
                  }
                })
              }
              {
                (this.props.addRecordOnClick && $(window).width() > 1024) ?
                      <span key={window.globals.guid()} style={{position: 'absolute', top: -15, left: this.state.addButtonOffset}}>
                      <IconButton title={this.props.addRecordOnClickToolTip} onClick={(e) => {
                        this.props.addRecordOnClick();
                      }} style={footerButtonCss.medium} iconStyle={footerButtonCss.mediumIcon}  ><ActionAdd color={blueGrey100}/>
                      </IconButton>
                        </span>
                 : null
              }
            </div>
          </TableHeaderColumn>
        </TableRow>
      );
      let NavRightProps = {};
      NavRightProps.disabled = (this.state.data && this.state.data.length == 0);
      NavRightProps.onTouchTap = () => this.buildSearchRequest({Page: this.nextPage()});
      let NavLeftProps = {};
      NavLeftProps.disabled = (this.previousPage(false) == this.state.Page);
      NavLeftProps.onTouchTap = () => this.buildSearchRequest({Page: this.previousPage()});

      var footer  = (
        <TableFooter className="widgetListFooter" style={{lineHeight: "20px"}}>
          <TableRow style={{height: 50, maxHeight: 50}}>
            {(this.props.searchEnabled) ?
            <span>
            <TableRowColumn colSpan={this.state.widthWindow > 1024 ? 1: 4} style={{height: 50, textAlign: 'left'}}>
              {/*Mobile pagination buttons need to be here*/}
              {(this.props.nextBackButtonsEnabled && this.state.widthWindow < 1024) ?
                <span>
                <RaisedButton
                  labelColor={window.materialColors["blueGrey900"]}
                  style={{minWidth: 40}}
                  icon={<HardwareKeyboardArrowLeft/>}
                  color={window.materialColors["blueGrey500"]}
                  {...NavLeftProps}
                />
                <RaisedButton
                  labelColor={window.materialColors["blueGrey900"]}
                  style={{minWidth: 40, marginRight: 15, marginLeft: 5}}
                  icon={<HardwareKeyboardArrowRight/>}
                  color={window.materialColors["blueGrey500"]}
                  {...NavRightProps}
                />
                </span>: null}
              <TextField
                ref={(c) => this.textFieldSearchRef = c}
                floatingLabelStyle={{color: blueGrey100}}
                hintStyle={{color: blueGrey100}}
                style={{color: blueGrey100, marginBottom: 10, width: 250}}
                inputStyle={{color: blueGrey100}}
                floatingLabelText={(this.props.searchByLabelText) ? this.props.searchByLabelText : window.appContent.WidgetListDefaultSearchBy}
                hintText={(this.props.searchByHintText) ? this.props.searchByHintText : window.appContent.WidgetListDefaultSearchBy}
                defaultValue={this.state.Criteria}
                fullWidth={true}
                onChange={this.state.widthWindow < 1024 ? this.handleSearchChange: this.handleSearchChangeStateOnly}
                onKeyPress={this.state.widthWindow > 1024 ? (event) => {
                  if (event.nativeEvent.key == "Enter"){
                    this.buildSearchRequest({Criteria: this.state.SearchVal, Page: 1})
                  }
                } : null}

                errorText={(this.state.data && this.state.data.length == 0) ? this.props.noDataMessage : null}
              />
              </TableRowColumn>
              {this.state.widthWindow > 1024 ?
            <TableRowColumn colSpan={1}>
              <RaisedButton
                    style={{minWidth: 40, marginTop:10, float: "left"}}
                    icon={<ActionSearch/>}
                    primary={true}
                    backgroundColor={window.materialColors["blueGrey700"]}
                    label={window.appContent.Search}

                    onTouchTap={() => {
                      this.buildSearchRequest({Criteria: this.state.SearchVal, Page: 1})
                    }}
                  /></TableRowColumn>: null}
            </span> : null}
              {(this.props.perPageSelectBoxEnabled && this.state.widthWindow > 1024) ?
            <TableRowColumn colSpan={1} style={{textAlign: 'right', textOverflow: "none"}}>
                  <span className="widgetListPerPage">
                    <SelectField
                        underlineStyle={{color: blueGrey100}}
                        underlineFocusStyle={{color: blueGrey100}}
                        floatingLabelStyle={{color: blueGrey100}}
                        labelStyle={{color: blueGrey100}}
                        style={{color: blueGrey100, width: 60, marginRight: 20}}
                        inputStyle={{color: blueGrey100}}
                        onChange={this.handlePerPageChange}
                        value={this.state.PerPage.toString()}
                      >
                        {perPageOptions}
                      </SelectField>
                  </span>
            </TableRowColumn> : null}

              {(this.props.nextBackButtonsEnabled && this.state.widthWindow > 1024) ?
            <TableRowColumn colSpan={1} style={{textAlign: 'right'}}>
                <span className="widgetListPagination">
                      <span style={{height: this.getRowHeight(), paddingTop: 11, color: blueGrey100}}>
                        <IconButton style={footerButtonCss.medium} iconStyle={footerButtonCss.mediumIcon}  {...NavLeftProps}>
                          <HardwareKeyboardArrowLeft/>
                        </IconButton>
                      </span>
                </span>
            </TableRowColumn>: null}

            {(this.props.nextBackButtonsEnabled && this.state.widthWindow > 1024) ?
            <TableRowColumn colSpan={1} style={{width: 10}}>
                      <span style={{height: this.getRowHeight(), minWidth: maxSpacer, maxWidth: maxSpacer, fontSize: '25px', color: blueGrey100}}>
                        <span className="Aligner">
                          <span className="Aligner-item Aligner-item--top"></span>
                          <span className="Aligner-item">{window.appContent.hasOwnProperty("Numbers" + this.state.Page) ? window.appContent["Numbers" + this.state.Page]: this.state.Page}</span>
                          <span className="Aligner-item Aligner-item--bottom"></span>
                        </span>
                        {(this.state.Page < 10) ? <span style={{minWidth: maxSpacer, maxWidth: maxSpacer}}>&#160;</span>: null}
                        {(this.state.Page >= 10) ? <span style={{minWidth: maxSpacer, maxWidth: maxSpacer}}>&#160;&#160;&#160;</span>: null}
                      </span>
            </TableRowColumn> : null}

            {(this.props.nextBackButtonsEnabled && this.state.widthWindow > 1024) ?
            <TableRowColumn colSpan={1}>
                      <span style={{height: 70, paddingTop: 11, color: blueGrey100}}>
                        <IconButton style={footerButtonCss.medium} iconStyle={footerButtonCss.mediumIcon} {...NavRightProps}>
                          <HardwareKeyboardArrowRight/>
                        </IconButton>
                      </span>
            </TableRowColumn> : null}

          </TableRow>
        </TableFooter>
      );

      var tableRows = null;
      if (this.state.data && this.state.data.length > 0) {
        tableRows = this.state.data.map( (row, index) => {
            let selected = (!row.hasOwnProperty('selected')) ? false : true;
            if (this.selectedRows.length > 0) {
              this.selectedRows.forEach((v) => {
                if (v == index) {
                  selected = true;
                }
              });
            }
            return (
            <TableRow className="widgetListRow" key={this.props.listViewModel.Page + "-" + index} selected={selected}>
              {
                this.state.fieldsWithButton.map((field, fieldindex) => {
                    if (field.stateKey == this.consts.paginationButtons) {
                      var content = (
                        <div  className="Aligner">
                          {
                            this.rowButtons.map((btnObj, buttonIndex) => {
                              if(btnObj.func(row)) {
                                var btnTypeName = this.props.rowButtonTypes;
                                if (btnTypeName == this.consts.paginationRowIcon || btnTypeName == this.consts.paginationRowRaisedButton || btnTypeName == this.consts.paginationRowFlatButton) {
                                var newProps = {};
                                Object.keys(btnObj.button.props).forEach((k) => {
                                  if (k == 'onTouchTap') {
                                    newProps['onTouchTap'] = () => {
                                      this.rowButtonClick(function(row) {
                                        (btnObj.button.props[k])(row);
                                      }, row);
                                    }
                                  } else {
                                    newProps[k] = btnObj.button.props[k];
                                  }
                                });

                                newProps.key = buttonIndex + '-' + index + '-' + fieldindex;
                                if (btnTypeName == this.consts.paginationRowIcon) {
                                  return (
                                    <span key={window.globals.guid()}>
                                      <span className="Aligner-item Aligner-item--top"></span>
                                      <span className="Aligner-item">
                                        <IconButton {...newProps}>
                                          {btnObj.button.props.children}
                                        </IconButton>
                                      </span>
                                      <span className="Aligner-item Aligner-item--bottom"></span>
                                    </span>
                                  )

                                } else if (btnTypeName == this.consts.paginationRowRaisedButton) {
                                  if (!newProps.hasOwnProperty("style")) {
                                    newProps.style = {};
                                  }
                                  newProps.style.minWidth = 40;
                                  newProps.style.marginRight = 15;
                                  return (
                                    <span key={window.globals.guid()}>
                                      <span className="Aligner-item Aligner-item--top"></span>
                                      <span className="Aligner-item">
                                        <RaisedButton {...newProps}>
                                          {btnObj.button.props.children}
                                        </RaisedButton>
                                      </span>
                                      <span className="Aligner-item Aligner-item--bottom"></span>
                                    </span>
                                  )

                                } else if (btnTypeName == this.consts.paginationRowFlatButton) {
                                  return (
                                    <span key={window.globals.guid()}>
                                      <span className="Aligner-item Aligner-item--top"></span>
                                      <span className="Aligner-item">
                                        <FlatButton {...newProps}>
                                          {btnObj.button.props.children}
                                        </FlatButton>
                                      </span>
                                      <span className="Aligner-item Aligner-item--bottom"></span>
                                    </span>
                                  )
                                }

                              } else {
                                console.log(btnTypeName, 'Pagination Icon button type Not implemented');
                              }
                            }
                          })
                        }
                      </div>
                    );
                      var style = {textAlign: 'right'};
                    } else {
                      var content = "";
                      var tooltipContent = "";
                      if (field.stateKey.indexOf(".") == -1) {
                        if (!row.hasOwnProperty(field.stateKey) && !this.warns.hasOwnProperty(field.stateKey)) {
                          console.error("Bad field pointer for field [" + field.stateKey + "].  Could not find in data.");
                          this.warns[field.stateKey] = true;
                        }
                        content = row[field.stateKey];
                        tooltipContent = row[field.tooltipKey];
                      } else {
                        // currently non recursive, just two and three levels deep
                        var pointers = field.stateKey.split(".");
                        var tooltipPointers = field.tooltipKey.split(".");

                        var currentContentPointer = row[pointers[0]];
                        var currentTooltipPointer = row[tooltipPointers[0]];
                        for(var i = 1; i < pointers.length; i++){
                          if (currentContentPointer == undefined){
                            currentContentPointer = "";
                            console.error("Bad field pointer for field [" + field.stateKey + "].  Could not find in data.");
                            this.warns[field.stateKey] = true;
                          } else {

                            if (currentContentPointer[pointers[i]] == undefined) {
                              currentContentPointer = "";
                              console.error("Bad field pointer for field [" + field.stateKey + "].  Could not find in data.");
                              this.warns[field.stateKey] = true;
                            } else {
                              currentContentPointer = currentContentPointer[pointers[i]];
                            }

                          }

                          if (currentTooltipPointer == undefined){
                            currentTooltipPointer = "";
                            console.error("Bad field pointer for field [" + field.tooltipKey + "].  Could not find in data.");
                            this.warns[field.tooltipKey] = true;
                          } else {

                              if (currentTooltipPointer[tooltipPointers[i]] == undefined) {
                                currentTooltipPointer = "";
                                console.error("Bad field pointer for field [" + field.tooltipKey + "].  Could not find in data.");
                                this.warns[field.tooltipKey] = true;
                              } else {
                                currentTooltipPointer = currentTooltipPointer[tooltipPointers[i]];
                              }

                          }
                        }

                        content = currentContentPointer;
                        tooltipContent = currentTooltipPointer;

                      }
                      content = field.func(row, content);
                      if (content == "") {
                        var newProps = Object.assign({},this.props.noFieldValue.props);
                        delete newProps.children;
                        content = <span key={window.globals.guid()} {...newProps} title={tooltipContent}>{this.props.noFieldValueTranslation}</span>
                      }
                      var style = {fontSize: '16px'};
                    }
                    return (
                      <TableRowColumn className="widgetListColumn" style={style} key={index + '-' + fieldindex}><span title={tooltipContent}>{content}</span></TableRowColumn>
                    )
                  }
                )
              }
            </TableRow>
            );
        });
      } else if (this.state.data && this.state.data.length == 0) {
        tableRows = ([
              <TableRow key="row1"></TableRow>,
              <TableRow key="row2"></TableRow>,
              <TableRow key="row3"></TableRow>,
              <TableRow key="row4"></TableRow>,
              <TableRow key="row5" style={{border: 'none'}}>
                <TableRowColumn><h3 className="center-xs">{this.props.noDataMessage}</h3></TableRowColumn>
              </TableRow>,
              <TableRow key="row6"></TableRow>,
              <TableRow key="row7"></TableRow>,
              <TableRow key="row8"></TableRow>,
              <TableRow key="row9"></TableRow>,
              <TableRow key="row10"></TableRow>
            ]
        );
      }

      var WidgetListTheme = getMuiTheme({
        palette: {
          accent1Color: deepOrange600,
          accent2Color: indigo600,
          accent3Color: blueGrey600,
          primary1Color: blueGrey900,
          primary2Color: indigo600,
          primary3Color: grey500,
          pickerHeaderColor: deepOrange600,
        }
      });

      WidgetListTheme.slider.trackColor = 'white';
      WidgetListTheme.slider.selectonColor = deepOrange600;
      WidgetListTheme.slider.trackSize = 4;
      WidgetListTheme.toggle.thumbOnColor = green500;
      WidgetListTheme.toggle.thumbOffColor = grey100;
      WidgetListTheme.toggle.trackOnColor = fade(WidgetListTheme.toggle.thumbOnColor, 0.5);
      WidgetListTheme.toggle.trackOffColor = grey100;
      WidgetListTheme.tableRow.height = this.props.rowHeight;
      WidgetListTheme.tableRowColumn.height = this.props.rowHeight;

      return (
        <MuiThemeProvider muiTheme={WidgetListTheme}>
          <div className={"widgetList " + this.props.name} style={{width: this.state.widthTable, boxShadow: 'rgba(0, 0, 0, 0.117647) 0px 1px 6px, rgba(0, 0, 0, 0.117647) 0px 1px 4px'}}>
          <input type="hidden" id="SortBy" value={this.state.SortBy}/>
          <input type="hidden" id="SortDirection" value={this.state.SortDirection}/>
          <Table
            height={this.state.tableHeight}
            fixedHeader={true}
            fixedFooter={true}
            selectable={this.showCheckBoxes()}
            multiSelectable={true}
            onRowSelection={this.rowSelected}
          >
            <TableHeader
              className="widgetListHeader"
              displaySelectAll={false}
              adjustForCheckbox={false}
              enableSelectAll={true}
            >
              {header}
              <TableRow style={{height: this.getRowHeight() + 'px'}}>
                {(this.props.showCheckboxes) ? <TableHeaderColumn className="widgetListColumn" style={{width: 24}}></TableHeaderColumn> : null}
                {this.state.fieldsWithButton.map((row, index) => (
                  <TableHeaderColumn key={index} tooltip={row.tooltip}>{this.buildHeaderColumnSortLink(row)}</TableHeaderColumn>
                ))}
              </TableRow>
            </TableHeader>
            <TableBody
              displayRowCheckbox={this.showCheckBoxes()}
              deselectOnClickaway={false}
              showRowHover={false}
            >
              {tableRows}
            </TableBody>
            {footer}
          </Table>

          <input type="hidden" name="rerender" value={this.props.listViewModel.PerPage + "-" + this.props.listViewModel.Page + "-" + this.props.listViewModel.SortBy + "-" + this.props.listViewModel.SortDirection + "-" + this.props.listViewModel.Criteria + "-" + this.props.listViewModel.CustomCriteria + "-" + this.props.listViewModel.IsDefaultFilter + "-" + this.props.listViewModel.ListTitle}/>
        </div>
        </MuiThemeProvider>
      );
    } catch(e) {
      return this.globs.ComponentError(this.getClassName(), e.message);
    }
  }
}


WidgetList.propTypes = {
  // Ultimately will be used to support two lists on the same page, but currently it wouldnt work 100%
  name: React.PropTypes.string.isRequired,
  listTitle: React.PropTypes.string.isRequired,

  // The WidgetList modelView pointer
  listViewModel: React.PropTypes.object.isRequired,
  // controller to run a GET on when user interacts with widgetlist
  controller: React.PropTypes.string,
  showCheckboxes: React.PropTypes.bool,
  fields: React.PropTypes.arrayOf(React.PropTypes.shape({
    sortable: React.PropTypes.bool,
    func: React.PropTypes.func, // Will pass your field and row data into it and pass the entire row to you in the first parameter so you can wrap and customize more.  The signature of the inputs is (row, currentValue) => currentValue;  Where you can use the entire row's data to build more nodes
    sortOn: React.PropTypes.string, //Exact mongo pointer to sort on if your view is in a different pointer
    tooltip: React.PropTypes.string, // will become headerDisplay if not passed
    responsiveKeep: React.PropTypes.bool, // default false
    headerDisplay: React.PropTypes.string.isRequired,
    stateKey: React.PropTypes.string.isRequired, // A dot syntax only three levels deep currently
  })),
  dataKey: React.PropTypes.string.isRequired, // The key in your View Model in which a []model.Object name exists

  rowButtonTypes: React.PropTypes.oneOf([paginationRowIcon, paginationRowRaisedButton, paginationRowFlatButton]),
  // Buttons on each row,  Will be resized small currently
  rowButtons: React.PropTypes.arrayOf(React.PropTypes.shape({
    func: React.PropTypes.func, // will result to a closure which will always return true.  This will pass rows to your callback with all the checked items
    button: React.PropTypes.node.isRequired // Pass in a <RaisedButton> <IconButton> or <FlatButton>  WidgetList will re-render and wrap your onTouchTap's so that it calls your function on your caller instead of handling inside the component.
  })),

  //Internal mostly, but can be overridden
  perPage: React.PropTypes.number,

  noFieldValue: React.PropTypes.node.isRequired, // Add a node with CSS applied when the string of a returned field callback is ""
  noFieldValueTranslation: React.PropTypes.string.isRequired,
  noDataMessage: React.PropTypes.string, // Message when no data is found
  searchByLabelText: React.PropTypes.string,
  searchByHintText: React.PropTypes.string,

  // You need to pass this because I cannot introspect on a minified react using this.
  checkboxButtonTypes: React.PropTypes.oneOf([paginationRowIcon, paginationRowRaisedButton, paginationRowFlatButton]),
  // Buttons upon checked
  checkboxButtons: React.PropTypes.arrayOf(React.PropTypes.shape({
    func: React.PropTypes.func, // will result to a closure which will always return true.  This will pass rows to your callback with all the checked items
    button: React.PropTypes.node.isRequired // Pass in a <RaisedButton> <IconButton> or <FlatButton>  WidgetList will re-render and wrap your onTouchTap's so that it calls your function on your caller instead of handling inside the component.
  })),

  // Controls on your end to show things how you want.
  checkboxButtonsShowAlways: React.PropTypes.bool, // Always show buttons
  searchEnabled: React.PropTypes.bool,
  nextBackButtonsEnabled: React.PropTypes.bool,
  perPageSelectBoxEnabled: React.PropTypes.bool,

  // a prop meant to bind data from the callers so external page changes outside of widgetlist can re-render the data through state and props
  data: React.PropTypes.array,

  // extra margin right if needed
  marginRight: React.PropTypes.number,

  // extra margin bottom if needed
  marginBottom: React.PropTypes.number,

  //How hight are your rows
  rowHeight: React.PropTypes.number,

  //Main outer page div to calculate size of widgetList (css selector)
  bodySelector: React.PropTypes.string,

  //You can pass an onclick and a magic + icon will always show up in a subtle manner
  addRecordOnClick: React.PropTypes.func,
  addRecordOnClickToolTip: React.PropTypes.string,

  // If there are elements above the widgetlist, pass in the appropriate pixels of items above so that it can calculate a full height needed to show the most rows.
  offsetHeightToList: React.PropTypes.number
};


WidgetList.defaultProps = {
  offsetHeightToList: 0,
  bodySelector: '#pagecontent',
  rowHeight: 60,
  marginRight: 40,
  marginBottom: 5,
  searchEnabled: true,
  perPage: 10,
  noFieldValue: noFieldValue,
  noFieldValueTranslation: "(none)",
  noDataMessage: "No Records Found",
  showCheckboxes: true,
  nextBackButtonsEnabled: true,
  perPageSelectBoxEnabled: true,
  checkboxButtonsShowAlways: false,
  fields: [],
  checkboxButtonTypes: paginationRowRaisedButton,
  checkboxButtons: [],
  rowButtonTypes: paginationRowRaisedButton,
  rowButtons: []
};


export default WidgetList;
