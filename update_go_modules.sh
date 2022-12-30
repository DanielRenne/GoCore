set -e
go get -u
echo -e "\n\n\n\n\n\n\n"
go mod tidy
dirname="buildCore"
cd $dirname
echo " Start of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go get -u
echo -e "\n\n\n\n\n\n\n"
echo " End of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go mod tidy
cd -

dirname="modelBuild"
cd $dirname
echo " Start of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go get -u
echo -e "\n\n\n\n\n\n\n"
echo " End of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"

go mod tidy
cd -

dirname="core"
cd $dirname
echo " Start of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go get -u
echo -e "\n\n\n\n\n\n\n"
echo " End of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"

go mod tidy
dirname="app"
cd $dirname
echo " Start of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go get -u
echo -e "\n\n\n\n\n\n\n"
echo " End of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go mod tidy
cd -

dirname="app/api"
cd $dirname
echo " Start of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go get -u
echo -e "\n\n\n\n\n\n\n"
echo " End of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go mod tidy
cd -

dirname="appGen"
cd $dirname
echo " Start of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go get -u
echo -e "\n\n\n\n\n\n\n"
echo " End of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go mod tidy
cd -

dirname="atomicTypes"
cd $dirname
echo " Start of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go get -u
echo -e "\n\n\n\n\n\n\n"
echo " End of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go mod tidy
cd -

dirname="channels"
cd $dirname
echo " Start of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go get -u
echo -e "\n\n\n\n\n\n\n"
echo " End of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go mod tidy
cd -

dirname="cron"
cd $dirname
echo " Start of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go get -u
echo -e "\n\n\n\n\n\n\n"
echo " End of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go get -u
go mod tidy
cd -

dirname="crypto"
cd $dirname
echo " Start of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go get -u
echo -e "\n\n\n\n\n\n\n"
echo " End of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go mod tidy
cd -

dirname="dbServices"
cd $dirname
echo " Start of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go get -u
echo -e "\n\n\n\n\n\n\n"
echo " End of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go mod tidy
cd -

dirname="dbServices/bolt/stubs"
cd $dirname
echo " Start of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go get -u
echo -e "\n\n\n\n\n\n\n"
echo " End of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go mod tidy
cd -

dirname="dbServices/common/stubs"
cd $dirname
echo " Start of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go get -u
echo -e "\n\n\n\n\n\n\n"
echo " End of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go mod tidy
cd -

dirname="dbServices/example"
cd $dirname
echo " Start of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go get -u
echo -e "\n\n\n\n\n\n\n"
echo " End of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go mod tidy
cd -

dirname="dbServices/example/modelsGenerate"
cd $dirname
echo " Start of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go get -u
echo -e "\n\n\n\n\n\n\n"
echo " End of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go mod tidy
cd -

dirname="dbServices/mongo/stubs"
cd $dirname
echo " Start of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go get -u
echo -e "\n\n\n\n\n\n\n"
echo " End of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go mod tidy
cd -

dirname="extensions"
cd $dirname
echo " Start of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go get -u
echo -e "\n\n\n\n\n\n\n"
echo " End of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go mod tidy
cd -

dirname="fileCache"
cd $dirname
echo " Start of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go get -u
echo -e "\n\n\n\n\n\n\n"
echo " End of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go mod tidy
cd -

dirname="ginServer"
cd $dirname
echo " Start of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go get -u
echo -e "\n\n\n\n\n\n\n"
echo " End of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go mod tidy
cd -

dirname="gitWebHooks"
cd $dirname
echo " Start of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go get -u
echo -e "\n\n\n\n\n\n\n"
echo " End of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go mod tidy
cd -

dirname="httpExtensions"
cd $dirname
echo " Start of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go get -u
echo -e "\n\n\n\n\n\n\n"
echo " End of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go mod tidy
cd -

dirname="logger"
cd $dirname
echo " Start of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go get -u
echo -e "\n\n\n\n\n\n\n"
echo " End of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go mod tidy
cd -

dirname="mongo"
cd $dirname
echo " Start of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go get -u
echo -e "\n\n\n\n\n\n\n"
echo " End of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go mod tidy
cd -

dirname="path"
cd $dirname
echo " Start of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go get -u
echo -e "\n\n\n\n\n\n\n"
echo " End of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go mod tidy
cd -

dirname="pubsub"
cd $dirname
echo " Start of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go get -u
echo -e "\n\n\n\n\n\n\n"
echo " End of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go mod tidy
cd -

dirname="serverSettings"
cd $dirname
echo " Start of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go get -u
echo -e "\n\n\n\n\n\n\n"
echo " End of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go mod tidy
cd -

dirname="store"
cd $dirname
echo " Start of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go get -u
echo -e "\n\n\n\n\n\n\n"
echo " End of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go mod tidy
cd -

dirname="utils"
cd $dirname
echo " Start of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go get -u
echo -e "\n\n\n\n\n\n\n"
echo " End of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go mod tidy
cd -

dirname="zip"
cd $dirname
echo " Start of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go get -u
echo -e "\n\n\n\n\n\n\n"
echo " End of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go mod tidy
cd -

dirname="workQueue"
cd $dirname
echo " Start of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go get -u
echo -e "\n\n\n\n\n\n\n"
echo " End of Updates from $dirname "
echo -e "\n\n\n\n\n\n\n"
go mod tidy
cd ..
cd ..
echo "Successfully updated all modules"
