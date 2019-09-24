import React, {Component} from "react";
import SocialPerson from "material-ui/svg-icons/social/person";
import BasePageComponent from "../../components/basePageComponent";
import InfoPopup from "../../components/infoPopup";
import MuiThemes from "../colors";
// import Slider from "react-slick";
import {HDMIIcon, NotesIcon, UserIcon, DocumentationIcon} from "../../globals/icons";
import {Avatar, IconButton, RaisedButton} from "../../globals/forms";

class Home extends BasePageComponent {
  constructor(props, context) {
    super(props, context);
    this.state.DevToolsShow = false;
    this.state.Width = $(window).width();
    this.state.ReleaseNotesOpen = false;

    this.resizeEvent = (e) => this.handleResize(e);
    this.createComponentEvents();
  }

  componentWillReceiveProps(nextProps) {
    if (window.appState.DeveloperMode) {
      this.createComponentEvents();
    }
    return true;
  }

  createComponentEvents() {

  }

  handleResize(e) {
    var changes = {};
    changes.Width = $(window).width();
    this.setComponentState(changes);
  }

  componentDidMount() {
    $('.page-content').css({paddingTop: 0});
    this.handleResize();
    window.addEventListener('resize', this.resizeEvent);
  }

  componentWillUnmount() {
    $('.page-content').css({paddingTop: 30});
    window.removeEventListener('resize', this.resizeEvent);
  }

  render() {
    this.logRender();
    let openReleaseNotes = () => this.setComponentState({ReleaseNotesOpen:true});
    let openDocumentation = () => {
      window.appState.Banner.HelpOpen = !window.appState.Banner.HelpOpen;
      window.goCore.setBannerStateFromExternal(window.appState.Banner);
    };
    let goCoreAppVersion = null;
    let accountUserRoomListItems = [];
    let accountUserRoomDesc = null;
    let hasAccountData = (this.state.Accounts != null && this.state.Accounts.length > 0);
    let hasUserData = (this.state.Users != null && this.state.Users.length > 0);

    if (hasAccountData) {
      accountUserRoomDesc = window.pageContent.Accounts;
    } else if (hasUserData) {
      accountUserRoomDesc = window.pageContent.Users;
    } else {
      accountUserRoomDesc = window.pageContent.TY;
    }

    let canSeeTechLink = this.state.RoomDeviceDownId != "" && this.globs.HasRole("TECHNOLOGY_VIEW");
    let canClickEquipment = this.globs.HasRole("EQUIPMENT_VIEW");
    let canClickBuildingsRooms = this.globs.HasRole("FEATURE_SITE_BUILDING_FLOOR_ROOM_VIEW");
    let canClickUsers = this.globs.HasRole("USER_VIEW");
    let appDocumentation = <div className="col-lg-4 col-md-12 col-xm-12">
        <h3 style={{marginTop: "1.5em", color:  window.materialColors["blueGrey700"]}}>{window.pageContent.ProdDoc}</h3>
        <div className="widget widget-shadow" style={{minHeight: 500}}>
          <div className="widget-header cover overlay" style={{cursor: "pointer"}} onClick={openDocumentation}>
            <div className="cover-background height-200" style={{border: "2px solid", backgroundImage: "url('/web/app/images/go_core_app_product_home.png')"}}></div>
          </div>
          <div className="widget-body" style={{height:"calc(100% - 250px)", paddingRight: 0, paddingLeft: 0, paddingTop: 20, paddingBottom: 0}}>
            <div className="margin-bottom-10" style={{marginTop: -70, cursor: "pointer"}} onClick={openDocumentation}>
              <a className="avatar avatar-100 bg-white img-bordered">
                <span>
                  <Avatar
                      icon={<DocumentationIcon />}
                      color={MuiThemes.default.palette.accent1Color}
                      backgroundColor={"white"}
                      size={90}
                  />
                </span>
              </a>
            </div>
            <div className="margin-bottom-20">
              <div className="font-size-20" style={{color:  window.materialColors["blueGrey700"]}}>{window.pageContent.ViewDocs}</div>
              <div className="font-size-14 grey-500">
                <span><a href="javascript:" style={{cursor: "pointer"}} onClick={() => {
                  window.open("/web/app/manual/???.pdf");
                }}>({window.pageContent.ViewDocumentation})</a></span>
              </div>
            </div>
            <div style={{textAlign: "left", maxHeight: 169, overflowY: "scroll", borderTop: "1px solid " + window.materialColors["grey200"]}}>
            </div>
          </div>
        </div>
      </div>;

    if (this.state.ReleaseDescriptionLines != null && this.state.ReleaseDescriptionLines.length > 0) {
      let description = null;

      if (this.state.ReleaseDescriptionLines.length == 1) {
        description = <p>
          {this.state.ReleaseDescriptionLines[0]}
        </p>;
      } else {
        description = <ul>
          {this.globs.map(this.state.ReleaseDescriptionLines, (v, k) => {
              if (k <= 3) {
                return <li key={k + "-li"}>{v}</li>
              }
            }
          )}
        </ul>;
      }
      goCoreAppVersion = <div className="col-lg-4 col-md-12 col-xm-12">
        <h3 style={{marginTop: "1.5em", color:  window.materialColors["blueGrey700"]}}>{window.pageContent.Updates}</h3>
        <div className="widget widget-shadow" style={{minHeight: 500}}>
          <div className="widget-header cover overlay"  style={{cursor: "pointer"}} onClick={openReleaseNotes}>
            <div className="cover-background height-200" style={{border: "2px solid",backgroundImage: "url('/web/app/images/go_core_company_home_hd.jpg')"}}></div>
          </div>
          <div className="widget-body padding-horizontal-30 padding-vertical-20" style={{height:"calc(100% - 250px)"}}>
            <div className="margin-bottom-10" style={{marginTop: -70, cursor: "pointer"}} onClick={openReleaseNotes}>
              <span>
                  <a className="avatar avatar-100 bg-white img-bordered">
                    <Avatar
                        icon={<NotesIcon/>}
                        color={MuiThemes.default.palette.accent1Color}
                        backgroundColor={"white"}
                        size={90}
                      />
                  </a>
              </span>
            </div>
            <div className="margin-bottom-20">
              <div className="font-size-20"  style={{color:  window.materialColors["blueGrey700"]}}>{window.pageContent.AppVersion}</div>
              <div className="font-size-14 grey-500">
                <span>{window.appState.displayVersion} <a href="javascript:" onClick={openReleaseNotes}>({window.pageContent.AppVersionNotes})</a></span>
              </div>
            </div>
            <span style={{textAlign: "left"}}>
              {description}
            </span>
          </div>
        </div>
        <InfoPopup open={this.state.ReleaseNotesOpen} parentStateKey="ReleaseNotesOpen" parent={this} onClose={() => {
            this.setComponentState({ReleaseNotesOpen:false})
          }}>
          <textarea style={{width:700,height:500}}>
            {this.state.ReleaseNotes}
          </textarea>
        </InfoPopup>
      </div>;
    }

    if (hasAccountData) {
        this.state.Accounts.forEach((item, k) => {
          let click = () => {
            return this.globs.HasRole("ACCOUNT_VIEW") ? window.api.get({action: "Root",uriParams: {}, controller: "accountList"}): null;
          }
          accountUserRoomListItems.push(
            <li key={k+"-accountli"} style={{paddingTop: 10}}  className="list-group-item">
              <div className="media">
                <div className="media-left">
                  <a className="avatar avatar-lg" href="javascript:" onClick={click}>
                    <Avatar backgroundColor={MuiThemes.default.palette.accent1Color} size={40}>
                      {item.AccountName.substr(0,1).toUpperCase()}
                    </Avatar>
                  </a>
                </div>
                <div className="media-body" style={{width: 237}}>
                  <h4 className="media-heading"><a href="javascript:" onClick={click}>{item.AccountName}</a></h4>
                  <small>{click ? <a href="javascript:" onClick={click}>{item.City + ", " + item.Region}</a>: <span>{item.City + ", " + item.Region}</span>}</small>
                </div>
                <div className="media-body" style={{width: 100}}>
                  <button type="button" style={{position: "relative"}} className="btn btn-primary" onClick={() => window.globals.handleAccountSwitch(item.Id)}>{window.pageContent.Access}</button>
                </div>
              </div>
            </li>);
        });
    } else if (hasUserData) {
      let click = () => {
        this.globs.HasRole("USER_VIEW") ? window.api.get({
          action: "Root", uriParams: {}, controller: "userList"
        }) : null;
      }
      this.state.Users.forEach((item, k) => {
        accountUserRoomListItems.push(
          <li  key={k+"-user2li"} style={{paddingTop: 10}}  className="list-group-item">
            <div className="media">
              <div className="media-left">
                <a className="avatar avatar-lg" href="javascript:" onClick={click}>
                  <Avatar backgroundColor={MuiThemes.default.palette.accent1Color} size={40}>
                    {item.Joins.User.First.substr(0,1).toUpperCase() + " " + item.Joins.User.Last.substr(0,1).toUpperCase()}
                  </Avatar>
                </a>
              </div>
              <div className="media-body">
                <h4 className="media-heading"><a href="javascript:" onClick={click}>{item.Joins.User.First + " " + item.Joins.User.Last}</a></h4>
                <small>{click ? <a href="javascript:" onClick={click}>{item.Joins.Account.AccountName}</a>: <span>{item.Joins.Account.AccountName}</span>}</small>
              </div>
            </div>
          </li>);
      });
    } else {
      accountUserRoomListItems = [<li key="empty" className="list-group-item">{window.pageContent.NoRecent}</li>];
    }

    let accountsUsers = <div className="col-lg-4 col-md-12 col-xm-12">
                <h3 style={{marginTop: "1.5em", color:  window.materialColors["blueGrey700"]}}>{accountUserRoomDesc}</h3>
                <div className="widget" id="widgetUserList" style={{backgroundColor: "white", minHeight: 500}}>
                  <div className="widget-header cover overlay">
                    <img className="cover-image height-200" src="/web/app/images/dashboard-header.jpg" alt="..." />
                    <div className="overlay-panel vertical-align overlay-background">
                      <div className="vertical-align-middle">
                        <a className="avatar avatar-100 pull-left margin-right-20" href="javascript:void(0)">
                          <Avatar backgroundColor={MuiThemes.default.palette.accent1Color} size={75}>
                            {window.appState.UserInitials}
                          </Avatar>
                        </a>

                        <div className="pull-left">
                          <div className="font-size-20">{window.pageContent.Welcome2}</div>
                          <div className="font-size-20">{window.appState.UserFirst} {window.appState.UserLast}</div>
                          <p className="margin-bottom-20 text-nowrap">
                            <span className="text-break">{window.appState.UserEmail}</span>
                          </p>
                        </div>
                      </div>
                    </div>
                  </div>

                  <div className="widget-body padding-horizontal-5 padding-vertical-20" style={{height:"calc(100% - 250px)"}}>
                    <div style={{marginTop: -70}}>
                      <a className="avatar avatar-100 bg-white img-bordered">
                        <Avatar
                            icon={<UserIcon/>}
                            color={MuiThemes.default.palette.accent1Color}
                            backgroundColor={"white"}
                            size={95}
                          />
                      </a>
                    </div>


                  <div className="widget-content padding-horizontal-20">
                    <ul className="list-group list-group-full list-group-dividered">
                      {accountUserRoomListItems}
                    </ul>
                  </div>
                  </div>
                </div>
              </div>;

    return (
        <div style={{textAlign: 'center'}}>
          <input type="hidden" value={this.state.Width}/>
          <div className="page">
            <div className="page-content container-fluid">
              {/*First Row*/}
              {accountsUsers}
              {appDocumentation}
              {goCoreAppVersion}
            </div>

            {/*End users see floors and rooms*/}
            {window.appState.AccountRoleId == "57c083a3dcba0f7a0be33eb1" || window.appState.AccountRoleId  == "57fbb8bc9f566eaa6f5e4e8b" ? this.globs.getFloors(this): null}
          </div>

          {(window.appState.DeveloperMode) ?
              <div className="page">
                <div className="page-content container-fluid">
                  <h3>Dev Tools</h3>
                  <div>
                    <div className="col-sm-6">
                      <div className="widget">
                        <div className="widget-content padding-30 bg-green-600">
                          <div className="widget-watermark darker font-size-60 margin-15"><i className="icon md-input-antenna" aria-hidden="true"></i></div>
                          <div className="counter counter-md counter-inverse text-left">
                            <div className="counter-number-group">
                              <span className="counter-number">{this.state.WebSocketConnectionsCount}</span>
                            </div>
                            <div className="counter-label text-capitalize">Web Sockets Connected</div>
                          </div>
                        </div>
                      </div>
                    </div>
                    <div className="col-sm-6" style={{cursor: "pointer"}} onClick={() => {
                      this.setComponentState({DevToolsShow: !this.state.DevToolsShow})
                    }}>
                      <div className="widget">
                        <div className="widget-content padding-30 bg-purple-600">
                          <div className="widget-watermark lighter font-size-60 margin-15"><i className="icon md-download" aria-hidden="true"></i></div>
                          <div className="counter counter-md counter-inverse text-left">
                            <div className="counter-number-wrap font-size-30">
                              <span className="counter-number">Open Sesame</span>
                            </div>
                            <div className="counter-label text-capitalize">Developer Tools (Click Here)</div>
                          </div>
                        </div>
                      </div>
                    </div>
                </div>
              </div>

              <div style={{display: (this.state.DevToolsShow) ? "block": "none"}}>
                <h4>Server Stuff</h4>
                <ul className="list-icons">
                  <li><i className="md-chevron-right" aria-hidden="true"/>
                    <a href="/#/appErrorList">View All Errors</a>
                  </li>
                  <li><i className="md-chevron-right" aria-hidden="true"/>
                    <a href="/#/transactionList">Transaction List</a>
                  </li>
                  <li><i className="md-chevron-right" aria-hidden="true"/>
                    <a href="/#/fileObjectList">File Object List</a>
                  </li>
                  <li><i className="md-chevron-right" aria-hidden="true"/>
                    <a href="javascript:" onClick={() => window.api.post({action: "ShutDownServer", state: {}, controller:"login"})}>os.Extt(0) graceful shutdown</a>
                  </li>
                </ul>

                <h4>Role Management</h4>
                <ul className="list-icons">
                  <li><i className="md-chevron-right" aria-hidden="true"/>
                    <a href="/#/roleFeatureList">Role Feature (Mapping)</a>
                  </li>
                  <li><i className="md-chevron-right" aria-hidden="true"/>
                    <a href="/#/featureList">Application Features</a>
                  </li>
                  <li><i className="md-chevron-right" aria-hidden="true"/>
                    <a href="/#/featureGroupList">Feature Pages Or Groups</a>
                  </li>
                  <li><i className="md-chevron-right" aria-hidden="true"/>
                    <a onClick={() => {
                      if (window.prompt("Type 'UNDERSTAND' if you know what this is going to do") == "UNDERSTAND") {
                        document.location = '#/roleFeatureList?action=Root&uriParams=eyJEdW1wQWxsRmVhdHVyZXNUb0FsbFJvbGVzIjp0cnVlfQ%3D%3D';
                      }
                    }} >Map all roles to all features</a>
                  </li>
                </ul>

              </div>
            </div>: null}
        </div>
    );
  }
}

export default Home;
