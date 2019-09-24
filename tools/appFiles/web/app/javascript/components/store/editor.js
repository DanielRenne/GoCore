/**
 * Created by Dan on 12/5/16.
 */
import {React,
    BaseComponent,
    Editor,
    EditorState,
    RichUtils,
    DraftJsStyleMap,
    DraftJsGetBlockStyle,
    DraftJsBlockStyleControls,
    DraftJsInlineStyleControls,
    ContentState,
    convertFromHTML,
    stateToHTML
  } from "../../globals/forms";
  import Loader from "./loader";

  class EditorStoreComponent extends BaseComponent {
    constructor(props, context) {
      super(props, context);

      this.unmounted = false;
      this.state = {
        value: undefined
      }

      this.state.editorState = EditorState.createEmpty();

      this.draftJSFocus = () => this.refs.editor.focus();
      this.draftJSHandleKeyCommand = (command) => this._DraftJsHandleKeyCommand(command);
      this.draftJSOnTab = (e) => this._DraftJsOnTab(e);
      this.draftJSToggleBlockType = (type) => this._DraftJstoggleBlockType(type);
      this.draftJSToggleInlineStyle = (style) => this._DraftJsToggleInlineStyle(style);

      this.handleDraftJsChange = (editorState) => {
        window.displayErrorAlert = false;
        this.setState({editorState}, () => {
          const {editorState} = this.state;
          var contentState = editorState.getCurrentContent();
          this.store.set(this.props.collection, this.props.id, "InfoPopup", stateToHTML(contentState));
        });
      };
    }

    componentDidMount() {
      this.subscriptionId = this.store.subscribe(this.props.collection, this.props.id, this.props.path,(data) => {this.handleValueChange(data)}, true);
    }

    componentWillUnmount() {
      this.unmounted = true;
      this.store.unsubscribe(this.subscriptionId);
    }

    handleValueChange(data) {
      if (this.state.value != data && !this.unmounted) {
        this.setState({value:data}, () => {
          if (!this.state.editorState) {
            if (this.state.value && this.state.value.length > 0) {
              const state = ContentState.createFromBlockArray(convertFromHTML(this.state.value));
              this.state.editorState = EditorState.createWithContent(state);
            } else {
              this.state.editorState = EditorState.createEmpty();
            }
          }
        });
      }
    }

    render() {
      try {

        const {editorState} = this.state;

        // If the user changes block type before entering any text, we can
        // either style the placeholder or hide it. Let's just hide it now.
        let className = 'RichEditor-editor';
        var contentState = editorState.getCurrentContent();
        if (!contentState.hasText()) {
          if (contentState.getBlockMap().first().getType() !== 'unstyled') {
            className += ' RichEditor-hidePlaceholder';
          }
        }

        return (
          <span>
            {this.state.value == undefined? <Loader/>: null}
            <span style={{display: this.state.value == undefined ? "none" : "block"}}>
              <div>
                <span style={{color: 'rgba(0,0,0,0.498039)', fontSize: 12}}>{window.appContent.RoomAddEditRoomInfo}:</span>
                <br/>
                <br/>
                <div className="RichEditor-root">
                  <DraftJsBlockStyleControls
                    editorState={editorState}
                    onToggle={this.draftJSToggleBlockType}
                  />
                  <DraftJsInlineStyleControls
                    editorState={editorState}
                    onToggle={this.draftJSToggleInlineStyle}
                  />
                  <div className={className} onClick={this.draftJSFocus}>
                    <Editor
                      blockStyleFn={DraftJsGetBlockStyle}
                      customStyleMap={DraftJsStyleMap}
                      editorState={editorState}
                      handleKeyCommand={this.draftJSHandleKeyCommand}
                      onChange={this.handleDraftJsChange}
                      onTab={this.draftJSOnTab}
                      ref="editor"
                      placeHolder={this.state.value != "" ? window.appContent.RoomAddEditDescribeRoomFeatures : ""}
                      spellCheck={true}
                    />
                  </div>
                </div>
              </div>
            </span>
          </span>
        );
      } catch(e) {
        return this.globs.ComponentError("EditorStore", e.message, e);
      }
    }
  }


  EditorStoreComponent.propTypes = {
    collection:React.PropTypes.string,
    id:React.PropTypes.string,
    path:React.PropTypes.string
  };

  export default EditorStoreComponent;
