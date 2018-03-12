var JSON = JSON || {};
// implement JSON.stringify serialization
JSON.stringify = JSON.stringify || function (obj) {
  var t = typeof (obj);
  if (t != "object" || obj === null) {

      // simple data type
      if (t == "string") obj = '"'+obj+'"';
      return String(obj);

  }
  else {

      // recurse array or object
      var n, v, json = [], arr = (obj && obj.constructor == Array);

      for (n in obj) {
          v = obj[n]; t = typeof(v);

          if (t == "string") v = '"'+v+'"';
          else if (t == "object" && v !== null) v = JSON.stringify(v);

          json.push((arr ? "" : '"' + n + '":') + String(v));
      }

      return (arr ? "[" : "{") + String(json) + (arr ? "]" : "}");
  }
};

// implement JSON.parse de-serialization
JSON.parse = JSON.parse || function (str) {
  if (str === "") str = '""';
  eval("var p=" + str + ";");
  return p;
};

!function(){function e(e){this.message=e}var t="undefined"!=typeof exports?exports:"undefined"!=typeof self?self:$.global,r="ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=";e.prototype=new Error,e.prototype.name="InvalidCharacterError",t.btoa||(t.btoa=function(t){for(var o,n,a=String(t),i=0,f=r,c="";a.charAt(0|i)||(f="=",i%1);c+=f.charAt(63&o>>8-i%1*8)){if(n=a.charCodeAt(i+=.75),n>255)throw new e("'btoa' failed: The string to be encoded contains characters outside of the Latin1 range.");o=o<<8|n}return c}),t.atob||(t.atob=function(t){var o=String(t).replace(/[=]+$/,"");if(o.length%4==1)throw new e("'atob' failed: The string to be decoded is not correctly encoded.");for(var n,a,i=0,f=0,c="";a=o.charAt(f++);~a&&(n=i%4?64*n+a:a,i++%4)?c+=String.fromCharCode(255&n>>(-2*i&6)):0)a=r.indexOf(a);return c})}();
