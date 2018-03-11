
import {React,
    BaseComponent,
    Checkbox
  } from "../../globals/forms";
  import Loader from "./loader";
import checkbox from "./checkbox";
  
  class CheckboxArrayStore extends BaseComponent {
    constructor(props, context) {
      super(props, context);
  
      this.state = {
        loaded: (this.props.items) ? true : false,
        items:this.props.items
      }
    }

    isChecked(key) {
      for(var i = 0; i < this.state.items.length; i++) {
        let item = this.state.items[i];
        if (item == key) {
          return true;
        }
      }
      return false;
    }
  
    componentDidMount() {
      this.subscriptionId = this.store.subscribe(this.props.collection, this.props.id, this.props.path,(data) => {this.handleValueChange(data)}, this.props.items ? false : true);
      if (!this.props.items) {
        this.store.getByPath({"collection":this.props.collection, 
                              "id":this.props.id, 
                              "path":this.props.path}, (data) => {
          this.setState({loaded:true, items:data});
        });
      }
    }
  
    componentWillUnmount() {
      this.store.unsubscribe(this.subscriptionId);
    }
  
    handleValueChange(data) {
      this.setState({items:data});
    }
  
    render() {
      try {

        if (!this.state.loaded) {
          return (<Loader/>);
        }

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

                    if (this.state.items == null) {
                      this.state.items = [];
                    }
                
                    for(var i = 0; i < this.state.items.length; i++) {
                      if (keyValue.key == this.state.items[i] && value == true) {
                        return;
                      }
                      if (keyValue.key == this.state.items[i] && value == false) {
                        this.state.items.splice(i, 1);
                        break;
                      }
                    }
                
                    if (value == true) {
                      this.state.items.push(keyValue.key);
                    }

                    this.store.set(this.props.collection, this.props.id, this.props.path, this.state.items);
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
            {checkboxes}
          </div>
        );
      } catch(e) {
        return this.globs.ComponentError("CheckboxArrayStore", e.message, e);
      }
    }
  }
  
  CheckboxArrayStore.propTypes = {
    collection:React.PropTypes.string,
    id:React.PropTypes.string,
    path:React.PropTypes.string,
    items:React.PropTypes.array,
    keyValues:React.PropTypes.array
  };
  
  export default CheckboxArrayStore;
  