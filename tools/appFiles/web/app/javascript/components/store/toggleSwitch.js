/**
 * Created by Dan on 12/5/16.
 */
import {React,
    BaseComponent,
    Toggle
  } from "../../globals/forms";
  import Loader from "./loader";

  class ToggleStoreComponent extends BaseComponent {
    constructor(props, context) {
      super(props, context);

      if (props.path == "DisableHomePage") {
        session_functions.VelocityDump("init")
      }
      this.unmounted = false;
      this.state = {
        value: undefined
      }
    }

    componentDidMount() {
      if (this.props.path == "DisableHomePage") {
        session_functions.VelocityDump("mounted")
      }
      this.subscriptionId = this.store.subscribe(this.props.collection,
                                                 this.props.id,
                                                 this.props.path,
                                                 this.handleValueChange,
                                                 true);
    }

    componentWillUnmount() {
      if (this.props.path == "DisableHomePage") {
        session_functions.VelocityDump("unmounted")
      }
      this.store.unsubscribe(this.subscriptionId);
      this.unmounted = true;
    }

    handleValueChange = (data) => {
      if (this.unmounted) {
        return;
      }
      if (data == null) {
        return;
      }
      if (data != this.state.value) {
        if (this.props.path == "DisableHomePage") {
          session_functions.VelocityDump("value set", data, this.state.value)
        }
        this.setState({value:data});
      }
    }

    render() {
      if (this.props.path == "DisableHomePage") {
        session_functions.VelocityDump("rendered home page", JSON.stringify(this.state, null, 2), JSON.stringify(this.props, null, 2))
      }
      try {

        return (
          <span>
              {this.state.value == undefined? <Loader/>: null}
              <span style={{display: this.state.value == undefined ? "none" : "block"}}>

                <Toggle
                    {...this.props}
                    toggled = {this.props.invert ? !this.state.value : this.state.value}
                    onToggle={(event, value) => {
                      this.store.set(this.props.collection, this.props.id, this.props.path, this.props.invert ? !value : value);
                      this.setState({value:this.props.invert ? !value : value}, () => {
                        if (this.props.onToggle !== undefined) {
                          this.props.onToggle(event, value);
                        }
                      });
                    }}
                />
              </span>
          </span>
        );
      } catch(e) {
        return this.globs.ComponentError("ToggleSwitchStore", e.message, e);
      }
    }
  }


  ToggleStoreComponent.propTypes = {
    collection:React.PropTypes.string,
    id:React.PropTypes.string,
    path:React.PropTypes.string,
    invert:React.PropTypes.bool,
    onToggle:React.PropTypes.func
  };

  export default ToggleStoreComponent;
