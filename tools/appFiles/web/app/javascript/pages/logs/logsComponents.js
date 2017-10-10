import React, {Component} from 'react';
import BasePageComponent from '../../components/basePageComponent';
import {
    RaisedButton
} from "../../globals/forms";
import ReactDOM from 'react-dom';
import {
  Play,
  Pause
} from "../../globals/icons";
class Logs extends BasePageComponent {
  constructor(props, context) {
    super(props, context);
    this.state.paused = false;
    this.state.pausedDueToScroll = false;
    this.state.logData = "";
    this.state.dataLength = "";
    this.state.dataCachedLengthSize = 0;
    this.state.dataCachedLength = "";
    this.state.pausedLogs = "";
    this.state.showClear = false;
    this.state.clearCountDown = 0;
    window.currentPauseAction = false;
    window.isAutoScrollTop = false;
  }

   clearPausedLogs() {
    let changes = {pausedLogs: "", dataLength: window.global.functions.humanFileSize(window.global.functions.roughSizeOfObject(this.state.logData))};
    if (window.global.functions.roughSizeOfObject(this.state.logData) > 1000000 && !this.state.showClear) {
      window.currentPauseAction = true;
      changes.paused = true;
      this.state.showClear = true;
      changes.showClear = true;
      changes.clearCountDown = 60;

      document.title = this.globs.__(window.pageContent.ClearingLog, {sec: changes.clearCountDown})
      // After 60 seconds clear it
      window.setTimeout(() => {
        document.title = "GoCoreAppHumanName | Logs";
        window.clearInterval(this.clearId);
        window.currentPauseAction = false;
        this.setComponentState({logData: "", pausedLogs: "", showClear: false, clearCountDown: 0, paused: window.currentPauseAction}, () => {
          this.clearPausedLogs();
        })
      }, 60000);
      //Count down to clear
      this.clearId = window.setInterval(() => {
        let newCountDown = this.state.clearCountDown - 1;
        document.title = this.globs.__(window.pageContent.ClearingLog, {sec: newCountDown })
        this.setComponentState({clearCountDown: newCountDown})
      }, 1000);
    }
    this.setComponentState(changes)
   }

   developerWarnings(data) {
     let searchErrors = ["WARNING: DATA RACE"]
     $.each(searchErrors,(k, v) => {
       if (data.indexOf(v) != -1) {
          alert(v);
       }
     });
   }

   componentDidMount() {
    if (this.state.Id == "") {
      alert("No Id");
      window.close();
    }

    var cb = (data) => {
      if (data.hasOwnProperty("Id")) {
        if (data.Id == this.state.Id) {
          if (window.appState.DeveloperMode) {
            this.developerWarnings(data.Data);
          }
          if (!this.state.paused) {
            this.setComponentState({logData: data.Data + this.state.pausedLogs + "\n" + this.state.logData }, () => this.clearPausedLogs())
          } else {
            this.setComponentState({paused: window.currentPauseAction, pausedLogs: data.Data + "\n" + this.state.pausedLogs}, () => {
              this.setComponentState({
                dataCachedLengthSize: window.global.functions.roughSizeOfObject(this.state.pausedLogs),
                dataCachedLength: window.global.functions.humanFileSize(window.global.functions.roughSizeOfObject(this.state.pausedLogs))
              })
            })
          }
        }
      }
    };
    this.logCallbackId = window.api.registerSocketCallback(cb, "LogData");
  }

  componentWillUnmount() {
    window.api.unRegisterSocketCallback(this.logCallbackId);
  }


  render() {
    return (
        <div>
          <h1><i className="icon md-traffic" aria-hidden="true"></i>{this.state.LongName}{(this.state.LongName != "") ? "(Device Id: " +  this.state.Id + ")": this.state.Id} Log</h1>
          <h3>{window.pageContent.MemUsage} {this.state.dataLength} {(this.state.showClear) ? <span style={{color: "red"}}>{this.globs.__(window.pageContent.LogClear2, {sec: this.state.clearCountDown})}</span>: null}. {this.state.paused && !this.state.pausedDueToScroll && this.state.dataCachedLengthSize > 0 ? <span style={{color: window.materialColors["green400"]}}>{this.globs.__(window.pageContent.AutoPausedNewData, {dataByte: this.state.dataCachedLength})}</span>: null} {this.state.pausedDueToScroll ? <span style={{color: window.materialColors["green400"]}}>{window.pageContent.AutoPaused} {this.state.dataCachedLengthSize > 0 ? <span>{this.globs.__(window.pageContent.AutoPausedNewData, {dataByte: this.state.dataCachedLength})}</span>: null}</span>: null}</h3>
          <div style={{marginBottom: 15}}>
            <RaisedButton
                label={window.appContent.ClearScreen}
                onTouchTap={() => {
                  this.setComponentState({logData: ""}, () => {
                    if (!this.state.showClear) {
                      if (this.state.paused) {
                        var logDiv = ReactDOM.findDOMNode(this.logRef);
                        if (logDiv != undefined) {
                          $(logDiv)[0].scrollTop = 0;
                        }
                        // set a flag that we should not run our pause event on this scroll
                        window.isAutoScrollTop = true;
                      } else {
                        window.isAutoScrollTop = false;
                      }
                      window.currentPauseAction = !this.state.paused;
                      this.setComponentState({pausedDueToScroll: false, paused: window.currentPauseAction}, () => {
                        window.setTimeout(() => {
                          // set back to false so next time we scroll we pause again
                          window.isAutoScrollTop = false;
                        }, 100)
                      })
                    } else {
                      alert(this.globs.__(window.pageContent.LogClear3, {sec: this.state.clearCountDown}))
                    }
                  })
                }}
                style={{marginRight: 15}}
             />
            <RaisedButton
                label={window.appContent.DlFull}
                onTouchTap={() => {
                window.api.post({
                  controller: "roomList",
                  action: "DownloadLog",
                  state: {DeviceId: this.state.Id, Name: this.state.Id + "log"},
                  leaveStateAlone: true,
                  callback: (vm) => {
                    if (vm.DeviceId == "") {
                        return;
                    }
                    window.globals.saveFile(vm.Data, vm.Name + ".log");
                  }
                });
              }}
                style={{marginRight: 15}}
                secondary={true}
            />
            <RaisedButton
                icon={this.state.paused ? <Play/>: <Pause/>}
                onTouchTap={() => {
                  if (!this.state.showClear) {
                    if (this.state.paused) {
                      var logDiv = ReactDOM.findDOMNode(this.logRef);
                      if (logDiv != undefined) {
                        $(logDiv)[0].scrollTop = 0;
                      }
                      // set a flag that we should not run our pause event on this scroll
                      window.isAutoScrollTop = true;
                    } else {
                      window.isAutoScrollTop = false;
                    }
                    window.currentPauseAction = !this.state.paused;
                    this.setComponentState({pausedDueToScroll: false, paused: window.currentPauseAction}, () => {
                      window.setTimeout(() => {
                        // set back to false so next time we scroll we pause again
                        window.isAutoScrollTop = false;
                      }, 100)
                    })
                  } else {
                    alert(this.globs.__(window.pageContent.LogClear3, {sec: this.state.clearCountDown}))
                  }
                }}

                secondary={true}
            />
          </div>
          <pre style={{maxHeight: 650, overflow: "auto"}} ref={(c) => this.logRef = c} onScroll={() => {
            if (!window.isAutoScrollTop) {
              window.currentPauseAction = true;
              this.setComponentState({
                pausedDueToScroll: window.currentPauseAction,
                paused: true,
                dataCachedLengthSize: window.global.functions.roughSizeOfObject(this.state.pausedLogs),
                dataCachedLength: window.global.functions.humanFileSize(window.global.functions.roughSizeOfObject(this.state.pausedLogs))
              });
            }
          }}>{this.state.logData}</pre>
        </div>
    );
    this.logRender();
  }
}
export default Logs;
