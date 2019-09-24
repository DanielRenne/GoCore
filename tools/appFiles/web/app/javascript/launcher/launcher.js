/**
 * Created by Dan on 11/3/16.
 */
import CoreRouter from "../router/routerES6";
import WSocket from "../webSocket/webSocket";

class Launcher {
  constructor() {
    this.webSocketCallbacks = [];
    this.content;
    this.currentPage = "";

    this.deadConnectionCount;
    this.reactPerf = (window.appState.DeveloperMode ? "&react_perf=1" : "");
    var router = new CoreRouter();
    var port = (window.location.protocol == "https:") ? 443 : window.appState.HTTPPort;
    if (port == undefined || port == "" || port == 0) {
      port = 80;
    }
    var ws = new WSocket(window.location.hostname, port, "ws");

    var websocketDataCallback = (jsonObj, pub) => {
      this.webSocketCallbacks.forEach((item) => {
        if (item.sub == "" && pub == "*") {
          item.callback(jsonObj, pub);
          return;
        }
        if (item.sub == pub || item.sub == "*") {
          item.callback(jsonObj, pub);
          return;
        }
      });
    };
    var webSocketOpen =() =>{
      this.deadConnectionCount = 0;
      this.handleWebSocketOpen();
    }

    var webSocketClosed = () => {
      this.handleWebSocketClosed();
    }

    var webSocketError = () => {
      this.handleWebSocketError();
    }

    ws.connect(webSocketOpen, websocketDataCallback, webSocketClosed, webSocketError);
    this.coreRouter = router;
    this.ws = ws;

    // window.onbeforeunload = function() {
    //   this.ws.onclose = function() {};
    //   this.ws.close();
    // }

    router.onLoad(() => this.loadPage());
    router.onStartRequest((typeRequest, information) => this.onStartRequest(typeRequest, information));
    router.onEndRequest(() => this.onEndRequest());
    router.onGetError((clientSide, message, e) => this.handleGetError(clientSide, message, e));
    router.onGetSuccess((data, status, xhr, Html, PageContent, State, Redirect, GlobalMessage, Trace, GlobalMessageType, getCallback, newPage, LeaveStateAlone) => this.handleGetSuccess(data, status, xhr, Html, PageContent, State, Redirect, GlobalMessage, Trace, GlobalMessageType, getCallback, newPage, LeaveStateAlone));
    router.onPostError((clientSide, message, e, controller, action, errorCallback) => this.handlePostError(clientSide, message, e, controller, action, errorCallback));
    router.onPostSuccess((data, status, xhr, State, Redirect, GlobalMessage, Trace, GlobalMessageType, TransactionId, LeaveStateAlone, postCallback, errorCallback) => this.handlePostSuccess(data, status, xhr, State, Redirect, GlobalMessage, Trace, GlobalMessageType, TransactionId, LeaveStateAlone, postCallback, errorCallback));
  }

  onStartRequest(typeRequest, information) {
    if ((window.location.href).indexOf("/#/home") == -1) {
      localStorage.setItem("uriRedirect", "~"+window.location.hash.substring(2));
    }
    let startSpinner = true;
    startSpinner &= window.location.href.indexOf("YouNeedToAddCertainPagesToStopSpinner") == -1 && window.location.href.indexOf("YouNeedToAddCertainPagesToStopSpinner2") == -1;
    if (typeRequest == "post") {
      let controllerActionSkips = ["roomList-TailLog", "serverSettingsModify-RestoreDatabase", "serverSettingsModify-FactoryReset", "equipmentCatalog-LearnIR"];
      startSpinner &= $.inArray(information.controller + "-" + information.action, controllerActionSkips) == -1;
    }
    if (information.action != undefined && $.inArray(information.action, ["YouNeedToAddCertainbActionsToStopSpinner"])) {
      startSpinner = false;
    }
    // if (window.hasOwnProperty("goCore") && window.goCore.setLoaderFromExternal != undefined && startSpinner) {
    //   window.loading = setTimeout(() => {
    //     window.goCore.setLoaderFromExternal({loading: true});
    //     core.Debug.Dump("100ms elapsed on spinner to show, you can further filter this from showing using this data:", typeRequest, information)
    //   },100);
    // }

    if (window.hasOwnProperty("goCore") && window.goCore.setLoaderFromExternal != undefined && startSpinner) {
      window.goCore.setLoaderFromExternal({loading: true});
    }

  }

  onEndRequest() {
    if (window.goCore.hasOwnProperty("setCurrentHistoryIdx")) {
      window.goCore.setCurrentHistoryIdx(window.history.length);
    }

    if (window.hasOwnProperty("goCore") && window.goCore.setLoaderFromExternal != undefined) {
      // clearTimeout(window.loading);
      window.goCore.setLoaderFromExternal({loading: false});
    }
  }

  loadPage() {
      this.content = $("#pagecontent");
      if (window.api == undefined) {
        window.api = {};
      }

      window.AppLoad();
      window.api.post = (param) => { return this.coreRouter.post(param)};
      window.api.get = (param) => { return this.coreRouter.get(param)};
      window.api.open = (param) => { return this.coreRouter.open(param)};
      window.api.download = (param) => {return this.coreRouter.download(param)};
      window.api.upload = (param) => {return this.coreRouter.upload(param)};
      window.api.uploadMultipart = (param) => {return this.coreRouter.uploadMultipart(param)};
      window.api.buildGetUrl = (param) => { return this.coreRouter.buildGetUrl(param)};
      window.api.newWindow = (param) => {
        if (!param.hasOwnProperty("uriParams")) {
          param.uriParams = {};
        }
        let url = window.api.buildGetUrl(param);
        window.open(url);
      };
      window.api.resize = () => {this.resizeContent()};
      window.api.writeSocket = (param) => { this.handleWriteSocket(param)};
      window.api.registerSocketCallback = (callback, subscription) => { return this.registerSocketCallback(callback, subscription)};
      window.api.unRegisterSocketCallback = (id) => { this.unRegisterSocketCallback(id)};
      window.api.closeSocket = () => {this.ws.close()};
      window.addEventListener("resize", () => this.resizeContent());


      window.LoadSiteBanner();
      window.LoadSideBarMenu();
      this.resizeContent();
      window.LoadSiteFooter();
      window.LoadSiteNotifications();
      window.store.init();
  }



  registerSocketCallback(callback, subscription) {
    var id = window.globals.guid();
    this.webSocketCallbacks.push({id:id, callback: callback, sub:subscription});
    return id;
  }

  unRegisterSocketCallback(id) {
    for (var i = 0; i < this.webSocketCallbacks.length; i++) {
      var s = this.webSocketCallbacks[i];
      if (s.id == id) {
        this.webSocketCallbacks.splice(i, 1);
        i--;
      }
    }
  }

  handleWriteSocket(param) {
    let objToSend = {};
        objToSend.action = param.action;
        objToSend.state = JSON.stringify(param.state);
        objToSend.controller = param.controller;

      let response = (data) => {

        let jsonObj = {};
        try{
          jsonObj = JSON.parse(data.State);
        }catch(e){
          this.handlePostError("Failed to parse pageState.", e, param.controller, param.action, param.error);
          return;
        }

        let Redirect = data.Redirect;
        let renderFooter = () => {
          this.renderFooterControls(data.GlobalMessage, data.Trace, data.GlobalMessageType, data.TransactionId, param.error);
        };

        let callback = () => {
          if (param.callback != undefined) {
            param.callback(jsonObj, data.GlobalMessage, data.Trace, data.GlobalMessageType, data.TransactionId, param.error);
          }
        };
        //Handle any redirects or refreshes
        if (Redirect == "refresh") {
          window.location.reload();
        } else if (Redirect == "back") {
          window.history.back();
        } else if (Redirect == "homeRefresh") {
          renderFooter();
          var customGet = window.location.origin + "/#/home" + this.reactPerf;
          window.location.assign(customGet);
          window.location.reload();
          callback();
          return;
        } else if (Redirect == "rerender") {
          renderFooter();
          this.currentPage = "";
          this.coreRouter.handlePartialPageLoad();
          callback();
          return;
        } else if (Redirect != "") {
          if (this.contains(Redirect,"~")) {
            $('body').css({display: "none"});
            renderFooter();
            var customGet = window.location.origin + "/#/" + this.replaceAll(Redirect, "~", "") + this.reactPerf;
            window.location.assign(customGet);
            window.location.reload();
            callback();
            return;
          } else if(this.contains(Redirect, "http") && this.contains(Redirect, "/web/custom")) {
            window.location.assign(Redirect + this.reactPerf);
            callback();
            return;
          }
          renderFooter();

          if (this.contains(Redirect, "http")) {
            window.open(Redirect + this.reactPerf);
          } else {
            var customGet = window.location.origin + "/#/" + Redirect + this.reactPerf;
            window.location.assign(customGet);
          }
          callback();
          return;
        }

        renderFooter();
        callback();
      };
      this.ws.send(objToSend, response, param.error, param.onTimeout, param.timeout, param.multiCallback);
  }

  handleGetError(clientSide, message, e){
    if (e != undefined) {
      if (e.hasOwnProperty("responseText")) {
        var e = e.responseText;
      } else if (typeof e === 'object') {
        var e = JSON.stringify(e, null, 2);
      }
    } else {
      var e = "uknown error";
    }
    if (e == "NOT AUTHORIZED") {
      try {
        window.global.functions.Popup(window.appContent.Unauthorized);
      } catch (err) {}
      document.location = '/#/home';
      return
    }

    if (clientSide) {
      var stack = new Error().stack;
      this.renderClientFooterControls(message, e + "\n\n" + stack, "Error", "", true);
    } else {
      this.renderFooterControls(message, e, "Error", "");
      if (window.history.length > 0) {
        // go back every one second there is an error.  If the next page errors and golang server is down keep going back until we get to a page.
        if (window.BackErrors == null || window.BackErrors == undefined) {
          window.BackErrors = 1;
        }
        if (window.BackErrors == 1) {
          window.history.back();
          window.BackErrors--;
        } else if (window.BackErrors == 0) {
          window.BackErrors = null;
        }
      }
    }
  }

  handlePostError(clientSide, message, e, controller, Action, postErrorCallback) {
    var err = "POST Error:\n\t500\n\nPath:\n\t"+"/api?path=" + controller+ "\n\nAction:\n\t"+ Action;
    if (e.hasOwnProperty("responseText")) {
      var e = e.responseText;
    } else if (typeof e === 'object') {
      var e = JSON.stringify(e, null, 2);
    }

    if (e == "NOT AUTHORIZED") {
      try {
        window.global.functions.Popup(window.appContent.Unauthorized);
      } catch (e) {}
      document.location = '/#/home';
      return
    }
    if (clientSide) {
      var stack = new Error().stack;
      this.renderClientFooterControls(message, e + "\n\n" + err + "\n\n" + stack, "Error", "", true);
    } else {
      this.renderFooterControls(message, err, "Error", "", postErrorCallback);
    }
  }

  handleGetSuccess(data, status, xhr, Html, PageContent, State, Redirect, GlobalMessage, Trace, GlobalMessageType, getCallback, newPage) {

    if (xhr.responseText == Launcher.HTTP_NOT_AUTHORIZED) {
      window.location.assign(document.location.origin);
      return;
    }

    $("." + "GoCore-components").css("visibility", "visible");

    window.pageContent = PageContent;
    window.pageState = State;

    this.resizeContent();

    if (State.SideBarMenu != undefined) {
      window.appState["SideBarMenu"] = State.SideBarMenu;
      window.unloadSideBarMenu();
      window.LoadSideBarMenu();
    }

    window.collapseSideBarMenu();
    window.collapseNavbarAvatar();

    var renderFooter = () => {
      this.renderFooterControls(GlobalMessage, Trace, GlobalMessageType, "");
    };

    let callback = () => {
      if (getCallback != "" && getCallback != undefined) {
        getCallback(State, GlobalMessage, Trace, GlobalMessageType, "");
      }
    };

    renderFooter();

    //Handle any redirects or refreshes
    if (Redirect == "refresh") {
      window.location.reload();
    } else if (Redirect == "back") {
      window.history.back();
    } else if (Redirect == "homeRefresh") {
      var customGet = window.location.origin + "/#/home" + (window.appState.DeveloperMode ? "?react_perf=1" : "");
      window.location.assign(customGet);
      window.location.reload();
      this.currentPage = "";
    } else if (Redirect == "rerender") {
      this.coreRouter.handlePartialPageLoad();
    } else if (Redirect != "") {
      if (this.contains(Redirect,"~")) {
        $('body').css({display:"none"});
        var baseUrl = this.replaceAll(Redirect, "~", "") + this.reactPerf;

        var partialPath = "";
        if (baseUrl.indexOf("/#/") == -1) {
          partialPath = "/#/";
        }

        var customGet = window.location.origin + partialPath + baseUrl;
        window.location.assign(customGet);
        window.location.reload();
      } else if(this.contains(Redirect, "http") && this.contains(Redirect, "/web/custom")) {
        window.location.assign(Redirect + this.reactPerf);
        callback();
        return;
      }
      renderFooter();

      var baseUrl = Redirect + this.reactPerf;

      var partialPath = "";
      if (baseUrl.indexOf("/#/") == -1) {
        partialPath = "/#/";
      }

      if (this.contains(Redirect, "http")) {
        window.open(Redirect + this.reactPerf);
      } else {
        var customGet = window.location.origin + partialPath + baseUrl;
        window.location.assign(customGet);
        this.currentPage = "";
      }
    }

      if (this.currentPage == "" || this.currentPage != newPage) {
        window.unloadAll();
        this.content.html(Html);
        this.currentPage = newPage;
        window.global.functions.pageStart();
        window["Load_" + newPage].call(this);
        window.global.functions.pageEnd();
        callback();
      } else if (!leaveStateAlone) {
        this.setPageState(State, callback);
      } else if (leaveStateAlone) {
        callback();
      } 
  }

  handlePostSuccess(data, status, xhr, State, Redirect, GlobalMessage, Trace, GlobalMessageType, TransactionId, LeaveStateAlone, postCallback, postErrorCallback) {

    if (xhr.responseText == Launcher.HTTP_NOT_AUTHORIZED) {
      window.location.assign(document.location.origin);
      return;
    }

    var renderFooter = () => {
      this.renderFooterControls(GlobalMessage, Trace, GlobalMessageType, TransactionId, postErrorCallback);
    };

    let callback = () => {
      if (postCallback != undefined) {
        postCallback(State, GlobalMessage, Trace, GlobalMessageType, TransactionId, postErrorCallback);
      }
    };

    //Handle any redirects or refreshes
    if (Redirect == "refresh") {
      window.location.reload();
    } else if (Redirect == "back") {
      window.history.back();
    } else if (Redirect == "homeRefresh") {
      renderFooter();
      var customGet = window.location.origin + "/#/home" + (window.appState.DeveloperMode ? "?react_perf=1" : "");
      window.location.assign(customGet);
      window.location.reload();
      callback();
      return;
    } else if (Redirect == "rerender") {
      renderFooter();
      this.currentPage = "";
      this.coreRouter.handlePartialPageLoad();
      callback();
      return;
    } else if (Redirect != "") {
      if (this.contains(Redirect, "~")) {
        $('body').css({display: "none"});
        renderFooter();
        var customGet = window.location.origin + "/#/" + this.replaceAll(Redirect, "~", "") + this.reactPerf;
        window.location.assign(customGet);
        window.location.reload();
        callback();
        return;
      } else if(this.contains(Redirect, "http") && this.contains(Redirect, "/web/custom")) {
        window.location.assign(Redirect + this.reactPerf);
        callback();
        return;
      }
      renderFooter();
      if (this.contains(Redirect, "http")) {
        window.open(Redirect + this.reactPerf);
      } else {
        var customGet = window.location.origin + "/#/" + Redirect + this.reactPerf;
        window.location.assign(customGet);
      }
      callback();
      return;
    }

    renderFooter();

    if (!LeaveStateAlone) {
      this.setPageState(State, callback);
    } else {
      callback();
    }
    if (window.appState.DeveloperMode && State.hasOwnProperty("DeveloperLog") && window.pageState &&  window.pageState.DeveloperLog != "") {
      console.log("<ServerLogs>");
      window.pageState.DeveloperLog = State.DeveloperLog;
      console.log(window.atob(window.pageState.DeveloperLog));
      console.log("</ServerLogs>");
    }
  }

  handleWebSocketOpen(){
    console.log("Opened Web Socket");

    var paramIndex = window.location.href.indexOf("?");
    var uri = window.location.href;

    if (paramIndex != -1) {
      uri = window.location.href.substring(0,paramIndex);
    }

    var partials = uri.split("#/");
    var newPage = "home";

    if (partials.length > 1) {
      var page = partials[1].split("/");
      newPage = page[page.length - 1];
      if (newPage == "") {
        newPage = partials[partials.length - 1];
      }
    }

    this.handleWriteSocket({controller:"App", action:"SetCurrentPage", state:{Page:newPage}});
    if (window.ReloadOnWebSocketReconnect) {
      window.location = "/";
    }
  }

  handleWebSocketPageSpecificHandler(){

  }

  handleWebSocketClosed(){
    console.log("Web Socket Closed");
    this.handleWebSocketPageSpecificHandler();
  }

  handleWebSocketError(){
    console.log("Web Socket Error");
    this.handleWebSocketPageSpecificHandler();
  }

  renderFooterControls(globalMessage, trace, globalMessageType, transactionId, postErrorCallback) {
    if (globalMessage != "" || trace != "") {
      window.appState[Launcher.VIEWMODEL_DIALOG_OPEN] = false;
      window.appState[Launcher.VIEWMODEL_SNACKBAR_TRANSACTION] = transactionId;
      var message = this.getTranslation(globalMessage);
      window.PopupServerMessage = "";

      if (message == "") {
        message = globalMessage;
        window.PopupServerMessage = globalMessage;
      }

      var snackBarMessageLength = message.length;

      if (message != "") {
        window.appState[Launcher.VIEWMODEL_SNACKBAR_MESSAGE] = message;
        window.appState[Launcher.VIEWMODEL_SNACKBAR_TYPE] = globalMessageType;
        window.appState[Launcher.VIEWMODEL_SNACKBAR_OPEN] = true;
      }

      if (trace != "") {

        var t = "";
        if (typeof(trace) == "string") {
          window.PopupServerMessage = message+"\n\n"+trace;
        } else {
          window.PopupServerMessage = trace;
        }
        window.appState[Launcher.VIEWMODEL_DIALOG_OPEN] = true;
      }

      if (globalMessageType == Launcher.PARAM_SNACKBAR_TYPE_ERROR || trace != "") {
        if (typeof(postErrorCallback) == "function") {
          postErrorCallback();
        }
      }
      if (message != "" || trace != "") {
        this.setFooterState(this.getAppState());
      }
    }
  }

  renderClientFooterControls(globalMessage, trace) {
   if (globalMessage != "" || trace != "") {

     var message = this.getTranslation(globalMessage);
     window.PopupClientMessage = "";

     if (message == "") {
       message = globalMessage;
       window.PopupClientMessage = globalMessage;
     }

     window.PopupClientMessage = message+"\n\n"+trace;
     window.appState[Launcher.VIEWMODEL_DIALOG_OPEN] = true;


     this.setFooterState(this.getAppState());
   }
  }

  showDialog() {

    if (window.PopupServerMessage != "") {
      window.appState[Launcher.VIEWMODEL_DIALOG_MESSAGE + "2"] = "Server Response: \n\n" + window.PopupServerMessage;
    } else {
      window.appState[Launcher.VIEWMODEL_DIALOG_MESSAGE + "2"] = "";
    }

    // Concatenate the second message
    window.appState["ShowDialogSubmitBug2"] = true;
    window.appState[Launcher.VIEWMODEL_DIALOG_TITLE + "2"] = this.getTranslation(Launcher.VIEWMODEL_DIALOG_ERROR_TITLE);
    if (window.PopupClientMessage != undefined) {
      window.appState[Launcher.VIEWMODEL_DIALOG_MESSAGE + "2"] += "\n\n" + window.PopupClientMessage;
    }
    window.appState[Launcher.VIEWMODEL_DIALOG_OPEN + "2"] = true;
    window.appState[Launcher.VIEWMODEL_POPUP_ERROR_SUBMIT + "2"] = true;


    var post = {AppError: {}};
    post.AppError.Message = window.appState[Launcher.VIEWMODEL_DIALOG_MESSAGE + "2"];
    post.AppError.StackShown = window.appState[Launcher.VIEWMODEL_DIALOG_MESSAGE + "2"];
    post.AppError.Url = window.location.href;
    window.api.post({action: "CreateAppError", state: post, controller:"appErrors"});

    this.setFooterState(this.getAppState());
  }

  showDialogMinimal() {
    if (window.PopupServerMessage != "") {
      window.appState[Launcher.VIEWMODEL_DIALOG_MESSAGE + "2"] = window.PopupServerMessage;
    } else {
      window.appState[Launcher.VIEWMODEL_DIALOG_MESSAGE + "2"] = "";
    }
    window.appState["ShowDialogSubmitBug2"] = false;
    window.appState[Launcher.VIEWMODEL_DIALOG_OPEN + "2"] = true;
    window.appState[Launcher.VIEWMODEL_POPUP_ERROR_SUBMIT + "2"] = true;
    this.setFooterState(this.getAppState());
  }

  contains(source, value) {
    if (source.indexOf(value) == -1) {
      return false;
    } else {
      return true;
    }
  }

  replaceAll(source, value, replacementValue) {
    return source.replace(new RegExp(value, "g"), replacementValue);
  }

  resizeContent(){
    var footerHeight = 0.0;
    var partials = window.location.href.split("#/");
    var page = "";
    var isNormalPage = true;
    try{
      if (partials[1] != undefined) {
        var pos = partials[1].indexOf("?");
        if (pos != -1) {
          page = partials[1].substr(0,pos)
        } else {
          page = partials[1];
        }
        isNormalPage = true;
        window.global.functions.forEach(window.global.functions.NonConformingPages(), (pageIter) => {
          isNormalPage &= page != pageIter;
        });
      }
    } catch(e) {}

    if (window.innerWidth > 767 && isNormalPage) {
      footerHeight = 35.0;
    }

    var headerElement = window.global.functions.GetHeaderElement()
    var headerHeight;

    if (isNormalPage) {
      if (window.innerWidth > 767) {
        headerHeight = headerElement.height() - 5;
      } else {
        headerHeight = headerElement.height() - 35;
      }
    } else {
      headerHeight = 0;
    }

    // window.global.functions.log(window.innerHeight);
    // window.global.functions.log(headerHeight );
    // window.global.functions.log(footerHeight);

    var scrollHeight = window.innerHeight - headerHeight - footerHeight;

    window.HeaderHeight = footerHeight;
    window.FooterHeight = headerHeight;
    window.ScrollHeight = scrollHeight;

    var innerScrollHeight = window.innerHeight - headerHeight - footerHeight;
    if (partials.length > 1 && isNormalPage) {
      innerScrollHeight -= (18 + 30) // margin top and bottom of .GoCore-content
    }
    var atlContent = $(".GoCore-content");
    atlContent.css("max-height", scrollHeight + "px").css("height", scrollHeight + "px").css("min-width", $(window).width()).css("width", $(window).width()-(($(window).width() > 767 && isNormalPage) ? 30: 0) );
    $("." + "GoCore-body").css("min-width", $(window).width());

    var alignDiv = $("." + "Align");
    if (alignDiv.length > 0) {
      alignDiv.css("height", innerScrollHeight + "px");
    }

    var content = $("#pagecontent");
    content.css("height",  innerScrollHeight + "px");
  }


  //Reactjs State Call Functions

  getAppState() {
    return window.appState;
  }

  getPageState() {
    return window.pageState;
  }

  refreshApp(data){
    this.setToolbarState(data);
    this.setMenuState(data);
    this.setFooterState(data);
  }

  refreshPage() {
    if (window.hasOwnProperty("goCore") && window.goCore.page != undefined) {
      window.goCore.page.forceUpdate();
    }
  }

  setPageState(data, cb) {
    if (Object.keys(data).length == 0) {
      if (typeof(cb) == "function") {
        cb();
      }
      return;
    }
    if (window.hasOwnProperty("goCore") && window.goCore.setStateFromExternal != undefined) {
      window.goCore.setStateFromExternal(data, cb);
    }
  }

  setLoaderState(data) {
    if (window.hasOwnProperty("goCore") && window.goCore.setLoaderFromExternal != undefined) {
      window.goCore.setLoaderFromExternal(data);
    }
  }

  setFooterState(data) {
    if (window.hasOwnProperty("goCore") && window.goCore.setFooterStateFromExternal != undefined) {
      window.goCore.setFooterStateFromExternal(data);
    }
  }

  setSideMenuBarState(data) {
    if (window.hasOwnProperty("goCore") && window.goCore.setSideBarMenuStateFromExternal != undefined) {
      window.goCore.setSideBarMenuStateFromExternal(data);
    }
  }

  setMenuState(data) {
    if (window.hasOwnProperty("goCore") && window.goCore.setMenuStateFromExternal != undefined) {
      window.goCore.setMenuStateFromExternal(data);
    }
  }

  setToolbarState(data) {
    if (window.hasOwnProperty("goCore") && window.goCore.setToolbarStateFromExternal != undefined) {
      window.goCore.setToolbarStateFromExternal(data);
    }
  }

  rerender() {
    this.currentPage = "";
    this.coreRouter.handlePartialPageLoad();
  }

  //Translation Functions
  getTranslation(key) {
    try {
      if (window.appContent != undefined && window.appContent[key] == undefined) {
        if (window.appContent != undefined && window.pageContent[key] == undefined) {
          return "";
        } else if (window.appContent == undefined) {
          return key;
        }
        return window.pageContent[key];
      }
    } catch (e) {
      return key;
    }
    return window.appContent[key];
  }
}

Launcher.HTTP_NOT_AUTHORIZED = "NOT AUTHORIZED";
Launcher.VIEWMODEL_SNACKBAR_TRANSACTION = "SnackBarUndoTransactionId";
Launcher.VIEWMODEL_SNACKBAR_MESSAGE     = "SnackbarMessage";
Launcher.VIEWMODEL_SNACKBAR_OPEN        = "SnackbarOpen";
Launcher.VIEWMODEL_SNACKBAR_TYPE        = "SnackbarType";

Launcher.PARAM_SNACKBAR_TYPE_SUCCESS = "";
Launcher.PARAM_SNACKBAR_TYPE_WARNING = "Warning";
Launcher.PARAM_SNACKBAR_TYPE_ERROR   = "Error";

Launcher.VIEWMODEL_DIALOG_MESSAGE = "DialogMessage";
Launcher.VIEWMODEL_DIALOG_OPEN    = "DialogOpen";
Launcher.VIEWMODEL_POPUP_ERROR_SUBMIT = "PopupErrorSubmit";
Launcher.VIEWMODEL_DIALOG_TITLE   = "DialogTitle";

Launcher.VIEWMODEL_DIALOG_SUCCESS_TITLE = "DialogSuccess";
Launcher.VIEWMODEL_DIALOG_WARNING_TITLE = "DialogWarning";
Launcher.VIEWMODEL_DIALOG_ERROR_TITLE   = "DialogError";
Launcher.SNACK_BAR_LENGTH_MAX = 50;

export default Launcher;
