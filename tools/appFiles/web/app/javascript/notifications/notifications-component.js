import React, {Component} from 'react';
import {deepOrange300,red300, green300} from 'material-ui/styles/colors';
import BaseComponent from '../components/base'
import InfoNotification from '../components/infoNotification'
import PagePopup from '../components/pagePopup'

class Notifications extends BaseComponent {
  constructor(props, context) {
    super(props, context);

    this.notificationCallbackId;
    this.InfoPopupContent;

    this.state = {
      Open:false,
      Title:"",
      LinkTitle:"",
      Link:"",
      InfoPopupOpen:false,
      InfoPopupTitle:""
    };

    window.showNotification = (id, type, title, linkTitle, link, IPAddress, User) => {
      this.setComponentState({Id: id, Type: type, InfoPopupOpen:false, Open:true,Title:title,LinkTitle:linkTitle,Link:link, IPAddress: IPAddress, User: User});
    };

    window.showPagePopup = (title, content) => {
      this.InfoPopupContent = content;
      this.setComponentState({InfoPopupOpen:true, Open:false, InfoPopupTitle:title});
    }
  }

  componentDidMount() {
    var cb = (data) => {
      this.handleSocketBroadcast(data);
    };

    this.notificationCallbackId = window.api.registerSocketCallback(cb, "Notification");

  }

  componentWillUnmount() {
    window.api.unRegisterSocketCallback(this.notificationCallbackId);
  }

  handleSocketBroadcast(data) {

    if (window.location.hash == "" || window.appState.loggedIn == false) {
      return;
    }

    if ((data.Title != undefined && data.Title == "RoomSyncRequested") && window.location.hash.indexOf("#/roomControl") != -1) {
      return;
    }

    window.showNotification(data.Id, data.Type, data.Title, data.URLTitle, data.URL, data.IPAddress, data.User);
  }

  render() {
    this.logRender();

    return (
      <div>
        <InfoNotification
          uniqueId={this.state.Id}
          type={this.state.Type}
          open={this.state.Open}
          title={this.state.Title}
          IPAddress={this.state.IPAddress}
          User={this.state.User}
          linkTitle={this.state.LinkTitle}
          link={this.state.Link}
        />
        <PagePopup
          open={this.state.InfoPopupOpen}
          title={this.state.InfoPopupTitle}>
          {
            (this.InfoPopupContent === undefined) ? null : this.InfoPopupContent
          }
        </PagePopup>
      </div>
    );
  }
}

export default Notifications;
