package br

import (
	"net"

	"sync"

	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/networks"
	"github.com/DanielRenne/goCoreAppTemplate/queries"
	"github.com/pkg/errors"
)

/*

  To call a custom store function like checking IPV4 addresses put something like this in your componentDidMount()

      this.registerSubscriptions([
        this.store.dbRegister("YourCollectionName", "MongoIdOfRecord", "FieldNameOfMongo", (response) => {
	    // So on change of the above field in the store, we are ensuring the value set is an IPv4 address
            // this value is on blur so then run our checks on the IP
            this.store.appSet("IPAddress.ValidateIPV4", {UserId: window.appState.UserId, IpAddress: response, MinIP: "", MaxIP: "", MinIP2: "225.0.0.1", MaxIP2: "239.255.255.254"}, (errors) => {
              if (!errors.Valid) {
                this.store.set("YourCollectionName", "MongoIdOfRecord", "FieldNameOfMongo", this.state.YourIPInState, () => {
                  this.globs.Popup(errors.Message);
                });
              } else {
                this.base.setState((s) => {
                  s.YourIPInState = response;
                  return s
                });
              }
            });
        }, true),
      ]);

	And for your ip address field import this


	import TextFieldStore from "../../components/store/textField";

	<TextFieldStore
		changeOnBlur={true}
		collection={"YourCollectionName"}
		id={"MongoIdOfRecord"}
		path={"FieldNameOfMongo"}
		value={this.state.YourIPInState}
		floatingLabelText={"* IPAddress"}
		hintText={"* IPAddress"}
		fullWidth={true}
	/>

	See serverSettingsModify for use case

*/

func init() {
	RegisterBr(&IPAddress{})
}

var lockNextIP sync.RWMutex

type IPAddress struct{}
type ValidateIPV4 struct {
	IpAddress string `json:"IpAddress"`
}

func GetResponse() map[string]interface{} {
	response := make(map[string]interface{}, 0)
	response["Valid"] = true
	response["Message"] = ""
	return response
}

func VerifyIp(ipAddress string, lowIp string, highIP string) error {
	ip := net.ParseIP(ipAddress)
	if ip.To4() == nil {
		return errors.New("Bad IPv4 Address")
	}
	if lowIp != "" {
		if !networks.IsGreaterThanOrEqualTo(ipAddress, lowIp) {
			return errors.New("Next Ip address (" + ipAddress + ") is not greater than or equal to: " + lowIp + " as per the requirements")
		}
	}
	if highIP != "" {
		if networks.IsGreaterThanOrEqualTo(ipAddress, highIP) {
			return errors.New("Next Ip address (" + ipAddress + ") is greater than: " + highIP + " as per the requirements")
		}
	}
	return nil
}

func (IPAddress) ValidateIPV4(x interface{}) interface{} {
	rawPost := x.(map[string]interface{})
	response := GetResponse()
	ipAddress := rawPost["IpAddress"].(string)
	userId := rawPost["UserId"].(string)
	if ipAddress != "" && userId != "" {
		var user model.User
		err := model.Users.Query().ById(userId, &user)
		if err != nil {
			response["Valid"] = false
			response["Message"] = "Error: " + err.Error()
			return response
		}
		if ipAddress == "" {
			response["Valid"] = false
			response["Message"] = constants.ERROR_REQUIRED_FIELD
			return response
		}

		ipObj := net.ParseIP(ipAddress)

		// if ipObj.To4()[3] == 0 || ipObj.To4()[3] == 255 {
		// 	response["Valid"] = false
		// 	response["Message"] = queries.AppContent.GetTranslationFromUser(user, "IPAddressCantStart")
		// 	return response
		// }

		if ipObj.To4() == nil {
			response["Valid"] = false
			response["Message"] = queries.AppContent.GetTranslationFromUser(user, "InInvalidIP")
			return response
		}
		// minIP := rawPost["MinIP"].(string)
		// maxIP := rawPost["MaxIP"].(string)

		// minIP2 := rawPost["MinIP2"].(string)
		// maxIP2 := rawPost["MaxIP2"].(string)
		// if minIP != "" && maxIP != "" {
		// 	if !networks.IsWithinRange(ipAddress, minIP, maxIP) {
		// 		response["Valid"] = false
		// 		replacements := queries.TagReplacements{
		// 			Tag1: queries.Q("start", minIP),
		// 			Tag2: queries.Q("end", maxIP),
		// 		}
		// 		response["Message"] = queries.AppContent.GetTranslationWithReplacementsFromUser(user, "IPAddressInRange", &replacements)
		// 		return response
		// 	}
		// }
		// if minIP2 != "" && maxIP2 != "" {
		// 	if !networks.IsWithinRange(ipAddress, minIP, maxIP) {
		// 		response["Valid"] = false
		// 		replacements := queries.TagReplacements{
		// 			Tag1: queries.Q("start", minIP),
		// 			Tag2: queries.Q("end", maxIP),
		// 		}
		// 		response["Message"] = queries.AppContent.GetTranslationWithReplacementsFromUser(user, "IPAddressInRange", &replacements)
		// 		return response
		// 	}
		// }
	}

	return response
}
