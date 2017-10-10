set -e
cd ..
webpack -p --config webpack-production.config.js

cp remark/material/global/fonts/material-design/Material-Design-Iconic-Font.eot dist/javascript/
cp remark/material/global/fonts/material-design/Material-Design-Iconic-Font.svg dist/javascript/
cp remark/material/global/fonts/material-design/Material-Design-Iconic-Font.ttf dist/javascript/
cp remark/material/global/fonts/material-design/Material-Design-Iconic-Font.woff dist/javascript/
cp remark/material/global/fonts/material-design/Material-Design-Iconic-Font.woff2 dist/javascript/
cp remark/material/global/fonts/brand-icons/brand-icons.svg dist/javascript/
cp remark/material/global/fonts/brand-icons/brand-icons.ttf dist/javascript/
cp remark/material/global/fonts/brand-icons/brand-icons.woff dist/javascript/
cp remark/material/global/fonts/brand-icons/brand-icons.woff2 dist/javascript/

cd javascript
gzip -f go-core-app.js
mv -f go-core-app.js.* ../dist/javascript/
cp ../node_modules/react-intl-tel-input/dist/flags.png ../dist/css/
cp ../node_modules/react-intl-tel-input/dist/flags@2x.png  ../dist/css/
cp ../node_modules/react-intl-tel-input/dist/libphonenumber.js  ../dist/javascript/
gzip -f ../dist/javascript/libphonenumber.js

cp polyfills/polyfills.js  ../dist/javascript/polyfills.js
cp RemarkCore.js ../dist/javascript/remark-core.js
cp json.js ../dist/javascript/json.js
cp RemarkExperimental.js ../dist/javascript/remark-experimental.js
cd ../dist/javascript/
gzip -f remark-core.js
gzip -f remark-experimental.js
