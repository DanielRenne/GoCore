import React, {Component} from 'react';
import BaseComponent from '../components/base';
import MuiThemes from '../pages/colors'
import Avatar from 'material-ui/Avatar';
import {ExitIcon} from '../icons/icons';
import IconButton from 'material-ui/IconButton';
import ContentBack from 'material-ui/svg-icons/hardware/keyboard-backspace';

class Banner extends BaseComponent {
  constructor(props, context) {
    super(props, context);
    this.state = window.appState.Banner;

    this.exitAccount = () => {
      window.globals.handleAccountSwitch(window.appState.UserPrimaryAccount);
    };

    window.goCore.setBannerStateFromExternal = (state) => {
      window.appState.Banner = state;
      this.setComponentState(window.appState.Banner);
    };

    window.goCore.setBannerCbFromExternal = (cb) => {
      this.base.setState(cb);
    };

    window.goCore.setCurrentHistoryIdx = (idx) => {
      if (idx != this.state.LatestHistoryLen) {
        this.setState({
          CurrentHistoryIdx: idx,
          LatestHistoryLen: idx
        });
      }
    };

    this.back = () => {
      if (window.history.length == 0 ) {
        return;
      }
      window.history.back();
    }
  }

  componentDidMount() {
    if (window.location.href.indexOf("roomControl") == -1) {
      $(".site-navbar").show();
    }
  }

  render() {
    this.logRender();

    var exitDoor = "none";
    var exitAccount = "";
    var companyPadding = 0;
    if (window.appState.Banner.IsSecondaryAccount) {
      exitDoor = "inherit";
      companyPadding = 50;

      exitAccount = (
        <li role="presentation">
          <a href='javascript:window.globals.handleAccountSwitch(window.appState.UserPrimaryAccount);' role="menuitem"><ExitIcon  color={"black"} width={15} height={15}/><span style={{paddingLeft:6}}>{window.appContent.NavbarAvatarExitAccount}</span></a>
        </li>
      );
    }

    var avatarBilling = "";

    // avatarBilling = (
    //   <li role="presentation">
    //     <a href="javascript:void(0)" role="menuitem"><i className="icon md-card" aria-hidden="true"></i>{window.appContent.NavbarAvatarBilling}</a>
    //   </li>
    // );

    return (
        <div>
        <nav className="site-navbar navbar navbar-inverse navbar-fixed-top navbar-mega hide" role="navigation" style={{backgroundColor:this.state.Color}}>

         {/* Mobile */}

          <div className="navbar-header" id="navbar-mobile">
            <button id="btnSideBarMenuCollapse" type="button" className="navbar-toggle hamburger hamburger-close navbar-toggle-left hided"
              data-toggle="menubar"><span className="sr-only">Toggle navigation</span>
              <span className="hamburger-bar"></span>
            </button>
            {window.history.length > 0 ? <IconButton style={{position:"absolute", marginTop:10}} onClick={this.back}>
              <ContentBack color={"white"} width={20} height={20}/>
            </IconButton>: null}
            <button  type="button" className="navbar-toggle collapsed" data-target="#site-navbar-collapse"
              data-toggle="collapse">
              <i className="icon md-more" aria-hidden="true"></i>
            </button>
            <div className="navbar-brand navbar-brand-center site-gridmenu-toggle" data-toggle="gridmenu">
              <a style={{width: 28, cursor: 'pointer'}} onClick={() => document.location = '/#/home'}>
                <img src="/web/app/images/go_core_company_icon.png"/>
              </a>
            </div>
            <button type="button" className="navbar-toggle collapsed" data-target="#site-navbar-search"
              data-toggle="collapse">
              <span className="sr-only">Toggle Search</span>
              <i className="icon md-search" aria-hidden="true"></i>
            </button>
          </div>


          <div className="navbar-container container-fluid" id="navbar-desktop">
            <div className="collapse navbar-collapse navbar-collapse-toolbar" id="site-navbar-collapse">
              <ul className="nav navbar-toolbar">
                <li className="hidden-float" id="toggleMenubar">
                  <a data-toggle="menubar" href="#" role="button">
                    <i className="icon hamburger hamburger-arrow-left">
                        <span className="sr-only">Toggle menubar</span>
                        <span className="hamburger-bar"></span>
                      </i>
                  </a>
                </li>
                {($(window).width() > 768) ? <li>
                  <IconButton style={{marginTop:10}} onClick={this.back}>
                    <ContentBack color={"white"} width={20} height={20}/>
                  </IconButton>
                </li>: null}
                <li className="hidden-xs" id="toggleFullscreen">
                  <a className="icon icon-fullscreen" data-toggle="fullscreen" href="#" role="button">
                    <span className="sr-only">Toggle fullscreen</span>
                  </a>
                </li>
                <li className="hidden-float">
                  <a className="icon md-search" data-toggle="collapse" href="#" data-target="#site-navbar-search"
                    role="button">
                    <span className="sr-only">Toggle Search</span>
                  </a>
                </li>
                <li className="dropdown dropdown-fw dropdown-mega">
                  <a className="dropdown-toggle" data-toggle="dropdown" href="#" aria-expanded="false"
                    data-animation="fade" role="button">Help <i className="icon md-chevron-down" aria-hidden="true"></i></a>
                  <ul className="dropdown-menu" role="menu">
                    <li role="presentation">
                      <div className="mega-content">
                        <div className="row">
                          <div className="col-sm-4">
                            <h5>Temporary Developer Widgets</h5>
                            <ul className="blocks-2">
                              <li className="mega-menu margin-0">
                                <ul className="list-icons">
                                  <li><i className="md-chevron-right" aria-hidden="true"></i>
                                    <a href="/#/widgets">Modals</a>
                                  </li>
                                  <li><i className="md-chevron-right" aria-hidden="true"></i>
                                    <a href="/#/widgets">Panels</a>
                                  </li>
                                  <li><i className="md-chevron-right" aria-hidden="true"></i>
                                    <a href="/#/widgets">Overlay</a>
                                  </li>
                                  <li><i className="md-chevron-right" aria-hidden="true"></i>
                                    <a href="/#/widgets">Tooltips</a>
                                  </li>
                                  <li><i className="md-chevron-right" aria-hidden="true"></i>
                                    <a href="/#/widgets">Scrollable</a>
                                  </li>
                                  <li><i className="md-chevron-right" aria-hidden="true"></i>
                                    <a href="/#/widgets">Typography</a>
                                  </li>
                                </ul>
                              </li>
                              <li className="mega-menu margin-0">
                                <ul className="list-icons">
                                  <li><i className="md-chevron-right" aria-hidden="true"></i>
                                    <a href="/#/widgets">Modals</a>
                                  </li>
                                  <li><i className="md-chevron-right" aria-hidden="true"></i>
                                    <a href="/#/widgets">Panels</a>
                                  </li>
                                  <li><i className="md-chevron-right" aria-hidden="true"></i>
                                    <a href="/#/widgets">Overlay</a>
                                  </li>
                                  <li><i className="md-chevron-right" aria-hidden="true"></i>
                                    <a href="/#/widgets">Tooltips</a>
                                  </li>
                                  <li><i className="md-chevron-right" aria-hidden="true"></i>
                                    <a href="/#/widgets">Scrollable</a>
                                  </li>
                                  <li><i className="md-chevron-right" aria-hidden="true"></i>
                                    <a href="/#/widgets">Typography</a>
                                  </li>
                                </ul>
                              </li>
                            </ul>
                          </div>
                          <div className="col-sm-4">
                            <h5>Media
                              <span className="badge badge-success">4</span>
                            </h5>
                            <ul className="blocks-3">
                              <li>
                                <a className="thumbnail margin-0" href="javascript:void(0)">
                                  <img className="width-full" src="/web/app/images/placeholder.png" alt="..."
                                  />
                                </a>
                              </li>
                              <li>
                                <a className="thumbnail margin-0" href="javascript:void(0)">
                                  <img className="width-full" src="/web/app/images/placeholder.png" alt="..."
                                  />
                                </a>
                              </li>
                              <li>
                                <a className="thumbnail margin-0" href="javascript:void(0)">
                                  <img className="width-full" src="/web/app/images/placeholder.png" alt="..."
                                  />
                                </a>
                              </li>
                              <li>
                                <a className="thumbnail margin-0" href="javascript:void(0)">
                                  <img className="width-full" src="/web/app/images/placeholder.png" alt="..."
                                  />
                                </a>
                              </li>
                              <li>
                                <a className="thumbnail margin-0" href="javascript:void(0)">
                                  <img className="width-full" src="/web/app/images/placeholder.png" alt="..."
                                  />
                                </a>
                              </li>
                              <li>
                                <a className="thumbnail margin-0" href="javascript:void(0)">
                                  <img className="width-full" src="/web/app/images/placeholder.png" alt="..."
                                  />
                                </a>
                              </li>
                            </ul>
                          </div>
                          <div className="col-sm-4">
                            <h5 className="margin-bottom-0">Accordion</h5>
                            {/* Accordion */}
                            <div className="panel-group panel-group-simple" id="siteMegaAccordion" aria-multiselectable="true"
                              role="tablist">
                              <div className="panel">
                                <div className="panel-heading" id="siteMegaAccordionHeadingOne" role="tab">
                                  <a className="panel-title" data-toggle="collapse" href="#siteMegaCollapseOne" data-parent="#siteMegaAccordion"
                                    aria-expanded="false" aria-controls="siteMegaCollapseOne">
                                      Collapsible Group Item #1
                                    </a>
                                </div>
                                <div className="panel-collapse collapse" id="siteMegaCollapseOne" aria-labelledby="siteMegaAccordionHeadingOne"
                                  role="tabpanel">
                                  <div className="panel-body">
                                    De moveat laudatur vestra parum doloribus labitur sentire partes, eripuit praesenti
                                    congressus ostendit alienae, voluptati ornateque
                                    accusamus clamat reperietur convicia albucius.
                                  </div>
                                </div>
                              </div>
                              <div className="panel">
                                <div className="panel-heading" id="siteMegaAccordionHeadingTwo" role="tab">
                                  <a className="panel-title collapsed" data-toggle="collapse" href="#siteMegaCollapseTwo"
                                    data-parent="#siteMegaAccordion" aria-expanded="false"
                                    aria-controls="siteMegaCollapseTwo">
                                      Collapsible Group Item #2
                                    </a>
                                </div>
                                <div className="panel-collapse collapse" id="siteMegaCollapseTwo" aria-labelledby="siteMegaAccordionHeadingTwo"
                                  role="tabpanel">
                                  <div className="panel-body">
                                    Praestabiliorem. Pellat excruciant legantur ullum leniter vacare foris voluptate
                                    loco ignavi, credo videretur multoque choro fatemur
                                    mortis animus adoptionem, bello statuat expediunt
                                    naturales.
                                  </div>
                                </div>
                              </div>

                              <div className="panel">
                                <div className="panel-heading" id="siteMegaAccordionHeadingThree" role="tab">
                                  <a className="panel-title collapsed" data-toggle="collapse" href="#siteMegaCollapseThree"
                                    data-parent="#siteMegaAccordion" aria-expanded="false"
                                    aria-controls="siteMegaCollapseThree">
                                      Collapsible Group Item #3
                                    </a>
                                </div>
                                <div className="panel-collapse collapse" id="siteMegaCollapseThree" aria-labelledby="siteMegaAccordionHeadingThree"
                                  role="tabpanel">
                                  <div className="panel-body">
                                    Horum, antiquitate perciperet d conspectum locus obruamus animumque perspici probabis
                                    suscipere. Desiderat magnum, contenta poena desiderant
                                    concederetur menandri damna disputandum corporum.
                                  </div>
                                </div>
                              </div>
                            </div>
                            {/* End Accordion */}
                          </div>
                        </div>
                      </div>
                    </li>
                  </ul>
                </li>
              </ul>
              {/* End Navbar Toolbar  */}

              {/* Navbar Toolbar Right  */}
              <ul className="nav navbar-toolbar navbar-right navbar-toolbar-right">

                <li className="dropdown" style={{marginBottom: 0}}>
                  <a data-toggle="dropdown" href="javascript:void(0)" title="Notifications" aria-expanded="false"
                    data-animation="scale-up" role="button">
                    <i className="icon md-notifications" aria-hidden="true"></i>
                    <span className="badge badge-danger up">5</span>
                  </a>
                  <ul className="dropdown-menu dropdown-menu-right dropdown-menu-media" role="menu">
                    <li className="dropdown-menu-header" role="presentation">
                      <h5>NOTIFICATIONS</h5>
                      <span className="label label-round label-danger">New 5</span>
                    </li>

                    <li className="list-group" role="presentation">
                      <div data-role="container">
                        <div data-role="content">
                          <a className="list-group-item" href="javascript:void(0)" role="menuitem">
                            <div className="media">
                              <div className="media-left padding-right-10">
                                <i className="icon md-receipt bg-red-600 white icon-circle" aria-hidden="true"></i>
                              </div>
                              <div className="media-body">
                                <h6 className="media-heading">A new order has been placed</h6>
                                <time className="media-meta">5 hours ago</time>
                              </div>
                            </div>
                          </a>
                          <a className="list-group-item" href="javascript:void(0)" role="menuitem">
                            <div className="media">
                              <div className="media-left padding-right-10">
                                <i className="icon md-account bg-green-600 white icon-circle" aria-hidden="true"></i>
                              </div>
                              <div className="media-body">
                                <h6 className="media-heading">Completed the task</h6>
                                <time className="media-meta">2 days ago</time>
                              </div>
                            </div>
                          </a>
                          <a className="list-group-item" href="javascript:void(0)" role="menuitem">
                            <div className="media">
                              <div className="media-left padding-right-10">
                                <i className="icon md-settings bg-red-600 white icon-circle" aria-hidden="true"></i>
                              </div>
                              <div className="media-body">
                                <h6 className="media-heading">Settings updated</h6>
                                <time className="media-meta">2 days ago</time>
                              </div>
                            </div>
                          </a>
                          <a className="list-group-item" href="javascript:void(0)" role="menuitem">
                            <div className="media">
                              <div className="media-left padding-right-10">
                                <i className="icon md-calendar bg-blue-600 white icon-circle" aria-hidden="true"></i>
                              </div>
                              <div className="media-body">
                                <h6 className="media-heading">Event started</h6>
                                <time className="media-meta">3 days ago</time>
                              </div>
                            </div>
                          </a>
                          <a className="list-group-item" href="javascript:void(0)" role="menuitem">
                            <div className="media">
                              <div className="media-left padding-right-10">
                                <i className="icon md-comment bg-orange-600 white icon-circle" aria-hidden="true"></i>
                              </div>
                              <div className="media-body">
                                <h6 className="media-heading">Message received</h6>
                                <time className="media-meta">3 days ago</time>
                              </div>
                            </div>
                          </a>
                        </div>
                      </div>
                    </li>

                    <li className="dropdown-menu-footer" role="presentation">
                      <a className="dropdown-menu-footer-btn" href="javascript:void(0)" role="button">
                        <i className="icon md-settings" aria-hidden="true"></i>
                      </a>
                      <a href="javascript:void(0)" role="menuitem">
                          See all notification
                        </a>
                    </li>
                  </ul>
                </li>
                <li className="dropdown" style={{marginBottom: 0}}>
                  <a data-toggle="dropdown" href="javascript:void(0)" title="Messages" aria-expanded="false"
                    data-animation="scale-up" role="button">
                    <i className="icon md-email" aria-hidden="true"></i>
                    <span className="badge badge-info up">3</span>
                  </a>
                  <ul className="dropdown-menu dropdown-menu-right dropdown-menu-media" role="menu">
                    <li className="dropdown-menu-header" role="presentation">
                      <h5>MESSAGES</h5>
                      <span className="label label-round label-info">New 3</span>
                    </li>

                    <li className="list-group" role="presentation">
                      <div data-role="container">
                        <div data-role="content">
                          <a className="list-group-item" href="javascript:void(0)" role="menuitem">
                            <div className="media">
                              <div className="media-left padding-right-10">
                                <span className="avatar avatar-sm avatar-online">
                                  <img src="" alt="..." />
                                  <i></i>
                                </span>
                              </div>
                              <div className="media-body">
                                <h6 className="media-heading">Mary Adams</h6>
                                <div className="media-meta">
                                  <time>30 minutes ago</time>
                                </div>
                                <div className="media-detail">Anyways, i would like just do it</div>
                              </div>
                            </div>
                          </a>
                          <a className="list-group-item" href="javascript:void(0)" role="menuitem">
                            <div className="media">
                              <div className="media-left padding-right-10">
                                <span className="avatar avatar-sm avatar-off">
                                  <img src="" alt="..." />
                                  <i></i>
                                </span>
                              </div>
                              <div className="media-body">
                                <h6 className="media-heading">Caleb Richards</h6>
                                <div className="media-meta">
                                  <time>12 hours ago</time>
                                </div>
                                <div className="media-detail">I checheck the document. But there seems</div>
                              </div>
                            </div>
                          </a>
                          <a className="list-group-item" href="javascript:void(0)" role="menuitem">
                            <div className="media">
                              <div className="media-left padding-right-10">
                                <span className="avatar avatar-sm avatar-busy">
                                  <img src="" alt="..." />
                                  <i></i>
                                </span>
                              </div>
                              <div className="media-body">
                                <h6 className="media-heading">June Lane</h6>
                                <div className="media-meta">
                                  <time>2 days ago</time>
                                </div>
                                <div className="media-detail">Lorem ipsum Id consectetur et minim</div>
                              </div>
                            </div>
                          </a>
                          <a className="list-group-item" href="javascript:void(0)" role="menuitem">
                            <div className="media">
                              <div className="media-left padding-right-10">
                                <span className="avatar avatar-sm avatar-away">
                                  <img src="" alt="..." />
                                  <i></i>
                                </span>
                              </div>
                              <div className="media-body">
                                <h6 className="media-heading">Edward Fletcher</h6>
                                <div className="media-meta">
                                  <time>3 days ago</time>
                                </div>
                                <div className="media-detail">Dolor et irure cupidatat commodo nostrud nostrud.</div>
                              </div>
                            </div>
                          </a>
                        </div>
                      </div>
                    </li>
                    <li className="dropdown-menu-footer" role="presentation">
                      <a className="dropdown-menu-footer-btn" href="javascript:void(0)" role="button">
                        <i className="icon md-settings" aria-hidden="true"></i>
                      </a>
                      <a href="javascript:void(0)" role="menuitem">
                          See all messages
                        </a>
                    </li>
                  </ul>
                </li>

                <li className="dropdown" style={{marginBottom: 0}}>
                      <a className="navbar-avatar dropdown-toggle" data-toggle="dropdown" href="#" aria-expanded="false"
                        data-animation="scale-up" role="button">
                        <Avatar color={MuiThemes.default.palette.accent2Color} size={30}>
                          {window.appState.UserInitials}
                        </Avatar>
                      </a>
                      <ul className="dropdown-menu" role="menu">
                        <li role="presentation">
                          <a href="/#/userProfile" role="menuitem"><i className="icon md-account" aria-hidden="true"></i>{window.appContent.NavbarAvatarProfile}</a></li>

                        {avatarBilling}
                        {(window.appState.AccountTypeShort == "cust" && this.globs.HasRole("USER_VIEW")) ?
                        <li role="presentation">
                          <a href="/#/userList" role="menuitem"><i className="icon md-settings" aria-hidden="true"></i>{window.appContent.NavbarAvatarUserList}</a>
                        </li>:
                            (this.globs.HasRole("USER_VIEW") || this.globs.HasRole("ACCOUNT_VIEW")) ?
                                <li role="presentation">
                                  <a href="/#/accountList" role="menuitem"><i className="icon md-settings" aria-hidden="true"></i>{window.appContent.NavbarAvatarAccountSettings}</a>
                                </li>
                                : null
                        }
                        {((this.globs.HasRole("MAINTENANCE_MODIFY") || this.globs.HasRole("SERVER_SETTING_MODIFY"))) ? <li role="presentation">
                          <a href="/#/serverSettingsModify" role="menuitem"><i className="icon md-storage" aria-hidden="true"></i>{window.appContent.NavbarAvatarServerSettings}</a>
                        </li>: null}
                        {(this.globs.HasRole("LOGS_VIEW")) ?
                        <li role="presentation">
                          <a href="javascript:" onClick={() => this.globs.ViewLog("app")} role="menuitem"><i className="icon md-traffic" aria-hidden="true"></i>{window.appContent.ViewLogs}</a>
                        </li>: null}
                        <li className="divider" role="presentation"></li>
                        {exitAccount}
                        <li role="presentation">
                          <a href='javascript:window.api.post({action: "Logout", state: {}, controller:"login"});' role="menuitem"><i className="icon md-power" aria-hidden="true"></i>{window.appContent.NavbarAvatarLogout}</a>
                        </li>
                      </ul>


                </li>
              </ul>
              {/* End Navbar Toolbar Right  */}

              <div className="navbar-brand navbar-brand-center">
                <a onClick={() => document.location = '/#/home'} style={{cursor: 'pointer'}}>
                  <img style={{width:200}} src="/web/app/images/go_core_app_product_white.png"/>
                </a>
              </div>
              {($(window).width() > 768) ?
              <div className="navbar-brand navbar-right">

                <span className="Aligner" style={{marginTop:15, marginLeft:companyPadding}}>
                  <div  onClick={() => document.location = '/#/home'}
                        className="navbar-brand-app-title Aligner-item--bottom"
                        style={{minWidth:30, minHeight:22, cursor: 'pointer'}}>
                    {((window.appState.AccountTypeShort == "cust") ? window.appContent.CustomerColon : window.appContent.AccountColon) + this.state.AccountName}
                  </div>
                  <IconButton className="Aligner-item--bottom"
                              title={window.appContent.ExitAccount}
                              style={{paddingTop:22, paddingLeft:10, display:exitDoor}}
                              onTouchTap={this.exitAccount}>
                    <ExitIcon  color={"white"} width={20} height={20}/>
                  </IconButton>
                </span>

              </div>: null}
            </div>
            {/* End Navbar Collapse  */}

            {/* Site Navbar Seach  */}
            <div className="collapse navbar-search-overlap" id="site-navbar-search">
              <form role="search">
                <div className="form-group">
                  <div className="input-search">
                    <i className="input-search-icon md-search" aria-hidden="true"></i>
                    <input type="text" className="form-control" name="site-search" placeholder="Search..."/>
                    <button type="button" className="input-search-close icon md-close" data-target="#site-navbar-search"
                      data-toggle="collapse" aria-label="Close"></button>
                  </div>
                </div>
              </form>
            </div>
            {/* End Site Navbar Seach  */}
          </div>
        </nav>



        </div>
    );
  }
}

Banner.propTypes = {
  page: React.PropTypes.element,
};

// AppBanner.defaultProps = {
//     showMenuPageTitle: true
// };

export default Banner;
