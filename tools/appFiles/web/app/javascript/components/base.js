import React, {Component} from "react";
// import ReactPerfAnalysis from "react-addons-perf";
import {RichUtils} from "draft-js";


// Most logging code stolen from https://www.npmjs.com/package/react-log-lifecycle

// Optional flags:
const flags = {
  // If logType is set to keys then the props of the object being logged
  // will be written out instead of the whole object. Remove logType or
  // set it to anything except keys to have the full object logged.
  logType: 'keys',
  // A list of the param "types" to be logged.
  // The example below has all the types.
  names: ['props', 'nextProps', 'nextState', 'prevProps', 'prevState']
};

class BaseComponent extends Component {
  constructor(props, flags) {
    super(props);
    this.updateParent(props);

    var showLog = true;
    if(window.appState.hasOwnProperty("DeveloperSuppressThesePages") && window.appState.DeveloperSuppressThesePages !== null && window.appState.DeveloperSuppressThesePages.length > 0) {
      window.appState.DeveloperSuppressThesePages.forEach((v) => {
        if (window.location.href.indexOf(v) > -1) {
          showLog = false;
        }
      });
      this.showLog = showLog;
    }
    if (showLog) {
      this.showLog = ((window.appState.hasOwnProperty("DeveloperLogTheseObjects") && window.appState.DeveloperLogTheseObjects && (window.appState.DeveloperLogTheseObjects.length == 0 || $.inArray(this.getClassName(), window.appState.DeveloperLogTheseObjects) != -1)) && $.inArray(this.getClassName(), window.appState.DeveloperSuppressTheseObjects) == -1);
    }
    this.cycleNum = 1;
    this.flags = flags || {};


    if (window.appState.DeveloperLogReact && this.showLog) {
      if (performance) {
        console.info(performance.now());
      }
      console.info(this.getClassName());
      console.info(`#${this.cycleNum}.1 constructor(props)
  - Start of cycle #${this.cycleNum}
  - replaces getInitialState()
  - assign to this.state to set initial state.
`);
    }

    this._log({props});
    this.__ = window.global.functions._;
    this.globs = window.global.functions;
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

  setComponentState(stateChanges, cb, setParent=false) {
    var callback = cb;

    if (window.appState.DeveloperLogState && this.showLog) {
      if (performance) {
        console.info(performance.now());
      }
      this.stack = new Error().stack;
      var line = "";
      if (this.stack != undefined && this.stack.split("\n") != undefined && this.stack.split("\n").length > 0) {
        if (this.stack.split("\n")[2].indexOf("makeAssimilatePrototype.js") > -1) {
          line = this.stack.split("\n")[3];
        } else {
          line = this.stack.split("\n")[2];
        }

        console.info("Stacktrace" + line.trim() + " called => setComponentState(" + stateChanges + ")");
        console.info(this.stack);
      }
      console.info("-----<setComponentState>----");
      console.info("");
      console.info("");
      console.info('State changes proposed:', stateChanges);
      // console.info("React Object:", this);
    }

    var ptr = null;
    if (!setParent) {
      ptr = this;
    } else if (setParent && typeof(this.parent) != "undefined") {
      ptr = this.parent;
    } else {
      ptr = this;
    }

    // let keys = Object.keys(stateChanges);
    // if (keys.length == 1 && ptr.state[keys[0]] === stateChanges[keys[0]]) {
    //   if (typeof(callback) == "function") {
    //     callback();
    //   }
    //   return
    // }

    var merge = this.globs.mergeDeep(ptr.state, stateChanges);
    if (window.appState.DeveloperLogState && this.showLog) {
      console.info("Merged State Sent To setState:", merge);
      console.info("");
      console.info("");
      console.info("-----</setComponentState>----");
    }

    if (Object.keys(stateChanges).length > 0) {

      ptr.setState(merge, ()=> {
        if (typeof(callback) == "function") {
          callback();
        }
        if (window.appState.DeveloperLogState && this.showLog) {
          var line = "";

          if (this.stack != undefined && this.stack.split("\n") != undefined && this.stack.split("\n").length > 0) {
            if (this.stack.split("\n")[2].indexOf("makeAssimilatePrototype.js") > -1) {
              line = this.stack.split("\n")[3];
            } else {
              line = this.stack.split("\n")[2];
            }
            console.info("Callback After React setState -> " + line.trim(), " called => setComponentState(", stateChanges, ")");
            if (performance) {
              console.info(performance.now());
            }
            console.info("Stacktrace");
            console.info(this.stack);
          }
          console.info("");
          console.info("React Is Done Mutating State From Requested Change", stateChanges);
          console.info("The True state of your component: ", this.state);
        }
      });
    } else {
      if (callback) {
        callback();
      }
    }

    // Maybe later mess with the update function and use the other methods....
    // this.setState(update(this.state, {$merge: stateChanges}), () => {
    //   if (window.appState.DeveloperLogState) {
    //     console.info("setComponentState called " + this.stack.split("\n")[1].trim(), this);
    //   }
    // });
  }

  _log(obj) {

    if (window.appState.DeveloperLogReact && this.showLog) {
      // obj should have a single property.
      // The name (key) of that property should be switched on in the flags
      // object if it should be logged. You swith it on by adding it to the
      // flags.names array.
      // The value of the single prop in obj is the object that is to be
      // logged. (Or the keys of the object to be logged.)
      const keys = Object.keys(obj);
      if (keys.length !== 1) {
        return;
      }
      const key = keys[0];

      if (this.flags.names && this.flags.names.indexOf(key) >= 0) {
        // The flags object can override logging the object by changing
        // the logType to 'keys' to just log out the keys of the object.
        const logObj = (this.flags.logType && this.flags.logType === 'keys' && obj[key])
            ? Object.keys(obj[key])
            : obj[key];

        console.info(key + ':', logObj);
      }
    }
  }

  componentWillMount() {
    if (window.appState.DeveloperLogReact && this.showLog) {
      if (performance) {
        console.info(performance.now());
      }
      console.info(this.getClassName());
      console.info(`#${this.cycleNum}.2 componentWillMount()
  - Invoked Once (client and server)
  - Can change state here with this.setState()  (will not trigger addition render)
  - Just before render()
`);
    }
  }

  componentDidMount() {
    this.updateParent(this.props);

    if (window.appState.DeveloperLogReact && this.showLog) {
      if (performance) {
        console.info(performance.now());
      }
      console.info(this.getClassName());
      console.info(`#${this.cycleNum}.3 componentDidMount()
  - Invoked Once (client only)
  - refs to children now available
  - integrate other JS frameworks, timers, ajax etc. here
  - Just after render()
  - End of Cycle #${this.cycleNum}
`);
    }
    this.cycleNum = 2;
  }

  componentWillReceiveProps(nextProps) {
    if (window.appState.DeveloperLogReact && this.showLog) {
      if (performance) {
        console.info(performance.now());
      }
      console.info(this.getClassName());
      console.info(`#${this.cycleNum}.1 componentWillReceiveProps(nextProps)
  - Start of cycle #${this.cycleNum}
  - invoked when component is receiving new props
  - not called in cycle #1
  - this.props is old props
  - parameter to this function is nextProps
  - can call this.setState() here (will not trigger addition render)
`);
    }
    this.updateParent(nextProps);

    if (window.appState.DeveloperLogState && this.showLog) {
      console.info("componentWillReceiveProps (" + this.getClassName() + ")");
      console.info("Old Props", this.props);
      console.info("New Props", nextProps);
    }

    this._log({nextProps});
  }

  shouldComponentUpdate(nextProps, nextState) {
    if (window.appState.DeveloperLogReact && this.showLog) {
      if (performance) {
        console.info(performance.now());
      }
      console.info(this.getClassName());
      console.info(`#${this.cycleNum}.2 shouldComponentUpdate(nextProps, nextState)
  - invoked when new props/state being received
  - not called on forceUpdate()
  - returning false from here prevents component update and the next 2 parts of the Lifecycle: componentWillUpdate() componentDidUpdate()
  - returns true by default;
`);
    }
    this.updateParent(nextProps);
    this._log({nextProps});
    this._log({nextState});
    return true;
  }

  componentWillUpdate(nextProps, nextState) {
    this.updateParent(nextProps);

    if (window.appState.DeveloperLogReact && this.showLog) {
      if (performance) {
        console.info(performance.now());
      }
      console.info(this.getClassName());
      console.info(`#${this.cycleNum}.3 componentWillUpdate(nextProps, nextState)
  - cannot use this.setState() (do that in componentWillReceiveProps() above)
  - Just before render()
`);
    }

    this._log({nextProps});
    this._log({nextState});
  }

  componentDidUpdate(prevProps, prevState) {
    if (window.appState.DeveloperLogReact && this.showLog) {
      if (performance) {
        console.info(performance.now());
      }
      console.info(this.getClassName());
      console.info(`#${this.cycleNum}.4 componentDidUpdate(prevProps, prevState)
  - Just after render()
`);
    }

    // if (window.appState.DeveloperLogStateChangePerformance && this.showLog) {
    //   ReactPerfAnalysis.stop();
    //   console.info("> React Inline Render Performance Using DeveloperLogStateChangePerformance (Expand me please!!) <");
    //   console.info("Time Wasted Below:");
    //   ReactPerfAnalysis.printWasted();
    //   console.info("Render Time In Your Code(Without react):");
    //   ReactPerfAnalysis.printExclusive();
    //   console.info("Render Time Complete Breakdown:");
    //   ReactPerfAnalysis.printInclusive();
    // }
    this._log({prevProps});
    this._log({prevState});

  }

  componentWillUnmount() {
    this.cycleNum = 3;
    if (window.appState.DeveloperLogReact && this.showLog) {
      if (performance) {
        console.info(performance.now());
      }
      console.info(this.getClassName());
      console.info(`#${this.cycleNum}.1 componentWillUnmount()
  - invoked immediately before a component is unmounted from DOM
  - do cleanup here. e.g. kill timers and unlisten to events such as flux store updates
`);
    }
  }

  getClassName() {
    return this.constructor.toString().split("function ")[1].split("(")[0];
  }

  logRender() {
    if ((window.appState.DeveloperLogState || window.appState.DeveloperLogReact) && this.showLog) {
      if (performance) {
        console.info(performance.now());
      }
      console.info('render invoked', this);
    }

    // if (window.appState.DeveloperLogStateChangePerformance && this.showLog) {
    //   try {
    //     ReactPerfAnalysis.start();
    //   } catch(e) {
    //     session_functions.Dump("caught", e.message)
    //
    //   }
    // }
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
