import React, {Component} from "react";
import BaseComponent from "../components/base";
import {IconButton, Avatar, grey900, blueGrey100, AppBar, Drawer, List, ListItem, deepOrange300} from "../globals/forms";
import {BackArrow, CommunicationCall, CommunicationCallEnd, AvMicOff, AvMic, NavigationClose, VolumeMute, VolumeUnMute, CameraIcon, DocumentationIcon, PhoneForwarded} from "../globals/icons";

class ContactList extends BaseComponent {
  constructor(props, context) {
    super(props, context);

		this.state = {
			ContactListIsOpen: false
    };

    this.open = () => {
      this.setComponentState({ContactListIsOpen: true})
    }

    this.close = () => {
      this.setComponentState({ContactListIsOpen: false})
    }

    this.handleSelectContact = (value) => {
    	this.props.DialContactNumber(value)
    }

  }

  render() {

    let defaultButton = {
    	ButtonType: ""
    }

    let contactList = [
      {First: "Joe", Last: "Alexander", Phone: "6468384627", Ext: "5634"},
      {First: "Margie", Last: "Baxter", Phone: "9256947262", Ext: ""},
      {First: "Bobby", Last: "Combs", Phone: "7167853937", Ext: "636"},
      {First: "Chad", Last: "Devito", Phone: "4325653634", Ext: ""},
      {First: "Tony", Last: "DiNozzo", Phone: "897475136", Ext: "12"},
      {First: "Donald", Last: "Duck", Phone: "3259782638", Ext: ""},
      {First: "Ziva", Last: "Gibbs", Phone: "5857397569", Ext: ""},
      {First: "Michelle", Last: "Holder", Phone: "6264953000", Ext: ""},
      {First: "Melanie", Last: "Louis", Phone: "6534543456", Ext: ""},
      {First: "Mary", Last: "Margaret", Phone: "6256356756", Ext: ""},
      {First: "Tim", Last: "Mallard", Phone: "3355757037", Ext: "444"},
      {First: "Abby", Last: "McGee", Phone: "6389562047", Ext: ""},
      {First: "Eric", Last: "Quinn", Phone: "6585893872", Ext: ""},
      {First: "Amy", Last: "Saprano", Phone: "2635659836", Ext: "234"},
      {First: "Katie", Last: "Smith", Phone: "5442345543", Ext: ""},
      {First: "Kevin", Last: "Spacey", Phone: "7546565454", Ext: ""},
      {First: "Elizabeth", Last: "Towney", Phone: "2421536387", Ext: ""},
      {First: "Daniel", Last: "Wahlberg", Phone: "4264537656", Ext: ""},
      {First: "Bob", Last: "Wellman", Phone: "5453245646", Ext: "1254"},
      {First: "John", Last: "Wick", Phone: "5183936586", Ext: ""}
    ];

// globals.buildButton = function(el, icon, iconText, label="", onClick=() => {}) {
  	let button1 = this.globs.buildButton(defaultButton, "", "1", "", () => this.appendToNumber("1"), this.state.DialerLabelColor)
  	let button2 = this.globs.buildButton(defaultButton, "", "2", "abc", () => this.appendToNumber("2"), this.state.DialerLabelColor)
  	let button3 = this.globs.buildButton(defaultButton, "", "3", "def", () => this.appendToNumber("3"), this.state.DialerLabelColor)
  	let buttonDel = this.globs.buildButton(defaultButton, <BackArrow/>, "", "del", () => this.appendToNumber("del"), this.state.DialerLabelColor)

  	try {
  		this.logRender();
  			return (
          <Drawer containerStyle={{zIndex:1501}} width={320} openSecondary={true} open={this.state.ContactListIsOpen}>
            <AppBar
                title={"Contact List"}
                iconElementLeft={<IconButton tooltip={window.appContent.Close} onClick={() => this.close()}><NavigationClose color={"white"}/></IconButton>}
            />
    				<div key={window.globals.guid()} className="Aligner">
    				<span className="Aligner" style={{minHeight:this.parentState.SourceSectionHeight}}>
              <div className="Aligner-item" style={{maxWidth:"none"}}>
    						{/*Search bar.*/}
                <List>
                {contactList.map((c, i) => {
                  let initials = c.First.charAt(0) + c.Last.charAt(0)
                    return (
                      <ListItem primaryText={c.Last + ", " + c.First}
                                secondaryText={c.Ext != "" ? c.Phone + " x" + c.Ext : c.Phone}
                                leftAvatar={<Avatar color={grey900} backgroundColor={deepOrange300} size={30} >{initials}</Avatar>}
                                onTouchTap={() => this.handleSelectContact(c.Phone + c.Ext)}
                                value = {c.Phone + c.Ext} />
                      )
                  })
                }
                </List>
    	      	</div>
    	      </span>
    	      </div>
          </Drawer>
  			)
    } catch(e) {
      session_functions.Dump(this.getClassName() + " errored with", e.message);
      return null;
      // return this.globs.ComponentError(this.getClassName(), e.message);
    }
  }


}

ContactList.propTypes = {
  parent: React.PropTypes.object
};

export default ContactList;