/**
 * Created by Dan on 11/28/16.
 */
class WSocket {
  constructor(domain, port, path) {
    this.path = path;
    this.port = port;
    this.domain = domain;
    this.socket = null;
    this.socketCallBacks = new Array();
    this.currentCallBackId = 0;
    this.hasConnected = false;
    this.closedExplicitly = false;
  }

  connect(onOpen, onMessage, onClose, onError) {
    try{
      this.onOpen = onOpen;
      this.onMessage = onMessage;
      this.onClose = onClose;
      var protocol = (window.location.protocol == "https:") ? "wss://" : "ws://";
      var host = protocol + this.domain + ":" + this.port + "/" + this.path;
      if ("WebSocket" in window) {
        this.socket = new WebSocket(host);
      }
      if ("MozWebSocket" in window) {
        this.socket = new MozWebSocket(host);
      }

      this.socket.onopen = () => {
        this.onOpen.call();
      };

      this.socket.onmessage = (msg) => {
        try {
          var jsonObj = JSON.parse(msg.data);
          if(jsonObj.callBackId != undefined){
            for(var i = 0; i < this.socketCallBacks.length; i++){
              var scb = this.socketCallBacks[i];
              if(scb.callBackId == jsonObj.callBackId){
                if(jsonObj.error != undefined){
                  if(scb.errorCallBack != undefined)
                    scb.errorCallBack(jsonObj.error);
                }
                else{
                  if(scb.messageCallBack != undefined)
                    scb.messageCallBack(jsonObj.data, jsonObj.pub);
                }
                if(!scb.multipleCallbacks)
                  this.socketCallBacks.splice(i, 1);
                if(jsonObj.endCallBacks != undefined)
                  this.socketCallBacks.splice(i, 1);
                return;
              }
            }
          }
          else {
            if (jsonObj.Key != undefined) {
              this.onMessage(jsonObj.Content, jsonObj.Key);
            } else {
              this.onMessage(jsonObj, "*");
            }
          }
        } catch (e) {
          console.warn("Websocket error", e)
        }
      };

      this.socket.onclose = () => {
        try {
          this.onClose.call();
        } catch (e) {
          console.warn("Websocket error", e)
        }
      }

      if(!this.hasConnected) {
        setInterval(() => {
          this.checkConnection(onOpen, onMessage, onClose, onError);
        }, 3000);
        setInterval(() => {
          this.sendPollMessage();
        }, 30000);
      }
      this.hasConnected = true;

    }
    catch(ex){
      console.error("Failed at connect:", ex);
      onError.call(ex);
    }
  }

  sendPollMessage() {
    if (this.socket.readyState == 1) { //Send a message for socket timeouts
      this.send({});
    }
  }

  checkConnection(onOpen, onMessage, onClose, onError) {
    if (this.closedExplicitly == true) {
      return;
    }

    if (this.socket == null || this.socket == undefined) {
        this.connect(onOpen, onMessage, onClose, onError);
        return;
    }
    if (this.socket.readyState == 1 || this.socket.readyState == 0) {  //Connecting or Open simply return.
      return;
    }
    this.connect(onOpen, onMessage, onClose, onError);
  }

  close() {
    this.closedExplicitly = true;
    this.socket.close();

  }

  send(data, onMessage, onError, onTimeout, timeoutSeconds, multipleCallbacks){
    try{

      if (this.socket.readyState == 2 || this.socket.readyState == 3) {
        return
      }

      this.currentCallBackId++;
      if(this.currentCallBackId == Number.MAX_VALUE) {
        this.currentCallBackId = 1;
      }
      var callBackId = this.currentCallBackId;
      var mcb = false;
      if(multipleCallbacks != undefined) {
        mcb = multipleCallbacks;
      }
      var obj = {callBackId : this.currentCallBackId, data : data};
      this.socketCallBacks.push({callBackId: callBackId, messageCallBack: onMessage, errorCallBack: onError, multipleCallbacks: mcb});
      this.socket.send(JSON.stringify(obj));
      if(onTimeout != undefined || onTimeout != null){
        setTimeout(() => {
          for(var i = 0; i < this.socketCallBacks.length; i++){
            var scb = this.socketCallBacks[i];
            if(scb.callBackId == callBackId){
              this.socketCallBacks.splice(i, 1);
              if(this.onTimeout != undefined)
                this.onTimeout.call();
              return;
            }
          }
        }, timeoutSeconds);
      }
    }
    catch(ex){
      if(onError != undefined)
        onError.call(ex);
    }
  }
}

export default WSocket;
