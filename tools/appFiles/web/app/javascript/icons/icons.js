import React, {Component} from 'react';
import BaseComponent from '../components/base';

class RoomIcon extends BaseComponent {
  constructor(props, context) {
    super(props, context);
  }

  render() {
    this.logRender();
    return (
      <svg style={{width:this.props.width, height:this.props.height}} viewBox="0 0 512 512" aria-labelledby="title">
          <title id="title">Room Icon</title>
            <path fill={this.props.color} d="M208,104.016C208,81.922,225.906,64,248,64s40,17.922,40,40.016C288,126.094,270.094,144,248,144   S208,126.094,208,104.016z M40,256c22.094,0,40-17.906,40-39.984C80,193.922,62.094,176,40,176S0,193.922,0,216.016   C0,238.094,17.906,256,40,256z M456,256c22.094,0,40-17.906,40-39.984C496,193.922,478.094,176,456,176s-40,17.922-40,40.016   C416,238.094,433.906,256,456,256z M496,295.266V386c0,7.734-6.266,14-14,14h-38.359c-6.75,0-12.53-4.813-13.766-11.438l-6.5-35   l-11.438,11.125c-1.906,1.858-4.22,3.266-6.75,4.077l-3.142,0.97l12.484,41.422c0.969,3.203,1.469,6.547,1.469,9.906v15.766   C416,443,411,448,404.828,448H91.172C85,448,80,443,80,436.828v-15.813c0-3.327,0.484-6.641,1.438-9.827l12.391-41.5l-3.016-0.923   c-2.531-0.813-4.844-2.219-6.75-4.077l-11.438-11.125l-6.5,35C64.891,395.188,59.109,400,52.359,400H14c-7.734,0-14-6.266-14-14   v-90.734C0,282.438,10.406,272,23.266,272H48c4.922,0,9.406,2.797,11.609,7.219l18.688,31.703l25.641,24.953l21.172-70.906   c1.594-5.328,6.484-8.969,12.047-8.969H149.5l-3.938-10.734c-2.344-5.141-2.031-11.109,0.844-15.984L181,176.703   C187.875,166.281,199.516,160,212,160l20.234,41.328l9.969-20.125c-2.109-1.469-5.391-5.563-5.391-14.063   c0-6.109,4.828-7.141,11.188-7.141s11.188,1.031,11.188,7.141c0,8.5-3.28,12.594-5.391,14.063l9.969,20.125L284,160   c12.484,0,24.125,6.281,31,16.703l34.594,52.578c2.875,4.875,3.188,10.844,0.844,15.984L346.484,256h11.983   c5.531,0,10.423,3.625,12.017,8.938l21.422,71.109l25.797-25.125l18.688-31.703C438.594,274.797,443.078,272,448,272h24.734   C485.594,272,496,282.438,496,295.266z M288,256h21.391l5.906-16.891l-19.703-28.625L288,256z M186.609,256H208l-7.594-45.516   l-19.703,28.625L186.609,256z M386.766,374.453l-28.233,8.719C356.781,383.734,355,384,353.25,384   c-7.297,0-14.078-4.672-16.422-11.984c-2.922-9.078,2.078-18.797,11.156-21.719l28.844-8.859L355.906,272h-15.297l-8.063,21.922   C329.641,300.25,323.375,304,316.844,304c-2.406,0-4.844-0.5-7.172-1.578c-8.672-3.953-12.484-14.203-8.516-22.875l2.641-7.547   H192.203l2.641,7.547c3.969,8.672,0.156,18.922-8.516,22.875C184,303.5,181.563,304,179.156,304   c-6.531,0-12.797-3.75-15.703-10.078L155.391,272h-15.672L119,341.375l29.016,8.922c9.078,2.922,14.078,12.641,11.156,21.719   C156.828,379.328,150.047,384,142.75,384c-1.75,0-3.531-0.266-5.281-0.828l-28.344-8.75L101.469,400h293L386.766,374.453z"/>
      </svg>
    );
  }
}

RoomIcon.propTypes = {
  color: React.PropTypes.string,
  width: React.PropTypes.number,
  height: React.PropTypes.number
};

RoomIcon.defaultProps = {
    color:"black",
    width:32,
    height:32
};

class ConferencingIcon2 extends BaseComponent {
  constructor(props, context) {
    super(props, context);
  }

  render() {
    this.logRender();
    return (
      <svg width={this.props.width} height={this.props.height} viewBox="0 0 512 512" aria-labelledby="title">
          <title id="title">Conferencing Icon</title>
            <path d="M368,32H144c-8.844,0-16,7.156-16,16v160c0,8.844,7.156,16,16,16h224c8.844,0,16-7.156,16-16V48   C384,39.156,376.844,32,368,32z M368,204.5c0,1.938-1.563,3.5-3.5,3.5h-217c-1.938,0-3.5-1.563-3.5-3.5v-153   c0-1.938,1.563-3.5,3.5-3.5h217c1.938,0,3.5,1.563,3.5,3.5V204.5z M224,96c0-17.672,14.328-32,32-32s32,14.328,32,32   s-14.328,32-32,32S224,113.672,224,96z M320,192c0-9.719,0-12.578,0-15.953c0-10.688-1.875-15.938-7-22.203   c-5.125-6.25-12.484-9.844-20.203-9.844h-4.594l-17.297,34.172l-9.313-13.359c1.938-1.453,5.017-5.453,5.017-13.781   c0-5.984-4.5-7.031-10.469-7.031c-5.953,0-10.875,1.047-10.875,7.031c0,8.328,2.938,12.328,4.875,13.781l-9.156,13.359L223.672,144   h-4.594c-7.719,0-15.063,3.594-20.203,9.844c-5.125,6.266-6.844,11.516-6.844,22.938c0,3.484,0,5.984-0.031,15.219L320,192L320,192   z M416,368c0-17.672,14.328-32,32-32s32,14.328,32,32s-14.328,32-32,32S416,385.672,416,368z M512,464c0-9.719,0-12.578,0-15.953   c0-10.688-1.875-15.938-7-22.203c-5.125-6.25-12.484-9.844-20.203-9.844h-4.594l-17.297,34.172l-9.313-13.359   c1.938-1.452,5.017-5.452,5.017-13.78c0-5.984-4.5-7.031-10.47-7.031c-5.952,0-10.875,1.047-10.875,7.031   c0,8.328,2.938,12.328,4.875,13.78l-9.155,13.359L415.672,416h-4.594c-7.719,0-15.063,3.594-20.203,9.844   c-5.125,6.267-6.844,11.517-6.844,22.938c0,3.483,0,5.983-0.031,15.219L512,464L512,464z M400,224c0-17.672,14.328-32,32-32   s32,14.328,32,32s-14.328,32-32,32S400,241.672,400,224z M496,320c0-9.719,0-12.578,0-15.953c0-10.688-1.875-15.938-7-22.203   c-5.125-6.25-12.484-9.844-20.203-9.844h-4.594l-17.297,34.172l-9.313-13.359c1.938-1.453,5.017-5.453,5.017-13.78   c0-5.984-4.5-7.031-10.47-7.031c-5.952,0-10.875,1.047-10.875,7.031c0,8.327,2.938,12.327,4.875,13.78l-9.155,13.359L399.672,272   h-4.594c-7.719,0-15.063,3.594-20.203,9.844c-5.125,6.266-6.844,11.516-6.844,22.938c0,3.483,0,5.983-0.031,15.219L496,320L496,320   z M64,400c-17.672,0-32-14.328-32-32s14.328-32,32-32s32,14.328,32,32S81.672,400,64,400z M128,464   c-0.031-9.234-0.031-11.734-0.031-15.219c0-11.422-1.719-16.672-6.844-22.938c-5.141-6.25-12.484-9.844-20.203-9.844h-4.594   l-17.313,34.172l-9.156-13.358c1.938-1.453,4.875-5.453,4.875-13.781c0-5.984-4.922-7.031-10.875-7.031   c-5.969,0-10.469,1.047-10.469,7.031c0,8.328,3.078,12.328,5.016,13.781l-9.313,13.358L31.797,416h-4.594   c-7.719,0-15.078,3.594-20.203,9.844c-5.125,6.267-7,11.517-7,22.203c0,3.375,0,6.234,0,15.953H128z M80,256   c-17.672,0-32-14.328-32-32s14.328-32,32-32s32,14.328,32,32S97.672,256,80,256z M144,320c-0.031-9.234-0.031-11.734-0.031-15.219   c0-11.422-1.719-16.672-6.844-22.938c-5.141-6.25-12.484-9.844-20.203-9.844h-4.594l-17.313,34.172l-9.156-13.358   c1.938-1.453,4.875-5.453,4.875-13.781c0-5.984-4.922-7.031-10.875-7.031c-5.969,0-10.469,1.047-10.469,7.031   c0,8.328,3.078,12.328,5.016,13.781l-9.313,13.358L47.797,272h-4.594c-7.719,0-15.078,3.594-20.203,9.844   c-5.125,6.266-7,11.516-7,22.203c0,3.375,0,6.234,0,15.953H144z M335.563,250.063c-0.844-5.781-5.875-10.063-11.78-10.063H188.219   c-5.906,0-10.938,4.281-11.781,10.063l-32.313,200.563c-0.5,3.375,0.531,6.781,2.781,9.344c2.266,2.563,5.547,4.031,9,4.031   h200.188c3.453,0,6.734-1.469,9-4.031c2.25-2.563,3.281-5.969,2.781-9.344L335.563,250.063z"/>
      </svg>
    );
  }
}

class WhiteboardIcon extends BaseComponent {
  constructor(props, context) {
    super(props, context);
  }

  render() {
    this.logRender();
    return (
      <svg width={this.props.width} height={this.props.height} version="1.1" viewBox="0 0 24 24">
        <g>
          <rect height="24" style={{fill:"none"}} width="24"/>
        </g>
        <g id="Line_Icons">
          <path d="M20,2H4C2.896,2,2,2.898,2,4v10c0,1.103,0.896,2,2,2h5.585l-5.999,5.999h2.828L11,17.414v4.585h2   v-4.585l4.586,4.585h2.828L14.415,16H20c1.104,0,2-0.897,2-2V4C22,2.898,21.104,2,20,2z M20,4v0.586l-5,5l-1.293-1.293   c-0.391-0.391-1.023-0.391-1.414,0L11,9.586L8.707,7.293c-0.391-0.391-1.023-0.391-1.414,0L5.586,9H4V4H20z M4,14v-3h2   c0.266,0,0.52-0.105,0.707-0.293L8,9.415l2.293,2.292c0.391,0.391,1.023,0.391,1.414,0L13,10.415l1.293,1.292   C14.488,11.902,14.744,12,15,12s0.512-0.098,0.707-0.293l4.292-4.292L19.997,14H4z" style={{fill:this.props.color}}/>
        </g>
      </svg>
    );
  }
}

WhiteboardIcon.propTypes = {
  color: React.PropTypes.string
};

WhiteboardIcon.defaultProps = {
    color:"black"
};


class ProjectorIcon extends BaseComponent {
  constructor(props, context) {
    super(props, context);
  }

  render() {
    this.logRender();
    return (
      <svg  width={this.props.width} height={this.props.height} id="Layer_1" version="1.1" viewBox="0 -8 64 64">
        <g id="projector">
          <path d="M42,32c-1.657,0-3,1.343-3,3s1.343,3,3,3s3-1.343,3-3S43.657,32,42,32z M42,36c-0.552,0-1-0.448-1-1   s0.448-1,1-1s1,0.448,1,1S42.552,36,42,36z" fill={this.props.color}/>
        <g>
          <rect fill={this.props.color} height="1" width="17" x="6" y="41"/>
          <rect fill={this.props.color} height="1" width="17" x="6" y="39"/>
          <rect fill={this.props.color} height="1" width="17" x="6" y="37"/>
          <rect fill={this.props.color} height="1" width="17" x="6" y="35"/>
          <rect fill={this.props.color} height="1" width="17" x="6" y="33"/>
          <rect fill={this.props.color} height="1" width="17" x="6" y="31"/>
          <rect fill={this.props.color} height="1" width="17" x="6" y="29"/>
          <path d="M59,24H46.778C45.313,23.362,43.7,23,42,23s-3.313,0.362-4.777,1H5c-1.654,0-3,1.346-3,3v16    c0,1.654,1.346,3,3,3h1v2h3v-2h28.223c1.464,0.638,3.077,1,4.777,1s3.313-0.362,4.778-1H55v2h3v-2h1c1.654,0,3-1.346,3-3V27    C62,25.346,60.654,24,59,24z M5,44c-0.552,0-1-0.449-1-1V27c0-0.551,0.448-1,1-1h30.687C32.854,27.99,31,31.275,31,35    s1.854,7.01,4.687,9H5z M42,45c-5.524,0-10-4.477-10-10s4.476-10,10-10s10,4.477,10,10S47.524,45,42,45z M60,43    c0,0.551-0.448,1-1,1H48.313C51.146,42.01,53,38.725,53,35s-1.854-7.01-4.687-9H59c0.552,0,1,0.449,1,1V43z" fill={this.props.color}/>
          <path d="M42,28c-3.867,0-7,3.135-7,7s3.133,7,7,7s7-3.135,7-7S45.867,28,42,28z M42,41c-3.314,0-6-2.687-6-6    s2.686-6,6-6s6,2.687,6,6S45.314,41,42,41z" fill={this.props.color}/>
          <path d="M57,27c-1.104,0-2,0.896-2,2s0.896,2,2,2s2-0.896,2-2S58.104,27,57,27z M57,30c-0.552,0-1-0.448-1-1    s0.448-1,1-1s1,0.448,1,1S57.552,30,57,30z" fill={this.props.color}/>
          </g>
        </g>
      </svg>
    );
  }
}

ProjectorIcon.propTypes = {
  color: React.PropTypes.string,
  width: React.PropTypes.number,
  height: React.PropTypes.number
};

ProjectorIcon.defaultProps = {
    color:"black",
    width:38,
    height:38
};

class RemoteIcon extends BaseComponent {
  constructor(props, context) {
    super(props, context);
  }
   //   width="128px"
  render() {
    this.logRender();
    let style= {};
    if (this.props.position) {
      style.position = this.props.position;
    }
    if (this.props.top) {
      style.top = this.props.top;
    }
    if (this.props.left) {
      style.left = this.props.left;
    }
    return (
      <svg enableBackground="new 0 0 128 128" viewBox="0 0 128 128" style={style}>
        <g>
          <path d="M51,125.5h26c6.3,0,11.5-5.2,11.5-11.5V14c0-6.3-5.2-11.5-11.5-11.5H51c-6.3,0-11.5,5.2-11.5,11.5v100   C39.5,120.3,44.7,125.5,51,125.5z M42.5,14c0-4.7,3.8-8.5,8.5-8.5h26c4.7,0,8.5,3.8,8.5,8.5v100c0,4.7-3.8,8.5-8.5,8.5H51   c-4.7,0-8.5-3.8-8.5-8.5V14z" fill={this.props.color}/>
          <path d="M57,24c2.8,0,5-2.2,5-5s-2.2-5-5-5s-5,2.2-5,5S54.2,24,57,24z M57,17c1.1,0,2,0.9,2,2s-0.9,2-2,2   s-2-0.9-2-2S55.9,17,57,17z" fill={this.props.color}/>
          <path d="M71,24c2.8,0,5-2.2,5-5s-2.2-5-5-5s-5,2.2-5,5S68.2,24,71,24z M71,17c1.1,0,2,0.9,2,2s-0.9,2-2,2   s-2-0.9-2-2S69.9,17,71,17z" fill={this.props.color}/>
          <path d="M57,39c2.8,0,5-2.2,5-5s-2.2-5-5-5s-5,2.2-5,5S54.2,39,57,39z M57,32c1.1,0,2,0.9,2,2s-0.9,2-2,2   s-2-0.9-2-2S55.9,32,57,32z" fill={this.props.color}/>
          <path d="M71,39c2.8,0,5-2.2,5-5s-2.2-5-5-5s-5,2.2-5,5S68.2,39,71,39z M71,32c1.1,0,2,0.9,2,2s-0.9,2-2,2   s-2-0.9-2-2S69.9,32,71,32z" fill={this.props.color}/>
          <path d="M57,54c2.8,0,5-2.2,5-5s-2.2-5-5-5s-5,2.2-5,5S54.2,54,57,54z M57,47c1.1,0,2,0.9,2,2s-0.9,2-2,2   s-2-0.9-2-2S55.9,47,57,47z" fill={this.props.color}/>
          <path d="M71,54c2.8,0,5-2.2,5-5s-2.2-5-5-5s-5,2.2-5,5S68.2,54,71,54z M71,47c1.1,0,2,0.9,2,2s-0.9,2-2,2   s-2-0.9-2-2S69.9,47,71,47z" fill={this.props.color}/>
          <path d="M56,84c2.8,0,5-2.2,5-5s-2.2-5-5-5s-5,2.2-5,5S53.2,84,56,84z M56,77c1.1,0,2,0.9,2,2s-0.9,2-2,2   s-2-0.9-2-2S54.9,77,56,77z" fill={this.props.color}/>
          <path d="M72,74c-2.8,0-5,2.2-5,5s2.2,5,5,5s5-2.2,5-5S74.8,74,72,74z M72,81c-1.1,0-2-0.9-2-2s0.9-2,2-2s2,0.9,2,2   S73.1,81,72,81z" fill={this.props.color}/>
          <path d="M64,82c-2.8,0-5,2.2-5,5s2.2,5,5,5s5-2.2,5-5S66.8,82,64,82z M64,89c-1.1,0-2-0.9-2-2s0.9-2,2-2s2,0.9,2,2   S65.1,89,64,89z" fill={this.props.color}/>
          <path d="M59,71c0,2.8,2.2,5,5,5s5-2.2,5-5s-2.2-5-5-5S59,68.2,59,71z M64,69c1.1,0,2,0.9,2,2s-0.9,2-2,2s-2-0.9-2-2   S62.9,69,64,69z" fill={this.props.color}/>
        </g>
      </svg>
    );
  }
}

RemoteIcon.propTypes = {
  color: React.PropTypes.string
};

RemoteIcon.defaultProps = {
    color:"black"
};

class RevokeUserIcon extends BaseComponent {
  constructor(props, context) {
    super(props, context);
  }

  render() {
    this.logRender();
    return (
      <svg width={this.props.width} height={this.props.height} viewBox="0 0 512 512" style={{position:this.props.position, top:this.props.top, left:this.props.left}}  >
          <g
             transform="translate(0,448)"
             id="layer1">
            <g
               transform="matrix(1.0239926,0,0,1.0239926,-513.92473,-8769.2473)"
               id="g5473">
              <path
                 style={{fill:this.props.color,fillOpacity:1,fillRule:"nonzero",stroke:"none"}}
                 id="rect3778-1-0-9-4-3-4"
                 d="m 886.8785,8360.3359 c -29.4309,0 -58.8625,11.2136 -81.3175,33.6686 -44.91,44.9098 -44.91,117.7391 0,162.649 44.91,44.9099 117.739,44.9099 162.649,0 l 0.4792,-0.4933 c 44.4463,-44.9877 44.2386,-117.4381 -0.4792,-162.1557 -22.455,-22.455 -51.9006,-33.6686 -81.3315,-33.6686 z m -1.1134,29.7929 c 0.4944,-0.01 0.9995,0 1.4939,0 21.6801,0.094 43.3254,8.4026 59.8677,24.9449 28.6006,28.6005 32.648,72.5145 12.0356,105.4309 L 841.7099,8403.0663 c 13.4444,-8.4183 28.7297,-12.7383 44.0552,-12.9375 z m -71.5791,40.7292 117.1705,117.1705 c -32.8187,20.1145 -76.3101,15.9717 -104.7262,-12.4443 -28.416,-28.416 -32.5609,-71.9097 -12.4443,-104.7262 z" />
              <path
                 transform="translate(0,-9449.9978)"
                 style={{fill:this.props.color,fillOpacity:1,stroke:"none"}}
                 id="path4863"
                 d="m 684.375,17612.219 c -45.601,0 -82.5625,45.521 -82.5625,101.656 0,56.135 36.9615,101.625 82.5625,101.625 45.601,0 82.5625,-45.49 82.5625,-101.625 0,-56.135 -36.9615,-101.656 -82.5625,-101.656 z m -82.96875,217.343 c -10.7311,-0.04 -22.20555,3.974 -33.34375,13.094 -30.61,25.063 -51.65025,58.078 -63.15625,94.594 -10.7675,34.177 8.5205,43.012 37.4375,46.812 23.703,3.118 39.6087,10.754 59.375,17.063 37.0891,11.846 59.58595,12.094 82.65625,12.094 23.0703,0 45.56705,-0.248 82.65625,-12.094 3.66859,-1.171 7.17733,-2.4 10.65625,-3.625 -15.48902,-23.64 -23.73569,-52.117 -21.6875,-80.531 1.19491,-30.425 14.60051,-58.967 35.1875,-80.844 -14.81509,-8.333 -29.67307,-8.323 -42.125,-2.313 -22.666,10.944 -42.7069,16.626 -64.6875,16.626 -21.9806,0 -42.0216,-5.682 -64.6875,-16.626 -5.6702,-2.736 -11.84255,-4.225 -18.28125,-4.25 z" />
            </g>
          </g>
      </svg>
    );
  }
}

RevokeUserIcon.propTypes = {
  color: React.PropTypes.string,
  width: React.PropTypes.number,
  height: React.PropTypes.number,
  top: React.PropTypes.number,
  left: React.PropTypes.number,
  position:React.PropTypes.string
};

RevokeUserIcon.defaultProps = {
    color:"black"
};

class ExitIcon extends BaseComponent {
  constructor(props, context) {
    super(props, context);
  }

  render() {
    this.logRender();
    return (
      <svg enableBackground="new 0 0 24 24" width={this.props.width} height={this.props.height} viewBox="0 0 24 24">
        <g>
          <path style={{fill:this.props.color}} d="M23,21h-3V1c0-0.6-0.4-1-1-1H5C4.4,0,4,0.4,4,1v20H1c-0.6,0-1,0.4-1,1v2h24v-2C24,21.4,23.6,21,23,21z M6,2h12v19H6V2z    M1,23v-1h22v1H1z"/>
          <circle style={{fill:this.props.color}} cx="16" cy="11" r="1"/>
        </g>
      </svg>
    );
  }
}

ExitIcon.propTypes = {
  color: React.PropTypes.string,
  width: React.PropTypes.number,
  height: React.PropTypes.number
};

ExitIcon.defaultProps = {
    color:"black"
};

class DisplayIcon extends BaseComponent {
  constructor(props, context) {
    super(props, context);
  }

  render() {
    this.logRender();
    return (
      <svg width={this.props.width} height={this.props.height} enableBackground="new 0 0 92.168 92.168" viewBox="0 0 92.168 92.168" >
        <rect clipRule="evenodd" fill="none" fillRule="evenodd" width={this.props.width} height={this.props.height} x="0"/>
        <path style={{fill:this.props.color}} d="M27.544,77.749c-0.855,0-1.527-0.67-1.527-1.523c0-0.833,0.672-1.523,1.527-1.523h17.013v-5.99H1.527  C0.672,68.711,0,68.019,0,67.188V15.943c0-0.853,0.672-1.523,1.527-1.523h89.113c0.835,0,1.527,0.67,1.527,1.523v51.245  c0,0.831-0.692,1.523-1.527,1.523H47.611v5.99h17.013c0.835,0,1.527,0.69,1.527,1.523c0,0.853-0.692,1.523-1.527,1.523H27.544z   M89.113,17.466H3.055v48.198h86.058V17.466z"/>
      </svg>
    );
  }
}

DisplayIcon.propTypes = {
  color: React.PropTypes.string,
  width: React.PropTypes.number,
  height: React.PropTypes.number
};

DisplayIcon.defaultProps = {
    color:"black"
};

class DVDIcon extends BaseComponent {
  constructor(props, context) {
    super(props, context);
  }

  render() {
    this.logRender();
    return (
      <svg width={this.props.width} height={this.props.height} enableBackground="new 0 0 512 512" viewBox="0 0 512 512" >
        <g>
          <polygon style={{fill:this.props.color}} points="456.025,134.681 460.727,134.681 460.727,107.092 472.313,107.092 472.313,103 444.439,103    444.439,107.092 456.025,107.092  "/>
            <path style={{fill:this.props.color}} d="M240.163,325.373C106.058,325.373,0,343.169,0,365.127c0,21.953,106.058,39.755,240.163,39.755   c134.103,0,244.804-17.802,244.804-39.755C484.967,343.169,374.266,325.373,240.163,325.373z M236.672,393.208   c-97.399,0-176.356-13.539-176.356-30.24c0-16.7,78.958-30.239,176.356-30.239c97.398,0,176.357,13.539,176.357,30.239   C413.029,379.669,334.07,393.208,236.672,393.208z"/>
            <path style={{fill:this.props.color}} d="M432.225,166.36H314.799l-62.694,72.549l-29.854-71.557h-185.1l-3.98,31.14h73.642   c0,0,33.412-1.941,31.845,22.529c0,0-0.225,40.058-49.757,41.741l-26.109-0.496l13.171-47.209H28.195l-21.89,78.848h86.574   c0,0,77.288,1.802,91.553-71.89c0,0,3.728-18.63,2.137-26.845c-0.245-1.267,0.532,0.665,0,0l44.636,116.29L342.66,198.827h84.588   c0,0,25.873,1.326,25.873,22.195c0,0,2.232,38.852-38.813,41.742h-34.826l8.957-47.708h-41.797l-23.885,78.848h79.611   c0,0,98.633,1.802,99.516-73.879C501.885,220.024,499.895,166.36,432.225,166.36z"/>
            <path style={{fill:this.props.color}} d="M502.848,103l-9.814,26.375L483.146,103h-6.67v31.681h4.447v-20.262c0-0.808-0.025,0.574-0.072-1.145   c-0.045-1.725-0.068-3.01-0.068-3.848v-1.049l9.908,26.304h4.613l9.814-26.304c0,1.85-0.016,3.588-0.049,5.21   c-0.029,1.619-0.043,0.131-0.043,0.831v20.262h4.443V103H502.848z"/>
            <path style={{fill:this.props.color}} d="M237.552,342.824c-82.607,0-149.574,9.121-149.574,20.373c0,11.251,66.967,20.372,149.574,20.372   c82.606,0,149.573-9.121,149.573-20.372C387.127,351.945,320.16,342.824,237.552,342.824z M237.292,369.58   c-19.315,0-34.974-2.967-34.974-6.627s15.659-6.627,34.974-6.627c19.315,0,34.974,2.967,34.974,6.627   S256.607,369.58,237.292,369.58z"/>
        </g>
      </svg>
    );
  }
}

DVDIcon.propTypes = {
  color: React.PropTypes.string,
  width: React.PropTypes.number,
  height: React.PropTypes.number
};

DVDIcon.defaultProps = {
  color:"black",
  width:512,
  height:512
};

class BluRayIcon extends BaseComponent {
  constructor(props, context) {
    super(props, context);
  }

  render() {
    this.logRender();
    return (
      <svg width={this.props.width} height={this.props.height} enableBackground="new 0 0 512 512"  viewBox="0 0 512 512" >
        <g>
          <polygon style={{fill:this.props.color}} points="364.489,395.516 380.7,395.516 380.7,425.074 387.281,425.074 387.281,395.516 403.491,395.516    403.491,390.914 364.489,390.914  "/>
            <path style={{fill:this.props.color}} d="M437.011,419.104l-13.835-28.189h-9.337v34.16h6.227v-22.801c0-0.916-0.036-2.342-0.102-4.277   c-0.008-0.34-0.016-0.58-0.025-0.891l13.794,27.969h6.457l13.733-28.117c0,2.09-0.02-0.428-0.063,1.393   c-0.044,1.826-0.065,3.137-0.065,3.926v22.801h6.222v-34.16h-9.27L437.011,419.104z"/>
            <path style={{fill:this.props.color}} d="M270.88,291.283c0,0,172.684-42.339,76.869-102.232c0,0-119.136-72.766-232.46,1.26l45-91.923L60.376,98.06   L7.489,207.258c0,0-44.598,83.646,82.303,83.838C257.544,299.083,270.88,291.283,270.88,291.283z M233.112,192.483   c41.975,0,76.003,13.838,76.003,30.91c0,17.071-34.028,30.909-76.003,30.909c-41.976,0-76.005-13.838-76.005-30.909   S191.136,192.483,233.112,192.483z"/>
            <path style={{fill:this.props.color}} d="M462.954,142.76c0,0-85.089-60.559-260.096-45.155l-0.002,7.821c0,0,200.156,12.782,226.64,96.724   c0,0,35.535,88.715-171.897,131.971c0,0-96.264,16.016-195.601,13.037l-0.64,18.607c0,0,283.522,6.029,364.619-40.08   C425.977,325.686,592.081,253.124,462.954,142.76z" />
        </g>
      </svg>
    );
  }
}

BluRayIcon.propTypes = {
  color: React.PropTypes.string,
  width: React.PropTypes.number,
  height: React.PropTypes.number
};

BluRayIcon.defaultProps = {
  color:"black",
  width:512,
  height:512
};

class AppleTVIcon extends BaseComponent {
  constructor(props, context) {
    super(props, context);
  }

  render() {
    this.logRender();
    return (
      <svg width={this.props.width} height={this.props.height} enableBackground="new 0 0 128 128" viewBox="0 0 128 128" >
        <g>
          <path style={{fill:this.props.color}} d="M97.882,0H30.118C13.553,0,0,13.553,0,30.116v67.767C0,114.446,13.553,128,30.118,128h67.764   C114.446,128,128,114.446,128,97.883V30.116C128,13.553,114.446,0,97.882,0z M41.874,46.39c1.318-1.548,3.597-2.727,5.425-2.805   c0.238,2.179-0.629,4.322-1.888,5.904c-1.317,1.55-3.423,2.739-5.476,2.584C39.667,49.982,40.707,47.748,41.874,46.39z    M52.4,75.844c-1.512,2.266-3.094,4.473-5.614,4.514c-2.45,0.055-3.27-1.439-6.077-1.439c-2.833,0-3.711,1.398-6.058,1.494   c-2.399,0.091-4.229-2.412-5.795-4.655c-3.135-4.576-5.572-12.901-2.305-18.567c1.586-2.778,4.478-4.564,7.567-4.612   c2.402-0.049,4.631,1.618,6.105,1.618c1.445,0,4.199-1.993,7.043-1.698c1.191,0.032,4.568,0.467,6.75,3.648   c-0.178,0.102-4.033,2.368-3.988,7.022c0.045,5.577,4.874,7.419,4.937,7.44C54.938,70.741,54.213,73.271,52.4,75.844z    M76.938,56.202h-7.057v14.619c0,3.362,0.951,5.264,3.697,5.264c1.344,0,2.125-0.112,2.856-0.335l0.223,3.75   c-0.954,0.339-2.465,0.674-4.368,0.674c-2.297,0-4.146-0.782-5.32-2.071c-1.345-1.514-1.906-3.919-1.906-7.115V56.202h-4.199   v-3.754h4.199v-5.04l4.818-1.458v6.497h7.057V56.202z M94.358,79.558h-4.706l-10.305-27.11h5.264l5.321,15.178   c0.897,2.521,1.626,4.762,2.187,7.003h0.165c0.618-2.241,1.401-4.481,2.298-7.003l5.266-15.178H105L94.358,79.558z"/>
        </g>
      </svg>
    );
  }
}

AppleTVIcon.propTypes = {
  color: React.PropTypes.string,
  width: React.PropTypes.number,
  height: React.PropTypes.number
};

AppleTVIcon.defaultProps = {
  color:"black",
  width:128,
  height:128
};


class PlayPause extends BaseComponent {
  constructor(props, context) {
    super(props, context);
  }

  render() {
    this.logRender();
    return (
      <svg width={this.props.width} height={this.props.height} enableBackground="new 0 0 72 72" viewBox="0 0 72 72" >
      <g>
        <path style={{fill:this.props.color}} d="M6.8,19v28l22-14L6.8,19z"/>
        <path style={{fill:this.props.color}} d="M35.3,47h8V19h-8V47z M51.3,19v28h8V19H51.3z"/>
      </g>
      </svg>
    );
  }
}

PlayPause.propTypes = {
  color: React.PropTypes.string,
  width: React.PropTypes.number,
  height: React.PropTypes.number
};

PlayPause.defaultProps = {
  color:"black",
  width:32,
  height:32
};

class VGA extends BaseComponent {
  constructor(props, context) {
    super(props, context);
  }

  render() {
    this.logRender();
    return (

      <svg width={this.props.width} height={this.props.height} viewBox={"-4 " + this.props.viewBoxBottom + " 56 56"} >
        <g>
          <path style={{fill:"none", stroke:this.props.color}} d="M38.6,35H9.4c-2.2,0-4.2-1.4-4.6-3.2L1.6,17.3   c-0.6-2.5,1.2-4.7,3.9-4.7h37c2.7,0,4.5,2.2,3.9,4.7l-3.2,14.5C42.8,33.6,40.8,35,38.6,35z" strokeMiterlimit="10" strokeWidth="2"/>
          <g>
            <line style={{fill:"none", stroke:this.props.color}} strokeMiterlimit="10" strokeWidth="2" x1="42.4" x2="40.4" y1="17.7" y2="17.7"/>
            <line style={{fill:"none", stroke:this.props.color}} strokeMiterlimit="10" strokeWidth="2" x1="38.1" x2="36" y1="17.7" y2="17.7"/>
            <line style={{fill:"none", stroke:this.props.color}} strokeMiterlimit="10" strokeWidth="2" x1="33.7" x2="31.7" y1="17.7" y2="17.7"/>
            <line style={{fill:"none", stroke:this.props.color}} strokeMiterlimit="10" strokeWidth="2" x1="29.4" x2="27.3" y1="17.7" y2="17.7"/>
            <line style={{fill:"none", stroke:this.props.color}} strokeMiterlimit="10" strokeWidth="2" x1="25" x2="23" y1="17.7" y2="17.7"/>
            <line style={{fill:"none", stroke:this.props.color}} strokeMiterlimit="10" strokeWidth="2" x1="20.7" x2="18.6" y1="17.7" y2="17.7"/>
            <line style={{fill:"none", stroke:this.props.color}} strokeMiterlimit="10" strokeWidth="2" x1="16.3" x2="14.3" y1="17.7" y2="17.7"/>
            <line style={{fill:"none", stroke:this.props.color}} strokeMiterlimit="10" strokeWidth="2" x1="12" x2="9.9" y1="17.7" y2="17.7"/>
            <line style={{fill:"none", stroke:this.props.color}} strokeMiterlimit="10" strokeWidth="2" x1="7.6" x2="5.6" y1="17.7" y2="17.7"/>
          </g>
          <g>
            <line style={{fill:"none", stroke:this.props.color}} strokeMiterlimit="10" strokeWidth="2" x1="38.1" x2="36" y1="23.8" y2="23.8"/>
            <line style={{fill:"none", stroke:this.props.color}} strokeMiterlimit="10" strokeWidth="2" x1="33.7" x2="31.7" y1="23.8" y2="23.8"/>
            <line style={{fill:"none", stroke:this.props.color}} strokeMiterlimit="10" strokeWidth="2" x1="29.4" x2="27.3" y1="23.8" y2="23.8"/>
            <line style={{fill:"none", stroke:this.props.color}} strokeMiterlimit="10" strokeWidth="2" x1="25" x2="23" y1="23.8" y2="23.8"/>
            <line style={{fill:"none", stroke:this.props.color}} strokeMiterlimit="10" strokeWidth="2" x1="20.7" x2="18.6" y1="23.8" y2="23.8"/>
            <line style={{fill:"none", stroke:this.props.color}} strokeMiterlimit="10" strokeWidth="2" x1="16.3" x2="14.3" y1="23.8" y2="23.8"/>
            <line style={{fill:"none", stroke:this.props.color}} strokeMiterlimit="10" strokeWidth="2" x1="12" x2="9.9" y1="23.8" y2="23.8"/>
          </g>
          <g>
            <line style={{fill:"none", stroke:this.props.color}} strokeMiterlimit="10" strokeWidth="2" x1="33.7" x2="31.7" y1="29.9" y2="29.9"/>
            <line style={{fill:"none", stroke:this.props.color}} strokeMiterlimit="10" strokeWidth="2" x1="29.4" x2="27.3" y1="29.9" y2="29.9"/>
            <line style={{fill:"none", stroke:this.props.color}} strokeMiterlimit="10" strokeWidth="2" x1="25" x2="23" y1="29.9" y2="29.9"/>
            <line style={{fill:"none", stroke:this.props.color}} strokeMiterlimit="10" strokeWidth="2" x1="20.7" x2="18.6" y1="29.9" y2="29.9"/>
            <line style={{fill:"none", stroke:this.props.color}} strokeMiterlimit="10" strokeWidth="2" x1="16.3" x2="14.3" y1="29.9" y2="29.9"/>
          </g>
        </g>
      </svg>
    );
  }
}

VGA.propTypes = {
  color: React.PropTypes.string,
  width: React.PropTypes.number,
  height: React.PropTypes.number,
  viewBoxBottom:React.PropTypes.number
};

VGA.defaultProps = {
  color:"black",
  width:48,
  height:48,
  viewBoxBottom:-5
};

class MuteVideo extends BaseComponent {
  constructor(props, context) {
    super(props, context);
  }

  render() {
    this.logRender();
    return (

      <svg width={this.props.width} height={this.props.height} viewBox="0 0 300 300" >

        <g/>
        <path
           style={{fill:this.props.color, strokeWidth:0.2262295}}
        />
        <g
           transform="translate(-0.22622951,15.609836)">
          <g
             transform="matrix(0.35526034,0,0,0.38856026,56.51111,37.554089)"
          >
            <g>
              <path
                 d="M 441.677,43.643 H 10.687 C 4.785,43.643 0,48.427 0,54.329 v 297.425 c 0,5.898 4.785,10.676 10.687,10.676 h 162.069 v 25.631 c 0,0.38 0.074,0.722 0.112,1.089 h -23.257 c -5.407,0 -9.796,4.389 -9.796,9.795 0,5.408 4.389,9.801 9.796,9.801 h 158.506 c 5.406,0 9.795,-4.389 9.795,-9.801 0,-5.406 -4.389,-9.795 -9.795,-9.795 h -23.256 c 0.032,-0.355 0.115,-0.709 0.115,-1.089 V 362.43 H 441.7 c 5.898,0 10.688,-4.782 10.688,-10.676 V 54.329 C 452.37,48.427 447.589,43.643 441.677,43.643 Z m -19.588,261.49 c 0,5.903 -4.784,10.687 -10.683,10.687 H 40.96 c -5.898,0 -10.684,-4.783 -10.684,-10.687 V 79.615 c 0,-5.898 4.786,-10.684 10.684,-10.684 h 370.446 c 5.898,0 10.683,4.785 10.683,10.684 z"
                 style={{fill:this.props.color}}
              />
            </g>
          </g>
          <path

             d="m 128.83771,259.6836 c -23.71859,-1.65779 -46.388608,-9.33991 -66.500488,-22.53479 -3.98658,-2.61551 -9.20286,-6.46442 -12.67984,-9.35604 -3.11723,-2.59243 -14.49502,-14.00148 -17.12606,-17.1731 C 14.440812,188.81219 3.4239725,162.2464 0.91485246,134.38032 c -0.39374,-4.37285 -0.52999,-14.49393 -0.26647,-19.79508 C 2.0994825,85.393495 13.407452,56.726873 32.596142,33.595083 c 2.53406,-3.054797 13.9457,-14.429793 17.17435,-17.119205 12.85516,-10.7081492 28.1897,-19.3730082 43.57609,-24.6228802 10.671898,-3.6412798 21.633188,-5.9710188 33.115718,-7.0384988 4.17357,-0.387999 19.79553,-0.392977 23.64098,-0.0075 18.57763,1.862103 34.43447,6.3865338 50.10984,14.29783282 9.84662,4.96955798 19.07736,11.06204318 26.92131,17.76865618 3.43964,2.940907 13.35498,12.849447 16.28033,16.269168 12.17721,14.235048 21.75539,32.041887 27.07547,50.336068 3.69887,12.719351 5.50678,25.541806 5.50235,39.024576 -0.0106,32.09259 -11.27847,63.0512 -32.07145,88.1164 -2.63521,3.17667 -14.01455,14.5861 -17.12605,17.17134 -22.95695,19.0742 -51.22011,30.28944 -80.42459,31.91367 -4.80503,0.26723 -13.55842,0.25671 -17.53278,-0.021 z m 13.36169,-27.56749 c 2.86959,-0.11576 6.63631,-0.37409 8.37049,-0.57406 17.1281,-1.97494 33.50225,-7.88021 47.83339,-17.25088 3.98982,-2.60883 9.27541,-6.51692 9.27541,-6.85811 0,-0.38892 -154.553968,-154.950279 -154.794328,-154.801728 -0.37599,0.232374 -4.83819,6.313084 -6.60933,9.006642 -16.78936,25.533311 -22.22673,57.077176 -14.9408,86.676256 10.44147,42.41853 44.54819,74.43887 87.549378,82.19379 5.13402,0.92587 8.97983,1.38371 13.34754,1.58899 1.99082,0.0936 3.87418,0.1835 4.18524,0.19986 0.31107,0.0163 2.91342,-0.065 5.78301,-0.18076 z m 82.58733,-41.91366 c 5.43347,-6.8172 10.61274,-15.52502 14.29933,-24.0412 12.88929,-29.77485 12.21782,-63.03375 -1.85892,-92.075258 -12.99227,-26.804174 -36.51385,-47.318415 -64.84025,-56.550122 -27.4662,-8.9513602 -56.96145,-6.819715 -82.913118,5.992202 -6.59423,3.255463 -11.70299,6.387513 -17.46166,10.705301 -2.15504,1.615829 -3.01211,2.390694 -3.01211,2.723227 0,0.31716 24.48757,24.951499 77.100188,77.56228 42.40511,42.40363 77.20356,77.09751 77.3299,77.09751 0.12634,0 0.73683,-0.63627 1.35664,-1.41394 z"
             style={{fill:this.props.color, strokeWidth:0.2262295}} />
        </g>
        <g/>
      </svg>
    );
  }
}

MuteVideo.propTypes = {
  color: React.PropTypes.string,
  width: React.PropTypes.number,
  height: React.PropTypes.number
};

MuteVideo.defaultProps = {
  color:"black",
  width:48,
  height:48
};

class UnMuteVideo extends BaseComponent {
  constructor(props, context) {
    super(props, context);
  }

  render() {
    this.logRender();
    return (

      <svg width={this.props.width} height={this.props.height} viewBox="-20 0 256 256" >

        <g/>
        <path
           style={{fill:this.props.color, strokeWidth:0.2262295}}
        />
        <g
           transform="translate(-0.22622951,15.609836)">
          <g
             transform="matrix(0.35526034,0,0,0.38856026,56.51111,37.554089)"
          >
            <g>
              <path
                 d="M 441.677,43.643 H 10.687 C 4.785,43.643 0,48.427 0,54.329 v 297.425 c 0,5.898 4.785,10.676 10.687,10.676 h 162.069 v 25.631 c 0,0.38 0.074,0.722 0.112,1.089 h -23.257 c -5.407,0 -9.796,4.389 -9.796,9.795 0,5.408 4.389,9.801 9.796,9.801 h 158.506 c 5.406,0 9.795,-4.389 9.795,-9.801 0,-5.406 -4.389,-9.795 -9.795,-9.795 h -23.256 c 0.032,-0.355 0.115,-0.709 0.115,-1.089 V 362.43 H 441.7 c 5.898,0 10.688,-4.782 10.688,-10.676 V 54.329 C 452.37,48.427 447.589,43.643 441.677,43.643 Z m -19.588,261.49 c 0,5.903 -4.784,10.687 -10.683,10.687 H 40.96 c -5.898,0 -10.684,-4.783 -10.684,-10.687 V 79.615 c 0,-5.898 4.786,-10.684 10.684,-10.684 h 370.446 c 5.898,0 10.683,4.785 10.683,10.684 z"
                 style={{fill:this.props.color}}
              />
            </g>
          </g>
        </g>
        <g/>
      </svg>
    );
  }
}

UnMuteVideo.propTypes = {
  color: React.PropTypes.string,
  width: React.PropTypes.number,
  height: React.PropTypes.number
};

UnMuteVideo.defaultProps = {
  color:"black",
  width:48,
  height:48
};


class DatabaseIcon extends BaseComponent {
  constructor(props, context) {
    super(props, context);
  }

  render() {
    this.logRender();
    return (

      <svg width={this.props.width} height={this.props.height} enable-background="new 0 0 48 48" viewBox="0 0 48 48">
      	<g id="Expanded">
      		<g>
      			<g>
      				<path style={{fill:this.props.color}} d="M24,20c-11.215,0-20-3.953-20-9s8.785-9,20-9s20,3.953,20,9S35.215,20,24,20z M24,4C15.486,4,6,6.875,6,11s9.486,7,18,7     s18-2.875,18-7S32.514,4,24,4z"/>
      			</g>
      			<g>
      				<path style={{fill:this.props.color}} d="M24,28c-11.215,0-20-3.953-20-9v-8c0-0.553,0.447-1,1-1s1,0.447,1,1v8c0,4.125,9.486,7,18,7s18-2.875,18-7v-8     c0-0.553,0.447-1,1-1s1,0.447,1,1v8C44,24.047,35.215,28,24,28z"/>
      			</g>
      			<g>
      				<path style={{fill:this.props.color}} d="M24,37c-11.215,0-20-3.953-20-9v-9c0-0.553,0.447-1,1-1s1,0.447,1,1v9c0,4.125,9.486,7,18,7s18-2.875,18-7v-9     c0-0.553,0.447-1,1-1s1,0.447,1,1v9C44,33.047,35.215,37,24,37z"/>
      			</g>
      			<g>
      				<path style={{fill:this.props.color}} d="M24,46c-11.215,0-20-3.953-20-9v-9c0-0.553,0.447-1,1-1s1,0.447,1,1v9c0,4.125,9.486,7,18,7s18-2.875,18-7v-9     c0-0.553,0.447-1,1-1s1,0.447,1,1v9C44,42.047,35.215,46,24,46z"/>
      			</g>
      		</g>
      	</g>
      </svg>
    );
  }
}

DatabaseIcon.propTypes = {
  color: React.PropTypes.string,
  width: React.PropTypes.number,
  height: React.PropTypes.number
};

DatabaseIcon.defaultProps = {
  color:"black",
  width:48,
  height:48
};

export {RoomIcon,
  ConferencingIcon2,
  WhiteboardIcon,
  ProjectorIcon,
  RemoteIcon,
  RevokeUserIcon,
  ExitIcon,
  DisplayIcon,
  DVDIcon,
  BluRayIcon,
  AppleTVIcon,
  PlayPause,
  VGA,
  MuteVideo,
  UnMuteVideo,
  DatabaseIcon
};
