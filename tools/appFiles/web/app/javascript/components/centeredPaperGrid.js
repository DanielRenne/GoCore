import React, {Component} from 'react';
import {Grid, Row, Col} from 'react-flexgrid-no-bs-conflict';
import Paper from 'material-ui/Paper';
import BaseComponent from './base';

class CenteredPaperGrid extends BaseComponent {
  constructor(props, context) {
    super(props, context);
  }
  render() {
    this.logRender();
    try {
      return (
          <Grid>
            <Row>
              <Col xs={12} md={12} sm={12} lg={6} smOffset={0} xsOffset={0} lgOffset={3} mdOffset={0}>
                <Paper style={{padding:'15px'}} zDepth={1}>
                  {this.props.children}
                </Paper>
              </Col>
            </Row>
          </Grid>
       );
    } catch(e) {
      return this.globs.ComponentError("CenteredPaperGrid", e.message, e);
    }
  }
}

export default CenteredPaperGrid;
