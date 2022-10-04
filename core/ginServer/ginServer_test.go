package ginServer_test

import (
	"github.com/DanielRenne/GoCore/core/app"
	"github.com/DanielRenne/GoCore/core/ginServer"
)

func ExampleConfigureGin() {
	/*
		import (
			"github.com/DanielRenne/GoCore/core/app"
			"github.com/DanielRenne/GoCore/core/ginServer"
		)
	*/
	csrf := "my crsf secret"
	cookieConfig := ginServer.SessionConfiguration{
		Enabled:               true,
		SessionKey:            "test",
		SessionName:           "test",
		SessionExpirationDays: 15,
		SessionSecureCookie:   false,
	}
	ginServer.ConfigureGin("debug", "", false, []string{}, csrf, cookieConfig, true)
	app.RunLite(9090)
}
