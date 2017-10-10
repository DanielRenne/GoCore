set -e
pwd
if [ $# -lt 1 ]; then
  #Not building in OSX due to node-gyp.  Dave builds flexbox examples and commits to source.
  cd $GOPATH/src/github.com/DanielRenne/goCoreAppTemplate/web/app/css/flexbox-examples
  gulp css
  cd $GOPATH/src/github.com/DanielRenne/goCoreAppTemplate/web/app/
  webpack -p --config webpack-production.config.js
fi

cd $GOPATH/src/github.com/DanielRenne/goCoreAppTemplate/web/app/javascript
gzip -f dist/css/go-core-app.css
mv -f dist/css/go-core-app.css.gz ../dist/css/
mv -f dist/css/go-core-app.css* ../dist/css/
rm -rf dist/

cd $GOPATH/src/github.com/DanielRenne/goCoreAppTemplate/web/app
cp css/RemarkCore.css dist/css/remark-core.css
cp css/RemarkExperimental.css dist/css/remark-experimental.css
cd dist/css/
gzip -f remark-core.css
gzip -f remark-experimental.css
