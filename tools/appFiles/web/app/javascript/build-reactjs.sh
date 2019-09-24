set -e
cd ..
webpack --config webpack-production.config.js

cp lib/Material-Design-Iconic-Font.eot dist/javascript/
cp lib/Material-Design-Iconic-Font.svg dist/javascript/
cp lib/Material-Design-Iconic-Font.ttf dist/javascript/
cp lib/Material-Design-Iconic-Font.woff dist/javascript/
cp lib/Material-Design-Iconic-Font.woff2 dist/javascript/
cp lib/brand-icons.svg dist/javascript/
cp lib/brand-icons.ttf dist/javascript/
cp lib/brand-icons.woff dist/javascript/
cp lib/brand-icons.woff2 dist/javascript/

cd javascript
gzip -f go-core-app.js
mv -f go-core-app.js.* ../dist/javascript/
set +e
rm ../dist/javascript/atlona-studio.js.map > /dev/null 2>&1
set -e

cp polyfills/polyfills.js  ../dist/javascript/polyfills.js
cp Core.js ../dist/javascript/core.js
cp json.js ../dist/javascript/json.js
cp Experimental.js ../dist/javascript/experimental.js
cd ../dist/javascript/
gzip -f core.js
gzip -f experimental.js
