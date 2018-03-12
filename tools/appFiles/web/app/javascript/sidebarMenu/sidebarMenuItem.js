import {
  React,
  BaseComponent,
  IconButton} from '../globals/forms';
import {deepOrange500, blueGrey500, deepOrange400, deepOrange300} from 'material-ui/styles/colors';
import {RoomIcon} from '../icons/icons'
import {ControlRoomIcon, Play} from '../globals/icons'

class SideBarMenuItem extends BaseComponent {
  constructor(props, context) {
    super(props, context);
    this.PageSize = 10;

    this.getItemsShow = (currentPage, items) => {
      return items !== undefined ? items.slice(currentPage*this.PageSize, (currentPage*this.PageSize)+this.PageSize) : [];
    }

    this.state = {
      Title: this.__(this.props.title),
      URL: this.props.url,
      Items: globals.sortByKey(this.props.items, "Title"),
      Expanded: this.props.expanded,
      Icon: this.props.icon,
      RightIcon: this.props.rightIcon,
      Hidden: this.props.hidden,
      RightIconLink: this.props.rightIconLink,
      Selected: this.props.selected,
      ShowPaging: false,
      ItemsShow: this.getItemsShow(0, globals.sortByKey(this.props.items, "Title")),
      Pages: this.props.items!==null? Math.ceil(this.props.items.length/this.PageSize) : 0,
      CurrentPage: 0,
      PrevDisabled: true,
      NextDisabled: false
    };

    this.toggleSubMenu = () => {
      if (this.state.Expanded) {
        this.setComponentState({Expanded: false});
      } else {
        this.setComponentState({Expanded: true});
      }
    };

    this.prevItem = () => {
      var currentPage = ((this.state.CurrentPage - 1) >= 0) ? (this.state.CurrentPage - 1):this.state.CurrentPage;
      var itemsShow = this.getItemsShow(currentPage, this.state.Items);
      var prevDisabled = currentPage > 0 ? false : true;

      this.setComponentState({
        CurrentPage: currentPage,
        ItemsShow: itemsShow,
        PrevDisabled: prevDisabled,
        NextDisabled: currentPage < this.state.Pages ? false : true,
      });
    };

    this.nextItem = () => {
      var currentPage = ((this.state.CurrentPage + 1) < this.state.Pages) ? (this.state.CurrentPage + 1):this.state.CurrentPage;
      var itemsShow = this.getItemsShow(currentPage, this.state.Items);
      var nextDisabled = currentPage !== this.state.Pages - 1 ? false : true;

      this.setComponentState({
        CurrentPage: currentPage,
        ItemsShow: itemsShow,
        PrevDisabled: currentPage > 0 ? false : true,
        NextDisabled: nextDisabled
      });
    };
  }

  render() {
    this.logRender();

    var displayStyle = "list-item";
    var icon = window.globals.fetchRemarkIcon(this.state.Icon);
    var iconClass = "site-menu-icon " + icon;
    var subItems = "";
    var subItemsArrow = "";
    var anchorItems = [];

    var itemTitle = (<span  key={"main-span" + this.state.Title} className="site-menu-title">{window.globals.translate( (this.state.Title == "-99999") ? "AllTitle" : this.state.Title  )}</span>);
    var itemIcon = (<i style={{paddingLeft:(this.props.depth > 1) ?  this.props.depth * 10 : 0, paddingRight:5}} key={"main-icon" + this.state.Title} className={iconClass} aria-hidden="true"></i>);

    if (this.state.Icon == "Room") {
      itemTitle = (<span style={{marginTop:-6}}  key={"main-span" + this.state.Title} className="site-menu-title">{window.globals.translate( (this.state.Title == "-99999") ? "AllTitle" : this.state.Title  )}</span>);
      itemIcon = (<i style={{paddingLeft:(this.props.depth > 1) ?  this.props.depth * 10 : 0, paddingRight:5, width:18}} key={"main-icon" + this.state.Title}  aria-hidden="true">
                    <RoomIcon width={18} height={18} color="#757575" />
                  </i>);
    }

    anchorItems.push(itemIcon);
    anchorItems.push(itemTitle);

    if (this.state.Icon == "Room" && this.state.RightIconLink != "") {
        var rightIcon = ( <span key={window.globals.guid()} style={{position:"absolute",
                                 right: 5,
                                 display: "inline-block",
                                 verticalAlign: "middle"}}>
                                 <IconButton style={{padding:0, width:44, height:40}} onClick={(event) =>  {eval(this.state.RightIconLink); event.preventDefault(); }}><ControlRoomIcon color={deepOrange300}/></IconButton>
                    </span>);
          anchorItems.push(rightIcon);
    }

    if (this.state.RightIcon != ""){
          var rightIconKey = window.globals.fetchRemarkIcon(this.state.RightIcon);
          var rightIcon = ( <span key={window.globals.guid()} style={{position:"absolute",
                                         right: 25,
                                         display: "inline-block",
                                         verticalAlign: "middle"}}>
                                         <i className={rightIconKey} aria-hidden="true"></i>
                            </span>);
          anchorItems.push(rightIcon);
    }

    if (this.state.Hidden){
      displayStyle = "none";
    }

    var anchorStyle = {};
    if (this.state.Selected) {
      anchorStyle = {color:deepOrange500};
    }

    var itemAnchor =
    (
      <a style={anchorStyle} href={this.state.URL} key={"a-" + this.state.Title}>{
        anchorItems.map((aItem) =>{
          return aItem;
        })
      }
      </a>
    );

    var expandSubItems = "none";

    if (this.state.Items != undefined && this.state.Items.length > 0){
      if (this.state.Expanded){
        expandSubItems = "block";
      }

      var itemArrow = (<span key={"site-menu-arrow-" + this.state.Title} className="site-menu-arrow"></span>);
      anchorItems.push(itemArrow);
      //
      itemAnchor = (
        <a  href="javascript:void(0)" key={"a-click-" + this.state.Title} onClick={this.toggleSubMenu}>{
          anchorItems.map((aItem) =>{
            return aItem;
          })
        }
        </a>
      );

      // var sortedSubItems = globals.sortByKey(this.state.Items, "Title");
      
      this.state.ShowPaging = this.state.Items.length > this.PageSize;
      var paging;
      if (this.state.ShowPaging){
        paging = (
          <li key={"ul-paging-"+this.state.CurrentPage} className="site-menu-item" style={{display:displayStyle}}>
            <div style={{paddingLeft:35}}>
              <IconButton style={{padding:0, width:30, height:30}} onClick={this.prevItem} disabled={this.state.PrevDisabled}>
                <Play color={blueGrey500} className="rotate180"/>
                </IconButton>

              <IconButton style={{padding:0, width:30, height:30}} onClick={this.nextItem} disabled={this.state.NextDisabled}>
                <Play color={blueGrey500}/>
                </IconButton>
            </div>
          </li>
        )
      };
      
      subItems = (
        <ul key={"ul-" + this.state.Title} className="site-menu-sub" style={{display:expandSubItems}}>{
          this.state.ItemsShow.map((item, k) => {
            var subItem = item;
            return <SideBarMenuItem key={this.state.CurrentPage+"-"+k}
                                    title={item.Title}
                                    url={item.URL}
                                    expanded={item.Expanded}
                                    items={item.Items}
                                    icon={item.Icon}
                                    hidden={item.Hidden}
                                    depth={this.props.depth + 1}
                                    rightIcon={item.RightIcon}
                                    rightIconLink={item.RightIconLink}
                                    rightIconFunc={item.RightIconFunc}
                                    selected={item.Selected}/>
          })
        }
        {paging}
        </ul>
      )
    }

    var liItems = [];
    liItems.push(itemAnchor);
    liItems.push(subItems);

    var key = this.state.Title + this.state.URL;

    return (
      <li key={key} className="site-menu-item" style={{display:displayStyle}}>
        {
          liItems.map((liItem) =>{
            return liItem;
          })
        }
      </li>
    );
  }
}

SideBarMenuItem.propTypes = {
  items: React.PropTypes.array,
  title: React.PropTypes.string,
  url: React.PropTypes.string,
  rightIcon: React.PropTypes.string,
  rightIconLink: React.PropTypes.string,
  expanded: React.PropTypes.bool,
  hidden: React.PropTypes.bool,
  selected: React.PropTypes.bool,
  depth: React.PropTypes.number
};

SideBarMenuItem.defaultProps = {
    items: []
};

export default SideBarMenuItem;
