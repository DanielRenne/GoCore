import {
    React,
    BasePageComponent,
    TextField,
    RaisedButton,
    Toggle,
    IconButton,
    ConfirmPopup,
} from "../../globals/forms";
import {deepOrange500} from 'material-ui/styles/colors';
import CircularProgress from 'material-ui/CircularProgress';
import Paper from 'material-ui/Paper';
import {Sync} from '../../globals/icons'


// WARNING!! Do not use pageContent and only use appContent because both the /home or the /login could be the endpoint needing two entries of page and app content

class Login extends BasePageComponent {
  constructor(props, context) {
    super(props, context);

    this.changeAlias;

    this.state = {
      MarginTop:  window.innerHeight / 2 - 230,
      ForgotPassword: false,
    };

    this.handleUserNameChange = (event) => {
      this.setComponentState({username: event.target.value, UserNameError: ""});
    };

    this.handlePasswordChange = (event) => {
      this.setComponentState({password: event.target.value, PasswordError: ""});
    };

    this.handleRoomIdChange = (event) => {
      this.setComponentState({RoomId: event.target.value, RoomIdError: ""});
    };

    this.handleRoomKeyChange = (event) => {
      this.setComponentState({RoomKey: event.target.value, RoomKeyError: ""});
    };

    this.handleRoomKeyDown = (event) => {
      if (event.nativeEvent.key == "Enter"){
        this.setComponentState({authMessage: "", RoomIdError: "", RoomKeyError: ""});
        window.api.post({action: "Authorize", state: this.state, controller:"login"});
      }
    };

    this.handlePasswordKeyDown = (event) => {
      if (event.nativeEvent.key == "Enter"){
        this.setComponentState({authMessage: "", UserNameError: "", PasswordError: ""});
        window.api.post({action: "Authorize", state: this.state, controller:"login"});
      }
    };

    this.handleForgotPassword = (event) => {
      if (!this.state.ForgotPassword) {
        this.setComponentState({UserNameError: "", ForgotPassword: true});
      } else {
        this.setComponentState({ForgotPassword: false}, () => {
          window.api.post({action: "ForgotPassword", state: this.state, controller:"login"});
        });
      }
    };


    this.handleForgotPasswordCancel = (event) => {
      this.setComponentState({UserNameError: "", ForgotPassword: false});
    };

    this.login = () => {
      this.setState({authMessage: "", UserNameError: "", PasswordError: ""});
      window.api.post({action: "Authorize", state: this.state, controller:"login"});
    };

    this.handleUseRoomKey = (event, value) => {
      this.setComponentState({LoginByKey: value});
    };

    this.resizeEvent = (e) => this.handleResize(e);
  }

  handleResize(e) {
    this.setComponentState({height: window.innerHeight});
  }


  componentDidMount() {
    this.handleResize();

    window.addEventListener('resize', this.resizeEvent);
    $('body').css({paddingTop: '0px'});

  }

  componentWillUnmount() {
    window.removeEventListener('resize', this.resizeEvent);
  }


  render() {

    var style = {
      width: 520,
      paddingBottom: 30,
      paddingTop: 30,
      minWidth: $(window).width() > 350 ? 350: $(window).width() - 30,
      margin: 20,
      textAlign: 'center',
      display: 'inline-block',
      position: 'relative'
    };

    this.logRender();

    var fields = (
      <div>
      <TextField
        fullWidth={true}
        hintText={window.appContent.LoginPageEmailAddress}
        floatingLabelText={window.appContent.LoginPageEmailAddress}
        id="username"
        autoComplete={this.state.username}
        onChange={this.handleUserNameChange}
        defaultValue={this.state.username}
        errorText={window.appContent[this.state.UserNameError]}
        style={{visibility:(this.state.LoginByKey) ? "collapse" : "visible"}}
        />
        {!this.state.ForgotPassword ?
          <span>
            <br />
            <TextField
              fullWidth={true}
              hintText={window.appContent.LoginPagePassword}
              floatingLabelText={window.appContent.LoginPagePassword}
              id="password"
              type="password"
              autoComplete={this.state.password}
              onChange={this.handlePasswordChange}
              onKeyDown={this.handlePasswordKeyDown}
              defaultValue={this.state.password}
              errorText={window.appContent[this.state.PasswordError]}
              style={{visibility:(this.state.LoginByKey) ? "collapse" : "visible"}}
            />
          </span> : null}
        </div>
      );

    if (this.state.LoginByKey) {
      fields = (
        <div>
          <TextField
            fullWidth={true}
            hintText={window.appContent.RoomIdLabel}
            floatingLabelText={window.appContent.RoomIdLabel}
            onChange={this.handleRoomIdChange}
            value={this.state.RoomId}
            errorText={window.appContent[this.state.RoomIdError]}
            style={{visibility:(this.state.LoginByKey) ? "visible" : "collapse"}}
          />
          <br />
          <TextField
            fullWidth={true}
            hintText={window.appContent.RoomKeyLabel}
            floatingLabelText={window.appContent.RoomKeyLabel}
            type="password"
            onChange={this.handleRoomKeyChange}
            value={this.state.RoomKey}
            onKeyDown={this.handleRoomKeyDown}
            errorText={window.appContent[this.state.RoomKeyError]}
            style={{visibility:(this.state.LoginByKey) ? "visible" : "collapse"}}
          />
        </div>
        );
    }

    return (
        <div>

          {this.state.height > 900 ? <div className="Aligner" style={{marginBottom: 50}}>
            <a href="#" target="_blank"><img style={{width:$(window).width() > 350 ? 350: $(window).width() - 100}} src="/web/app/images/go_core_company_logo.png"/></a>
          </div>: null}

          {this.state.height > 490 ? <div className="Aligner">
            <img style={{width:$(window).width() > 350 ? this.state.height < 900 ? 350: 200: $(window).width() - 100}} src={"/web/app/images/go_core_app_product.png"}/>
          </div>: null}

          <div className="Align" style={{maxHeight:450}}>

            <Paper style={style} className={($(window).width() > 400) ? 'Align-center' : ''} rounded={true} zDepth={3}>

              <h1 className="Paper-Title" style={{paddingTop:0}}>{window.appContent.Login}</h1>
              <div style={{paddingLeft: 30, paddingRight: 30}}>
                  {fields}
                  <span className="Aligner" style={{marginBottom: 10, marginTop: 10}}>
                    <RaisedButton
                      label={!this.state.ForgotPassword ? appContent.btnLogin: window.appContent.ResetPass}
                      onTouchTap={!this.state.ForgotPassword ? this.login: this.handleForgotPassword}
                      secondary={true}
                    />
                    {this.state.ForgotPassword ?
                      <RaisedButton
                        style={{marginLeft:40}}
                        label={window.appContent.Cancel}
                        onTouchTap={this.handleForgotPasswordCancel}
                        secondary={true}
                      /> : null}
                  </span>
                  <span className="Aligner">
                    {!this.state.ForgotPassword ?
                      <span>
                        <span style={{clear: "both"}} className="Aligner">
                          <a href="javascript:" onClick={this.handleForgotPassword}>{window.appContent.ForgotPassword}</a>
                        </span>
                      </span>: null}
                  </span>
                </div>
            </Paper>
          </div>
        </div>
    );
  }
}

export default Login;
