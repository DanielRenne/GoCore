mkdir -p $GOPATH/src/github.com/DanielRenne/goCoreAppTemplate/web/app/dist/css/
mkdir -p $GOPATH/src/github.com/DanielRenne/goCoreAppTemplate/web/app/dist/javascript/
set -e
bash build-react.sh
bash build-css.sh 1