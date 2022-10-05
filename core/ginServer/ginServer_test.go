package ginServer_test

import (
	"github.com/DanielRenne/GoCore/core/ginServer"
)

// Simple example using csrf security, session cookies, and a custom port.
func ExampleConfigureGin() {
	csrf := "my crsf secret"
	cookieConfig := ginServer.SessionConfiguration{
		Enabled:               true,
		SessionKey:            "test",
		SessionName:           "test",
		SessionExpirationDays: 15,
		SessionSecureCookie:   false,
	}
	ginServer.ConfigureGin("debug", "", false, []string{}, csrf, cookieConfig, true)
	// Run this if you want to use the gin server after configuring it how you want
	// import "github.com/DanielRenne/GoCore/core/app"
	// app.RunLite(9090)
}
