import React, { PropTypes, Component } from 'react';
import IntlTelInput from 'react-intl-tel-input';
import '../../node_modules/react-intl-tel-input/dist/main.css';
import BaseComponent from './base';

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

  onChange(event, phoneWithoutDialCode, countryData, numberFormatted) {
      var state = {
        Value: numberFormatted,
        Numeric: (countryData.hasOwnProperty("dialCode") ? numberFormatted.replace("+" + countryData.dialCode + " ", "").replace(/[^0-9]/g, ""): numberFormatted),
        DialCode: (countryData.hasOwnProperty("dialCode") ? countryData.dialCode: null),
        CountryISO: (countryData.hasOwnProperty("iso2") ? countryData.iso2: null),
      };
      this.setComponentState(state);
      this.props.OnChange(state);
  }

  render() {
    try {
      this.logRender();
      return (
        <div>
          <label htmlFor={this.Id} style={{height: 35, paddingTop: 10, fontSize: 12, display: 'block', cursor: 'text', color: 'rgba(0, 0, 0, 0.498039)'}}>{this.props.Label}</label>
          <IntlTelInput
            ref={(c) => this.intlInput = c}
            id={this.Id}
            style={{border: (this.props.ErrorText) ? "1px solid rgb(244, 67, 54)": "none"}}
            defaultCountry={this.state.DefaultCountry}
            value={this.state.Value}
            disabled={this.props.Disabled}
            separateDialCode={false}
            nationalMode={false}
            autoHideDialCode={false}
            css={['intl-tel-input', 'form-control']}
            utilsScript={'/dist/javascript/libphonenumber.js.gz'}
            onPhoneNumberChange={(...params) => this.onChange(...params)}
          />
          {(this.props.ErrorText) ? <div style={{height: 35, paddingTop: 10, fontSize: 12, display: 'block', cursor: 'text', color: 'rgb(244, 67, 54)'}}>{this.props.ErrorText}</div>: null}
        </div>
      );
    } catch(e) {
      return this.globs.ComponentError(this.getClassName(), e.message);
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
