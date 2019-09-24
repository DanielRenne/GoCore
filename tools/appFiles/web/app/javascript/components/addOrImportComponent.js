import React, {Component} from "react";
import ActionButton from "./actionButton";
import {ConfirmPopup, DropZone, orange900, green900} from "../globals/forms";


export default function AddOrImportPage(Component, Action, Tooltip, DownloadTemplate, Controller, ControllerAction, MiddleColumnFunc, ShowUpload=true, ShowViaForm=true, TitleOption=1, ImportController="", JsonCallback=null, JsonCallbackUriParams=null, DisplayIcon=true, ExecHasOnlyOneOption=true) {
  if (JsonCallback == null) {
    var JsonCallback = (input) => {
      return input
    }
  }
  if (JsonCallbackUriParams == null) {
    var JsonCallbackUriParams = (input) => {
      return input
    }
  }
  if (ImportController == "") {
    ImportController = Controller;
  }
  class AddOrImportPageComponent extends Component {
    constructor(props, context) {
      super(props, context);
      this.page = null;
      this.state = {
        files: [],
        dragging: false,
        addPageAction: Action,
        addPageController: Controller,
        buttonVisible: Tooltip == "AddRoom" ? false:true,
        uploadPageAction: ControllerAction,
        uploadPageController: ImportController,
        jsonCallback: JsonCallback,
        jsonCallbackUriParams: JsonCallbackUriParams,
        title: ""
      };

      if (Action == "CONTROLLER_USERADD" && !this.globs.HasRole("USER_ADD")) {
        ShowViaForm = false;
      }

      this.hasOnlyOneOption = (ExecHasOnlyOneOption && MiddleColumnFunc && MiddleColumnFunc() == null && ((!ShowUpload && ShowViaForm) || (ShowUpload && !ShowViaForm)));
      this.addOrImportRef;
      this.completionPopup;
    }

    componentDidMount() {
      var title = "";
      if (TitleOption == 1) {
        title = window.appContent.ImportAndUploadTitle;
      } else if (TitleOption == 2) {
        title = window.appContent.InviteOrAddViaForm;
      } else if (TitleOption == 3) {
        title = window.appContent.ImportAndUploadTitleJSON;
      } else if (TitleOption == 4) {
         title = window.appContent.ImportFile;
      }

      this.setComponentState({
        title: title
      });

      window.clickCurrentAddOrImportActionButton = () => {
        this.open();
      }
    }

    open() {
      if (this.hasOnlyOneOption) {
        window.global.functions.redirect(this.state.addPageAction);
      } else {
        this.addOrImportRef.handleOpen();
      }
    }

    openImport() {
      this.addOrImportRef.handleOpen();
    }

    onDrop(files) {
      this.setState({
        files: files,
        dragging: false
      });

      var reader = new FileReader();
      reader.onload = (ev) => {
          try {
              var jsonObj = {};
              jsonObj.FileUpload = {};
              jsonObj.FileUpload.Name = files[0].name;
              jsonObj.FileUpload.Content = window.atob(ev.target.result.substr(ev.target.result.indexOf("base64,") + 7));
              jsonObj.FileUpload.Size = files[0].size;
              jsonObj.FileUpload.Type = files[0].type;
              jsonObj.FileUpload.Modified = files[0].lastModifiedDate;
              jsonObj.FileUpload.ModifiedUnix = files[0].lastModified;
              jsonObj = this.state.jsonCallback(jsonObj);
              window.api.post({action: this.state.uploadPageAction, state: jsonObj, controller: this.state.uploadPageController, leaveStateAlone: true, callback: (vm) => {
                this.setComponentState({FileUpload: vm.FileUpload}, () => {
                  this.onClosePopup();

                  if (vm.FileUpload && vm.FileUpload.Meta.CompleteFailure) {
                    this.completionPopup.handleClose();
                    if (this.state.uploadPageController == "siteImport") {
                      this.globs.PopupWindow(window.appContent.ImportFailedDescription + "\n\n" + vm.FileUpload.Meta.FileErrors.join("\n") + "\n\nId Row Information:\n\n" + vm.FileUpload.Meta.RowsSkippedDetails.join("\n"), "ImportFailed");
                    }
                  } else {
                    this.completionPopup.handleOpen();
                  }
                });
              }});
          }
          catch (ex) {
            core.Debug.Dump(ex);
            try {
              this.close();
            } catch (e) {}
          }
      };
      //And now, read the image and base64
      reader.readAsDataURL(files[0]);
    }

    onClosePopup(files) {
      this.setState({
        dragging: false,
        files: []
      });
      this.addOrImportRef.setComponentState({open: false});
    }

    onDragEnter() {
      this.setState({
        dragging: true
      });
    }

    onDragLeave() {
      this.setState({
        dragging: false
      });
    }

    download(e) {
      e.stopPropagation();
      window.api.download({action: "GetImportCSVTemplate", state: {}, controller: this.state.addPageController, fileName: DownloadTemplate});
    }

    render(){
      try {
        if (Action == "CONTROLLER_ACCOUNTADD" && !this.globs.HasRole("ACCOUNT_ADD")) {
          return <div><Component/></div>;
        }
        if (Action == "CONTROLLER_USERADD" && !this.globs.HasRole("USER_ADD") && !this.globs.HasRole("ACCOUNT_INVITE")) {
          return <div><Component/></div>;
        }
        if (Action == "CONTROLLER_ROLEADD" && !this.globs.HasRole("ROLE_ADD")) {
          return <div><Component/></div>;
        }
        if (Action == "CONTROLLER_FILEOBJECTADD" && !this.globs.HasRole("FILEOBJECT_ADD")) {
          return <div><Component/></div>;
        }
        if (Action == "CONTROLLER_EQUIPMENTADD" && !this.globs.HasRole("EQUIPMENT_ADD")) {
          return <div><Component/></div>;
        }
        if (Action == "CONTROLLER_LICENSEADD" && !this.globs.HasRole("LICENSE_ADD")) {
          return <div><Component/></div>;
        }
        if (Action == "CONTROLLER_SITEADD" && !this.globs.HasRole("SITE_ADD")) {
          return <div><Component/></div>;
        }
        if (Action == "CONTROLLER_CUSTOMSRCADD" && !this.globs.HasRole("CUSTOMSRC_ADD")) {
          return <div><Component/></div>;
        }
        //AdditionalPages


        var uploadHTML = null;
        if (window.File && window.FileReader && window.FileList && window.Blob) {
          if (this.state.dragging) {
            uploadHTML =
                <div>
                  <div className="counter-icon margin-bottom-5"><i className="icon md-plus-circle" aria-hidden="true"></i>
                  </div>
                  <span className="counter-number">
                    {window.appContent.ImportAndUploadDrag}
                    <div style={{marginTop: 40}}></div>
                  </span>
                </div>
          } else {
            uploadHTML = this.state.files && this.state.files.length > 0 ? <div>
              <h2 style={{color: 'white'}}>{window.appContent.ImportAndUploadStatus}</h2>
              <div>{this.globs.map(this.state.files, (file) => <div key={file.name}>{file.name}</div>)}</div>
            </div> :
                <div>
                  <div className="counter-icon margin-bottom-5"><i className="icon md-upload" aria-hidden="true"></i>
                  </div>
                  <span className="counter-number">{window.appContent.ImportAndUpload}
                    {DownloadTemplate !== false ? <div style={{marginTop: 20}}>
                      <a style={{color: 'white'}} onClick={this.download}>({window.appContent.ImportAndUploadDownload})</a>
                    </div>: null}
                  </span>
                </div>
          }
        } else {
          uploadHTML =  <div>
                  <div className="counter-icon margin-bottom-5"><i className="icon md-warning" aria-hidden="true"></i>
                  </div>
                  <span className="counter-number">
                    {window.appContent.ImportAndUploadNotSupported}
                    <div style={{marginTop: 40}}></div>
                  </span>
                </div>
        }

        var colSize = 0;

        if (MiddleColumnFunc && MiddleColumnFunc() != null) {
          colSize += 1;
        }

        if (ShowUpload) {
          colSize += 1;
        }

        if (ShowViaForm) {
          colSize += 1;
        }
        if (colSize == 1) {
          var colSmSize = 8;
        } else if (colSize == 2) {
          var colSmSize = 4;
        } else {
          var colSmSize = 3;
        }

        if (MiddleColumnFunc && ShowUpload && ShowViaForm) {
          var width = 900;
        } else {
          var width = 500;
        }


        //the icons for the below counter counter-lg classes for icons in the <i> tag are found here:
        //https://materialdesignicons.com/
        let msg = "";
        if (this.state.method == "ImportRoomJson") {
          msg = window.appContent.ImportFile;
        } else {
          msg = window.appContent.ImportAndUploadAddViaForm;
        }

        return (
          <div>
            <Component parent={this} ref={(c) => this.page = c}/>
            {!this.hasOnlyOneOption ?
                <div>
                  <ConfirmPopup
                      autoClose={false}
                      showActionButtons={true}
                      showActionSubmit={false}
                      actionCancelLabel={window.appContent.ImportClose}
                      onSubmit={() => {}}
                      title={window.appContent.ImportCompleted}
                      popupHTML={<div style={{marginLeft: 15}}>
                        {(this.state.FileUpload) ? <div style={{color: green900, fontSize: 20}}>{this.state.FileUpload.Meta.RowsCommittedInfo}</div>: null}
                        {(this.state.FileUpload && this.state.FileUpload.Meta.RowsSkipped > 0) ? <div style={{color: orange900, fontSize: 20, marginTop: 12}}>{this.state.FileUpload.Meta.RowsSkippedInfo}</div>: null}
                        {(this.state.FileUpload && this.state.FileUpload.Meta.RowsSkipped > 0) ? <div style={{marginTop: 20, overflowY: 'scroll', height: $(window).height() > 500 ? $(window).height() - 400 : 250}}>
                          <ul>
                            {this.globs.map(this.state.FileUpload.Meta.FileErrors, (e) =>
                                <li key={e} style={{fontSize: 16}}><strong>{e}</strong></li>)}
                          </ul>
                          <div style={{height: 75}}></div>
                        </div>: null}
                      </div>}
                    ref={(component) => this.completionPopup = component}
                  />
                  <ConfirmPopup
                      autoClose={true}
                      showActionButtons={false}
                      onClose={() => this.onClosePopup()}
                      onSubmit={() => {}}
                      title={this.state.title}
                      width={width}
                      popupHTML={<div className="widget">
                      <div className="widget-content">
                        <div className="row center">
                          {ShowViaForm ?
                            <div className={"col-sm-" + colSmSize} style={{cursor: 'pointer', marginRight: 26, marginLeft: 52}} onClick={() => window.global.functions.redirect(this.state.addPageAction)}>
                              <div className="counter counter-lg counter-inverse bg-purple-600 vertical-align height-150">
                                <div className="vertical-align-middle">
                                  <div className="counter-icon margin-bottom-5"><i className="icon md-edit" aria-hidden="true"/></div>
                                  <span className="counter-number">
                                    {window.appContent.ImportAndUploadAddViaForm}
                                    <div style={{marginTop: 40}}></div>
                                  </span>
                                </div>
                              </div>
                            </div>: null}
                          {(MiddleColumnFunc) ? MiddleColumnFunc(this.handleInvite) : null}
                          {ShowUpload ?
                            <div className={"col-sm-" + colSmSize} style={ShowViaForm ? {cursor: 'pointer'}: {cursor: 'pointer', marginRight: 26, marginLeft: 52}}>
                              <DropZone onDragEnter={() => this.onDragEnter()} onDragLeave={() => this.onDragLeave()} onDrop={(files) => this.onDrop(files)} style={{cursor:'pointer'}} inputProps={{name:"file_upload", id:"file_upload"}}>
                                <div className="counter counter-lg counter-inverse bg-blue-600 vertical-align height-150">
                                  <div className="vertical-align-middle">
                                   {uploadHTML}
                                  </div>
                                </div>
                              </DropZone>
                            </div>: null}
                        </div>
                      </div>
                    </div>}
                    ref={(component) => this.addOrImportRef = component}
                  />
                </div>: null}
            {DisplayIcon ? <ActionButton action={this.state.addPageAction} visible={this.state.buttonVisible}  onClick={() => this.open()} tooltip={window.pageContent[Tooltip]} type="AddImport"/>: null}
          </div>

        );
      } catch(e) {
        return this.globs.ComponentError(this.getClassName(), e.message);
      }
    }
  }

  return AddOrImportPageComponent;
}
