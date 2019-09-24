/**
 * Created by Dan on 12/5/16.
 */
import {React,
    BaseComponent,
    RaisedButton,
    deepOrange600,
    FileUpload,
    Row,
    Col
  } from "../../globals/forms";
  import Loader from "./loader";

  class ImageUploaderStoreComponent extends BaseComponent {
    constructor(props, context) {
      super(props, context);

      this.changedImage = false;
      this.unmounted = false;

      this.state = {
        cannotSubmit:false,
        value: undefined
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

        var image = (
          <Col xs={12} md={6} style={{padding:10}}>
            <img src={"/fileObject/" + this.state.value + ((this.changedImage === true) ? "?" + new Date().getTime() : "")}
                 style={{width:this.props.width, height:this.props.height, objectFit:"contain"}} />
            <br/>
            <RaisedButton
              label={window.appContent.ClearImage}
              secondary={true}
              onClick={this.handleClearFile}
            />
          </Col>
        );

        if (this.state.value == "") {
          image = (this.props.defaultImage === undefined) ? null : this.props.defaultImage;
        }

        var htmlUpload = (
        <span key={window.globals.guid()}>
          <span key={window.globals.guid()} className="AlignerLeft" style={{display: !this.state.cannotSubmit ? "flex": "none"}}>
            <Row>
              <Col xs={12} md={6} style={{padding:10}}>
                <FileUpload fileId={this.state.value}
                            onComplete={this.handleFileChange}
                            onFileChange={this.handleFileChangeStart}
                            onError={this.handleFileChangeError}
                            width={this.props.width}
                            height={this.props.height}>
                  <div style={{height: this.props.height, width: this.props.width, paddingLeft:20, paddingRight:20, border:"5px #EBEBEB dotted"}}>
                    <div>
                      <RaisedButton
                          style={{marginLeft: 30, marginTop: 68}}
                          label={"* " + window.appContent.UploadOrDrag}
                          secondary={true}
                      />
                      {/* {this.state.Room.Errors.BackgroundImage ? <div style={{color:'red', marginLeft: 15, marginTop: 6}}>{window.appContent.GenericFileUploadError}</div> : null} */}
                    </div>
                  </div>
                </FileUpload>
              </Col>
              {image}
            </Row>
          </span>
          </span>
        );

        return (

          <div>
            <span>
              {this.state.value == undefined? <Loader/>: null}
              <span style={{display: this.state.value == undefined ? "none" : "block"}}>
                {htmlUpload}
                {this.state.cannotSubmit ? <LinearProgress color={deepOrange600} mode="indeterminate" /> : null}
              </span>
            </span>
          </div>
        );
      } catch(e) {
        return this.globs.ComponentError("ImageUploaderStore", e.message, e);
      }
    }
  }


  ImageUploaderStoreComponent.propTypes = {
    collection:React.PropTypes.string,
    id:React.PropTypes.string,
    path:React.PropTypes.string,
    height:React.PropTypes.any,
    width:React.PropTypes.any,
    defaultImage:React.PropTypes.any
  };

  export default ImageUploaderStoreComponent;
