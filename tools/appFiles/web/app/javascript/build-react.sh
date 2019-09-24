set -e
cd ..
#npm install --unsafe-perm
#npm install --no-optional --save --save-exact immutability-helper@2.0.0
#npm install --no-optional --save --save-exact classnames@2.2.5 normalize.css@5.0.0 react@15.4.2 react-addons-css-transition-group@15.4.2 react-dom@15.4.2
#npm install --no-optional --save --save-exact react-flexgrid@0.8.0
#npm install --no-optional --save --save-exact
npm install --no-optional -g webpack@2.6.1
cd javascript
bash build-reactjs.sh
