
import {React,
    BaseComponent,
    TextField
  } from "../../globals/forms";
  import Loader from "./loader";

  
  class TextFieldStore extends BaseComponent {
    constructor(props, context) {
      super(props, context);
  
      this.state = {
        loaded: (this.props.value) ? true : false,
        value:this.props.value,
        errorText:""
      };

      this.changing = false;
      this.errorTimeout;
      this.changingTimeout;
    }
  
    componentDidMount() {
      this.subscriptionId = this.store.subscribe(this.props.collection, this.props.id, this.props.path,(data) => {this.handleValueChange(data)}, this.props.value ? false : true)
      if (this.props.validateErrorMessage) {
        this.subscriptionErrorId = this.store.subscribe(this.props.collection, this.props.id, "Errors." + this.props.path, (data) => {this.handleErrorValueChange(data)}, false);
      }
      
      if (!this.props.value) {
        this.store.getByPath({"collection":this.props.collection, 
                              "id":this.props.id, 
                              "path":this.props.path}, (data) => {
          this.setState({loaded:true, value:data});
        });
      }
    }
  
    componentWillUnmount() {
      this.store.unsubscribe(this.subscriptionId);
      if (this.props.validateErrorMessage) {
        this.store.unsubscribe(this.subscriptionErrorId);
      }
    }
  
    handleValueChange(data) {
      if (this.changing) {
        return;
      }
      if (data == null) {
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

        if (!this.state.loaded) {
          return (<Loader/>);
        }
        return (
            <TextField
              {...this.props}
              value={this.state.value.toString()}
              onChange={(event) => {
                clearTimeout(this.changingTimeout);

                let t = typeof(this.props.value);
                let value = event.target.value;
                if (t == "number" || this.props.isNumber === true) {
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

                let t = typeof(this.props.value);
                let value = event.target.value;
                if (t == "number" || this.props.isNumber === true) {
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
        );
      } catch(e) {
        return this.globs.ComponentError("TextField", e.message, e);
      }
    }
  }
  
  
  TextFieldStore.propTypes = {
    collection:React.PropTypes.string,
    id:React.PropTypes.string,
    path:React.PropTypes.string,
    value:React.PropTypes.oneOfType([
      React.PropTypes.string,
      React.PropTypes.number
    ]),
    changeOnBlur:React.PropTypes.bool,
    isNumber:React.PropTypes.number
  };

  TextFieldStore.defaultProps = {
    changeOnBlur: true
  };

  export default TextFieldStore;
  