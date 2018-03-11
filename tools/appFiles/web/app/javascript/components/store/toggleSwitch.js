
import {React,
    BaseComponent,
    Toggle
  } from "../../globals/forms";
  import Loader from "./loader";
  
  class ToggleSwitchStore extends BaseComponent {
    constructor(props, context) {
      super(props, context);
  
      this.state = {
        loaded: (this.props.value !== undefined) ? true : false,
        value:this.props.value
      }
    }
  
    componentDidMount() {
      this.subscriptionId = this.store.subscribe(this.props.collection, 
                                                 this.props.id, 
                                                 this.props.path,
                                                 this.handleValueChange, 
                                                 this.props.value !== undefined ? false : true);
    }
  
    componentWillUnmount() {
      this.store.unsubscribe(this.subscriptionId);
    }

    handleValueChange = (data) => {
      if (data == null) {
        return;
      }
      if (data != this.state.value) {
        this.setState({loaded:true, value:data});
      }
    }
  
    render() {
      try {

        if (!this.state.loaded) {
          return (<Loader/>);
        }

        return (
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
        );
      } catch(e) {
        return this.globs.ComponentError("ToggleSwitchStore", e.message, e);
      }
    }
  }
  
  
  ToggleSwitchStore.propTypes = {
    collection:React.PropTypes.string,
    id:React.PropTypes.string,
    path:React.PropTypes.string,
    value:React.PropTypes.bool,
    invert:React.PropTypes.bool,
    onToggle:React.PropTypes.func
  };
  
  export default ToggleSwitchStore;
  