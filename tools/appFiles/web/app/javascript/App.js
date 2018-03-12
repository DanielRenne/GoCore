import "babel-polyfill";
import injectTapEventPlugin from "react-tap-event-plugin";
import ReactDom from "react-dom";
import {
  red50,
  red100,
  red200,
  red300,
  red400,
  red500,
  red600,
  red700,
  red800,
  red900,
  pink50,
  pink100,
  pink200,
  pink300,
  pink400,
  pink500,
  pink600,
  pink700,
  pink800,
  pink900,
  purple50,
  purple100,
  purple200,
  purple300,
  purple400,
  purple500,
  purple600,
  purple700,
  purple800,
  purple900,
  deepPurple50,
  deepPurple100,
  deepPurple200,
  deepPurple300,
  deepPurple400,
  deepPurple500,
  deepPurple600,
  deepPurple700,
  deepPurple800,
  deepPurple900,
  indigo50,
  indigo100,
  indigo200,
  indigo300,
  indigo400,
  indigo500,
  indigo600,
  indigo700,
  indigo800,
  indigo900,
  blue50,
  blue100,
  blue200,
  blue300,
  blue400,
  blue500,
  blue600,
  blue700,
  blue800,
  blue900,
  lightBlue50,
  lightBlue100,
  lightBlue200,
  lightBlue300,
  lightBlue400,
  lightBlue500,
  lightBlue600,
  lightBlue700,
  lightBlue800,
  lightBlue900,
  cyan50,
  cyan100,
  cyan200,
  cyan300,
  cyan400,
  cyan500,
  cyan600,
  cyan700,
  cyan800,
  cyan900,
  teal50,
  teal100,
  teal200,
  teal300,
  teal400,
  teal500,
  teal600,
  teal700,
  teal800,
  teal900,
  green50,
  green100,
  green200,
  green300,
  green400,
  green500,
  green600,
  green700,
  green800,
  green900,
  lightGreen50,
  lightGreen100,
  lightGreen200,
  lightGreen300,
  lightGreen400,
  lightGreen500,
  lightGreen600,
  lightGreen700,
  lightGreen800,
  lightGreen900,
  lime50,
  lime100,
  lime200,
  lime300,
  lime400,
  lime500,
  lime600,
  lime700,
  lime800,
  lime900,
  yellow50,
  yellow100,
  yellow200,
  yellow300,
  yellow400,
  yellow500,
  yellow600,
  yellow700,
  yellow800,
  yellow900,
  amber50,
  amber100,
  amber200,
  amber300,
  amber400,
  amber500,
  amber600,
  amber700,
  amber800,
  amber900,
  orange50,
  orange100,
  orange200,
  orange300,
  orange400,
  orange500,
  orange600,
  orange700,
  orange800,
  orange900,
  deepOrange50,
  deepOrange100,
  deepOrange200,
  deepOrange300,
  deepOrange400,
  deepOrange500,
  deepOrange600,
  deepOrange700,
  deepOrange800,
  deepOrange900,
  brown50,
  brown100,
  brown200,
  brown300,
  brown400,
  brown500,
  brown600,
  brown700,
  brown800,
  brown900,
  grey50,
  grey100,
  grey200,
  grey300,
  grey400,
  grey500,
  grey600,
  grey700,
  grey800,
  grey900,
  blueGrey50,
  blueGrey100,
  blueGrey200,
  blueGrey300,
  blueGrey400,
  blueGrey500,
  blueGrey600,
  blueGrey700,
  blueGrey800,
  blueGrey900
} from "./globals/forms";
//Pages currently stored in UI !attic folders
//import Users from './pages/users/users'
//import Notifications from './pages/notifications/notifications'
import Launcher from "./launcher/launcher";
import Store from "./components/store/store";
import BaseLogger from "./components/baseLogger";


window.materialColors = {
  transparent: "transparent", red50: red50, red100: red100, red200: red200, red300: red300, red400: red400, red500: red500, red600: red600, red700: red700, red800: red800, red900:  red900, pink50: pink50, pink100: pink100, pink200: pink200, pink300: pink300, pink400: pink400, pink500: pink500, pink600: pink600, pink700: pink700, pink800: pink800, pink900: pink900, purple50: purple50, purple100: purple100, purple200: purple200, purple300: purple300, purple400: purple400, purple500: purple500, purple600: purple600, purple700: purple700, purple800: purple800, purple900: purple900, deepPurple50: deepPurple50, deepPurple100: deepPurple100, deepPurple200: deepPurple200, deepPurple300: deepPurple300, deepPurple400: deepPurple400, deepPurple500: deepPurple500, deepPurple600: deepPurple600, deepPurple700: deepPurple700, deepPurple800: deepPurple800, deepPurple900: deepPurple900, indigo50: indigo50, indigo100: indigo100, indigo200: indigo200, indigo300: indigo300, indigo400: indigo400, indigo500: indigo500, indigo600: indigo600, indigo700: indigo700, indigo800: indigo800, indigo900: indigo900, blue50: blue50, blue100: blue100, blue200: blue200, blue300: blue300, blue400: blue400, blue500: blue500, blue600: blue600, blue700: blue700, blue800: blue800, blue900: blue900, lightBlue50: lightBlue50, lightBlue100: lightBlue100, lightBlue200: lightBlue200, lightBlue300: lightBlue300, lightBlue400: lightBlue400, lightBlue500: lightBlue500, lightBlue600: lightBlue600, lightBlue700: lightBlue700, lightBlue800: lightBlue800, lightBlue900: lightBlue900, cyan50: cyan50, cyan100: cyan100, cyan200: cyan200, cyan300: cyan300, cyan400: cyan400, cyan500: cyan500, cyan600: cyan600, cyan700: cyan700, cyan800: cyan800, cyan900: cyan900,  teal50: teal50, teal100: teal100, teal200: teal200, teal300: teal300, teal400: teal400, teal500: teal500, teal600: teal600, teal700: teal700, teal800: teal800, teal900: teal900, green50: green50, green100: green100, green200: green200, green300: green300, green400: green400, green500: green500, green600: green600, green700: green700, green800: green800, green900: green900, lightGreen50: lightGreen50, lightGreen100: lightGreen100, lightGreen200: lightGreen200, lightGreen300: lightGreen300, lightGreen400: lightGreen400, lightGreen500: lightGreen500, lightGreen600: lightGreen600, lightGreen700: lightGreen700, lightGreen800: lightGreen800, lightGreen900: lightGreen900, lime50: lime50, lime100: lime100, lime200: lime200, lime300: lime300, lime400: lime400, lime500: lime500, lime600: lime600, lime700: lime700, lime800: lime800, lime900: lime900, yellow50: yellow50, yellow100: yellow100, yellow200: yellow200, yellow300: yellow300, yellow400: yellow400, yellow500: yellow500, yellow600: yellow600, yellow700: yellow700, yellow800: yellow800, yellow900: yellow900, amber50: amber50, amber100: amber100, amber200: amber200, amber300: amber300, amber400: amber400, amber500: amber500, amber600: amber600, amber700: amber700, amber800: amber800, amber900: amber900, orange50: orange50, orange100: orange100, orange200: orange200, orange300: orange300, orange400: orange400, orange500: orange500, orange600: orange600, orange700: orange700, orange800: orange800, orange900: orange900, deepOrange50: deepOrange50, deepOrange100: deepOrange100, deepOrange200: deepOrange200, deepOrange300: deepOrange300, deepOrange400: deepOrange400, deepOrange500: deepOrange500, deepOrange600: deepOrange600, deepOrange700: deepOrange700, deepOrange800: deepOrange800, deepOrange900: deepOrange900, brown50: brown50, brown100: brown100, brown200: brown200, brown300: brown300, brown400: brown400, brown500: brown500, brown600: brown600, brown700: brown700, brown800: brown800, brown900: brown900, grey50: grey50, grey100: grey100, grey200: grey200, grey300: grey300, grey400: grey400, grey500: grey500, grey600: grey600, grey700: grey700, grey800: grey800, grey900: grey900, blueGrey50: blueGrey50, blueGrey100: blueGrey100, blueGrey200: blueGrey200, blueGrey300: blueGrey300, blueGrey400: blueGrey400, blueGrey500: blueGrey500, blueGrey600: blueGrey600, blueGrey700: blueGrey700, blueGrey800: blueGrey800, blueGrey900: blueGrey900
};


window.TimeoutCallbacks = [];
window.global = require('./globals/global_functions');
const path = require('path');
Array.prototype.getUnique = function(){
   var u = {}, a = [];
   for(var i = 0, l = this.length; i < l; ++i){
      if(u.hasOwnProperty(this[i])) {
         continue;
      }
      a.push(this[i]);
      u[this[i]] = 1;
   }
   return a;
};

window.performance = window.performance || false;
if (window.performance) {
  performance.now = (function() {
      return performance.now       ||
          performance.mozNow    ||
          performance.msNow     ||
          performance.oNow      ||
          performance.webkitNow ||
          Date.now  /*none found - fallback to browser default */
  })();
}


window.session_functions = {};
window.core = {};
window.core.Debug = {};
window.core.Debug.Dump = function() {
  var stack = new Error().stack;
  var line = "";
  if (stack) {
    if (stack.split("\n")[2].indexOf("makeAssimilatePrototype.js") > -1) {
      line = stack.split("\n")[3];
    } else {
      line = stack.split("\n")[2];
    }
  }
  console.debug("");
  var colors = "background: #222; color: #bada55";
  console.debug("%c!!!!!!!!!!!!! DEBUG (" + line + ") !!!!!!!!!!!!!", colors);
  //console.debug(stack);
  for (var i = 0; i < arguments.length; ++i) {
    console.debug("%cParm" + i + ":", colors, arguments[i]);
  }
  console.debug("%c!!!!!!!!!!!!! ENDDEBUG !!!!!!!!!!!!!", colors);
  console.debug("");
};

window.session_functions.Dump = window.core.Debug.Dump;

let UnloadFunctions = [];
let ReactIds = ['page'];
ReactIds.forEach(function(val){
    UnloadFunctions.push(function() {
        if (document.getElementById(val)) {
            ReactDom.unmountComponentAtNode(document.getElementById(val));
        }
    });
});

window.unloadAll = function() {
    UnloadFunctions.forEach(function(val){
        (val)();
    });
};

window.SpinnerTimedout = false;
window.ReloadOnWebSocketReconnect = false;
window.HeaderHeight = 0;
window.FooterHeight = 0;
window.ScrollHeight = 0;
window.InFullScreen = false;

window.unloadSideBarMenu = function() {
  if (document.getElementById('GoCore-sidebarMenu')) {
      ReactDom.unmountComponentAtNode(document.getElementById('GoCore-sidebarMenu'));
  }
};

window.displayErrorAlert = true;

window.AppLoad = function() {
    Object.defineProperty(Array.prototype, 'chunk', {
        value: function(chunkSize) {
            var array=this;
            return [].concat.apply([],
                array.map(function(elem,i) {
                    return i%chunkSize ? [] : [array.slice(i,i+chunkSize)];
                })
            );
        }
    });
};


window.onerror = (msg, url, line, col, error) => {
  if (window.displayErrorAlert) {
     var extra = !col ? '' : '\n' + window.appContent.WindowOnErrorWordColumn + ': ' + col;
     extra += !error ? '' : '\n' + window.appContent.WindowOnErrorWordError + ': ' + error;
     if (msg != "") {
       var stack = "\n\n" + new Error().stack;

       // if (error != null && error.hasOwnProperty("message") && typeof(error.message) == "string" && error.message.indexOf("replaceChild") != -1 && error.message.indexOf("Cannot") != -1 && error.message.indexOf("property") != -1 && window.erroredOut == -1) {
       //   //React mount issue because of dirty error on last page.  I dont think we can recover from this easily due to how we mount things.
       //   window.location.reload();
       //   return;
       // }

       window.globals.PopupJSError(window.appContent.WindowOnErrorClientSide + ": " + msg, "\n" + window.appContent.WindowOnErrorWordURL + ": " + url + "\n" + window.appContent.WindowOnErrorWordLine + ": " + line + extra + stack, "", "Error", true, {ClientSide: true});
     }
  }
  return false; // False here allows the error to propagate back to console
};

window.launcher = new Launcher();
window.store = new Store();
window.baseLogger = new BaseLogger();

// Needed for onTouchTap
// http://stackoverflow.com/a/34015469/988941
injectTapEventPlugin();

import Notifications from './notifications/notifications'
import Home from './pages/home/home'
import Settings from './pages/settings/settings'
import Footer from './footer/footer'
import Loader from './loader/loader'
import Banner from './banner/banner'
import SideBarMenu from './sidebarMenu/sidebarMenu'
import AccountList from './pages/accountList/accountList'
import AccountModify from './pages/accountModify/accountModify'
import AccountAdd from './pages/accountAdd/accountAdd'
import UserList from './pages/userList/userList'
import UserModify from './pages/userModify/userModify'
import UserAdd from './pages/userAdd/userAdd'
import UserProfile from './pages/userProfile/userProfile'
import PasswordReset from './pages/passwordReset/passwordReset'
import Transactions from './pages/transactions/transactions'
import TransactionModify from './pages/transactionModify/transactionModify'
import TransactionList from './pages/transactionList/transactionList'
import TransactionAdd from './pages/transactionAdd/transactionAdd'
import AppErrors from './pages/appErrors/appErrors'
import AppErrorModify from './pages/appErrorModify/appErrorModify'
import AppErrorList from './pages/appErrorList/appErrorList'
import AppErrorAdd from './pages/appErrorAdd/appErrorAdd'
import Features from './pages/features/features'
import FeatureModify from './pages/featureModify/featureModify'
import FeatureList from './pages/featureList/featureList'
import FeatureAdd from './pages/featureAdd/featureAdd'
import RoleFeatures from './pages/roleFeatures/roleFeatures'
import RoleFeatureModify from './pages/roleFeatureModify/roleFeatureModify'
import RoleFeatureList from './pages/roleFeatureList/roleFeatureList'
import RoleFeatureAdd from './pages/roleFeatureAdd/roleFeatureAdd'
import FeatureGroups from './pages/featureGroups/featureGroups'
import FeatureGroupModify from './pages/featureGroupModify/featureGroupModify'
import FeatureGroupList from './pages/featureGroupList/featureGroupList'
import FeatureGroupAdd from './pages/featureGroupAdd/featureGroupAdd'
import Roles from './pages/roles/roles'
import RoleModify from './pages/roleModify/roleModify'
import RoleList from './pages/roleList/roleList'
import RoleAdd from './pages/roleAdd/roleAdd'
import FileObjects from './pages/fileObjects/fileObjects'
import FileObjectModify from './pages/fileObjectModify/fileObjectModify'
import FileObjectList from './pages/fileObjectList/fileObjectList'
import FileObjectAdd from './pages/fileObjectAdd/fileObjectAdd'
import Logs from './pages/logs/logs'
import ServerSettingsModify from './pages/serverSettingsModify/serverSettingsModify'