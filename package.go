package goCore

import (
	// Required package dependencies for mongo models
	_ "github.com/DanielRenne/GoCore/buildCore"
	_ "github.com/DanielRenne/GoCore/core"
	_ "github.com/DanielRenne/GoCore/core/app"
	_ "github.com/DanielRenne/GoCore/core/app/api"
	_ "github.com/DanielRenne/GoCore/core/appGen"
	_ "github.com/DanielRenne/GoCore/core/atomicTypes"
	_ "github.com/DanielRenne/GoCore/core/channels"
	_ "github.com/DanielRenne/GoCore/core/cmdExec"
	_ "github.com/DanielRenne/GoCore/core/crypto"
	_ "github.com/DanielRenne/GoCore/core/dbServices"
	_ "github.com/DanielRenne/GoCore/core/dbServices/bolt/stubs"
	_ "github.com/DanielRenne/GoCore/core/dbServices/common/stubs"
	_ "github.com/DanielRenne/GoCore/core/dbServices/mongo/stubs"
	_ "github.com/DanielRenne/GoCore/core/extensions"
	_ "github.com/DanielRenne/GoCore/core/fileCache"
	_ "github.com/DanielRenne/GoCore/core/ginServer"
	_ "github.com/DanielRenne/GoCore/core/gitWebHooks"
	_ "github.com/DanielRenne/GoCore/core/httpExtensions"
	_ "github.com/DanielRenne/GoCore/core/logger"
	_ "github.com/DanielRenne/GoCore/core/mongo"
	_ "github.com/DanielRenne/GoCore/core/pubsub"
	_ "github.com/DanielRenne/GoCore/core/serverSettings"
	_ "github.com/DanielRenne/GoCore/core/store"
	_ "github.com/DanielRenne/GoCore/core/utils"
	_ "github.com/DanielRenne/GoCore/core/workQueue"
	_ "github.com/DanielRenne/GoCore/core/zip"
	_ "github.com/DanielRenne/GoCore/modelBuild"
	_ "github.com/altipla-consulting/i18n-dateformatter"
	_ "github.com/asaskevich/govalidator"
	_ "github.com/fatih/camelcase"
)
