/**
 * Created by Dan on 12/5/16.
 */
import {React,
    BaseComponent,
    TextField
  } from "../../globals/forms";
  import Loader from "./loader";


  class TextFieldStoreComponent extends BaseComponent {
    constructor(props, context) {
      super(props, context);

      this.state = {
        value: undefined,
        errorText:""
      };

      this.changing = false;
      this.unmounted = false;
      this.errorTimeout;
      this.changingTimeout;
    }

    componentDidMount() {
      this.subscriptionId = this.store.subscribe(this.props.collection, this.props.id, this.props.path,(data) => {this.handleValueChange(data)}, true)
      if (this.props.validateErrorMessage) {
        this.subscriptionErrorId = this.store.subscribe(this.props.collection, this.props.id, "Errors." + this.props.path, (data) => {this.handleErrorValueChange(data)}, false);
      }
    }

    componentWillUnmount() {
      this.unmounted = true;
      this.store.unsubscribe(this.subscriptionId);
      if (this.props.validateErrorMessage) {
        this.store.unsubscribe(this.subscriptionErrorId);
      }
    }

    handleValueChange(data) {
      if (this.unmounted) {
        return;
      }
      if (data == null) {
        return;
      }
      if (this.changing) {
        return;
      }
      if (data != this.state.value) {
        this.setState((state) => {
          state.value = data;
          return state;
        });
        window.setTimeout(() => {
          this.setState((state) => {
            state.errorText = "";
            return state;
          });
        }, 10000);
      }
    }

    handleErrorValueChange(data) {
      this.setState((state) => {
        state.errorText = data;
        return state;
      });
    }

    render() {
      try {
        return (
          <span>
              {this.state.value == undefined? <Loader/>: null}
              <span style={{display: this.state.value == undefined ? "none" : "block"}}>
                <TextField
                  {...this.props}
                  value={this.state.value != undefined ? this.state.value.toString(): this.state.value}
                  onChange={(event) => {
                    clearTimeout(this.changingTimeout);

                    let value = event.target.value;
                    if (this.props.isNumber === true) {
                      value = Number(event.target.value)
                      if (isNaN(value)) {
                        this.changing = false;
                        return;
                      }
                    }



                    this.setState({value:value});
                    if (this.props.changeOnBlur === false) {
                      this.changing = true;
                      this.store.set(this.props.collection, this.props.id, this.props.path, value);
                      this.changingTimeout = window.setTimeout(() => {
                        this.changing = false;
                      }, 1500);
                    }
                  }}
                  onBlur={(event) => {
                    if (this.props.changeOnBlur === false) {
                      return;
                    }

                    let value = event.target.value;
                    if (this.props.isNumber === true) {
                      value = Number(event.target.value)
                      if (isNaN(value)) {
                        return;
                      }
                    }

                    this.store.set(this.props.collection, this.props.id, this.props.path, value);
                  }

                  }
                  errorText={this.globs.translate(this.state.errorText)}
              />
            </span>
          </span>
        );
      } catch(e) {
        return this.globs.ComponentError("TextField", e.message, e);
      }
    }
  }


  TextFieldStoreComponent.propTypes = {
    collection:React.PropTypes.string,
    id:React.PropTypes.string,
    path:React.PropTypes.string,
    changeOnBlur:React.PropTypes.bool,
    isNumber:React.PropTypes.bool
  };

  TextFieldStoreComponent.defaultProps = {
    changeOnBlur: true
  };

  export default TextFieldStoreComponent;
