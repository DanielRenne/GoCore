import React from 'react';

import CenteredPaperGrid from '../components/centeredPaperGrid';
import BasePageComponent from '../components/basePageComponent';
import BaseComponent from '../components/base';
import WidgetList from '../components/widgetList';
import AddRecordPage from '../components/addRecordPageComponent';
import BackPage from '../components/backPageComponent';
import AddOrImportPage from '../components/addOrImportComponent';
import ConfirmPopup from '../components/confirmPopup';
import ConfirmDelete from '../components/confirmDelete';
import CenterAnything from '../components/centerAnything';
import Dropzone from 'react-dropzone';
import MegaIcon from '../components/megaButton';
import PhoneInput from '../components/phoneInput';
import FileUpload from '../components/fileUpload';
import FileUploadWithState from '../components/fileUploadWithState';
import {Grid, Row, Col} from 'react-flexgrid-no-bs-conflict';
import {Editor, EditorState, RichUtils, convertFromHTML, convertFromRaw, convertToRaw, ContentState} from 'draft-js';
import TextField from 'material-ui/TextField';
import RaisedButton from 'material-ui/RaisedButton';
import FlatButton from 'material-ui/FlatButton';
import IconButton from 'material-ui/IconButton';
import {RadioButton, RadioButtonGroup} from 'material-ui/RadioButton';
import Toggle from 'material-ui/Toggle';
import FloatingActionButton from 'material-ui/FloatingActionButton';
import ContentAdd from 'material-ui/svg-icons/content/add';
import Avatar from 'material-ui/Avatar';
import Dialog from 'material-ui/Dialog';
import SelectField from 'material-ui/SelectField';
import MenuItem from 'material-ui/MenuItem';
import {List, ListItem} from 'material-ui/List';
import Checkbox from 'material-ui/Checkbox';
import Divider from 'material-ui/Divider';
import AutoComplete from 'material-ui/AutoComplete';
import {red50, red100, red200, red300, red400, red500, red600, red700, red800, red900, pink50, pink100, pink200, pink300, pink400, pink500, pink600, pink700, pink800, pink900, purple50, purple100, purple200, purple300, purple400, purple500, purple600, purple700, purple800, purple900,deepPurple50, deepPurple100, deepPurple200, deepPurple300, deepPurple400, deepPurple500, deepPurple600, deepPurple700, deepPurple800, deepPurple900, indigo50, indigo100, indigo200, indigo300, indigo400, indigo500, indigo600, indigo700, indigo800, indigo900, blue50, blue100, blue200, blue300, blue400, blue500, blue600, blue700, blue800, blue900, lightBlue50, lightBlue100, lightBlue200, lightBlue300, lightBlue400, lightBlue500, lightBlue600, lightBlue700, lightBlue800, lightBlue900, cyan50, cyan100, cyan200, cyan300, cyan400, cyan500, cyan600, cyan700, cyan800, cyan900,  teal50, teal100, teal200, teal300, teal400, teal500, teal600, teal700, teal800, teal900, green50, green100, green200, green300, green400, green500, green600, green700, green800, green900, lightGreen50, lightGreen100, lightGreen200, lightGreen300, lightGreen400, lightGreen500, lightGreen600, lightGreen700, lightGreen800, lightGreen900, lime50, lime100, lime200, lime300, lime400, lime500, lime600, lime700, lime800, lime900, yellow50, yellow100, yellow200, yellow300, yellow400, yellow500, yellow600, yellow700, yellow800, yellow900, amber50, amber100, amber200, amber300, amber400, amber500, amber600, amber700, amber800, amber900,orange50, orange100, orange200, orange300, orange400, orange500, orange600, orange700, orange800, orange900, deepOrange50, deepOrange100, deepOrange200, deepOrange300, deepOrange400, deepOrange500, deepOrange600, deepOrange700, deepOrange800, deepOrange900, brown50, brown100, brown200, brown300, brown400, brown500, brown600, brown700, brown800, brown900, grey50, grey100, grey200, grey300, grey400, grey500, grey600, grey700, grey800, grey900, blueGrey50, blueGrey100, blueGrey200, blueGrey300, blueGrey400, blueGrey500, blueGrey600, blueGrey700, blueGrey800, blueGrey900 } from 'material-ui/styles/colors';
import AppBar from 'material-ui/AppBar'
import Drawer from 'material-ui/Drawer'
import CircularProgress from 'material-ui/CircularProgress';
import LinearProgress from 'material-ui/LinearProgress';
import TimePicker from 'material-ui/TimePicker';
import Slider from 'material-ui/Slider';
import {convertToHTML} from 'draft-convert';

let stateToHTML = convertToHTML;

// Custom overrides for "code" style.
const DraftJsStyleMap = {
  CODE: {
    backgroundColor: 'rgba(0, 0, 0, 0.05)',
    fontFamily: '"Inconsolata", "Menlo", "Consolas", monospace',
    fontSize: 16,
    padding: 2,
  },
};


function DraftJsGetBlockStyle(block) {
  switch (block.getType()) {
    case 'blockquote': return 'RichEditor-blockquote';
    default: return null;
  }
}


class DraftJsStyleButton extends React.Component {
  constructor() {
    super();
    this.onToggle = (e) => {
      e.preventDefault();
      this.props.onToggle(this.props.style);
    };
  }

  render() {
    let className = 'RichEditor-styleButton';
    if (this.props.active) {
      className += ' RichEditor-activeButton';
    }

    return (
      <span className={className} onMouseDown={this.onToggle}>
        {this.props.label}
      </span>
    );
  }
}

const DRAFT_JS_BLOCK_TYPES = [
  {label: 'H1', style: 'header-one'},
  {label: 'H2', style: 'header-two'},
  {label: 'H3', style: 'header-three'},
  {label: 'H4', style: 'header-four'},
  {label: 'H5', style: 'header-five'},
  {label: 'H6', style: 'header-six'},
  {label: 'Blockquote', style: 'blockquote'},
  {label: 'UL', style: 'unordered-list-item'},
  {label: 'OL', style: 'ordered-list-item'},
];
  // {label: 'Code Block', style: 'code-block'},

  
const DraftJsBlockStyleControls = (props) => {
  const {editorState} = props;
  const selection = editorState.getSelection();
  const blockType = editorState
    .getCurrentContent()
    .getBlockForKey(selection.getStartKey())
    .getType();

  return (
    <div className="RichEditor-controls">
      {DRAFT_JS_BLOCK_TYPES.map((type) =>
        <DraftJsStyleButton
          key={type.label}
          active={type.style === blockType}
          label={type.label}
          onToggle={props.onToggle}
          style={type.style}
        />
      )}
    </div>
  );
};

var DRAFT_JS_INLINE_STYLES = [
  {label: 'Bold', style: 'BOLD'},
  {label: 'Italic', style: 'ITALIC'},
  {label: 'Underline', style: 'UNDERLINE'},
  {label: 'Monospace', style: 'CODE'},
];


const DraftJsInlineStyleControls = (props) => {
  var currentStyle = props.editorState.getCurrentInlineStyle();
  return (
    <div className="RichEditor-controls">
      {DRAFT_JS_INLINE_STYLES.map(type =>
        <DraftJsStyleButton
          key={type.label}
          active={currentStyle.has(type.style)}
          label={type.label}
          onToggle={props.onToggle}
          style={type.style}
        />
      )}
    </div>
  );
};


var DropZone = Dropzone;

export {
  React,
  DropZone,
  CenteredPaperGrid,
  BasePageComponent,
  BaseComponent,
  WidgetList,
  AddRecordPage,
  AddOrImportPage,
  BackPage,
  ConfirmDelete,
  TextField,
  RaisedButton,
  FlatButton,
  IconButton,
  RadioButton,
  RadioButtonGroup,
  Toggle,
  FloatingActionButton,
  ContentAdd,
  Avatar,
  SelectField,
  MenuItem,
  List,
  ListItem,
  Checkbox,
  Divider,
  ConfirmPopup,
  Grid,
  Row,
  Col,
  AutoComplete,
  PhoneInput,
  red50, red100, red200, red300, red400, red500, red600, red700, red800, red900, pink50, pink100, pink200, pink300, pink400, pink500, pink600, pink700, pink800, pink900, purple50, purple100, purple200, purple300, purple400, purple500, purple600, purple700, purple800, purple900,deepPurple50, deepPurple100, deepPurple200, deepPurple300, deepPurple400, deepPurple500, deepPurple600, deepPurple700, deepPurple800, deepPurple900, indigo50, indigo100, indigo200, indigo300, indigo400, indigo500, indigo600, indigo700, indigo800, indigo900, blue50, blue100, blue200, blue300, blue400, blue500, blue600, blue700, blue800, blue900, lightBlue50, lightBlue100, lightBlue200, lightBlue300, lightBlue400, lightBlue500, lightBlue600, lightBlue700, lightBlue800, lightBlue900, cyan50, cyan100, cyan200, cyan300, cyan400, cyan500, cyan600, cyan700, cyan800, cyan900,  teal50, teal100, teal200, teal300, teal400, teal500, teal600, teal700, teal800, teal900, green50, green100, green200, green300, green400, green500, green600, green700, green800, green900, lightGreen50, lightGreen100, lightGreen200, lightGreen300, lightGreen400, lightGreen500, lightGreen600, lightGreen700, lightGreen800, lightGreen900, lime50, lime100, lime200, lime300, lime400, lime500, lime600, lime700, lime800, lime900, yellow50, yellow100, yellow200, yellow300, yellow400, yellow500, yellow600, yellow700, yellow800, yellow900, amber50, amber100, amber200, amber300, amber400, amber500, amber600, amber700, amber800, amber900,orange50, orange100, orange200, orange300, orange400, orange500, orange600, orange700, orange800, orange900, deepOrange50, deepOrange100, deepOrange200, deepOrange300, deepOrange400, deepOrange500, deepOrange600, deepOrange700, deepOrange800, deepOrange900, brown50, brown100, brown200, brown300, brown400, brown500, brown600, brown700, brown800, brown900, grey50, grey100, grey200, grey300, grey400, grey500, grey600, grey700, grey800, grey900, blueGrey50, blueGrey100, blueGrey200, blueGrey300, blueGrey400, blueGrey500, blueGrey600, blueGrey700, blueGrey800, blueGrey900,
  AppBar,
  Drawer,
  FileUpload,
  CircularProgress,
  LinearProgress,
  Editor,
  EditorState,
  RichUtils,
  DraftJsStyleMap,
  DraftJsGetBlockStyle,
  DraftJsStyleButton,
  DraftJsBlockStyleControls,
  DraftJsInlineStyleControls,
  stateToHTML,
  FileUploadWithState,
  ContentState,
  convertFromHTML,
  convertFromRaw,
  convertToRaw,
  TimePicker,
  Slider,
  MegaIcon,
  Dialog,
  CenterAnything
};
