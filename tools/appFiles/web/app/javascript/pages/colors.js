import getMuiTheme from 'material-ui/styles/getMuiTheme';
import {fade} from 'material-ui/utils/colorManipulator';

import {deepOrange600, blueGrey900, blueGrey400, blueGrey600, indigo600, grey500, grey100, green500, greenA400} from 'material-ui/styles/colors';
var MuiThemes = {};

MuiThemes.default = getMuiTheme({
  palette: {
    accent1Color: "#f06837",
    accent2Color: indigo600,
    accent3Color: blueGrey600,
    primary1Color: blueGrey900,
    primary2Color: indigo600,
    primary3Color: grey500,
    pickerHeaderColor: "#f06837",
    textColor: blueGrey900,
  },
  textField: {
    textColor: blueGrey900,
    hintColor: blueGrey400,
    floatingLabelColor: blueGrey900,
    focusColor: blueGrey900,
  }
});

MuiThemes.default.tableRow.height = 60;
MuiThemes.default.tableRowColumn.height = 60;
MuiThemes.default.slider.trackColor = 'white';
MuiThemes.default.slider.selectonColor = deepOrange600;
MuiThemes.default.slider.trackSize = 4;
MuiThemes.default.toggle.thumbOnColor = green500;
MuiThemes.default.toggle.thumbOffColor = grey100;
MuiThemes.default.toggle.trackOnColor = fade(MuiThemes.default.toggle.thumbOnColor, 0.5);
MuiThemes.default.toggle.trackOffColor = grey100;


MuiThemes.opposite = getMuiTheme({
  palette: {
    accent1Color: blueGrey900,
    accent2Color: indigo600,
    accent3Color: blueGrey600,
    primary1Color: "#f06837",
    primary2Color: indigo600,
    primary3Color: grey500,
    pickerHeaderColor: "#f06837",
  },
  textField: {
    textColor: blueGrey900,
    hintColor: blueGrey400,
    floatingLabelColor: blueGrey900,
    focusColor: blueGrey900,
  }
});
MuiThemes.opposite.tableRow.height = 60;
MuiThemes.opposite.tableRowColumn.height = 60;
MuiThemes.opposite.slider.trackColor = 'white';
MuiThemes.opposite.slider.selectonColor = "#f06837";
MuiThemes.opposite.slider.trackSize = 4;
MuiThemes.opposite.toggle.thumbOnColor = green500;
MuiThemes.opposite.toggle.thumbOffColor = grey100;
MuiThemes.opposite.toggle.trackOnColor = fade(MuiThemes.default.toggle.thumbOnColor, 0.5);
MuiThemes.opposite.toggle.trackOffColor = grey100;

export default MuiThemes;