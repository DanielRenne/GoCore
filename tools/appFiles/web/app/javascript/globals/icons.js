import React from 'react';
import Icon from './icon';

// Action items: delete, add, edit, copy
import ActionBuild from 'material-ui/svg-icons/action/build';
import ContentUndo from 'material-ui/svg-icons/content/undo';
import ActionDelete from 'material-ui/svg-icons/action/delete';
import ActionAdd from 'material-ui/svg-icons/content/add';
import EditIcon from 'material-ui/svg-icons/image/edit';
import ImageAddAPhoto from 'material-ui/svg-icons/image/add-a-photo';
import ContentContentCopy from 'material-ui/svg-icons/content/content-copy';

import CommunicationCallMade from 'material-ui/svg-icons/communication/call-made';
import SpeakerPhone from 'material-ui/svg-icons/communication/speaker-phone';
// Navigation, menu, volume, mute, lock, hamburger/menu icons
import AvVolumeUp from 'material-ui/svg-icons/av/volume-up';
import AudioIcon from 'material-ui/svg-icons/av/equalizer';
import ActionLock from 'material-ui/svg-icons/action/lock';
import VideoDistributionIcon from 'material-ui/svg-icons/action/settings-input-hdmi';
import DisplayIcon from 'material-ui/svg-icons/hardware/tv';
import TouchPadIcon from 'material-ui/svg-icons/hardware/tablet-mac';
import ConferencingIcon from 'material-ui/svg-icons/av/video-call';
import NewRelease from 'material-ui/svg-icons/av/new-releases';
import CameraIcon from 'material-ui/svg-icons/av/videocam';
import ControlIcon from 'material-ui/svg-icons/image/transform';
import MediaIcon from 'material-ui/svg-icons/action/theaters';
import LinkIcon from 'material-ui/svg-icons/content/link';
import ReportProblemIcon from 'material-ui/svg-icons/action/report-problem';
import PowerIcon from 'material-ui/svg-icons/action/power-settings-new';
import ActionSettingsApplications from 'material-ui/svg-icons/action/settings-applications';
import ActionSettingsInputAntenna from 'material-ui/svg-icons/action/settings-input-antenna';
import ActionSettingsInputComponent from 'material-ui/svg-icons/action/settings-input-component';
import ActionSettingsInputHdmi from 'material-ui/svg-icons/action/settings-input-hdmi';
import ActionSettingsInputSvideo from 'material-ui/svg-icons/action/settings-input-svideo';
import LockOpen from 'material-ui/svg-icons/action/lock-open';
import LockClosed from 'material-ui/svg-icons/action/lock-outline';
import VolumeUp from 'material-ui/svg-icons/av/volume-up';
import VolumeDown from 'material-ui/svg-icons/av/volume-down';
import VolumeMute from 'material-ui/svg-icons/av/volume-off';
import VolumeUnMute from 'material-ui/svg-icons/av/volume-mute';
import HamburgerBar from 'material-ui/svg-icons/image/dehaze';
import ArrowBack from 'material-ui/svg-icons/navigation/arrow-back';
import ArrowForward from 'material-ui/svg-icons/navigation/arrow-forward';
import Circle from 'material-ui/svg-icons/image/brightness-1';
import SpaceBar from 'material-ui/svg-icons/editor/space-bar';
import HorizontalLine from 'material-ui/svg-icons/content/remove';
import AccountEnterIcon from 'material-ui/svg-icons/av/library-books';
import NavigationClose from 'material-ui/svg-icons/navigation/close';
import Explore from 'material-ui/svg-icons/action/explore';
import Presets from 'material-ui/svg-icons/device/widgets'

import AvVideocam from 'material-ui/svg-icons/av/videocam';
import AvVideocamOff from 'material-ui/svg-icons/av/videocam-off';
import AvMic from 'material-ui/svg-icons/av/mic';
import AvMicOff from 'material-ui/svg-icons/av/mic-off';
import AvMicNone from 'material-ui/svg-icons/av/mic-none';
import AvMovie from 'material-ui/svg-icons/av/movie';
import AvHighQuality from 'material-ui/svg-icons/av/high-quality';
import AvHd from 'material-ui/svg-icons/av/hd';
import AvGames from 'material-ui/svg-icons/av/games';
import AvForward5 from 'material-ui/svg-icons/av/forward-5';
import AvForward10 from 'material-ui/svg-icons/av/forward-10';
import AvForward30 from 'material-ui/svg-icons/av/forward-30';
import AvFiberDvr from 'material-ui/svg-icons/av/fiber-dvr';
import AvClosedCaption from 'material-ui/svg-icons/av/closed-caption';
import AvEqualizer from 'material-ui/svg-icons/av/equalizer';
import HardwareSpeakerGroup from 'material-ui/svg-icons/hardware/speaker-group'
import HardwareSpeaker from 'material-ui/svg-icons/hardware/speaker'

// Bluray/cable box Button icons: playback, replays, directionals, select, select all,
import Pause from 'material-ui/svg-icons/av/pause';
import Play from 'material-ui/svg-icons/av/play-arrow';
import Stop from 'material-ui/svg-icons/av/stop';
import FastForward from 'material-ui/svg-icons/av/fast-forward';
import Rewind from 'material-ui/svg-icons/av/fast-rewind';
import SkipNext from 'material-ui/svg-icons/av/skip-next';
import SkipPrevious from 'material-ui/svg-icons/av/skip-previous';
import ShuffleIcon from 'material-ui/svg-icons/av/shuffle';
import LoopIcon from 'material-ui/svg-icons/av/loop';
import Record from 'material-ui/svg-icons/av/fiber-manual-record';
import MicOn from 'material-ui/svg-icons/av/mic';
import MicHollow from 'material-ui/svg-icons/av/mic-none';
import MicOff from 'material-ui/svg-icons/av/mic-off';
import Replay from 'material-ui/svg-icons/av/replay';
import Replay5 from 'material-ui/svg-icons/av/replay-5';
import Replay10 from 'material-ui/svg-icons/av/replay-10';
import Replay30 from 'material-ui/svg-icons/av/replay-30';
import Left from 'material-ui/svg-icons/navigation/chevron-left';
import Right from 'material-ui/svg-icons/navigation/chevron-right';
import Up from 'material-ui/svg-icons/navigation/expand-less';
import Down from 'material-ui/svg-icons/navigation/expand-more';
import Eject from 'material-ui/svg-icons/action/eject';
import Select from 'material-ui/svg-icons/device/gps-not-fixed';
import SelectAll from 'material-ui/svg-icons/content/select-all';

// Ethernet
import Ethernet from 'material-ui/svg-icons/action/settings-ethernet';
import CommunicationImportExport from 'material-ui/svg-icons/communication/import-export';
import Reconnect from 'material-ui/svg-icons/av/loop';

// directionals
import UpArrow from 'material-ui/svg-icons/hardware/keyboard-arrow-up';
import DownArrow from 'material-ui/svg-icons/hardware/keyboard-arrow-down';
import RightArrow from 'material-ui/svg-icons/hardware/keyboard-arrow-right';
import Keyboard from 'material-ui/svg-icons/hardware/keyboard';
import ActionLightbulbOutline from 'material-ui/svg-icons/action/lightbulb-outline';

// sources
import ContentBlock from 'material-ui/svg-icons/content/block';
import Computer from 'material-ui/svg-icons/hardware/computer';
import VideoGame from 'material-ui/svg-icons/hardware/videogame-asset';
import ControlRoomIcon from 'material-ui/svg-icons/action/touch-app';
import {
  UserRevokeIcon,
  DVDIcon,
  BluRayIcon,
  AppleTVIcon,
  PlayPause,
  VGA,
  RemoteIcon,
  ExitIcon,
  ConferencingIcon2,
  ProjectorIcon,
  WhiteboardIcon,
  MuteVideo,
  UnMuteVideo
} from '../icons/icons';

// Plus/minus
import Add from 'material-ui/svg-icons/content/add';
import Minus from 'material-ui/svg-icons/content/remove';
import FileCloudDone from 'material-ui/svg-icons/file/cloud-done';

// Document Camera Buttons: Brightness, zoom, pan, rotate, input, menu
import ZoomIn from 'material-ui/svg-icons/action/zoom-in';
import ZoomOut from 'material-ui/svg-icons/action/zoom-out';
import BrightnessDown from 'material-ui/svg-icons/image/brightness-7';
import BrightnessUp from 'material-ui/svg-icons/image/brightness-5';
import Rotate from 'material-ui/svg-icons/image/rotate-90-degrees-ccw';
import PhotoCamera from 'material-ui/svg-icons/image/photo-camera';
import Pan from 'material-ui/svg-icons/action/open-with';
import Mask from 'material-ui/svg-icons/action/chrome-reader-mode';
import Input from 'material-ui/svg-icons/action/input';
import Menu from 'material-ui/svg-icons/navigation/menu';
import PIP from 'material-ui/svg-icons/action/picture-in-picture-alt';

// back arrow
import BackArrow from 'material-ui/svg-icons/hardware/keyboard-backspace';
import Home from 'material-ui/svg-icons/action/home';

import FileFileDownload from 'material-ui/svg-icons/file/file-download';
import ActionSearch from 'material-ui/svg-icons/action/search';
import ActionInfoOutline from 'material-ui/svg-icons/action/info-outline';
import ImageAudiotrack from 'material-ui/svg-icons/image/audiotrack';

//Screen icons
import Screens from 'material-ui/svg-icons/hardware/devices-other';
import ImageBackground from 'material-ui/svg-icons/image/image';
import Save from 'material-ui/svg-icons/content/save';

import ScreenShadeIcon from 'material-ui/svg-icons/editor/format-line-spacing';
import SocialPerson from 'material-ui/svg-icons/social/person';
import MountIcon from 'material-ui/svg-icons/editor/vertical-align-center';

import DeviceDevices from 'material-ui/svg-icons/device/devices';
import AvFiberSmartRecord from 'material-ui/svg-icons/av/fiber-smart-record';
import HardwareLaptopChromebook from 'material-ui/svg-icons/hardware/laptop-chromebook';
import HardwareMemory from 'material-ui/svg-icons/hardware/memory';
import AvRadio from 'material-ui/svg-icons/av/radio';
import AvSurroundSound from 'material-ui/svg-icons/av/surround-sound';
import HardwareToys from 'material-ui/svg-icons/hardware/toys';
import NotificationDiscFull from 'material-ui/svg-icons/notification/disc-full';
import AvLibraryMusic from 'material-ui/svg-icons/av/library-music';
import AvMusicVideo from 'material-ui/svg-icons/av/music-video';
import AvQueueMusic from 'material-ui/svg-icons/av/queue-music';
import HardwareDeveloperBoard from 'material-ui/svg-icons/hardware/developer-board';
import AvSubscriptions from 'material-ui/svg-icons/av/subscriptions';
import DeviceDvr from 'material-ui/svg-icons/device/dvr';
import HardwareDeviceHub from 'material-ui/svg-icons/hardware/device-hub';
import NotificationLiveTv from 'material-ui/svg-icons/notification/live-tv';
import HardwareDesktopMac from 'material-ui/svg-icons/hardware/desktop-mac';
import HardwareDesktopWindows from 'material-ui/svg-icons/hardware/desktop-windows';
import Security from 'material-ui/svg-icons/hardware/security';
import Email from 'material-ui/svg-icons/communication/mail-outline';
import Settings from 'material-ui/svg-icons/action/settings';

// Extras
import ActionDescription from 'material-ui/svg-icons/action/description';
import CommunicationImportContacts from 'material-ui/svg-icons/communication/import-contacts';
import Freeze from 'material-ui/svg-icons/places/ac-unit';
import ScreenShare from 'material-ui/svg-icons/communication/screen-share';
import StopScreenShare from 'material-ui/svg-icons/communication/stop-screen-share';
import EmptyCheckBox from 'material-ui/svg-icons/image/crop-square';
import ImageNegative from 'material-ui/svg-icons/image/crop-original';
import Image from 'material-ui/svg-icons/image/image';
import Sync from 'material-ui/svg-icons/notification/sync';
import FileDownload from 'material-ui/svg-icons/file/file-download';
import FileUpload from 'material-ui/svg-icons/file/file-upload';
import Fullscreen from 'material-ui/svg-icons/navigation/fullscreen';
import FullscreenExit from 'material-ui/svg-icons/navigation/fullscreen-exit';
import CD from 'material-ui/svg-icons/av/album';
import Mouse from 'material-ui/svg-icons/hardware/mouse';
import Favorite from 'material-ui/svg-icons/action/favorite';
import List from 'material-ui/svg-icons/action/list';
import USB from 'material-ui/svg-icons/device/usb';
import Snooze from 'material-ui/svg-icons/av/snooze';
import Help from 'material-ui/svg-icons/action/help';
import SimCard from 'material-ui/svg-icons/hardware/sim-card';
import LooksOne from 'material-ui/svg-icons/image/looks-one';
import LooksTwo from 'material-ui/svg-icons/image/looks-two';
import CommunicationCall from 'material-ui/svg-icons/communication/call';
import License from 'material-ui/svg-icons/communication/vpn-key';
import SettingsRemote from 'material-ui/svg-icons/action/settings-remote';
import Cloud from 'material-ui/svg-icons/file/cloud';
import CommunicationCallEnd from 'material-ui/svg-icons/communication/call-end';
import PhoneForwarded from 'material-ui/svg-icons/notification/phone-forwarded';
import CheckCircle from 'material-ui/svg-icons/action/check-circle';
import Rotate360 from 'material-ui/svg-icons/action/autorenew';
import WorldIcon from 'material-ui/svg-icons/social/public';


// with all of these icons. do not put the category in the name of the icon.  It messes up the padding and does strange stuff
var icon = "add_a_photo";
const addPhotoIcon = <Icon iconTag={icon}/>;
const addPhotoIconLarge = <Icon large={true} iconTag={icon}/>;

var icon = "power_settings_new";
const serverSettingsIcon = <Icon iconTag={icon}/>;
const serverSettingsIconLarge = <Icon large={true} iconTag={icon}/>;

icon = "supervisor_account";
const userListIcon = <Icon iconTag={icon}/>;
const userListIconLarge = <Icon large={true} iconTag={icon}/>;

icon = "assignment_ind";
const accountsIcon = <Icon iconTag={icon}/>;
const accountsIconLarge = <Icon large={true} iconTag={icon}/>;

icon = "security";
const blockIcon = <Icon iconTag={icon}/>;
const blockIconLarge = <Icon large={true} iconTag={icon}/>;

icon = "perm_identity";
const userProfileIcon = <Icon iconTag={icon}/>;
const userProfileIconLarge = <Icon large={true} iconTag={icon}/>;

icon = "new_releases";
const notificationsIcon = <Icon iconTag={icon}/>;
const notificationsIconLarge = <Icon large={true} iconTag={icon}/>;

icon = "business";
const businessIcon = <Icon iconTag={icon}/>;
const businessIconLarge = <Icon large={true} iconTag={icon}/>;

icon = "lock_outline";
const maintenanceIcon = <Icon iconTag={icon}/>;
const maintenanceIconLarge = <Icon large={true} iconTag={icon}/>;
const CPUIcon = HardwareDeveloperBoard;
const DatabaseIcon = FileCloudDone;
const NotesIcon = ActionDescription;
const DocumentationIcon = CommunicationImportContacts;
const ExportIcon = CommunicationImportExport;
const DeleteIcon = ActionDelete;
const AddIcon = ActionAdd;
const CopyIcon = ContentContentCopy;
const AddImageIcon = ImageAddAPhoto;
const LockIcon = ActionLock;
const DownloadIcon = FileFileDownload;
const SearchIcon = ActionSearch;
const InfoIcon = ActionInfoOutline;
const AntennaIcon = ActionSettingsInputAntenna;
const ComponentIcon = ActionSettingsInputComponent;
const SvideoIcon = ActionSettingsInputSvideo;
const HDMIIcon = ActionSettingsInputHdmi;
const UndoIcon = ContentUndo;
const ToolsIcon = ActionBuild;
const LightIcon = ActionLightbulbOutline;
const DeviceIcon = DeviceDevices;
const RecordIcon = AvFiberSmartRecord;
const LaptopIcon = HardwareLaptopChromebook;
const MemoryIcon = HardwareMemory;
const AudioNoteIcon = ImageAudiotrack;
const RadioIcon = AvRadio;
const SurroundIcon = AvSurroundSound;
const SpeakerGroupIcon = HardwareSpeakerGroup;
const SpeakerIcon = HardwareSpeaker;
const FanIcon = HardwareToys;
const DiscIcon = NotificationDiscFull;
const MusicLibraryIcon = AvLibraryMusic;
const MusicIcon = AvMusicVideo;
const UserIcon = SocialPerson;
const QueueMusicIcon = AvQueueMusic;
const SubscriptionsIcon = AvSubscriptions;
const DVRIcon = DeviceDvr;
const HubIcon = HardwareDeviceHub;
const AntennaTVIcon = NotificationLiveTv;
const ComputerMonitorIcon2 = HardwareDesktopWindows;
const ComputerMonitorIcon = HardwareDesktopMac;
const LoginAs = CommunicationCallMade;
export {
  LoginAs,
  UserIcon,
  ComputerMonitorIcon2,
  ComputerMonitorIcon,
  AntennaTVIcon,
  HubIcon,
  DVRIcon,
  MusicIcon,
  CPUIcon,
  QueueMusicIcon,
  MusicLibraryIcon,
  DiscIcon,
  FanIcon,
  SpeakerGroupIcon,
  SpeakerIcon,
  SpeakerPhone,
  SurroundIcon,
  AudioNoteIcon,
  LaptopIcon,
  MemoryIcon,
  SubscriptionsIcon,
  LightIcon,
  DeviceIcon,
  ToolsIcon,
  UndoIcon,
  ComponentIcon,
  RadioIcon,
  SvideoIcon,
  HDMIIcon,
  RecordIcon,
  AntennaIcon,
  addPhotoIcon,
  addPhotoIconLarge,
  serverSettingsIcon,
  DatabaseIcon,
  serverSettingsIconLarge,
  userListIcon,
  userListIconLarge,
  accountsIcon,
  accountsIconLarge,
  userProfileIcon,
  userProfileIconLarge,
  notificationsIcon,
  notificationsIconLarge,
  maintenanceIcon,
  maintenanceIconLarge,
  businessIcon,
  businessIconLarge,
  DeleteIcon,
  AddIcon,
  EditIcon,
  LockIcon,
  AddImageIcon,
  CopyIcon,
  ExportIcon,
  AudioIcon,
  VideoDistributionIcon,
  DisplayIcon,
  TouchPadIcon,
  ConferencingIcon,
  ControlIcon,
  AvVolumeUp,
  MediaIcon,
  CameraIcon,
  PowerIcon,
  LockOpen,
  LockClosed,
  ActionSettingsApplications,
  VolumeUp,
  VolumeDown,
  VolumeMute,
  VolumeUnMute,
  HamburgerBar,
  ArrowBack,
  ArrowForward,
  Circle,
  SpaceBar,
  HorizontalLine,
  AccountEnterIcon,
  NavigationClose,
  Pause,
  Play,
  Stop,
  FastForward,
  Rewind,
  SkipNext,
  SkipPrevious,
  ShuffleIcon,
  LoopIcon,
  Record,
  MicOn,
  MicHollow,
  MicOff,
  Replay,
  Replay5,
  Replay10,
  Replay30,
  Left,
  Right,
  Up,
  Down,
  Eject,
  Select,
  UserRevokeIcon,
  SelectAll,
  Ethernet,
  UpArrow,
  DownArrow,
  Keyboard,
  RightArrow,
  Computer,
  VideoGame,
  DVDIcon,
  BluRayIcon,
  AppleTVIcon,
  ControlRoomIcon,
  Add,
  Minus,
  PlayPause,
  ZoomIn,
  ZoomOut,
  BrightnessDown,
  BrightnessUp,
  Rotate,
  PhotoCamera,
  Pan,
  Mask,
  Input,
  Menu,
  PIP,
  BackArrow,
  Home,
  DownloadIcon,
  ContentBlock,
  blockIcon,
  blockIconLarge,
  SearchIcon,
  InfoIcon,
  Explore,
  Presets,
  AvRadio,
  AvVideocam,
  AvVideocamOff,
  AvMic,
  AvMicOff,
  AvMicNone,
  AvMovie,
  AvHighQuality,
  AvHd,
  AvGames,
  AvForward5,
  AvForward10,
  AvForward30,
  AvFiberDvr,
  AvClosedCaption,
  AvEqualizer,
  VGA,
  MuteVideo,
  UnMuteVideo,
  Screens,
  ImageBackground,
  Save,
  ScreenShadeIcon,
  MountIcon,
  Security,
  Email,
  Settings,
  Reconnect,
  Freeze,
  ScreenShare,
  StopScreenShare,
  EmptyCheckBox,
  ImageNegative,
  Image,
  Sync,
  FileUpload,
  FileDownload,
  Fullscreen,
  FullscreenExit,
  CD,
  Mouse,
  Favorite,
  List,
  USB,
  Snooze,
  SimCard,
  Help,
  LooksOne,
  RemoteIcon,
  ExitIcon,
  ProjectorIcon,
  WhiteboardIcon,
  LooksTwo,
  CommunicationCall,
  ConferencingIcon2,
  License,
  SettingsRemote,
  LinkIcon,
  ReportProblemIcon,
  NewRelease,
  NotesIcon,
  DocumentationIcon,
  Cloud,
  CommunicationCallEnd,
  PhoneForwarded,
  CheckCircle,
  Rotate360,
  WorldIcon
   };
