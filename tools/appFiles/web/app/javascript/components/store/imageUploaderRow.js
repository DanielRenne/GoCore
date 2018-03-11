
import {React,
    BaseComponent,
    RaisedButton,
    deepOrange600,
    red400,
    FileUpload,
    LinearProgress,
    Row,
    Col,
    IconButton
  } from "../../globals/forms";
  import {DeleteIcon} from "../../globals/icons";
  import Loader from "./loader";
  
  class ImageUploaderRowStore extends BaseComponent {
    constructor(props, context) {
      super(props, context);

      this.changedImage = false;

      this.state = {
        loaded: (this.props.value) ? true : false,
        cannotSubmit:false,
        value:this.props.value
      }

      this.handleClearFile = () => {
        this.changedImage = true;
        this.store.set(this.props.collection, this.props.id, this.props.path, "");
      };

      this.handleFileChange = (FileObject) => {
        window.displayErrorAlert = true;
        this.changedImage = true;
        this.store.set(this.props.collection, this.props.id, this.props.path, FileObject.Id);
        this.setState({cannotSubmit:false});
      };

      this.handleFileChangeStart = () => {
        window.displayErrorAlert = true;
        this.setState({cannotSubmit: true});
      };

      this.handleFileChangeError = () => {
        window.displayErrorAlert = true;
        this.setState({cannotSubmit: false});
      };

    }
  
    componentDidMount() {
      this.subscriptionId = this.store.subscribe(this.props.collection, 
                                                 this.props.id, 
                                                 this.props.path,
                                                 (data) => {
        this.handleValueChange(data)
      }, this.props.value ? false : true);
    }
  
    componentWillUnmount() {
      this.store.unsubscribe(this.subscriptionId);
    }
  
    handleValueChange(data) {
      this.setState({loaded:true, value:""}, () => {
        this.setState({loaded:true, value:data});
      });
    }
  
    render() {
      try {

        if (!this.state.loaded) {
          return (<Loader/>);
        }

        var image = (
          <div key={window.globals.guid()}>
            <img src={"/fileObject/" + this.state.value + ((this.changedImage === true) ? "?" + new Date().getTime() : "")} 
            style={{marginLeft:5, width:36, objectFit:"contain", marginTop:-5}} />
            <IconButton onClick={this.handleClearFile} tooltip={window.appContent.ClearImage}>
              <DeleteIcon width={20} height={20} color={red400}  />
            </IconButton>
          </div>
        );

        if (this.state.value == "") {
          image = (this.props.defaultImage === undefined) ? null : this.props.defaultImage;
        }

        var htmlUpload = (
        <span key={window.globals.guid()}>
          <span key={window.globals.guid()} className="AlignerLeft" style={{display: !this.state.cannotSubmit ? "flex": "none"}}>
            <FileUpload fileId={this.state.value}
                        onComplete={this.handleFileChange}
                        onFileChange={this.handleFileChangeStart}
                        onError={this.handleFileChangeError}
                        width={this.props.width} 
                        height={this.props.height}
                        disableSpinner={true}> 
              <div style={{height: this.props.height, width: this.props.width, border:"5px #EBEBEB dotted"}}>
                <div>
                  <RaisedButton
                      style={{marginLeft: 5, marginTop: 5, height:30}}
                      label={"* " + window.appContent.Upload}
                      secondary={true}
                  />
                  {/* {this.state.Room.Errors.BackgroundImage ? <div style={{color:'red', marginLeft: 15, marginTop: 6}}>{window.appContent.GenericFileUploadError}</div> : null} */}
                </div>
              </div>
            </FileUpload>
            {image}              
            </span>
          </span>               
        );

        return (

          <Row
            style={this.props.style}
          >
            {htmlUpload}
            {this.state.cannotSubmit ? <LinearProgress color={deepOrange600} mode="indeterminate" /> : null}
          </Row>
        );
      } catch(e) {
        return this.globs.ComponentError("ImageUploaderRowStore", e.message, e);
      }
    }
  }
  
  
  ImageUploaderRowStore.propTypes = {
    collection:React.PropTypes.string,
    id:React.PropTypes.string,
    path:React.PropTypes.string,
    value:React.PropTypes.string,
    height:React.PropTypes.any,
    width:React.PropTypes.any,
    defaultImage:React.PropTypes.any
  };
  
  export default ImageUploaderRowStore;
  