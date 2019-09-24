/**
 * Created by Dan on 12/5/16.
 */
import {React,
    BaseComponent,
    Checkbox
  } from "../../globals/forms";
  import Loader from "./loader";
import checkbox from "./checkbox";

  class CheckboxArrayStoreComponent extends BaseComponent {
    constructor(props, context) {
      super(props, context);

      this.state = {
        value: undefined
      }
      this.unmounted = false;
    }

    isChecked(key) {
      if (this.state.value != undefined) {
        for(var i = 0; i < this.state.value.length; i++) {
          let item = this.state.value[i];
          if (item == key) {
            return true;
          }
        }
      }
      return false;
    }

    componentDidMount() {
      this.subscriptionId = this.store.subscribe(this.props.collection, this.props.id, this.props.path,(data) => {this.handleValueChange(data)}, true);
    }

    componentWillUnmount() {
      this.unmounted = true;
      this.store.unsubscribe(this.subscriptionId);
    }

    handleValueChange(data) {
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

        let checkboxes = [];

        if (this.props.keyValues) {
          this.props.keyValues.map((kv, i) => {
            let keyValue = this.props.keyValues[i];
            checkboxes.push(
              <span key={window.globals.guid()}>
                <Checkbox
                  {...this.props}
                  checked = {this.isChecked(keyValue.key)}
                  label = {keyValue.value}
                  onCheck={(event, value) => {

                    if (this.state.value == undefined || this.state.value == null) {
                      this.state.value = [];
                    }

                    for(var i = 0; i < this.state.value.length; i++) {
                      if (keyValue.key == this.state.value[i] && value == true) {
                        return;
                      }
                      if (keyValue.key == this.state.value[i] && value == false) {
                        this.state.value.splice(i, 1);
                        break;
                      }
                    }

                    if (value == true) {
                      this.state.value.push(keyValue.key);
                    }

                    this.store.set(this.props.collection, this.props.id, this.props.path, this.state.value);
                  }}
                />
                <div style={{height: 10}}/>
              </span>
            );
          });
        } else {
          return null;
        }


        return (
          <div>
            <span>
              {this.state.value == undefined? <Loader/>: null}
              <span style={{display: this.state.value == undefined ? "none" : "block"}}>
                {checkboxes}
              </span>
            </span>
          </div>
        );
      } catch(e) {
        return this.globs.ComponentError("CheckboxArrayStore", e.message, e);
      }
    }
  }

  CheckboxArrayStoreComponent.propTypes = {
    collection:React.PropTypes.string,
    id:React.PropTypes.string,
    path:React.PropTypes.string,
    keyValues:React.PropTypes.array
  };

  export default CheckboxArrayStoreComponent;
