import {
    React,
    CenteredPaperGrid,
    BasePageComponent,
    TextField,
    RaisedButton,
    grey900,
    red500,
    green900
} from "../../globals/forms";
import Paper from 'material-ui/Paper';

class PasswordReset extends BasePageComponent {
  constructor(props, context) {
    super(props, context);
    //this.eventHandlerExample = (event) => {
      //window.api.post({action: "Warning", state: this.state});
    //};
  }

  componentDidUpdate() {
    if (window.appState.DeveloperLogState) {
      console.log("componentDidUpdate", this.state);
    }
  }

  render() {
    this.logRender();

    var pass8Chars = red500
    var passInteger = red500
    var passLower = red500
    var passUpper = red500
    var passSpecial = red500
    var passMatch = red500

    if(this.state.Password != undefined && this.state.Password != "") {
      this.state.Password.length >= 8 ? (pass8Chars = green900) : (pass8Chars = red500)
      this.state.Password.search(/\d/) != -1 ? (passInteger = green900) : (passInteger = red500)
      this.state.Password.search(/[a-z]/) != -1 ? (passLower = green900) : (passLower = red500)
      this.state.Password.search(/[A-Z]/) != -1 ? (passUpper = green900) : (passUpper = red500)
      this.state.Password.search(/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]+/) != -1 ? (passSpecial = green900) : (passSpecial = red500)
      this.state.Password === this.state.ConfirmPassword ? (passMatch = green900) : (passMatch = red500)
    }

    return (
      <div>
      <div className="Align" style={{maxHeight:600}}>
        <Paper style={this.style} className={($(window).width() > 400) ? 'Align-center' : ''} rounded={true} zDepth={3}>
          <span className="AlignerRight">
            <div></div>
            <h1 className="Paper-Title-Alt">{window.pageContent.Title}</h1>
            <div></div>
          </span>
          <div style={{paddingLeft: 30, paddingRight: 30, marginTop:10}}>
          {
            (this.state.PasswordReset.Complete === true) ?
            <div>
              <br />
              {window.pageContent.ResetCompleted}
              <br />
              <br />
              <span className="AlignerRight">
              <div></div>
              <RaisedButton
                label={window.appContent.OK}
                onTouchTap={() => {
                  window.api.post({action: "Reset", state: this.state});
                }}
                secondary={true}
              />
              <div></div>
              </span>
              <br />
              <br />
            </div>
            :
            <div>
            <br />
            <span className ="AlignerLeft">
              <span>
              <TextField
                floatingLabelText={"* " + window.pageContent.NewPassword}
                hintText={"* " + window.pageContent.NewPassword}
                fullWidth={false}
                type="password"
                onChange={(e) => {this.setComponentState({Password:e.target.value, PasswordErrors:""})}}
                errorText={window.pageContent[this.state.PasswordErrors]}
              />
              <br />
              <TextField
                floatingLabelText={"* " + window.pageContent.ConfirmPassword}
                hintText={"* " + window.pageContent.ConfirmPassword}
                fullWidth={false}
                type="password"
                onChange={(e) => {this.setComponentState({ConfirmPassword:e.target.value, ConfirmPasswordErrors:""})}}
                onKeyDown={(event) => {
                  if (event.nativeEvent.key == "Enter"){
                    window.api.post({action: "Reset", state: this.state});
                  }
                }}
                errorText={window.pageContent[this.state.ConfirmPasswordErrors]}
              />
              </span>
              <span style={{marginLeft: 10}}>
                <label style={{color: grey900}}>Password criteria:</label><br />
                <label style={{color: pass8Chars, marginLeft: 40, marginRight: 100}}>8 Characters</label><br />
                <label style={{color: passInteger, marginLeft: 40}}>1 Integer</label><br />
                <label style={{color: passLower, marginLeft: 40}}>1 Lowercase</label><br />
                <label style={{color: passUpper, marginLeft: 40}}>1 Uppercase</label><br />
                <label style={{color: passSpecial, marginLeft: 40}}>1 Special Character</label><br />
                <label style={{color: passMatch, marginLeft: 40}}>{this.state.Password != "" ? this.state.Password === this.state.ConfirmPassword ? "Passwords Match!" : "Passwords do not match.":""}</label><br />
              </span>
            </span>
            <br />
            <br />
            <span className="AlignerRight">
            <div></div>
            <RaisedButton
              label={window.appContent.OK}
              onTouchTap={() => {
                window.api.post({action: "Reset", state: this.state});
              }}
              secondary={true}
            />
            <div></div>
            </span>
              <br />
            </div>
          }



          </div>
        </Paper>
      </div>
      {
        (window.appState.AccountName == "") ?
        <div className="Aligner">
          <img style={{width:400}} src="/web/app/images/go_core_app_product.png"/>
        </div>
        : null
      }
    </div>
    );
  }
}

export default PasswordReset;
