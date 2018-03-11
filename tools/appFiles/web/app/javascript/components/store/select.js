
import {React,
    BaseComponent,
    SelectField
  } from "../../globals/forms";
  import Loader from "./loader";
  
  class SelectStore extends BaseComponent {
    constructor(props, context) {
      super(props, context);

      this.state = {
        loaded: (this.props.value) ? true : false,
        value:this.props.value,
        errorText:""
      }
    }
  
    componentDidMount() {
      this.subscriptionId = this.store.subscribe(this.props.collection, this.props.id, this.props.path,(data) => {this.handleValueChange(data)}, this.props.value ? false : true);
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
      this.setState({errorText:"", value:data});
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
        );
      } catch(e) {
        return this.globs.ComponentError("SelectStore", e.message, e);
      }
    }
  }
  
  
  SelectStore.propTypes = {
    collection:React.PropTypes.string,
    id:React.PropTypes.string,
    path:React.PropTypes.string,
    value:React.PropTypes.any,
    emptyValue:React.PropTypes.any,
    onChange:React.PropTypes.func
  };
  
  export default SelectStore;
  