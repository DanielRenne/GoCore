import React from 'react';
import ReactDom from 'react-dom';
import Loadable from 'react-loading-overlay';
import {BaseComponent, CenterAnything} from "../globals/forms";

class Loader extends BaseComponent {

  constructor(props, context) {
    super(props, context);
    this.state = {};
    this.state.DisplayHeight = window.innerHeight;
    this.state.DisplayWidth = window.innerWidth;
    this.state.loading = false;
    this.state.text = window.appContent.LoadingPage;
    window.goCore.setLoaderFromExternal = (state, cb) => {
      if (this.state.loading == false && state.loading == false) {
        return;
      }
      if (this.state.loading == true && state.loading == true) {
        return;
      }
      if (state.loading == false) {
        // set back to default
        state.text = window.appContent.LoadingPage;
      }
      this.setComponentState(state, () => {
        if (typeof(cb) == "function") {
          cb();
        }
      });
    };
    this.resizeEvent = (e) => this.handleResize(e);
  }

  handleResize(e) {
    // e == null is a case where screenchange happens
    if (e == null || window.innerWidth != this.state.DisplayWidth && window.innerHeight != this.state.DisplayHeight) {
      this.setComponentState({
        DisplayHeight: window.innerHeight,
        DisplayWidth: window.innerWidth,
      });
    }
  }

  componentDidMount() {
    window.addEventListener('resize', this.resizeEvent);
  }

  componentWillUnmount() {
    window.removeEventListener('resize', this.resizeEvent);
  }

  render() {
    let img = <img style={{width:this.state.DisplayWidth * .5}} src={"/web/app/images/go_core_app_product.png"}/>;
    return <span>
     <Loadable
        zIndex={999999}
        active={this.state.loading}
        spinner={true}
        text={this.state.text}
        >
       <span style={{display:!this.state.loading ? "inline": "none"}}>{this.props.children}</span>
       {this.state.loading ? <CenterAnything>{img}</CenterAnything> : null}
     </Loadable>
      {this.state.loading ? <span style={{position: 'fixed', top: 100, left: 15, zIndex: 9999999}}><a href="javascript:" onClick={() => {
        this.setComponentState({loading: false})
      }}><h2 style={{color: "white"}}>{window.appContent.CloseAjax}</h2></a></span>: null}
    </span>

  }
}
export default Loader;
