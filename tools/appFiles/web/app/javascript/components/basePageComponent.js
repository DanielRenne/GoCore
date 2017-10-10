import BaseComponent from './base'

class BasePageComponent extends BaseComponent {
  constructor(props, context) {
    super(props, context);
    if (window.appState.DeveloperMode && window.pageState.hasOwnProperty("DeveloperLog")) {
      this.state = {};
      $.each(window.pageState, (k, v) => {
        if (k != "DeveloperLog") {
          this.state[k] = v;
        }
      });
    } else {
      this.state = window.pageState;
    }
    this.uriParams = this.globs.GetUriParams();
    if (window.appState.DeveloperLogState) {
      console.info("initialPageState", window.pageState);
    }
    window.goCore.setStateFromExternal = (state, cb) => {
      window.pageState = state;
      this.stack = new Error().stack;
      this.setComponentState(window.pageState, () => {
        if (typeof(cb) == "function") {
          cb();
        }
        if (window.appState.DeveloperLogState) {
          if (this.stack) {
            console.info("setStateFromExternal " + this.stack.split("\n")[1].trim(), this);
          }
          console.info("state ==>> ", this.state);
        }
      });
    };
  }
}

export default BasePageComponent;
