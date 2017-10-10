// import ReactPerfAnalysis from "react-addons-perf";
import {React, SelectField, IconButton, Avatar, MenuItem, Checkbox} from "./forms";
import ReactDom from "react-dom";
import {grey900, blueGrey100} from "material-ui/styles/colors";
import Loader from "../loader/loader";
import {
  CameraIcon,
  DisplayIcon,
  Select,
  ZoomIn,
  ZoomOut,
  BrightnessDown,
  BrightnessUp,
  Rotate,
  PhotoCamera,
  Pan,
  Mask,
  Input,
  Menu,
  PIP,
  BackArrow,
  Home,
  Up,
  Down,
  Left,
  Right,
  Replay,
  ConferencingIcon,
  Rewind,
  PlayPause,
  FastForward,
  ProjectorIcon,
  ToolsIcon,
  AntennaTVIcon,
  Stop,
  Pause,
  Play,
  SkipNext,
  SkipPrevious,
  ShuffleIcon,
  LoopIcon,
  Record,
  Add,
  Minus,
  ComponentIcon,
  SvideoIcon,
  HDMIIcon,
  PowerIcon,
  VolumeMute,
  VolumeUnMute,
  SearchIcon,
  AvRadio,
  AvVideocam,
  AvVideocamOff,
  AvMic,
  AvMicOff,
  AvMicNone,
  AvMovie,
  AvHighQuality,
  AvHd,
  AvGames,
  AvForward5,
  AvForward10,
  AvForward30,
  AvFiberDvr,
  AvClosedCaption,
  AvEqualizer,
  VGA,
  UndoIcon,
  Freeze,
  EmptyCheckBox,
  ImageNegative,
  Image,
  Sync,
  FileUpload,
  VideoDistributionIcon,
  HubIcon,
  CPUIcon,
  VideoGame,
  FileDownload,
  Fullscreen,
  FullscreenExit,
  Ethernet,
  CD,
  Eject,
  Favorite,
  List,
  USB,
  InfoIcon,
  Snooze,
  NavigationClose,
  SimCard,
  LooksOne,
  LooksTwo,
  CommunicationCall,
  Help,
  AntennaIcon,
  MusicIcon,
  QueueMusicIcon,
  SubscriptionsIcon,
  RecordIcon,
  TouchPadIcon,
  DVRIcon,
  AudioNoteIcon,
  AppleTVIcon,
  ComputerMonitorIcon,
  DVDIcon,
  LightIcon,
  FanIcon,
  BluRayIcon,
  SpeakerGroupIcon,
  SpeakerIcon,
  MemoryIcon,
  RadioIcon,
  SurroundIcon,
  DeviceIcon,
  MountIcon,
  DiscIcon,
  ArrowForward
} from "./icons";
import ScreenShadeIcon from "material-ui/svg-icons/editor/format-line-spacing";
import KeyboardIcon from "material-ui/svg-icons/hardware/keyboard";
import SpeakerPhone from "material-ui/svg-icons/communication/speaker-phone";
const objectAssign = require('object-assign');


var globals = {};

globals.quickSpinnerDelayedOff = function (extraDelay=null) {
  let delay = 200;
  if (extraDelay != null) {
    delay = extraDelay;
  }
  window.setTimeout(() => {
    globals.quickSpinner(false);
  }, delay);
}

globals.quickSpinner = function (onOff=false, cb=null) {
  if (onOff === true) {
    window.SpinnerTimedout = false
    window.setTimeout(() => {
      window.SpinnerTimedout = true;
    }, 2000)
    $("#GoCore-loader").addClass("go-core-loading");
  } else {
    $("#GoCore-loader").removeClass("go-core-loading");
  }
  if (cb != null) {
    setTimeout(cb, 200);
  }
};

globals.buildUrl = function(controller, action, jsonObject) {
  return window.location.origin + "/#/" + controller+ "?action=" + action + "&uriParams=" + window.btoa(JSON.stringify(jsonObject))
};

globals.isArray = function(pointer) {
  return (pointer != undefined && pointer != null && Array.isArray(pointer))
};

globals.length = function(pointer) {
  if (globals.isArray(pointer)) {
    return pointer.length;
  } else {
    return 0;
  }
};

globals.map = function(pointer, mapFunc, elseFunc) {
  if (elseFunc == undefined) {
    elseFunc = () => {
      return []
    }
  }
  if (globals.isArray(pointer)) {
    return pointer.map(mapFunc);
  } else {
    return elseFunc();
  }
};

globals.filter = function(pointer, mapFunc, elseFunc) {
  if (elseFunc == undefined) {
    elseFunc = () => {
      return []
    }
  }
  if (globals.isArray(pointer)) {
    return pointer.filter(mapFunc);
  } else {
    return elseFunc();
  }
};

globals.reduce = function(pointer, mapFunc, elseFunc) {
  if (elseFunc == undefined) {
    elseFunc = () => {
      return []
    }
  }
  if (globals.isArray(pointer)) {
    return pointer.reduce(mapFunc);
  } else {
    return elseFunc();
  }
};

globals.forEach = function(pointer, mapFunc, elseFunc) {
  if (elseFunc == undefined) {
    elseFunc = () => {
      return []
    }
  }
  if (globals.isArray(pointer)) {
    pointer.forEach(mapFunc);
  }
};


globals.filter = function(pointer, mapFunc, elseFunc) {
  if (elseFunc == undefined) {
    elseFunc = () => {
      return []
    }
  }
  if (globals.isArray(pointer)) {
    return pointer.filter(mapFunc);
  } else {
    return elseFunc();
  }
};

globals.filter = function(pointer, mapFunc, elseFunc) {
  if (elseFunc == undefined) {
    elseFunc = () => {
      return []
    }
  }
  if (globals.isArray(pointer)) {
    return pointer.filter(mapFunc);
  } else {
    return elseFunc();
  }
};

globals.findIndex = function(pointer, mapFunc, elseFunc) {
  if (elseFunc == undefined) {
    elseFunc = () => {
      return []
    }
  }
  if (globals.isArray(pointer)) {
    return pointer.findIndex(mapFunc);
  } else {
    return elseFunc();
  }
};

globals.reduceRight = function(pointer, mapFunc, elseFunc) {
  if (elseFunc == undefined) {
    elseFunc = () => {
      return []
    }
  }
  if (globals.isArray(pointer)) {
    return pointer.reduceRight(mapFunc);
  } else {
    return elseFunc();
  }
};

globals.some = function(pointer, mapFunc, elseFunc) {
  if (elseFunc == undefined) {
    elseFunc = () => {
      return []
    }
  }
  if (globals.isArray(pointer)) {
    return pointer.some(mapFunc);
  } else {
    return elseFunc();
  }
};

globals.__ = function(tpl, replacements) {
  return globals._(tpl, replacements);
};

globals._ = function(tpl, replacements) {
  if (typeof(tpl) == "undefined") {
    return tpl;
  }
  $.each(replacements, function(k,v) {
    tpl = tpl.replace("{" + k + "}",v);
  });
  return tpl;
};

globals.ComponentError = function(classname, message) {
  session_functions.Dump(classname + " errored with", message);
  return <div className="example" style={{margin: 15 , backgroundColor: window.materialColors["deepOrange400"]}}>
    <div className="step error" style={{ backgroundColor: window.materialColors["deepOrange400"]}}>
      <span className="step-number">!</span>
      <div className="step-desc">
        <span className="step-title">Error</span>
        <p>{globals._(window.appContent.ErrorJs2, {classname: classname, message: message})}</p>
      </div>
    </div>
  </div>
};


globals.HasRole = function(role) {
  if (!window.appState.HasRole.hasOwnProperty(role)) {
    return false;
  }
  return window.appState.HasRole[role];
};


globals.GetUriParams = function() {
  var ret = {};
  var value = globals.GetVar("uriParams");
  if (value != null) {
    try {
      ret = JSON.parse(window.atob(value));
    } catch (e) {}
  }
  return ret;
};

globals.ViewLog = function(logName) {
  window.api.newWindow({"controller": "logs", uriParams:{Id: logName}})
};

globals.GetHeaderElement = function() {
  var headerElement;
  if (window.innerWidth < 767) {
    headerElement = $("#navbar-mobile");
  } else {
    headerElement = $("#navbar-desktop");
  }
  return headerElement;
};

globals.GetHeaderElement = function() {
  var headerElement;
  if (window.innerWidth < 767) {
    headerElement = $("#navbar-mobile");
  } else {
    headerElement = $("#navbar-desktop");
  }
  return headerElement;
};


globals.FullScreenContent = function() {
  window.InFullScreen = true;
  window.paddingTopOnly = 0;
  window.paddingLeftOnly = 0;
  window.paddingRightOnly = 0;
  window.paddingBottomOnly = 0;
};

globals.MostlyFullScreenContent = function() {
  window.InFullScreen = false;
  window.paddingTopOnly = 18;
  window.paddingLeftOnly = 15;
  window.paddingRightOnly = 0;
  window.paddingBottomOnly = 0;
};

globals.NormalScreenContent = function() {
  window.InFullScreen = false;
  window.paddingTopOnly = 0;
  window.paddingLeftOnly = 0;
  window.paddingRightOnly = 0;
  window.paddingBottomOnly = 0;
};

globals.GetVar = function(name, url)
{
  if (!url) {
    url = window.location.href;
  }
  name = name.replace(/[\[\]]/g, "\\$&");
  var regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)"),
      results = regex.exec(url);
  if (!results) return null;
  if (!results[2]) return '';
  return decodeURIComponent(results[2].replace(/\+/g, " "));
};

globals.IsRoomConfiguredFully = function(roomVM, renderFooterError=false) {
  if (roomVM) {
    if (roomVM.Joins.RoomDevices.Count == 0) {
      if (window.appState.SynchronizationId == "") {
        if (renderFooterError) {
          window.api.get({action:"Load", uriParams:{Id:roomVM.Id}, controller:"roomModifyDevices"});
          window.appState.SnackbarMessage = window.appContent.RoomControlComponentsRoomRequiresTechnology;
          window.appState.SnackbarOpen = true;
          window.appState.SnackbarType = "Error";
          window.goCore.setFooterStateFromExternal(window.appState);
        }
      }
      return false;
    }
  }

  // var displays = roomVM.Joins.RoomDevices.Items.reduce((roomDevice) => roomDevice.Joins && roomDevice.Joins.Device.Joins.Equipment.IsDisplayProjector);

  var displays = 0;
  $.each(roomVM.Joins.RoomDevices.Items, (k, roomDevice) => {
    if (roomDevice.Joins.Device != undefined && roomDevice.Joins.Device.Joins.Equipment != undefined && roomDevice.Joins.Device.Joins.Equipment.IsDisplayProjector) {
      displays++;
    }
  });

  if (displays == 0) {
      if (window.appState.SynchronizationId == "") {
        if (renderFooterError) {
          window.api.get({action: "Load", uriParams: {Id: roomVM.Id}, controller: "roomModifyDevices"});
          window.appState.SnackbarMessage = window.appContent.RoomControlComponentsRoomRequiresDisplay;
          window.appState.SnackbarOpen = true;
          window.appState.SnackbarType = "Error";
          window.goCore.setFooterStateFromExternal(window.appState);
        }
      }
      return false;
  }

  // var inputs_enabled = roomVM.Joins.RoomDevices.Items.reduce((roomDevice) =>  roomDevice && roomDevice.Inputs.reduce((input, i) => input.Enabled));
  var inputs_enabled = 0;
  $.each(roomVM.Joins.RoomDevices.Items, (k, roomDevice) => {
    $.each(roomDevice.Inputs, (k2, input) => {
      if (input.Enabled) {
        inputs_enabled++;
      }
    })
  });

  if (inputs_enabled == 0) {
      if (window.appState.SynchronizationId == "") {
        if (renderFooterError) {
          window.api.get({action: "Load", uriParams: {Id: roomVM.Id}, controller: "roomModifyDevices"});
          window.appState.SnackbarMessage = window.appContent.RoomControlComponentsRoomRequiresSource;
          window.appState.SnackbarOpen = true;
          window.appState.SnackbarType = "Error";
          window.goCore.setFooterStateFromExternal(window.appState);
        }
      }
      return false;
  }
  return true;
};

globals.PopupJSError = function(globalMessage, trace="") {
  window.launcher.renderClientFooterControls(globalMessage, trace);
};

globals.getValidMasterAudioDevices = function(room) {
  if (room.Joins.RoomDevices.Count > 0) {
    var items = room.Joins.RoomDevices.Items.filter((rd) => {
      if (rd.Joins.hasOwnProperty("Device")) {
        if (rd.Joins.Device.Joins.hasOwnProperty("Equipment")) {
          if (rd.Joins.Device.Joins.Equipment.Drivers.Commands != null && rd.Joins.Device.Joins.Equipment.Drivers.Commands.length > 0) {
            var hasCommand = false;
            if (rd.Joins.Device.Joins.Equipment.DeviceDriverId == "580eb20449520212e233b51a") {
              // Must have all of them
              let hasVolume = false;
              let hasMute = false;
              rd.Joins.Device.Joins.Equipment.Drivers.Commands.forEach((command) => {
                if ($.inArray(command.Type, ["VOLUME UP"]) != -1 || $.inArray(command.Type, ["VOLUME DOWN"]) != -1) {
                  hasVolume = true;
                }
                if ($.inArray(command.Type, ["MUTE ON"]) != -1 || $.inArray(command.Type, ["MUTE OFF"]) != -1 ) {
                  hasMute = true;
                }
              });
              if (hasVolume && hasMute) {
                return true;
              }
            } else {
              rd.Joins.Device.Joins.Equipment.Drivers.Commands.forEach((command) => {
                if ($.inArray(command.Type, ["MUTE TOGGLE", "MUTE ON", "VOLUME UP", "VOLUME DOWN"]) != -1) {
                  hasCommand = true;
                }
              });
              if (hasCommand) {
                return true;
              }
            }
          }
        }
      }
    });
    if (globals.length(items) > 0) {
      items.unshift({key:"", value: window.appContent.None});
    }
    return {items: items};
  }
  return {items: []}
};

globals.ObjectId = function () {
    var timestamp = (new Date().getTime() / 1000 | 0).toString(16);
    return timestamp + 'xxxxxxxxxxxxxxxx'.replace(/[x]/g, function() {
        return (Math.random() * 16 | 0).toString(16);
    }).toLowerCase();
};

globals.Popup = function(globalMessage, trace="", globalMessageType="", transactionId="") {
  window.launcher.renderFooterControls(globalMessage, trace, globalMessageType, transactionId);
};

globals.PopupWindow = function(info="", title="Information") {
  if (title != "") {
    window.appState["DialogTitle2"] = globals.translate(title);
  }
  window.launcher.renderFooterControls("", info);
  window.launcher.showDialogMinimal();
};

globals.RemoteControl = function(remoteCommand, cb=()=>{}) {
  remoteCommand.CurrentActionUser = window.appState.UserId;
  console.log(remoteCommand);
  if (window.appState.DeveloperMode) {
    core.Debug.Dump("RemoteCommand Call", JSON.stringify(remoteCommand));
  }
  window.api.post({action:"RemoteControl", controller: "roomControl", state: {RemoteCommand: remoteCommand}, leaveStateAlone: true, callback: cb});
};

globals.log = function(log) {
  if (window.appState.DeveloperMode) {
    core.Debug.Dump(log);
  }
};

globals.render = function(children) {
  window.goCore.page = ReactDom.render(<Loader>{children}</Loader>, document.getElementById("page"))
};

globals.handleUnSynchronizeBroadcast = function(data) {
  if (window.hasOwnProperty("SynchronizationId") && window.SynchronizationId == data.Synchronization.Id) {
    // setTimeout(() => {
       window.api.post({action: "Logout", state: {}, controller:"login"});
    //  },2000);

  }
};



globals.basePageStyle = function() {
  if ($(window).width() > 768) {
    if (!window.InFullScreen) {
      if (window.paddingTopOnly > 0 || window.paddingBottomOnly > 0 || window.paddingLeftOnly > 0 || window.paddingRightOnly > 0) {
        let newPadding = {};
        newPadding.paddingTop = (window.paddingTopOnly > 0) ? window.paddingTopOnly: 0;
        newPadding.paddingRight = (window.paddingRightOnly > 0) ? window.paddingRightOnly: 0;
        newPadding.paddingBottom = (window.paddingBottomOnly > 0) ? window.paddingBottomOnly: 0;
        newPadding.paddingLeft = (window.paddingLeftOnly > 0) ? window.paddingLeftOnly: 0;
        return newPadding
      } else {
        return {paddingTop: 18, paddingRight: 30, paddingBottom: 30, paddingLeft: 30}
      }
    } else {
      return {}
    }
  } else {
    return {}
  }
};

globals.pageEnd = function() {
  if (window.appState.DeveloperMode) {
    if (window.pageState.DeveloperLog != "") {
      console.debug("<ServerLogs>");
      try {
        console.debug(window.atob(window.pageState.DeveloperLog));
      } catch(e) {}
      console.debug("</ServerLogs>");
    }

    let bytesPage = window.global.functions.roughSizeOfObject(window.pageState);
    let bytesApp = window.global.functions.roughSizeOfObject(window.appState);
    let bytesContent = window.global.functions.roughSizeOfObject(window.appContent);
    let bytesContentPage = window.global.functions.roughSizeOfObject(window.pageContent);
    let totalBytes = bytesPage+bytesApp+bytesContent+bytesContentPage;
    core.Debug.Dump("!!! MEMORY USAGE !!!", "pageState size: " + window.global.functions.humanFileSize(bytesPage), "appState size: " + window.global.functions.humanFileSize(bytesApp), "appContent size: " + window.global.functions.humanFileSize(bytesContent), "pageContent size: " + window.global.functions.humanFileSize(bytesContentPage), "total memory: " + window.global.functions.humanFileSize(totalBytes),  "!!! MEMORY USAGE !!!");

    console.log(".................................React Load Ended", new Date());
    try {
      if (window.appState.DeveloperLogStateChangePerformance && window.performance) {
        globals.createMark("react-load-done");
        window.performance.measure('react-load-final', 'react-load', 'react-load-done');
        if (window.goCore.page.constructor.name == 'ButtonBarPage' || window.goCore.page.constructor.name == 'PlainPage') {
          var pageComponent = window.goCore.page.props.page.type.prototype.constructor.name;
        } else {
          var pageComponent = window.goCore.page.props.children.type.name;
        }
        if (pageComponent == 'AddRecordPageComponent' || pageComponent == 'BackPageComponent' ) {
          var pageComponent = window.goCore.page.props.page.type.prototype.constructor.prototype.__proto__.constructor.name;
        }
        if (window.performance) {
          console.info("pageLoaded for <" + pageComponent + "/> in " + window.performance.getEntriesByName('react-load-final')[0].duration + 'ms');
        }
        if (totalBytes > 1000000 && window.appState.DeveloperMode) {
          //alert("Warning, your total state and content exceeds 1MB (Current at around " + window.global.functions.humanFileSize(totalBytes) + ").  Trim it down to help with page load time (" + (performance.now()/1000) + " seconds if you did a full page refresh).")
        }
        performance.clearMeasures("react-load-final");

        // commented out... see https://github.com/facebook/react/commit/cbc450860e895a8bae0591ef8eaed1dabdf30fa6#commitcomment-22144714
        // ReactPerfAnalysis.stop();
        // console.info(">>>>>>>>>>>>>>>>>>>>>>>  React Render Performance (Expand me please!!) <<<<<<<<<<<<<<<<<<<<<");
        // console.info("Time Wasted Below:");
        // ReactPerfAnalysis.printWasted();
        // console.info("Render Time In Your Code(Without react):");
        // ReactPerfAnalysis.printExclusive();
        // console.info("Render Time Complete Breakdown:");
        // ReactPerfAnalysis.printInclusive();
      }
    } catch (e) {
      console.log("Error in developer mode pageEnd", e);
    }
  }
};

globals.createMark = function(name) {
  if (performance.mark === undefined) {
    console.log("performance.mark Not supported");
    return;
  }
  // Create the performance mark
  performance.mark(name);
};


globals.pageStart = function() {
  if (window.appState.DeveloperMode && window.appState.DeveloperLogStateChangePerformance) {
    console.log(".................................React Load Started", new Date());
    // commented out... see https://github.com/facebook/react/commit/cbc450860e895a8bae0591ef8eaed1dabdf30fa6#commitcomment-22144714
    // ReactPerfAnalysis.start();
    globals.createMark("react-load");
  }
};

globals.IsMasterAccount = function() {
  return window.appState.UserPrimaryAccount == "58405718f94c671b05350857";
};

globals.IsSystemAccount = function() {
  return window.appState.IsSystemAccount;
};



globals.getRoomImagePath = function(room) {
  let roomImageVectorPath = globals.getRoomImageInfo().path;
  let roomImageStockPath = globals.getRoomStockImageInfo().path;
  let image = null;
  if (room.CustomFileType == "custom") {
    image = "/fileObject/" + room.ImageCustom + "?" + new Date().getTime();
  } else if (room.CustomFileType == "real") {
    image = roomImageStockPath + room.ImageFileNameNonVector + ".jpg";
  } else if (room.CustomFileType == "vector") {
    image = roomImageVectorPath + room.ImageFileName + ".jpg";
  }
  return image
};

globals.getBuildingImageInfo = function() {
  return {
    images: [
      {
        Group: true,
        Name: window.appContent.BuildingTypeAssembly,
        Id: "id-assembly",
      },
      {
        Name: window.appContent.BuildingCafe,
        Id: "assembly-cafe",
      },
      {
        Name: window.appContent.BuildingChurch,
        Id: "assembly-church-alt",
      },
      {
        Name: window.appContent.BuildingPub,
        Id: "assembly-pub",
      },
      {
        Name: window.appContent.BuildingRestauraunt2,
        Id: "assembly-restaurant",
      },
      {
        Name: window.appContent.BuildingRestauraunt,
        Id: "assembly-small-restaurant-house",
      },
      {
        Group: true,
        Name: window.appContent.BuildingTypeBusiness,
        Id: "id-business",
      },
      {
        Name: window.appContent.BuildingBankBranch,
        Id: "business-bank-alt",
      },
      {
        Name: window.appContent.BuildingBankHeadquarters,
        Id: "business-bank-alt2",
      },
      {
        Name: window.appContent.BuildingCinema,
        Id: "business-cinema",
      },
      {
        Name: window.appContent.BuildingCinemaAlt,
        Id: "business-cinema-alt2",
      },
      {
        Name: window.appContent.BuildingConventionCenter,
        Id: "business-convention-center",
      },
      {
        Name: window.appContent.BuildingGym,
        Id: "business-gym",
      },
      {
        Name: window.appContent.BuildingGymAlt,
        Id: "business-gym-alt",
      },
      {
        Name: window.appContent.BuildingMuseum,
        Id: "business-palace-museum",
      },
      {
        Name: window.appContent.BuildingTheater,
        Id: "business-theater",
      },
      {
        Name: window.appContent.BuildingMuseumAlt,
        Id: "business-museum-alt",
      },
      {
        Group: true,
        Name: window.appContent.BuildingTypeEducational,
        Id: "id-education",
      },
      {
        Name: window.appContent.BuildingElementarySchool,
        Id: "educational-elementary-school",
      },
      {
        Name: window.appContent.BuildingHighSchoolAlt,
        Id: "educational-high-school-alt",
      },
      {
        Name: window.appContent.BuildingHighSchool,
        Id: "educational-high-school",
      },
      {
        Name: window.appContent.BuildingSchoolLarge,
        Id: "educational-large-school",
      },
      {
        Name: window.appContent.BuildingUniversity,
        Id: "educational-university-alt",
      },
      {
        Group: true,
        Name: window.appContent.BuildingTypeEnterprise,
        Id: "id-enterprise",
      },
      {
        Name: window.appContent.BuildingMostlyGlassBuilding,
        Id: "enterprise-7-story-windowed-office",
      },
      {
        Name: window.appContent.BuildingMostlyBrownSkyscraper,
        Id: "enterprise-brown-skyscraper",
      },
      {
        Name: window.appContent.BuildingDarkSkyscraper1,
        Id: "enterprise-detailed-high-rise-1",
      },
      {
        Name: window.appContent.BuildingDarkSkyscraper2,
        Id: "enterprise-detailed-high-rise-2",
      },
      {
        Name: window.appContent.BuildingDarkSkyscraper3,
        Id: "enterprise-detailed-high-rise-3",
      },
      {
        Name: window.appContent.BuildingDarkSkyscraper4,
        Id: "enterprise-detailed-high-rise-4",
      },
      {
        Name: window.appContent.BuildingDarkSkyscraper5,
        Id: "enterprise-detailed-high-rise-5",
      },
      {
        Name: window.appContent.BuildingDarkSkyscraper6,
        Id: "enterprise-three-detailed-high-rises",
      },
      {
        Name: window.appContent.BuildingLightSkyscraper1,
        Id: "enterprise-high-rise-1",
      },
      {
        Name: window.appContent.BuildingLightSkyscraper2,
        Id: "enterprise-high-rise-2",
      },
      {
        Name: window.appContent.BuildingLightSkyscraper3,
        Id: "enterprise-high-rise-3",
      },
      {
        Name: window.appContent.BuildingLightSkyscraper4,
        Id: "enterprise-high-rise-4",
      },
      {
        Name: window.appContent.BuildingLightSkyscraper5,
        Id: "enterprise-high-rise-multi-5",
      },
      // -- could not get grey background easily on this.
      // {
      //   Name: window.appContent.BuildingLightSkyscraper6,
      //   Id: "enterprise-detailed-high-rise-6",
      // },
      {
        Name: window.appContent.BuildingTealSkyscraper,
        Id: "enterprise-teal-skyscraper",
      },
      {
        Group: true,
        Name: window.appContent.BuildingTypeGovernment,
        Id: "id-government",
      },
      {
        Name: window.appContent.BuildingCapitol,
        Id: "government-capitol",
      },
      {
        Name: window.appContent.BuildingCityHall,
        Id: "government-city-hall",
      },
      {
        Name: window.appContent.BuildingCityCourt,
        Id: "government-court-alt",
      },
      {
        Name: window.appContent.BuildingEmbassy,
        Id: "government-embassy",
      },
      {
        Name: window.appContent.BuildingFireStation,
        Id: "government-fire-station-alt",
      },
      {
        Name: window.appContent.BuildingGenericBuilding,
        Id: "government-generic-building",
      },
      {
        Name: window.appContent.BuildingGenericBuilding2,
        Id: "government-parliament",
      },
      {
        Name: window.appContent.BuildingGenericLibrary,
        Id: "government-library-alt",
      },
      {
        Name: window.appContent.BuildingPostOffice,
        Id: "government-post-office-alt",
      },
      {
        Name: window.appContent.BuildingPoliceStation,
        Id: "government-police-alt",
      },
      {
        Group: true,
        Name: window.appContent.BuildingTypeInstitution,
        Id: "id-institution",
      },
      {
        Name: window.appContent.BuildingHospital,
        Id: "institutional-hospital",
      },
      {
        Name: window.appContent.BuildingHospitalAlt,
        Id: "institutional-hospital-alt",
      },
      {
        Name: window.appContent.BuildingPrison,
        Id: "institutional-prison-alt",
      },
      {
        Name: window.appContent.BuildingClinic,
        Id: "institutional-clinic",
      },
      {
        Group: true,
        Name: window.appContent.BuildingTypeMercantile,
        Id: "id-mercantile",
      },
      {
        Name: window.appContent.BuildingConvenienceStore,
        Id: "mercantile-convenience-store",
      },
      {
        Name: window.appContent.BuildingConvenienceStoreHighRise,
        Id: "mercantile-convenience-store-high-rise",
      },
      {
        Name: window.appContent.BuildingInternationalStore,
        Id: "mercantile-international-shop",
      },
      {
        Name: window.appContent.BuildingRetailBuilding,
        Id: "mercantile-retail",
      },
      {
        Name: window.appContent.BuildingRetailBuildingAlt,
        Id: "mercantile-retail-alt",
      },
      {
        Name: window.appContent.BuildingShoppingMall,
        Id: "mercantile-shopping-mall",
      },
      {
        Name: window.appContent.BuildingShoppingMallAlt,
        Id: "mercantile-shopping-mall-alt",
      },
      {
        Name: window.appContent.BuildingSuperMarket,
        Id: "mercantile-supermarket",
      },
      {
        Name: window.appContent.BuildingSuperMarketAlt,
        Id: "mercantile-supermarket-alt",
      },
      {
        Group: true,
        Name: window.appContent.BuildingTypeResidential,
        Id: "id-residential",
      },
      {
        Name: window.appContent.BuildingApartment,
        Id: "residential-apartment",
      },
      {
        Name: window.appContent.BuildingApartmentPink,
        Id: "residential-apartment-pink",
      },
      {
        Name: window.appContent.BuildingColorfulTownHome,
        Id: "residential-colorful-town-home",
      },
      {
        Name: window.appContent.BuildingGatedMansion,
        Id: "residential-gated-mansion",
      },
      {
        Name: window.appContent.BuildingHistoricDowntownBuilding,
        Id: "residential-historic-downtown-building",
      },
      {
        Name: window.appContent.BuildingHistoricTownHome,
        Id: "residential-historic-town-home",
      },
      {
        Name: window.appContent.BuildingHistoricHotel,
        Id: "residential-hotel-alt",
      },
      {
        Name: window.appContent.BuildingHistoricHotelAlt,
        Id: "residential-pink-hotel",
      },
      {
        Name: window.appContent.BuildingTanHotel,
        Id: "residential-hotel-tan",
      },
      {
        Name: window.appContent.BuildingSingleFamilyHouse,
        Id: "residential-house",
      },
      {
        Name: window.appContent.BuildingInternationalHouse,
        Id: "residential-international-house",
      },
      {
        Name: window.appContent.BuildingThreeStoryBrownApartment,
        Id: "residential-three-story-apartment-brown",
      },
      {
        Name: window.appContent.BuildingTwoStoryTanLofts,
        Id: "residential-two-story-tan-building",
      },
    ],
    path: "/web/app/images/buildings/vector-image-source/exports/"
  };
};


globals.controlButtonSelectBox = function(value, errorTxt, onChange, width) {
  var buttons = globals.getControlButtons();
  let filtered = buttons.filter((v) => !v.IsInput).filter((v) => v.GlobalIcon);
  let exists = filtered.filter((v) => value == v.Id);

  if (exists.length == 0) {
    value = "Other";
  }
  return <SelectField floatingLabelText={"* " + window.appContent.CommandIconText}
               hintText={"* " + window.appContent.CommandIconText}
               value={value}
               onChange={onChange}
               style={{width: width}}
               errorText={errorTxt}>
    {filtered.map((v) => {
      let tx = "";
      if (!v.Icon) {
        tx = globals.translate(v.Text);
        if (tx != undefined && tx.length > 5) {
          tx = tx.substr(0, 4);
        }
      }
      var additionalProps = {};
      if (v.Group) {
        additionalProps.disabled = true;
      } else {
        additionalProps.leftIcon = <div style={(v.hasOwnProperty("ButtonColor")) ? {backgroundColor: window.materialColors[v.ButtonColor]}: null}>{(v.Icon !== false) ? v.Icon: tx}</div>;
      }
      return <MenuItem {...additionalProps} key={v.hasOwnProperty("UniqueKey") ? v.UniqueKey : v.Id} value={v.Id} primaryText={v.Name}/>;
    })}
    <MenuItem leftIcon={<span>?</span>} key={"Other"} value={"Other"} primaryText={"Other"}/>
  </SelectField>

};

globals.controlButtonInputSelectBox = function(value, errorTxt, onChange, width) {
  var buttons = globals.getControlButtons();
  return <div>
      <SelectField floatingLabelText={"* " + window.appContent.CommandTypeBtnType}
                   hintText={"* " + window.appContent.CommandTypeBtnType}
                   value={value}
                   onChange={onChange}
                   style={{width: width}}
                   errorText={errorTxt}>
        {buttons.filter((v) => v.IsInput).map((v) => {
          var additionalProps = {};
          if (v.Group) {
            additionalProps.disabled = true;
          } else {
            additionalProps.leftIcon = <div style={(v.hasOwnProperty("ButtonColor")) ? {backgroundColor: window.materialColors[v.ButtonColor]}: null}>{(v.Icon !== false) ? v.Icon: globals.translate(v.Text).substr(0,4)}</div>;
          }
          return <MenuItem {...additionalProps} key={v.hasOwnProperty("UniqueKey") ? v.UniqueKey : v.Id} value={v.Id} primaryText={v.Name}/>;
        })}
      </SelectField>
  </div>;
};

globals.colorPickerSelect = function(value, errorTxt, onChange, width, label, bkgFontColor, allowTransparent=true, required=true) {
  let defaultOption = null;
  if (bkgFontColor == "fontColor") {
    defaultOption = <MenuItem {...iconTextProps} key={"default" + "grey900"} value={"grey900"} primaryText={window.appContent.Default}/>;
  } else if (bkgFontColor == "backgroundColor") {
    defaultOption = <MenuItem {...bkgndProps} key={"default" + blueGrey100} value={"blueGrey100"} primaryText={window.appContent.Default}/>;
  }
  let bkgndProps = {leftIcon: <div style={{width: 25, height: 25, backgroundColor: "#ECEFF1"}}>&nbsp;</div>}
  let iconTextProps = {leftIcon: <div style={{width: 25, height: 25, backgroundColor: "#212121", fill: "grey900"}}>&nbsp;</div>}
  return <div>
      <SelectField floatingLabelText={(required ? "* " : "") + label}
                   hintText={(required ? "* " : "") + label}
                   value={value}
                   onChange={onChange}
                   style={{width: width}}
                   errorText={errorTxt}>
        {defaultOption}
        {Object.keys(window.materialColors).filter((colorKey) => {
          if (colorKey == "transparent" && !allowTransparent) {
            return false;
          }
          return true;
        }).map((colorKey) => {
          let additionalProps = {
            leftIcon: <div style={{width: 25, height: 25, backgroundColor: window.materialColors[colorKey]}}>&nbsp;</div>
          };
          return <MenuItem {...additionalProps} key={colorKey} value={colorKey} primaryText={colorKey}/>;
        })}
      </SelectField>
  </div>;
};

globals.color = function(str) {
  return window.materialColors[str];
};

globals.isEmpty = function(str) {
  return globals.isBlank(str);
};

globals.isBlank = function(str) {
  return (typeof(str) == "undefined" || str == null || str === "" || /^\s*$/.test(str));
};

globals.smallIcon = function() {
  return {
    width: 40,
    height: 40,
  }
};

globals.smallShadow = function() {
  return {
    width: 64,
    height: 64,
    padding: 7,
  }
};
globals.toTitleCase = function (str) {
    return str.replace(/(?:^|\s)\w/g, function(match) {
        return match.toUpperCase();
    });
};

globals.camelize = function (str) {
  return str.replace(/(?:^\w|[A-Z]|\b\w)/g, function(letter, index) {
    return index == 0 ? letter.toLowerCase() : letter.toUpperCase();
  }).replace(/\s+/g, '').replace(/[^\w\s]/gi, '');
};

globals.repeatEvery = function(func, interval, globalIntervalKey) {
    // Check current time and calculate the delay until next interval
    var now = new Date,
        delay = interval - now % interval;

    function start() {
        // Execute function now...
        func();
        // ... and every interval
        window[globalIntervalKey] = setInterval(func, interval);
    }

    // Delay execution until it's an even interval
    setTimeout(start, delay);
};

globals.updateUserPreferences = function(key, value) {
  var foundMatch = false;
  if (!window.appState.UserPreferences) {
    window.appState.UserPreferences = [];
  }
  $.each(window.appState.UserPreferences, (k, v)=> {
    if (v.hasOwnProperty("Key") && v.Key == key) {
      window.appState.UserPreferences[k].Value = value;
      foundMatch = true;
    }
  });
  if (!foundMatch) {
    window.appState.UserPreferences.push({Key: key, Value: value});
  }
  return window.appState.UserPreferences;
};

globals.SubmitUserPreferenceChange = function(key, value) {
  window.api.post({action: "UpdatePreferences", state: {UserPreferences: globals.updateUserPreferences(key, value)}, controller:"users", leaveStateAlone: true});
};

globals.getUserPreferences = function(key) {
  var foundMatch = null;
  if (!window.appState.UserPreferences) {
    window.appState.UserPreferences = [];
  }
  $.each(window.appState.UserPreferences, (k, v)=> {
    if (v.hasOwnProperty("Key") && v.Key == key) {
      foundMatch = v.Value;
    }
  });
  if (!foundMatch) {
    return "<none>";
  } else {
    return foundMatch;
  }
};

globals.guid = function() {
  function s4() {
    return Math.floor((1 + Math.random()) * 0x10000)
      .toString(16)
      .substring(1);
  }
  return s4() + s4() + s4() + s4() + s4() + s4();
};

globals.uuid = function() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
   var r = Math.random()*16|0, v = c === 'x' ? r : (r&0x3|0x8);
   return v.toString(16);
  });
};

globals.roughSizeOfObject = function( object ) {
  var objectList = [];
  var stack = [ object ];
  var bytes = 0;
  var loops = 0;

  while ( stack.length && loops < 100000 ) {
    var value = stack.pop();
    if ( typeof value === 'boolean' ) {
      bytes += 4;
    } else if ( typeof value === 'string' ) {
      bytes += value.length * 2;
    } else if ( typeof value === 'number' ) {
      bytes += 8;
    } else if (typeof value === 'object' && objectList.indexOf(value) === -1) {
      objectList.push(value);
      for(var i in value) {
        stack.push( value[ i ] );
      }
    }
    loops++;
  }
  return bytes;
};

globals.humanFileSize = function(bytes, si) {
    var thresh = si ? 1000 : 1024;
    if(Math.abs(bytes) < thresh) {
        return bytes + ' B';
    }
    var units = si
        ? ['kB','MB','GB','TB','PB','EB','ZB','YB']
        : ['KiB','MiB','GiB','TiB','PiB','EiB','ZiB','YiB'];
    var u = -1;
    do {
        bytes /= thresh;
        ++u;
    } while(Math.abs(bytes) >= thresh && u < units.length - 1);
    return bytes.toFixed(1)+' '+units[u];
};

globals.toggleSideMenu = function(event){
  var target = $( event.target );
  var liParent = target.parent(".site-menu-item");
  var ulSubMenu = liParent.find(".site-menu-sub");
  ulSubMenu.toggle();
};

globals.fetchRemarkIcon = function(key){
  //https://zavoloklom.github.io/material-design-iconic-font/icons.html
  switch(key){
    case "Dashboard":
      return "md-view-dashboard";
    case "Account":
      return "md-assignment-account";
    case "AccessAccount":
      return "md-collection-text";
    case "GlobeLock":
      return "md-globe-lock";
    case "Users":
      return "md-accounts-alt";
  }
  return "";
};

globals.productTitle = function(){
  return "GoCoreAppHumanName | ";
};

globals.translate = function(key) {

  if (window.pageContent != undefined && window.pageContent.hasOwnProperty(key)) {
    return window.pageContent[key];
  }

  if (window.appContent != undefined && window.appContent.hasOwnProperty(key)) {
    return window.appContent[key];
  }

  return key;
};

globals.translateTimeUnits = function(units, timeType, translate=true) {
  var plural = "";
  if (units > 1) {
    var plural = "Plural";
  }
  if (translate) {
    return globals.translate("TimeUnits" + timeType.toUpperCase() + plural);
  } else {
    return "TimeUnits" + timeType.toUpperCase() + plural;
  }
};

globals.getUrlParameter = function(sParam) {
    var sPageURL = decodeURIComponent(window.location.search.substring(1)),
        sURLVariables = sPageURL.split('&'),
        sParameterName,
        i;

    for (i = 0; i < sURLVariables.length; i++) {
        sParameterName = sURLVariables[i].split('=');

        if (sParameterName[0] === sParam) {
            return sParameterName[1] === undefined ? true : sParameterName[1];
        }
    }
};

globals.redirect = function(controllerUpperCase, action='Root', uriParms={}) {
    var parms = {action:action, uriParams:uriParms};
    if (controllerUpperCase.indexOf("http") != -1) {
      document.location = controllerUpperCase;
    }
    if (controllerUpperCase.indexOf("CONTROLLER") == -1) {
      controllerUpperCase = "CONTROLLER_" + controllerUpperCase;
    }
    parms.controller = window.appState.routes.Paths[controllerUpperCase];
    if (parms.controller != null) {
      window.api.get(parms);
    }
};

globals.resolveController = function(controllerUpperCase) {
    return window.appState.routes.Paths[controllerUpperCase];
};

globals.times = x => f => {
  if (x > 0) {
    f();
    globals.times (x - 1) (f)
  }
};

globals.widgetListDefaults = () => {
  return {
    noFieldValueTranslation:window.appContent.WidgetListNoFieldValue,
    noDataMessage:window.appContent.WidgetListNoRecordsFound
  }
};

globals.widgetListButtonBarOffset = () => {
  return {
    marginRight: 50,
    rowHeight: $(window).height() > 750 ? 70: 50,
    offsetHeightToList: $(window).height() > 750 ? 226: 150,
  }
};

globals.clickCurrentAddOrImportActionButton = ()  => {
  window.clickCurrentAddOrImportActionButton();
};

globals.FloatingActionButtonClick = (e, customOnClick, typeObject, action) => {
  var redirectFunc = () => {
    window.global.functions.redirect(action);
  };
  if (typeObject == "Back"){
    redirectFunc = () => window.history.back();
  }
  if (customOnClick) {
    customOnClick();
  } else {
    redirectFunc();
  }
};

globals.isObject = (item) => {
  return (item && typeof item === 'object' && !Array.isArray(item) && item !== null);
};

globals.mergeDeep = (target, source, log=null) => {
  let output = objectAssign({}, target);
  if (target === null) {
    target = {};
  }
  log = false;
  // if (log == null) {
  //   log = target.hasOwnProperty("currentButton") || target.hasOwnProperty("VideoRemoteState");
  // }

  if (globals.isObject(target) && globals.isObject(source)) {
    Object.keys(source).forEach(key => {
      if (globals.isObject(source[key])) {
        if (!(key in target)) {
          if (log) {core.Debug.Dump("in not", key, source[key], output);}
          objectAssign(output, { [key]: source[key] });
        } else {
          if (Object.keys(source[key]).length === 0) {
            objectAssign(output, { [key]: source[key] });
            if (log) {core.Debug.Dump("in if", key, source[key], output);}
          } else {
            if (log) {core.Debug.Dump("in deep", target[key], source[key]);}
            output[key] = globals.mergeDeep(target[key], source[key], log);
          }
        }
      } else {
        objectAssign(output, { [key]: source[key] });
        if (log) {core.Debug.Dump("elssss", key, source[key], output);}
      }
    });
  }
  return output;
};

globals.sortByKey = (array, key) => {
    if (array instanceof Array) {
        return array.sort(function (a, b) {
            var x = a[key]; var y = b[key];
            return ((x < y) ? -1 : ((x > y) ? 1 : 0));
        });
    }
}

globals.saveFile = function(data, fileName, type) {
    try {
        var textFileAsBlob = new Blob([data], { type: type  });
        var downloadLink = document.createElement("a");
        downloadLink.download = fileName;
        downloadLink.innerHTML = "Download File";
        if (window.webkitURL != null) {
            // Chrome allows the link to be clicked
            // without actually adding it to the DOM.
            downloadLink.href = window.webkitURL.createObjectURL(textFileAsBlob);
        } else {
            // Firefox requires the link to be added to the DOM
            // before it can be clicked.
            downloadLink.href = window.URL.createObjectURL(textFileAsBlob);
            downloadLink.onclick = globals.destroyClickedElement;
            downloadLink.style.display = "none";
            document.body.appendChild(downloadLink);
        }
        downloadLink.click();
    }
    catch (ex) {
        core.Debug.Dump("Error at jQueryDom.saveFile:  " + ex);
    }
};

globals.destroyClickedElement = function(event) {
  // remove the link from the DOM
  document.body.removeChild(event.target);
}

globals.formatAMPM = function(date) {
  var hours = date.getHours();
  var minutes = date.getMinutes();
  var ampm = hours >= 12 ? 'pm' : 'am';
  hours = hours % 12;
  hours = hours ? hours : 12; // the hour '0' should be '12'
  minutes = minutes < 10 ? '0'+minutes : minutes;
  var strTime = hours + ':' + minutes + ' ' + ampm;
  return strTime;
}

globals.formatDate = function(date) {
  var yyyy = date.getFullYear().toString();
    var mm = (date.getMonth()+1).toString();
    var dd  = date.getDate().toString();

    var mmChars = mm.split('');
    var ddChars = dd.split('');

    return yyyy + '-' + (mmChars[1]?mm:"0"+mmChars[0]) + '-' + (ddChars[1]?dd:"0"+ddChars[0]);
}

globals.getCookie = function(cname) {
    var name = cname + "=";
    var decodedCookie = decodeURIComponent(document.cookie);
    var ca = decodedCookie.split(';');
    for(var i = 0; i <ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0) == ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) == 0) {
            return c.substring(name.length, c.length);
        }
    }
    return "";
}


export let functions = globals;
window.globals = globals;
