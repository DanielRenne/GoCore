/**
 * Created by Dan on 12/5/16.
 */
import {React,
    BaseComponent,
    TextField,
  } from "../../globals/forms";
  import Loader from "./loader";

  class TextInputStoreComponent extends BaseComponent {
    constructor(props, context) {
      super(props, context);
      this.unmounted = false;

      this.state = {
        value: undefined,
        errorText:""
      }

      this.changing = false;
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
      if (this.changing) {
        return;
      }
      if (data == null) {
        return;
      }
      if (data != this.state.value) {
        this.setState({errorText:"", value:data});
      }
    }

    handleErrorValueChange(data) {
      this.setState({errorText:data});
    }

    render() {
      try {
        return (
          <span>
            {this.state.value == undefined ? <Loader/>: null}

            <span style={{display: this.state.value == undefined ? "none" : "block"}}>
              <TextField
                onChange={(event) => {
                  clearTimeout(this.changingTimeout);
                  this.changing = true;

                  let value = event.target.value;
                  if (this.props.isNumber === true) {
                    value = Number(event.target.value)
                    if (isNaN(value)) {
                      this.changing = false;
                      return
                    }
                  }

                  this.store.set(this.props.collection, this.props.id, this.props.path, value);
                  this.setState({value:value});
                  this.changingTimeout = setTimeout(() => {
                    this.changing = false;
                  }, 1500);
                }}
                value={this.state.value}
                {...this.props.property}
              />
            </span>
          </span>
        );

        {/* <input
              type="text"
              className="form-control"
              value = {this.state.value.toString()}
              onChange={(event) => {
                clearTimeout(this.changingTimeout);
                this.changing = true;

                let t = typeof(this.props.value);
                let value = event.target.value;
                if (t == "number") {
                  value = Number(event.target.value)
                  if (isNaN(value)) {
                    this.changing = false;
                    return
                  }
                }

                this.store.set(this.props.collection, this.props.id, this.props.path, value);
                this.setState({value:value});
                this.changingTimeout = setTimeout(() => {
                  this.changing = false;
                }, 1500);
              }}
          /> */}

      } catch(e) {
        return this.globs.ComponentError("TextInputStore", e.message, e);
      }
    }
  }


  TextInputStoreComponent.propTypes = {
    collection:React.PropTypes.string,
    id:React.PropTypes.string,
    isNumber:React.PropTypes.bool,
    path:React.PropTypes.string
  };

  export default TextInputStoreComponent;
