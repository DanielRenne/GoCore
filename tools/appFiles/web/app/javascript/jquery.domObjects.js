var jDom = new jQueryDom();

function jQueryDom () {
	// body...
}

jQueryDom.prototype.center = function() {
	return $("<center></center>");
};

jQueryDom.prototype.empty = function() {
	return $();
};

jQueryDom.prototype.header = function(num, attributes) {
	if(num == undefined)
		num = 3;
	if(attributes == undefined)
		attributes = "";
	return $("<h" + num + " " + attributes + "></h" + num + ">");
};

jQueryDom.prototype.font = function(attributes) {
	if(attributes == undefined)
		attributes = "";
	return $("<font " + attributes + "></font>");
};

jQueryDom.prototype.div = function(attributes) {
	if(attributes == undefined)
		attributes = "";
	return $("<div " + attributes + "></div>");
};

jQueryDom.prototype.nav = function(attributes) {
	if(attributes == undefined)
		attributes = "";
	return $("<nav " + attributes + "></nav>");
};

jQueryDom.prototype.table = function(attributes) {
	if(attributes == undefined)
		attributes = "";
	return $("<table " + attributes + "></table>");
};

jQueryDom.prototype.row = function(attributes) {
	if(attributes == undefined)
		attributes = "";
	return $("<tr " + attributes + "></tr>");
};

jQueryDom.prototype.column = function(attributes) {
	if(attributes == undefined)
		attributes = "";
	return $("<td " + attributes + "></td>");
};

jQueryDom.prototype.tableHead = function(attributes) {
	if(attributes == undefined)
		attributes = "";
	return $("<thead " + attributes + "></thead>");
};

jQueryDom.prototype.tableHeader = function(attributes) {
	if(attributes == undefined)
		attributes = "";
	return $("<th " + attributes + "></th>");
};

jQueryDom.prototype.tableBody = function(attributes) {
	if(attributes == undefined)
		attributes = "";
	return $("<tbody " + attributes + "></tbody>");
};

jQueryDom.prototype.button = function(attributes) {
	//to set button text pass it in as value attribute
	if(attributes == undefined)
		attributes = "";
	return $("<input type='button' " + attributes + "></button>");
};

jQueryDom.prototype.dropdown = function(attributes, options, selectedValue) {
	if(attributes == undefined)
		attributes = "";
	if(options == undefined)
		return $("<select " + attributes + "></select>");
	else{
		var select = $("<select " + attributes + "></select>");
		options.forEach(function(option){
			var selected = "";
	        if (option.value == selectedValue)
	            selected = "selected";
	        select.append("<option value = '" + option.value + "' " + selected + ">" + option.display + "</option>");
		});
		return select;
	}
};

jQueryDom.prototype.textbox = function(attributes) {
	if(attributes == undefined)
		attributes = "";
	return $("<input type='text' " + attributes + "/>");
};

jQueryDom.prototype.textarea = function(attributes) {
	if(attributes == undefined)
		attributes = "";
	return $("<textarea " + attributes + "/>");
};

jQueryDom.prototype.password = function(attributes) {
	if(attributes == undefined)
		attributes = "";
	return $("<input type='password' " + attributes + "/>");
};

jQueryDom.prototype.checkbox = function(attributes) {
	if(attributes == undefined)
		attributes = "";
	return $("<input type='checkbox' " + attributes + "/>");
};

jQueryDom.prototype.icon = function(attributes) {
	if(attributes == undefined)
		attributes = "";
	return $("<i " + attributes + ">");
};

jQueryDom.prototype.span = function(attributes) {
	if(attributes == undefined)
		attributes = "";
	return $("<span " + attributes + "></span>");
};

jQueryDom.prototype.link = function(attributes) {
	if(attributes == undefined)
		attributes = "";
	return $("<a " + attributes + "></a>");
};

jQueryDom.prototype.image = function(attributes) {
	if(attributes == undefined)
		attributes = "";
	return $("<img " + attributes + "/>");
};

jQueryDom.prototype.ul = function(attributes) {
	if(attributes == undefined)
		attributes = "";
	return $("<ul " + attributes + "></ul>");
};

jQueryDom.prototype.li = function(attributes) {
	if(attributes == undefined)
		attributes = "";
	return $("<li " + attributes + "></li>");
};

jQueryDom.prototype.option = function(attributes) {
	if(attributes == undefined)
		attributes = "";
	return $("<option " + attributes + "></option>");
};

jQueryDom.prototype.iFrame = function(attributes) {
	if(attributes == undefined)
		attributes = "";
	return $("<iframe " + attributes + "></iframe>");
};

//HTML Forms
jQueryDom.prototype.form = function(attributes) {
	if(attributes == undefined)
		attributes = "";
	return $("<form " + attributes + "></form>");
}; 

jQueryDom.prototype.fieldset = function(attributes) {
	if(attributes == undefined)
		attributes = "";
	return $("<fieldset " + attributes + "></fieldset>");
}; 

jQueryDom.prototype.label = function(attributes) {
	if(attributes == undefined)
		attributes = "";
	return $("<label " + attributes + "></label>");
}; 

jQueryDom.prototype.submit = function(attributes) {
	if(attributes == undefined)
		attributes = "";
	return $("<input type='submit' " + attributes + "></button>");
};


jQueryDom.prototype.bindOffFocus = function(obj, callback, params) {
	params.jQueryObj = obj;
	obj.bind("focusout", function(e){
	    callback(params);
	    e.preventDefault();
	});
	return obj;
	
};


jQueryDom.prototype.bindOnChangeOffFocus = function(obj, callback, params) {
	params.jQueryObj = obj;
	obj.bind("focusout change", function(e){
	    callback(params);
	    e.preventDefault();
	});
	return obj;
	
};

jQueryDom.prototype.bindTouchClick = function(obj, callback, params) {
	params.jQueryObj = obj;
	if (window.navigator.pointerEnabled) {
        obj.get(0).addEventListener("pointerup", function(e){
            callback(params);
        e.preventDefault();
        }, false);
    }
    else{
	    obj.bind("touchend mouseup", function(e){
	        callback(params);
	        e.preventDefault();
	    });
	}
	return obj;
};

jQueryDom.prototype.bindTouch = function(obj, callback, params) {
	params.jQueryObj = obj;
	if (window.navigator.pointerEnabled) {
        obj.get(0).addEventListener("pointerup", function(e){
            callback(params);
        e.preventDefault();
        }, false);
    }
    else{
	    obj.bind("touchend", function(e){
	        callback(params);
	        e.preventDefault();
	    });
	}
	return obj;
};

jQueryDom.prototype.bindEnterKey = function(obj, callback, params) {
	params.jQueryObj = obj;
	obj.on('keydown', function (e) {
        e.stopPropagation();
        if (e.keyCode == 13) {
            callback(params);
        }
    });
    return obj;
};

jQueryDom.prototype.bindKeyDown = function(obj, callback, params) {
	params.jQueryObj = obj;
	obj.on('keydown', function (e) {
        e.stopPropagation();
			callback(e, params);
    });
    return obj;
};

jQueryDom.prototype.dialog = function(obj, title, width, height, closeCallback){

	var dialogDiv = this.div('syle="z-index="10000;"').append(obj); 

    $("body").append(dialogDiv);

    dialogDiv.dialog({
        width: width, height: height, title: title,close: function(event, ui) { $(this).remove(); if(closeCallback)closeCallback();}
    });
	
	return dialogDiv;
}

jQueryDom.prototype.openFile = function(callback, accept) {
	// Check for the various File API support.
	if (window.File && window.FileReader && window.FileList && window.Blob) {
	  // Great success! All the File APIs are supported.
	  var acceptInsert = "";
	  if(accept)
	  	acceptInsert =  "accept='" + accept + "'";
	  var fileInput = $("<input type='file' name='files[]' " + acceptInsert + "/>");

	  $("body").append(fileInput);

	  fileInput.trigger("click");

	  fileInput.on('change', function(evt){
	  	console.log(evt.target.files[0]);

        var reader = new FileReader();
        var fileName = evt.target.files[0].name;
        //When the file has been read...
        reader.onload = function (ev) {
            try {
                var jsonObj = new Object();
                jsonObj.fileName = fileName;
                jsonObj.data = window.atob(ev.target.result.substr(ev.target.result.indexOf("base64,") + 7));
                callback(jsonObj);
				fileInput.remove();
            }
            catch (ex) {
                alert(ex);
            }
        };
        //And now, read the image and base64
        reader.readAsDataURL(evt.target.files[0]);


	  	
	  });
	  
	} else {
	  alert('The File APIs are not fully supported in this browser.');
	}

};

jQueryDom.prototype.saveFile = function(data, fileName, type) {
    try {
        var textFileAsBlob = new Blob([data], { type: type  });
        var downloadLink = document.createElement("a");
        downloadLink.download = fileName;
        downloadLink.innerHTML = "Download File";
        if (window.webkitURL != null) {
            // Chrome allows the link to be clicked
            // without actually adding it to the DOM.
            downloadLink.href = window.webkitURL.createObjectURL(textFileAsBlob);
        }
        else {
            // Firefox requires the link to be added to the DOM
            // before it can be clicked.
            downloadLink.href = window.URL.createObjectURL(textFileAsBlob);
            downloadLink.onclick = destroyClickedElement;
            downloadLink.style.display = "none";
            document.body.appendChild(downloadLink);
        }

        downloadLink.click();
    }
    catch (ex) {
        alert("Error at jQueryDom.saveFile:  " + ex);
    }

};