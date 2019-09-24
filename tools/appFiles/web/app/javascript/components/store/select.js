/**
 * Created by Dan on 12/5/16.
 */
import {React,
    BaseComponent,
    SelectField
  } from "../../globals/forms";
  import Loader from "./loader";

  class SelectStoreComponent extends BaseComponent {
    constructor(props, context) {
      super(props, context);
      this.unmounted = false;

      this.state = {
        value: undefined,
        errorText:""
      }
    }

    componentDidMount() {
      this.subscriptionId = this.store.subscribe(this.props.collection, this.props.id, this.props.path,(data) => {this.handleValueChange(data)}, true);
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
            {this.state.value == undefined? <Loader/>: null}
            <span style={{display: this.state.value == undefined ? "none" : "block"}}>
              <select
                {...this.props}
                value={(this.props.emptyValue) ? (this.state.value) ? this.state.value : this.props.emptyValue : this.state.value}
                onChange={(e) => {
                  this.store.set(this.props.collection, this.props.id, this.props.path, e.target.value);
                  if (this.props.onChange !== undefined) {
                    this.props.onChange(e);
                  }
                }}
                errorText={this.globs.translate(this.state.errorText)}
              >
                {this.props.children}
              </select>
            </span>
          </span>
        );
      } catch(e) {
        return this.globs.ComponentError("SelectStore", e.message, e);
      }
    }
  }


  SelectStoreComponent.propTypes = {
    collection:React.PropTypes.string,
    id:React.PropTypes.string,
    path:React.PropTypes.string,
    emptyValue:React.PropTypes.any,
    onChange:React.PropTypes.func
  };

  export default SelectStoreComponent;
