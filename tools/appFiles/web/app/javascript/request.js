
var token = "";

 function request(domain, port, path){
 	this.path = path;
 	this.port = port;
 	this.domain = domain;
 	this.socket = null;
 	this.socketCallBacks = new Array();
 	this.currentCallBackId = 0;
 	this.hasConnected = false;
 }

 request.prototype.connectSocket = function(onOpen, onMessage, onClose, onError) {
 	try{
	    this.onOpen = onOpen;
	    this.onMessage = onMessage;
	    this.onClose = onClose;
	    var r = this;
	    var host = "wss://" + this.domain + ":" + this.port + "/" + this.path;
	    if ("WebSocket" in window)
	        this.socket = new WebSocket(host);
	    if ("MozWebSocket" in window)
	        this.socket = new MozWebSocket(host);

	    this.socket.onopen = function () {
	    	r.onOpen.call();
	    }
        this.socket.onmessage = function (msg) {
			var jsonObj = JSON.parse(msg.data);
			if(jsonObj.callBackId != undefined){
				for(var i = 0; i < r.socketCallBacks.length; i++){
					var scb = r.socketCallBacks[i];
					if(scb.callBackId == jsonObj.callBackId){
						if(jsonObj.error != undefined){
							if(scb.errorCallBack != undefined)
								scb.errorCallBack(jsonObj.error);
						}
						else{
							if(scb.messageCallBack != undefined)
								scb.messageCallBack(jsonObj.data);
						}
						if(!scb.multipleCallbacks)
							r.socketCallBacks.splice(i, 1);
						if(jsonObj.endCallBacks != undefined)
							r.socketCallBacks.splice(i, 1);
						return;
					}
				}
			}
			else
				r.onMessage(jsonObj);
        }
	    this.socket.onclose = function() {
	    	r.onClose.call();
	    }

	    if(!this.hasConnected)
			setInterval(function(){r.checkConnection(onOpen, onMessage, onClose, onError);}, 3000);
		this.hasConnected = true;

	}
	catch(ex){
		onError.call(ex);
	}

};

request.prototype.checkConnection = function(onOpen, onMessage, onClose, onError) {
		    if (this.socket == null || this.socket == undefined) {
		        this.connectSocket(onOpen, onMessage, onClose, onError);
		        return;
		    }
		    if (this.socket.readyState == 1)
		        return;
			this.connectSocket(onOpen, onMessage, onClose, onError);
};

request.prototype.sendSocketRequest = function(data, onMessage, onError, onTimeout, timeoutSeconds, multipleCallbacks) {

	try{
		var r = this;
		this.currentCallBackId++;
		if(this.currentCallBackId == Number.MAX_VALUE)
			this.currentCallBackId = 1;
		var callBackId = this.currentCallBackId;
		var mcb = false;
		if(multipleCallbacks != undefined)
			mcb = multipleCallbacks;
		var t = token;
		var obj = {callBackId : this.currentCallBackId, token: t, data : data};
		this.socketCallBacks.push({callBackId: callBackId, messageCallBack: onMessage, errorCallBack: onError, multipleCallbacks: mcb});
		this.socket.send(JSON.stringify(obj));
		if(onTimeout != undefined || onTimeout != null){
			setTimeout(function(){
				for(var i = 0; i < r.socketCallBacks.length; i++){
					var scb = r.socketCallBacks[i];
					if(scb.callBackId == callBackId){
						r.socketCallBacks.splice(i, 1);
						if(r.onTimeout != undefined)
							r.onTimeout.call();
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
};