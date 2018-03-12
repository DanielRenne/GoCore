import React, {Component} from "react";
import BaseComponent from "../components/base";
import SideBarMenuItem from "./sidebarMenuItem";
import MuiThemes from "../pages/colors";
import Avatar from "material-ui/Avatar";


class SideBarMenu extends BaseComponent {
  constructor(props, context) {
    super(props, context);

    this.state = window.appState.SideBarMenu;
  }

  render() {
    this.logRender();
    if (window.location.hash.indexOf("#/roomControl") != -1) {
      return null;
    }

    if (this.state.Items == null && window.appState.DeveloperMode) {
      alert("Dude, you haven't implemented the sidebar viewModel!\n\n#devfail");
    }

    var sidebarMapping = (this.state.Items == null) ? console.error("No sidebar items passed in viewModel!") :
      (window.appState.AccountTypeShort != "cust" ?
        (this.state.Items.map((item, k) => {
              if (item.Title != "SitesTitle") {
                return <SideBarMenuItem key={k}
                                        title={item.Title}
                                        url={item.URL}
                                        expanded={item.Expanded}
                                        items={item.Items}
                                        icon={item.Icon}
                                        rightIcon={item.RightIcon}
                                        rightIconFunc={item.RightIconFunc}
                                        hidden={item.Hidden}
                                        depth={0}
                                        rightIconLink={item.RightIconLink}
                                        selected={item.Selected}/>
              }
            })) : (this.state.Items.map((item, k) => {
                return <SideBarMenuItem key={k}
                                        title={item.Title}
                                        url={item.URL}
                                        expanded={item.Expanded}
                                        items={item.Items}
                                        icon={item.Icon}
                                        rightIcon={item.RightIcon}
                                        rightIconFunc={item.RightIconFunc}
                                        hidden={item.Hidden}
                                        depth={0}
                                        rightIconLink={item.RightIconLink}
                                        selected={item.Selected}/>
              }))
            )

    return (
      <div className="site-menubar" style={{overflowY:"auto"}}>
        <div className="site-menubar-header">
          <div className="cover overlay">
            <img className="cover-image" src="/web/app/images/dashboard-header.jpg"
              alt="..."/>
            <div className="overlay-panel vertical-align overlay-background">
              <div className="vertical-align-middle">
                <div className="row">
                  <div className="col-md-3" style={{marginLeft: 5}}>
                    <Avatar backgroundColor={MuiThemes.default.palette.accent1Color} size={50}>
                      {window.appState.UserInitials}
                    </Avatar>
                  </div>
                  <div className="col-md-8">
                    <div className="site-menubar-info">
                      <h5 className="site-menubar-user">{window.appState.UserFirst} {window.appState.UserLast}</h5>
                      <p className="site-menubar-email">{window.appState.UserEmail}</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      <div className="site-menubar-body">
        <div>
          <div>
            <ul className="site-menu">
            {
              sidebarMapping
            }
            </ul>
          </div>
        </div>
      </div>
  </div>
    );
  }
}

SideBarMenu.propTypes = {
  items: React.PropTypes.array
};

SideBarMenu.defaultProps = {
    items: []
};

export default SideBarMenu;
