/**
 * Created by Dan on 11/3/16.
 */

class CoreRouter {
  constructor() {
    this.loadCallback;
    this.onGetErrorCallback;
    this.onStartRequestCallback;
    this.onEndRequestCallback;
    this.onPostErrorCallback;
    this.onGetSuccessCallback;
    this.onPostSuccessCallback;
    this.getCallback;

    var routes = {
      '/': () => {
        this.handlePartialPageLoad()
      },
      '/*': () => this.handlePartialPageLoad(),
    };

    var router = window.Router(routes);


    $(document).ready(() => {
      router.init("/");
      core.Debug.Dump(window.location.href == document.location.origin + "/")

      if (window.location.href == document.location.origin + "/") {
        this.handlePartialPageLoad();
      }

      if (this.loadCallback != undefined) {
        this.loadCallback.call(this);
      }

    });

    this.onLoad = (callback) => {
      this.loadCallback = callback;
    };

    this.onGetError = (callback) =>  {
      this.onGetErrorCallback = callback;
    };

    this.onEndRequest = (callback) =>  {
      this.onEndRequestCallback = callback;
    };

    this.onStartRequest = (callback) =>  {
      this.onStartRequestCallback = callback;
    };

    this.onPostError = (callback) =>  {
      this.onPostErrorCallback = callback;
    };

    this.onGetSuccess = (callback) =>  {
      this.onGetSuccessCallback = callback;
    };

    this.onPostSuccess = (callback) =>  {
      this.onPostSuccessCallback = callback;
    }
  }

  getVar(name, url) {
    if (!url) {
      url = window.location.href;
    }
    name = name.replace(/[\[\]]/g, "\\$&");
    var regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)"),
        results = regex.exec(url);
    if (!results) {
      return null;
    }
    if (!results[2]) {
      return '';
    }
    return decodeURIComponent(results[2].replace(/\+/g, " "));
  }

  handlePartialPageLoad() {

    var paramIndex = window.location.href.indexOf("?");
    var uri = window.location.href;
    var additionalParams = "";

    if (paramIndex != -1) {
      uri = window.location.href.substring(0,paramIndex);
      additionalParams = window.location.href.substring(paramIndex).replace("?", "&");
    }

    var partials = uri.split("#/");
    core.Debug.Dump(partials);
    let info = {};
    var newPage = "home";
    var partialPath = "home";

    if (partials.length > 1 && partials[1] != "") {
      var page = partials[1].split("/");
      newPage = page[page.length - 1];
      partialPath = partials[1];
      if (newPage == "") {
        newPage = partials[partials.length - 1];
        partialPath = partials[partials.length - 1];
      }
      info.page = newPage;
      info.partialPath = partialPath;
    }

    this.clientSide = false;

    if (this.onStartRequestCallback != undefined) {
      this.onStartRequestCallback("handlePartialPageLoad", info);
    }
    $.ajax({
      url: "/dist/markup?path=" + partialPath + "&file=" + newPage + additionalParams,
      type: 'GET',
      dataType: "json",
      success: (data, status, xhr) => {
        var contentObj;
        try{
          contentObj = JSON.parse(data.PageContent);
        }catch(e){
          if (this.onEndRequestCallback != undefined) {
            this.onEndRequestCallback();
          }
          if (this.onGetErrorCallback != undefined) {
            this.onGetErrorCallback("Failed to parse pageContent.", e);
          }
          return;
        }

        var jsonObj;
        try{
          jsonObj = JSON.parse(data.Json);
        }catch(e){
          if (this.onEndRequestCallback != undefined) {
            this.onEndRequestCallback();
          }
          if (this.onGetErrorCallback != undefined) {
            this.onGetErrorCallback("Failed to parse pageState.", e);
          }
          return;
        }

        if (this.onEndRequestCallback != undefined) {
          this.onEndRequestCallback();
        }

        if (this.onGetSuccessCallback != undefined) {
          this.clientSide = true;
          let leaveStateAlone = false;
          if (this.getVar("leaveStateAlone") == "1") {
            leaveStateAlone = true;
          }
          this.onGetSuccessCallback( data, status, xhr, data.Html, contentObj, jsonObj, data.Redirect, data.GlobalMessage, data.Trace, data.GlobalMessageType, this.getCallback, newPage, leaveStateAlone);
          this.getCallback = undefined;

          //Tell the websocket what page I am on.
          var gatewayId = window.globals.getCookie("GatewayId");
          window.api.writeSocket({controller:"App", action:"SetCurrentPage", state:{Page:newPage, GatewayId: gatewayId}});

          return;
        }
      },
      error: (e) => {
        if (this.onEndRequestCallback != undefined) {
          this.onEndRequestCallback();
        }
        if (this.onGetErrorCallback != undefined) {
          if (this.clientSide == true) {
            this.onGetErrorCallback(this.clientSide, CoreRouter.ROUTER_A_CLIENT_SIDE_ERROR_OCCURED, e);
          } else {
            this.onGetErrorCallback(this.clientSide, CoreRouter.ROUTER_A_SERVER_ERROR_OCCURED, e);
          }
        }
      }
    })

  }

  post(param) {

    var controller = param.controller;
    var disableSpinner = param.disableSpinner;
    if (controller == undefined || controller == "") {
      var partials = window.location.href.split("#/");
      controller = "home";
      if (partials.length > 1) {
        var paramIndex = partials[1].indexOf("?");

        if (paramIndex != -1) {
          controller = partials[1].substring(0, paramIndex);
        } else {
          controller = partials[1];
        }
      }
    }

    this.clientSide = false;

    var url = "/api?path=" + controller;
    if (param.urlOverride !== undefined) {
      url = param.urlOverride;
    }

    var apiPayload = {};
    apiPayload.action = param.action;
    apiPayload.state = JSON.stringify(param.state);
    var apiPayloadString = JSON.stringify(apiPayload);
    if (this.onStartRequestCallback != undefined && (disableSpinner == undefined || disableSpinner == false)) {
      this.onStartRequestCallback("post", {state: param.state, action: param.action, controller: controller});
    }
    $.ajax({
        beforeSend: (xhrObj) => {
            xhrObj.setRequestHeader("Content-Type","application/json");
            xhrObj.setRequestHeader("Accept","application/json");
        },
        url: url,
        type: 'POST',
        async: !param.hasOwnProperty("async") ? true: param.async,
        data: apiPayloadString,
        success: (data, status, xhr) => {

          var jsonObj;
          try{
            jsonObj = JSON.parse(data.State);
            if (window.appState.DeveloperMode) {
              // core.Debug.Dump("window.api.post(JSON.parse(\"" + JSON.stringify(param) + "\"))");
            }
          } catch(e) {
            if (this.onEndRequestCallback != undefined) {
              this.onEndRequestCallback();
            }
            if (this.onPostErrorCallback != undefined) {
              this.onPostErrorCallback("Failed to parse pageState.", e, controller, apiPayload.action, param.error);
            }
            return;
          }

          if (this.onPostSuccessCallback != undefined) {
            if (this.onEndRequestCallback != undefined) {
              this.onEndRequestCallback();
            }
            this.clientSide = true;
            this.onPostSuccessCallback(data, status, xhr, jsonObj, data.Redirect, data.GlobalMessage, data.Trace, data.GlobalMessageType, data.TransactionId, param.leaveStateAlone, param.callback, param.error);
            return;
          }

        },
        error: (e) => {
          if (this.onEndRequestCallback != undefined) {
            this.onEndRequestCallback();
          }
          if (this.onPostErrorCallback != undefined) {
            if (this.clientSide == true) {
              this.onPostErrorCallback(this.clientSide, CoreRouter.ROUTER_A_CLIENT_SIDE_ERROR_OCCURED, e, controller, apiPayload.action, param.error);
            } else {
              this.onPostErrorCallback(this.clientSide, CoreRouter.ROUTER_A_SERVER_ERROR_OCCURED, e, controller, apiPayload.action, param.error);
            }
          }
        }
    });
  }

  upload(param) {
    if (this.onStartRequestCallback != undefined) {
      this.onStartRequestCallback("upload", param);
    }
    var reader = new FileReader();
    reader.onload = (ev) => {
        try {
            var FileUpload = {};
            FileUpload.Id = param.fileId;
            FileUpload.Name = param.file.name;
            FileUpload.Content = ev.target.result.substr(ev.target.result.indexOf("base64,") + 7);
            FileUpload.Size = param.file.size;
            FileUpload.Type = param.file.type;
            FileUpload.Modified = param.file.lastModifiedDate;
            FileUpload.ModifiedUnix = param.file.lastModified;
            window.api.post({action: "Save",
                             state: {FileObject: FileUpload,
                                     Width:param.width,
                                     Height:param.height},
                             controller: "fileUpload",
                             leaveStateAlone: true,
                             disableSpinner: param.disableSpinner,
                             callback: (a, b, c, d, e, f, g, h, i, j) => {
              if (this.onEndRequestCallback != undefined) {
                this.onEndRequestCallback();
              }
              if (param.callback != undefined) {
                param.callback(a, b, c, d, e, f, g, h, i, j)
              }

            }, error: () => {
                if (this.onEndRequestCallback != undefined) {
                  this.onEndRequestCallback();
                }
                if (param.error != undefined) {
                  param.error()
                }
              }
            });
        }
        catch (ex) {
          if (this.onEndRequestCallback != undefined) {
            this.onEndRequestCallback();
          }
          if (typeof param.error == "function") {
            param.error(ex);
          }
        }
    };
    try {
      //And now, read the image and base64
      reader.readAsDataURL(param.file);
    }
    catch (ex) {
      if (this.onEndRequestCallback != undefined) {
        this.onEndRequestCallback();
      }
      if (typeof param.error == "function") {
        param.error(ex);
      }
    }

  }

  uploadMultipart(param) {
    if (this.onStartRequestCallback != undefined) {
      this.onStartRequestCallback("uploadMultipart", param);
    }
    var reader = new FileReader();
    reader.onload = (ev) => {
        try {
            var FileUpload = {};
            FileUpload.Id = param.fileId;
            FileUpload.Name = param.file.name;
            FileUpload.Content = ev.target.result.substr(ev.target.result.indexOf("base64,") + 7);
            FileUpload.Size = param.file.size;
            FileUpload.Type = param.file.type;
            FileUpload.Modified = param.file.lastModifiedDate;
            FileUpload.ModifiedUnix = param.file.lastModified;

            // var jForm = new FormData();
            // jForm.append("file", $(‘#file’).get(0).files[0]);

            $.ajax({
              url: param.url,
              type: "POST",
              data: FileUpload.Content,
              mimeType: "multipart/form-data",
              contentType: false,
              cache: false,
              processData: false,
              success: (data, textStatus, jqXHR) => {
                if (this.onEndRequestCallback != undefined) {
                  this.onEndRequestCallback();
                }
              },

              error: (jqXHR, textStatus, errorThrown) => {
                if (this.onEndRequestCallback != undefined) {
                  this.onEndRequestCallback();
                }
              }

              });


        }
        catch (ex) {
          if (this.onEndRequestCallback != undefined) {
            this.onEndRequestCallback();
          }
          if (typeof param.error == "function") {
            param.error(ex);
          }
        }
    };
    try {
      //And now, read the image and base64
      reader.readAsDataURL(param.file);
    }
    catch (ex) {
      if (this.onEndRequestCallback != undefined) {
        this.onEndRequestCallback();
      }
      if (typeof param.error == "function") {
        param.error(ex);
      }
    }

  }

  download(param) {
    var controller = param.controller;

    if (window.URL==null){
      alert(window.appContent.ErrorDownload);
      if (this.onEndRequestCallback != undefined) {
        this.onEndRequestCallback();
      }
    }else{
      if (controller == undefined || controller == "") {
        var partials = window.location.href.split("#/");
        controller = "home";
        if (partials.length > 1) {
          var paramIndex = partials[1].indexOf("?");

          if (paramIndex != -1) {
            controller = partials[1].substring(0, paramIndex);
          } else {
            controller = partials[1];
          }
        }
      }

      this.clientSide = false;
      var apiPayload = {};
      apiPayload.action = param.action;
      apiPayload.state = JSON.stringify(param.state);
      var apiPayloadString = JSON.stringify(apiPayload);

      if (param.fileObjectId != undefined) {

        var downloadLink = document.createElement("a");
        downloadLink.download = param.fileName;
        downloadLink.innerHTML = "Download File";
        if (window.appState.UserAgent.Name != "Firefox") {
            // Chrome allows the link to be clicked
            // without actually adding it to the DOM.
            downloadLink.href = "/fileObject/" + param.fileObjectId;
        }
        else {
            // Firefox requires the link to be added to the DOM
            // before it can be clicked.
            downloadLink.href = "/fileObject/" + param.fileObjectId;
            downloadLink.onclick = this.destroyClickedElement;
            downloadLink.style.display = "none";
            document.body.appendChild(downloadLink);
        }

        downloadLink.click();

      } else {
        if (this.onStartRequestCallback != undefined) {
          this.onStartRequestCallback("download", {controller: controller, state: param.state, action: param.action});
        }
        $.ajax({
            url: "/api?path=" + controller,
            type: 'POST',
            data: apiPayloadString,
            success: (data, status, xhr) => {

              var blob = this.b64toBlob(data);
              var downloadLink = document.createElement("a");
              downloadLink.download = param.fileName;
              downloadLink.innerHTML = "Download File";
              if (window.URL != null) {
                  // Chrome allows the link to be clicked
                  // without actually adding it to the DOM.
                  downloadLink.href = window.URL.createObjectURL(blob);
              }
              else {
                  // Firefox requires the link to be added to the DOM
                  // before it can be clicked.
                  downloadLink.href = window.URL.createObjectURL(blob);
                  downloadLink.onclick = this.destroyClickedElement;
                  downloadLink.style.display = "none";
                  document.body.appendChild(downloadLink);
              }

              downloadLink.click();
              if (this.onEndRequestCallback != undefined) {
                this.onEndRequestCallback();
              }
            }
        });
      }
    }
  }

  destroyClickedElement(event) {
      document.body.removeChild(event.target);
  }

  b64toBlob(b64Data, contentType, sliceSize) {
    contentType = contentType || '';
    sliceSize = sliceSize || 512;

    var byteCharacters = atob(b64Data);
    var byteArrays = [];

    for (var offset = 0; offset < byteCharacters.length; offset += sliceSize) {
      var slice = byteCharacters.slice(offset, offset + sliceSize);

      var byteNumbers = new Array(slice.length);
      for (var i = 0; i < slice.length; i++) {
        byteNumbers[i] = slice.charCodeAt(i);
      }

      var byteArray = new Uint8Array(byteNumbers);

      byteArrays.push(byteArray);
    }

    var blob = new Blob(byteArrays, {type: contentType});
    return blob;
  }


  buildGetUrl(param) {

    var controller = param.controller;

    if (controller == undefined || controller == "") {
      var paramIndex = window.location.href.indexOf("?");
      var uri = window.location.href;

      if (paramIndex != -1) {
        uri = window.location.href.substring(0,paramIndex);
      }

      var partials = uri.split("#/");
      var controller = "home";
      if (partials.length > 1) {
        var paramIndex = partials[1].indexOf("?");

        if (paramIndex != -1) {
          controller = partials[1].substring(0, paramIndex);
        } else {
          controller = partials[1];
        }
      }
    }

    if (controller.indexOf("CONTROLLER") != -1) {
      controller = window.appState.routes.Paths[controller];
    }

    var action = param.action;
    var uriParams = window.btoa(JSON.stringify(param.uriParams));
    if (action == undefined) {
      action = "Root";
    }

    var partialPath = "";
    if (controller.indexOf("/#/") == -1) {
      partialPath = "/#/";
    }

    //fullMount basically passes a unique token so that the entire react tree is rebuilt and you lose all of local this and replace with whatever the state of the server returns to the page/request
    let leaveState = "0";
    if (param.hasOwnProperty("leaveStateAlone") && param.leaveStateAlone) {
      leaveState = "1";
    }
    //(window.appState.DeveloperMode ? "&react_perf=1" : "") +
    var customGet = window.location.origin + partialPath + controller + "?action=" + action + "&uriParams=" + encodeURIComponent(uriParams) + "&leaveStateAlone=" + leaveState + ((param.hasOwnProperty("fullMount") && param.fullMount) ? "&token=" + window.globals.guid() : "");
    this.getCallback = param.callback;

    console.log("URI Params:  ", param.uriParams);

    return customGet;
  }

  get(param) {
    window.location.assign(this.buildGetUrl(param));
  }

  open(param) {
    window.open(this.buildGetUrl(param));
  }
}

CoreRouter.ROUTER_A_SERVER_ERROR_OCCURED = "RouterAServerErrorOccurred";
CoreRouter.ROUTER_A_CLIENT_SIDE_ERROR_OCCURED = "RouterAClientSideErrorOccurred";

export default CoreRouter;
