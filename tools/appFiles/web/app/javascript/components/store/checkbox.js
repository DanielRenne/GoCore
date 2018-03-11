
import {React,
    BaseComponent,
    Checkbox
  } from "../../globals/forms";
  import Loader from "./loader";
  
  class CheckboxStore extends BaseComponent {
    constructor(props, context) {
      super(props, context);
  
      this.state = {
        loaded: (this.props.value) ? true : false,
        value:this.props.value
      }
    }
  
    componentDidMount() {
      this.subscriptionId = this.store.subscribe(this.props.collection, this.props.id, this.props.path, this.handleValueChange, this.props.value ? false : true);
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
    }
  
    handleValueChange = (data) => {
      if (data == null) {
        return;
      }
      if (data != this.state.value) {
        this.setState({value:data});
      }
    }
  
    render() {
      try {

        if (!this.state.loaded) {
          return (<Loader/>);
        }

        return (
            <Checkbox
                {...this.props}
                checked = {this.props.invert ? !this.state.value : this.state.value}
                onCheck={(event, value) => {
                  this.store.set(this.props.collection, this.props.id, this.props.path, this.props.invert ? !value : value);
                  this.setState({value:this.props.invert ? !value : value});
                }}
            />
        );
      } catch(e) {
        return this.globs.ComponentError("CheckboxStore", e.message, e);
      }
    }
  }
  
  
  CheckboxStore.propTypes = {
    collection:React.PropTypes.string,
    id:React.PropTypes.string,
    path:React.PropTypes.string,
    value:React.PropTypes.bool,
    invert:React.PropTypes.bool
  };
  
  export default CheckboxStore;
  