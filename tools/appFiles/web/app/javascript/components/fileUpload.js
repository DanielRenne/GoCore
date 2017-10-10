/**
 * Created by Dan on 11/14/16.
 */
import React, {Component} from "react";
import BaseComponent from "../components/base";
import {DropZone} from "../globals/forms";


class FileUpload extends BaseComponent {
  constructor(props, context) {
    super(props, context);

    this.state = {
      files: [],
      dragging: false,
      fileId: this.props.fileId
    };

    if (!this.props.hasOwnProperty("width") && this.props.hasOwnProperty("imageWidth")) {
      this.state.width = this.props.imageWidth;
    }
    if (!this.props.hasOwnProperty("height") && this.props.hasOwnProperty("imageHeight")) {
      this.state.height = this.props.imageHeight;
    }

    this.onComplete = this.props.onComplete;
    this.onFileChange = this.props.onFileChange;
    this.onError = this.props.onError;

    this.handleOnComplete = (fileObject) => {
      this.setComponentState({fileId:fileObject.FileObject.Id});
      if (this.onComplete != undefined){
        this.onComplete(fileObject.FileObject);
      }
    };

    this.handleOnError = () => {
      if (this.onError != undefined) {
        this.onError();
      }
    }

  }

  onDrop(files) {
    if (this.onFileChange != undefined){
      this.onFileChange();
    }
    window.api.upload({fileId:this.state.fileId, width: this.props.imageWidth, height: this.props.imageHeight, file: files[0], callback:this.handleOnComplete, error: this.handleOnError})
  }

  onDragEnter() {
    this.setComponentState({dragging:true});
  }

  onDragLeave() {
    this.setComponentState({dragging:false});
  }

  componentWillReceiveProps(nextProps) {
      this.setComponentState({fileId:nextProps.fileId});
  }

	render() {
    try {
      var uploaderStyle = {
        cursor: "pointer",
        display: "inline-block",
        width: this.state.width,
        height: this.state.height
      }

      this.logRender();

      return (
        <DropZone onDragEnter={() => this.onDragEnter()} onDragLeave={() => this.onDragLeave()}
                  onDrop={(files) => this.onDrop(files)} style={uploaderStyle}
                  inputProps={{name: "file_upload", id: "file_upload"}}>
          {this.props.children}
        </DropZone>
      );
    } catch(e) {
      return this.globs.ComponentError(this.getClassName(), e.message);
    }
	}
}

FileUpload.propTypes = {
  fileId: React.PropTypes.string,
  onComplete: React.PropTypes.func,
  onError: React.PropTypes.func,
  onFileChange: React.PropTypes.func,
  width: React.PropTypes.number,
	height: React.PropTypes.number,
	imageWidth: React.PropTypes.number,
	imageHeight: React.PropTypes.number
};

FileUpload.defaultProps = {
  width:0,
  height:0,
	imageWidth:0,
	imageHeight:0
};

export default FileUpload;
