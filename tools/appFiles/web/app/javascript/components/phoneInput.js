import React, { PropTypes, Component } from 'react';
import BaseComponent from './base';
import TextField from 'material-ui/TextField';

class PhoneInput extends BaseComponent {
    constructor(props, context) {
      super(props, context);
      this.intlInput;
      if (this.props.Label) {
        this.Id = this.props.Label.replace(/ /gi, "_");
      } else {
        this.Id = "errr";
      }
      this.state = {
        DefaultCountry: "us",
        Value: props.InitialValue,
        Numeric: null,
        DialCode: null,
        CountryISO: null
      };
  }

  changeCountry(iso) {
    this.setComponentState({DefaultCountry: iso}, () => this.clarValue());
  }

  clarValue() {
    this.setComponentState({Value:"", Numeric: null, DialCode: null, CountryISO: null});
  }

  getRef() {
    return this.intlInput;
  }

  render() {
    try {
      this.logRender();
      return (
        <div>
          <TextField
            ref={(c) => this.intlInput = c}
            floatingLabelText={this.props.Label}
            hintText={this.props.Label}
            defaultValue={this.state.Value}
            fullWidth={true}
            errorText={this.props.ErrorText}
            onChange={(event) => {
              let state = {
                Value: event.target.value,
                Numeric: event.target.value,
                DialCode: null,
                CountryISO: null
              };
              this.setComponentState(state);
              this.props.OnChange(state);
            }}
          />
        </div>
      );
    } catch(e) {
      return this.globs.ComponentError("PhoneInput", e.message, e);
    }
  }
}

PhoneInput.propTypes = {
  InitialValue: React.PropTypes.string,
  Disabled: React.PropTypes.bool,
  OnChange: React.PropTypes.func.isRequired,
  ErrorText: React.PropTypes.string,
  Label: React.PropTypes.string.isRequired,
};

PhoneInput.defaultProps = {
  Disabled: false,
  InitialValue: "",
  ErrorText: "",
};


export default PhoneInput;
