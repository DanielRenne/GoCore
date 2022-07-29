package main

import (
	_ "github.com/DanielRenne/GoCore/core"
	_ "github.com/DanielRenne/GoCore/core/app"
	_ "github.com/DanielRenne/GoCore/core/app/api"
	_ "github.com/DanielRenne/GoCore/core/appGen"
	_ "github.com/DanielRenne/GoCore/core/atomicTypes"
	_ "github.com/DanielRenne/GoCore/core/channels"
	_ "github.com/DanielRenne/GoCore/core/crypto"
	_ "github.com/DanielRenne/GoCore/core/dbServices"
	_ "github.com/DanielRenne/GoCore/core/dbServices/bolt/stubs"
	_ "github.com/DanielRenne/GoCore/core/dbServices/common/stubs"
	_ "github.com/DanielRenne/GoCore/core/dbServices/mongo/acct"
	_ "github.com/DanielRenne/GoCore/core/dbServices/mongo/stubs"
	_ "github.com/DanielRenne/GoCore/core/extensions"
	_ "github.com/DanielRenne/GoCore/core/fileCache"
	_ "github.com/DanielRenne/GoCore/core/ginServer"
	_ "github.com/DanielRenne/GoCore/core/gitWebHooks"
	_ "github.com/DanielRenne/GoCore/core/httpExtensions"
	_ "github.com/DanielRenne/GoCore/core/logger"
	_ "github.com/DanielRenne/GoCore/core/pubsub"
	_ "github.com/DanielRenne/GoCore/core/serverSettings"
	_ "github.com/DanielRenne/GoCore/core/store"
	_ "github.com/DanielRenne/GoCore/core/syncAtomic"
	_ "github.com/DanielRenne/GoCore/core/utils"
	_ "github.com/DanielRenne/GoCore/core/zip"
)

func main() {

}
