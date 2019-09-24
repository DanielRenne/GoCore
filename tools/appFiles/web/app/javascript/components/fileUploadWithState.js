import React from "react";
import BaseComponent from "../components/base";
import {deepOrange600, FileUpload, LinearProgress, RaisedButton} from "../globals/forms";


class FileUploadWithState extends BaseComponent {
	constructor(props, context) {
		super(props, context);

		this.state = {ImagePicked: (this.parentState[this.props.rowKey][this.props.fieldKey] != "") ? true:false};

		this.handleClear = () => {
			let changes = {};
			changes[this.props.rowKey] = {};
			changes[this.props.rowKey][this.props.fieldKey] = "";
			this.setComponentState({ImagePicked: false});
			this.setParentState(changes);
		};

		this.handleFileChange = (FileObject) => {
			window.displayErrorAlert = true;
			let changes = {};
			changes["cannotSubmit"] = false;
			changes[this.props.rowKey] = {};
			changes[this.props.rowKey][this.props.fieldKey] = FileObject.Id;
			this.setComponentState({ImagePicked: true});
			this.setParentState(changes);
		};

		this.handleImageChangeStart = () => {
			window.displayErrorAlert = true;
			this.setParentState({cannotSubmit: true});
		};

		this.handleFileChangeError = () => {
			window.displayErrorAlert = true;
			this.setParentState({cannotSubmit: false});
		};
	}

	render() {
    try {
      this.logRender();

      var image = (
        <div>
          <div style={{width: '100%', padding: 5, border: "5px #EBEBEB dotted"}}>
            <div className="text-center">
              <img
                src={"/fileObject/" + this.parentState[this.props.rowKey][this.props.fieldKey] + "?" + new Date().getTime()}
                style={{maxWidth: '100%', maxHeight: '100%', objectFit:"contain"}}/>
            </div>
          </div>
          <div>
              <RaisedButton
              style={{display: 'block', marginTop: 5}}
              label={window.appContent.ClearImage}
              secondary={true}
              onClick={this.handleClear}
            />
          </div>
        </div>
      );

      if (this.parentState[this.props.rowKey][this.props.fieldKey] == "") {
        image = "";
      }

      var htmlUpload = <span>
          <span style={{display: !this.parentState.cannotSubmit ? "inline" : "none"}}>
            {!this.state.ImagePicked &&
            <FileUpload
              width={this.props.width}
              height={this.props.height}
              imageWidth={this.props.imageWidth}
              imageHeight={this.props.imageHeight}
              fileId={this.parentState[this.props.rowKey][this.props.fieldKey]}
              onComplete={this.handleFileChange}
              onFileChange={this.handleImageChangeStart}
              onError={this.handleFileChangeError}>
              <div style={{textAlign: 'center', width: '100%', marginTop: 10, border: "5px #EBEBEB dotted"}}>
                <RaisedButton
                  style={{marginTop: 68, marginBottom: 68}}
                  label={"* " + window.appContent.UploadOrDrag}
                  secondary={true}
                />
                {this.parentState[this.props.rowKey]["Errors"][this.props.fieldKey] ? <div style={{
                  color: 'red',
                  marginLeft: 15,
                  marginTop: 6
                }}>{window.appContent.GenericFileUploadError}</div> : null}
              </div>
            </FileUpload>
            }
            {this.state.ImagePicked && image}
          </span>
        </span>;

      return (
        <div style={{marginBottom: 25}}>
          {htmlUpload}
          {this.parentState.cannotSubmit ? <LinearProgress color={deepOrange600} mode="indeterminate"/> : null}
        </div>
      );
    } catch(e) {
      return this.globs.ComponentError("FileUploadWithState", e.message. e);
    }
	}
}

FileUploadWithState.propTypes = {
	parent: React.PropTypes.object,
	rowKey: React.PropTypes.string,
	width: React.PropTypes.number,
	height: React.PropTypes.number,
	fieldKey: React.PropTypes.string,
	imageWidth: React.PropTypes.number,
	imageHeight: React.PropTypes.number
};

FileUploadWithState.defaultProps = {
	width: 0,
	height: 0,
	imageWidth:0,
	imageHeight:0
};

export default FileUploadWithState;
