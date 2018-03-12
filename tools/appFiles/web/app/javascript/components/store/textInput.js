
import {React,
    BaseComponent,
    TextField,
  } from "../../globals/forms";
  import Loader from "./loader";
  
  class TextInputStore extends BaseComponent {
    constructor(props, context) {
      super(props, context);
  
      this.state = {
        loaded: (this.props.value) ? true : false,
        value:this.props.value,
        errorText:""
      }

      this.changing = false;
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
        this.setState({errorText:"", value:data});
      }
    }

    handleErrorValueChange(data) {
      this.setState({errorText:data});
    }
  
    render() {
      try {

        if (!this.state.loaded) {
          return (<Loader/>);
        }

        return (
          <TextField
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
            defaultValue={this.state.value}
            {...this.props.property}
          />
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
  
  
  TextInputStore.propTypes = {
    collection:React.PropTypes.string,
    id:React.PropTypes.string,
    path:React.PropTypes.string,
    value:React.PropTypes.oneOfType([
      React.PropTypes.string,
      React.PropTypes.number
    ])
  };
  
  export default TextInputStore;
  