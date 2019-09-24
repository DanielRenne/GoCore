import React, {Component} from 'react';
import UserAddModify from './userModifyComponents'
import {BackPage} from '../../globals/forms';

class UserModify extends UserAddModify {
  constructor(props, context) {
    super(props, context);
  }
}

export default BackPage(UserModify);
