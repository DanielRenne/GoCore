import util from 'util';

class BaseLogger {
  constructor() {

    this.logObject = (params) => {
      return;
    };

    this.logTheseObjects = (params) => {
      this.logObject(params);
    };

    this.supressPages = (params) => {
      this.logTheseObjects(params);
    };

    this.logReactLifecycle = (getClassName, message) => {
      this.supressPages({getClassName:getClassName, message:message, logReactLifecycle:true});
    };

    this.log = (getClassName, obj) => {
      this.supressPages({getClassName:getClassName, obj:obj, log:true});
    };

    this.logRecieveProps = (getClassName, props, nextProps) => {
      this.supressPages({getClassName:getClassName, props:props, nextProps:nextProps, recieveProps:true});
    };

    this.logRender = (getClassName, obj) => {
      this.supressPages({getClassName:getClassName, obj:obj, logRender:true});
    };

    this.logSetStateStart = (getClassName) => {
      this.supressPages({getClassName:getClassName, logSetStateStart:true});
    };

    this.logSetStateEnd = (getClassName, obj) => {
      this.supressPages({getClassName:getClassName, obj:obj, logSetStateEnd:true});
    };

    this.logSetComponentStateStart = (getClassName, obj) => {
      this.supressPages({getClassName:getClassName, obj:obj, logSetComponentStateStart:true});
    };

    this.logSetComponentStateMerge = (getClassName, obj) => {
      this.supressPages({getClassName:getClassName, obj:obj, logSetComponentStateMerge:true});
    };

    this.logSetComponentStateEnd = (getClassName, obj, obj2) => {
      this.supressPages({getClassName:getClassName, obj:obj, obj2:obj2,  logSetComponentStateEnd:true});
    };

    if (window.appState.DeveloperMode) {

      this.supressPages = (params) => {
        if(window.appState.hasOwnProperty("DeveloperSuppressThesePages") && window.appState.DeveloperSuppressThesePages !== null && window.appState.DeveloperSuppressThesePages.length > 0) {
          window.appState.DeveloperSuppressThesePages.forEach((v) => {
            if (window.location.href.indexOf(v) > -1) {
              return;
            }
          });
        }
        var className = "";
        if (params.getClassName != undefined) {
          className = params.getClassName();
        }
        if (window.appState.hasOwnProperty("DeveloperSuppressTheseObjects") && window.appState.DeveloperSuppressTheseObjects != null && window.appState.DeveloperSuppressTheseObjects.length > 0 && $.inArray(className, window.appState.DeveloperSuppressTheseObjects) != -1) {
          return;
        }
        this.logTheseObjects(params);
      };

      if (window.appState.hasOwnProperty("DeveloperLogTheseObjects") && window.appState.DeveloperLogTheseObjects && window.appState.DeveloperLogTheseObjects.length > 0) {
        this.logTheseObjects = (params) => {
          var className = "";
          if (params.getClassName != undefined) {
            className = params.getClassName();
          }

          if ($.inArray(className, window.appState.DeveloperLogTheseObjects) == -1) {
            return;
          }
          this.logObject(params);
        }
      }

      if (window.appState.DeveloperLogState || window.appState.DeveloperLogReact) {
        this.logObject = (params) => {
          if (params.logReactLifecycle === true && window.appState.DeveloperLogReact) {
            var className = "";
            if (params.getClassName != undefined) {
              className = params.getClassName();
            }

            if (performance) {
              console.info(performance.now());
            }
            console.info(className);
            console.info(params.message);
          } else if (params.logSetStateEnd === true && window.appState.DeveloperLogState) {
            var line = "";
            if (window.appState.UserAgent.Name != "Safari" && this.stack != undefined && this.stack.split("\n") != undefined && this.stack.split("\n").length > 0) {
              if (this.stack.split("\n")[2].indexOf("makeAssimilatePrototype.js") > -1) {
                line = this.stack.split("\n")[3];
              } else {
                line = this.stack.split("\n")[2];
              }
              console.info("Callback After React setState -> " + line.trim(), " called => base.setState()");
              if (performance) {
                console.info(performance.now());
              }
              console.info("Stacktrace");
              console.info(this.stack);
            }
            console.info("");
            console.info("The True state of your component: ", params.obj);
          } else if (params.logSetStateStart === true && window.appState.DeveloperLogState) {
            if (performance) {
              console.info(performance.now());
            }
            this.stack = new Error().stack;
            var line = "";
            if (window.appState.UserAgent.Name != "Safari" && this.stack != undefined && this.stack.split("\n") != undefined && this.stack.split("\n").length > 0) {
              if (this.stack.split("\n")[2].indexOf("makeAssimilatePrototype.js") > -1) {
                line = this.stack.split("\n")[3];
              } else {
                line = this.stack.split("\n")[2];
              }
              console.info("Stacktrace " + this.stack + " called => base.setState()");
              console.info(this.stack);
            }
            console.info("-----<setComponentState>----");
            console.info("");
            console.info("");
            // console.info("React Object:", this);
          } else if (params.logSetComponentStateMerge === true && window.appState.DeveloperLogState) {
            console.info("Merged State Sent To setState:", params.obj);
            console.info("");
            console.info("");
            console.info("-----</setComponentState>----");
          } else if (params.logSetComponentStateEnd === true && window.appState.DeveloperLogState) {
            let line = "";
            if (window.appState.UserAgent.Name != "Safari" && this.stack != undefined && this.stack.split("\n") != undefined && this.stack.split("\n").length > 0) {
              if (this.stack.split("\n")[2].indexOf("makeAssimilatePrototype.js") > -1) {
                line = this.stack.split("\n")[3];
              } else {
                line = this.stack.split("\n")[2];
              }
              console.info("Callback After React setState -> " + line.trim(), " called => base.deepMergeState(", params.obj, ")");
              if (performance) {
                console.info(performance.now());
              }
              console.info("Stacktrace");
              console.info(this.stack);
            }
            console.info("");
            console.info("React Is Done Mutating State From Requested Change", params.obj);
            console.info("The True state of your component: ", params.obj2);
          } else if (params.logSetComponentStateStart === true && window.appState.DeveloperLogState) {
            if (performance) {
              console.info(performance.now());
            }
            this.stack = new Error().stack;
            var line = "";
            if (window.appState.UserAgent.Name != "Safari" && this.stack != undefined && this.stack.split("\n") != undefined && this.stack.split("\n").length > 0) {
              if (this.stack.split("\n")[2].indexOf("makeAssimilatePrototype.js") > -1) {
                line = this.stack.split("\n")[3];
              } else {
                line = this.stack.split("\n")[2];
              }
              console.info("Stacktrace " + this.stack + " called => base.deepMergeState(" + params.obj + ")");
              console.info(this.stack);
            }
            console.info("-----<setComponentState>----");
            console.info("");
            console.info("");
            console.info('State changes proposed:', params.obj);
            // console.info("React Object:", this);
          } else if (params.recieveProps === true && window.appState.DeveloperLogState) {
            var className = "";
            if (params.getClassName != undefined) {
              className = params.getClassName();
            }
            console.info("componentWillReceiveProps (" + className + ")");
            console.info("Old Props", params.props);
            console.info("New Props", params.nextProps);
          } else if (params.logRender === true && window.appState.DeveloperLogState) {
            if (performance) {
              console.info(performance.now());
            }
            if (window.appState.DeveloperLogState) {
              console.info("State Compare:", util.inspect(params.obj.state));
              console.info("Props Compare:", util.inspect(params.obj.props));
            }
            console.info('render invoked full object', params.obj);
          }
        }
      }
    }
  }
}



export default BaseLogger;