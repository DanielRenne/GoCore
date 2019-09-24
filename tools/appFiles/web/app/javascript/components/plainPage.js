import React, {Component} from "react";
import BaseComponent from "./base";
import GoCoreParentNode from './goCoreParentNode';
import MuiThemeProvider from "material-ui/styles/MuiThemeProvider";
import MuiThemes from "../pages/colors";

class PlainPage extends BaseComponent {
  constructor(props, context) {
    super(props, context);
  }

  render() {
    try {
      this.logRender();
      var title = (this.props.headerElement == undefined) ?
                <h1 style={{margin: '0em 0 .5em'}}><MuiThemeProvider muiTheme={MuiThemes.opposite}>{this.props.icon}</MuiThemeProvider> {this.props.title}</h1>
                :
                <span className="AlignerRight">
                  <h1 style={{margin: '0em 0 .5em'}}><MuiThemeProvider muiTheme={MuiThemes.opposite}>{this.props.icon}</MuiThemeProvider> {this.props.title}</h1>
                  <span>{this.props.headerElement}</span>
                </span>;

      var breadCrumbs = "";
      if (this.props.breadCrumbs != undefined && this.props.breadCrumbs.length > 0){

        breadCrumbs = (
          <ol className="breadcrumb breadcrumb-arrow" style={{marginLeft:30, display:"inline"}}>
            {
              this.globs.map(this.props.breadCrumbs, (crumb, k) => {
                var props = {};
                if (typeof(crumb.Link) == 'function') {
                  props.style = {cursor: 'pointer'};
                  props.onClick = crumb.Link;
                } else {
                  props.href = crumb.Link;
                }
                if (this.props.breadCrumbs.length != 1 && this.props.breadCrumbs[this.props.breadCrumbs.length-1] == crumb) {
                  return <li key={k}>{crumb.Title}</li>
                } else {
                  return <li key={k}><a {...props}>{crumb.Title}</a></li>
                }
              })
            }
          </ol>
        );

        title = <div><span className="h1"><MuiThemeProvider muiTheme={MuiThemes.opposite}>{this.props.icon}</MuiThemeProvider> {this.props.title} </span>
                  {
                    (this.props.headerElement == undefined) ?
                    <span style={{fontSize:22}}>{breadCrumbs}</span>
                    :
                    <span className="Aligner">
                      <span style={{fontSize:22}}>{breadCrumbs}</span>
                      <span>{this.props.headerElement}</span>
                    </span>
                  }
                </div>
        ;
      }

      return (
          <GoCoreParentNode>
            {title}
            <div>
               <hr style={{marginBottom: 15, marginTop: 8}}/>
               <MuiThemeProvider muiTheme={MuiThemes.default}>
                 {this.props.page}
               </MuiThemeProvider>
            </div>
          </GoCoreParentNode>
      );
    } catch(e) {
      return this.globs.ComponentError("PlainPage", e.message, e);
    }
  }
}

PlainPage.propTypes = {
  page: React.PropTypes.element,
  title: React.PropTypes.string,
  showMenuPageTitle: React.PropTypes.bool,
  icon: React.PropTypes.element,
  breadCrumbs: React.PropTypes.array,
  headerElement: React.PropTypes.element
};

PlainPage.defaultProps = {
    showMenuPageTitle: true
};

export default PlainPage;
