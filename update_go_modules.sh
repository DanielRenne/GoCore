set -e
go get -u
go mod tidy
cd core
echo "Updating core"
go get -u
go mod tidy
cd app
echo "Updating app"
go get -u
go mod tidy
cd -
echo "Updating app/api"
cd app/api
go get -u
go mod tidy
cd -
echo "Updating appGen"
cd appGen
go get -u
go mod tidy
cd -
echo "Updating atomicTypes"
cd atomicTypes
go get -u
go mod tidy
cd -
echo "Updating channels"
cd channels
go get -u
go mod tidy
cd -
echo "Updating cron"
cd cron
go get -u
go mod tidy
cd -
echo "Updating crypto"
cd crypto
go get -u
go mod tidy
cd -
echo "Updating dbServices"
cd dbServices
go get -u
go mod tidy
cd -
echo "Updating dbServices/bolt/stubs"
cd dbServices/bolt/stubs
go get -u
go mod tidy
cd -
echo "Updating dbServices/common/stubs"
cd dbServices/common/stubs
go get -u
go mod tidy
cd -
echo "Updating dbServices/example"
cd dbServices/example
go get -u
go mod tidy
cd -
echo "Updating dbServices/example/modelsGenerate"
cd dbServices/example/modelsGenerate
go get -u
go mod tidy
cd -
echo "Updating dbServices/mongo/stubs"
cd dbServices/mongo/stubs
go get -u
go mod tidy
cd -
echo "Updating extensions"
cd extensions
go get -u
go mod tidy
cd -
echo "Updating fileCache"
cd fileCache
go get -u
go mod tidy
cd -
echo "Updating ginServer"
cd ginServer
go get -u
go mod tidy
cd -
echo "Updating gitWebHooks"
cd gitWebHooks
go get -u
go mod tidy
cd -
echo "Updating httpExtensions"
cd httpExtensions
go get -u
go mod tidy
cd -
echo "Updating logger"
cd logger
go get -u
go mod tidy
cd -
echo "Updating mongo"
cd mongo
go get -u
go mod tidy
cd -
echo "Updating path"
cd path
go get -u
go mod tidy
cd -
echo "Updating pubsub"
cd pubsub
go get -u
go mod tidy
cd -
echo "Updating serverSettings"
cd serverSettings
go get -u
go mod tidy
cd -
echo "Updating store"
cd store
go get -u
go mod tidy
cd -
echo "Updating utils"
cd utils
go get -u
go mod tidy
cd -
echo "Updating zip"
cd zip
go get -u
go mod tidy
cd -
echo "Updating workQueue"
cd workQueue
go get -u
go mod tidy
cd ..
cd ..
echo "Successfully updated all modules"
