import React, {Component} from "react";
// import ReactPerfAnalysis from "react-addons-perf";
import {RichUtils} from "draft-js";
import Store from "./store/store";


// Most logging code stolen from https://www.npmjs.com/package/react-log-lifecycle

class BaseComponent extends Component {
  constructor(props, context) {
    super(props, context);

    this.logger = window.baseLogger;

    this.updateParent(props);
    this.cycleNum = 1;

    this.logger.logReactLifecycle(() => {return this.getClassName();}, `#${this.cycleNum}.1 constructor(props)
      - Start of cycle #${this.cycleNum}
      - replaces getInitialState()
      - assign to this.state to set initial state.
    `);

    this.logger.log(() => {return this.getClassName();}, {props});
    this.__ = window.global.functions._;
    this.globs = window.global.functions;
    this.store = window.store;

    //https://medium.freecodecamp.org/functional-setstate-is-the-future-of-react-374f30401b6b
    this.setFunctionalState = (setStateCb, cb, setParent=false) => {
      this.logger.logSetStateStart(() => {return this.getClassName();});
      var ptr = null;
      if (!setParent) {
        ptr = this;
      } else if (setParent && typeof(this.parent) != "undefined") {
        ptr = this.parent;
      } else {
        ptr = this;
      }
      ptr.setState(setStateCb, ()=> {
        if (typeof(cb) == "function") {
          try {
            cb();
          } catch(e) {
            console.error("Error in setFunctionalState callback", e);
          }
        }
        this.logger.logSetStateEnd(() => {return this.getClassName();}, this.state);
      });
    };

    // legacy setComponentState probably has I/O issues with large payloads passing to this.globs.mergeDeep
    this.setComponentState = (stateChanges, cb, setParent=false) => {
      var callback = cb;
      this.logger.logSetComponentStateStart(() => {return this.getClassName();}, stateChanges);
      var ptr = null;
      if (!setParent) {
        ptr = this;
      } else if (setParent && typeof(this.parent) != "undefined") {
        ptr = this.parent;
      } else {
        ptr = this;
      }
      var merge = this.globs.mergeDeep(ptr.state, stateChanges);
      this.logger.logSetComponentStateMerge(() => {return this.getClassName();}, {merge});
      if (Object.keys(stateChanges).length > 0) {

        ptr.setState(merge, ()=> {
          if (typeof(callback) == "function") {
            try {
              callback();
            } catch(e) {
              console.error("Error in setComponentState callback", e);
            }
          }
          this.logger.logSetComponentStateStart(() => {return this.getClassName();}, stateChanges, this.state);
        });
      } else {
        if (callback) {
          callback();
        }
      }
    };

    this.base = {};
    this.base.setState = this.setFunctionalState;
    this.base.deepMergeState = this.setComponentState;
  }

  // CSS shortcuts
  typeof(o) {
    let typeInfo = "";
    if (o != null && o != undefined) {
      typeInfo = o.constructor.name; // returns Array, Boolean, Object, String, etc...
    } else if (o == null) {
      typeInfo = "Null";
    } else if (o == undefined) {
      typeInfo = "Undefined";
    }
    return typeInfo;
  }

  merge(...params) {
    var merge = {};
    params.forEach((p) => {
      if (this.typeof(p) == "Object") {
        merge = this.globs.mergeDeep(merge, p);
      }
    });
    return merge;
  }

  overFlowVertical(attrs=null) {
    return this.merge(attrs, {overflowY: "auto", overflowX: "hidden"});
  }

  overFlowHorizontal(attrs=null) {
    return this.merge(attrs, {overflowY: "hidden", overflowX: "auto"});
  }

  overFlowBoth(attrs=null) {
    return this.merge(attrs, {overflowY: "auto", overflowX: "auto"});
  }

  updateParent(props) {
    if (props.hasOwnProperty("parent")) {
      this.parent = props.parent;
      this.parentState = props.parent.state;
      // core.Debug.Dump("this.parent", this.parent, "this.parentState", this.parentState);
    }
  }

  replaceAll(source, value, replacementValue) {
    if (source == undefined) {
      return source;
    }
    return source.replace(new RegExp(value, "g"), replacementValue);
  }

  openFile(callback, accept){
    // Check for the various File API support.
    if (window.File && window.FileReader && window.FileList && window.Blob) {
      // Great success! All the File APIs are supported.
      var acceptInsert = "";
      if(accept)
        acceptInsert =  "accept='" + accept + "'";
      var fileInput = $("<input type='file' name='files[]' " + acceptInsert + "/>");
      $("body").append(fileInput);
      fileInput.trigger("click");
      fileInput.on('change', function(evt){
        console.info(evt.target.files[0]);

          var reader = new FileReader();
          var fileName = evt.target.files[0].name;
          //When the file has been read...
          reader.onload = function (ev) {
              try {
                  var jsonObj = new Object();
                  jsonObj.fileName = fileName;
                  jsonObj.data = window.atob(ev.target.result.substr(ev.target.result.indexOf("base64,") + 7));
                  callback(jsonObj);
                  fileInput.remove();
              }
              catch (ex) {
                  window.global.functions.PopupWindow(ex);
              }
          };
          //And now, read the image and base64
          reader.readAsDataURL(evt.target.files[0]);
      });
    } else {
      window.global.functions.PopupWindow('The File APIs are not fully supported in this browser.');
    }
  }

  setParentState(stateChanges, cb) {
    this.setComponentState(stateChanges, cb, true)
  }

  componentWillMount() {

    this.logger.logReactLifecycle(() => {return this.getClassName();}, `#${this.cycleNum}.2 componentWillMount()
    - Invoked Once (client and server)
    - Can change state here with this.setState()  (will not trigger addition render)
    - Just before render()
  `);

  }

  componentDidMount() {
    this.updateParent(this.props);

    this.logger.logReactLifecycle(() => {return this.getClassName();}, `#${this.cycleNum}.3 componentDidMount()
    - Invoked Once (client only)
    - refs to children now available
    - integrate other JS frameworks, timers, ajax etc. here
    - Just after render()
    - End of Cycle #${this.cycleNum}
  `);

    this.cycleNum = 2;
  }

  componentWillReceiveProps(nextProps) {

    this.logger.logReactLifecycle(() => {return this.getClassName();}, `#${this.cycleNum}.1 componentWillReceiveProps(nextProps)
    - Start of cycle #${this.cycleNum}
    - invoked when component is receiving new props
    - not called in cycle #1
    - this.props is old props
    - parameter to this function is nextProps
    - can call this.setState() here (will not trigger addition render)
  `);

    this.updateParent(nextProps);

    this.logger.logRecieveProps(() => {return this.getClassName();}, this.props, nextProps);
    this.logger.log(() => {return this.getClassName();}, {nextProps});
  }

  shouldComponentUpdate(nextProps, nextState) {

    this.logger.logReactLifecycle(() => {return this.getClassName();}, `#${this.cycleNum}.2 shouldComponentUpdate(nextProps, nextState)
    - invoked when new props/state being received
    - not called on forceUpdate()
    - returning false from here prevents component update and the next 2 parts of the Lifecycle: componentWillUpdate() componentDidUpdate()
    - returns true by default;
  `);

    this.updateParent(nextProps);
    this.logger.log(() => {return this.getClassName();}, {nextProps});
    this.logger.log(() => {return this.getClassName();}, {nextState});

    return true;
  }

  componentWillUpdate(nextProps, nextState) {
    this.updateParent(nextProps);

    this.logger.logReactLifecycle(() => {return this.getClassName();}, `#${this.cycleNum}.3 componentWillUpdate(nextProps, nextState)
    - cannot use this.setState() (do that in componentWillReceiveProps() above)
    - Just before render()
  `);

    this.logger.log(() => {return this.getClassName();}, {nextProps});
    this.logger.log(() => {return this.getClassName();}, {nextState});
  }

  componentDidUpdate(prevProps, prevState) {

    this.logger.logReactLifecycle(() => {return this.getClassName();}, `#${this.cycleNum}.4 componentDidUpdate(prevProps, prevState)
    - Just after render()
    `);

    this.logger.log(() => {return this.getClassName();}, {prevProps});
    this.logger.log(() => {return this.getClassName();}, {prevState});
  }

  registerSubscriptions(subscriptions) {
    this.subscriptions = subscriptions;
    this.store.registerAll(this.subscriptions);
  }

  unmountSubscriptions() {
    if (this.hasOwnProperty("subscriptions")) {
      this.store.unsubscribeAll(this.subscriptions);
    }
  }

  componentWillUnmount() {
    if (this.hasOwnProperty("unmount")) {
      this.unmount();
    }
    this.unmountSubscriptions();
    this.cycleNum = 3;
    this.logger.logReactLifecycle(() => {return this.getClassName();}, `#${this.cycleNum}.1 componentWillUnmount()
    - invoked immediately before a component is unmounted from DOM
    - do cleanup here. e.g. kill timers and unlisten to events such as flux store updates
  `);
  }

  getClassName() {
    //warning. use of this in production is frowned upon.  dev only!! compression will make classnames into simple letters which is useless
    return this.constructor.toString().split("function ")[1].split("(")[0];
  }

  logRender() {
    this.logger.logRender(() => {return this.getClassName();}, this);
  }

  _DraftJsHandleKeyCommand(command) {
    const {editorState} = this.state;
    const newState = RichUtils.handleKeyCommand(editorState, command);
    if (newState) {
      this.handleDraftJsChange(newState);
      return true;
    }
    return false;
  }

  _DraftJsOnTab(e) {
    const maxDepth = 4;
    this.handleDraftJsChange(RichUtils.onTab(e, this.state.editorState, maxDepth));
  }

  _DraftJstoggleBlockType(blockType) {
    this.handleDraftJsChange(
      RichUtils.toggleBlockType(
        this.state.editorState,
        blockType
      )
    );
  }

  _DraftJsToggleInlineStyle(inlineStyle) {
    this.handleDraftJsChange(
      RichUtils.toggleInlineStyle(
        this.state.editorState,
        inlineStyle
      )
    );
  }

}

export default BaseComponent;
