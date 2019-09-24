/**
 * Created by Dan on 12/5/16.
 */
import {React,
    BaseComponent,
  } from "../../globals/forms";
import {deepOrange500} from "material-ui/styles/colors";
import CircularProgress from "material-ui/CircularProgress";
  
  
  class Loader extends BaseComponent {
    constructor(props, context) {
      super(props, context);

    }
  
    render() {
      try {
        return (
          <div><CircularProgress style={{width:15, height:15}} thickness={3.5} color={deepOrange500}/></div>
        );
      } catch(e) {
        return this.globs.ComponentError("Template", e.message, e);
      }
    }
  }

  
  export default Loader;
  