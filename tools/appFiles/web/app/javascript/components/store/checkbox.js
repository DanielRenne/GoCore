/**
 * Created by Dan on 12/5/16.
 */
import {React,
    BaseComponent,
    Checkbox
  } from "../../globals/forms";
  import Loader from "./loader";

  class CheckboxStoreComponent extends BaseComponent {
    constructor(props, context) {
      super(props, context);

      this.state = {
        value: undefined
      }
      this.unmounted = false;
    }

    componentDidMount() {
      this.subscriptionId = this.store.subscribe(this.props.collection, this.props.id, this.props.path, this.handleValueChange, true);
    }

    componentWillUnmount() {
      this.unmounted = true;
      this.store.unsubscribe(this.subscriptionId);
    }

    handleValueChange = (data) => {
      if (this.unmounted) {
        return;
      }
      if (data == null) {
        return;
      }
      if (data != this.state.value) {
        this.setState({value:data});
      }
    }

    render() {
      try {
        return (
            <span>
              {this.state.value == undefined? <Loader/>: null}
              <span style={{display: this.state.value == undefined ? "none" : "block"}}>
                <Checkbox
                    {...this.props}
                    checked = {this.props.invert ? !this.state.value : this.state.value}
                    onCheck={(event, value) => {
                      this.store.set(this.props.collection, this.props.id, this.props.path, this.props.invert ? !value : value);
                      this.setState({value:this.props.invert ? !value : value});
                    }}
                />
              </span>
            </span>
        );
      } catch(e) {
        return this.globs.ComponentError("CheckboxStore", e.message, e);
      }
    }
  }


  CheckboxStoreComponent.propTypes = {
    collection:React.PropTypes.string,
    id:React.PropTypes.string,
    path:React.PropTypes.string,
    invert:React.PropTypes.bool
  };

  export default CheckboxStoreComponent;
