echo "New version number (with out the v)"
read version
git tag "v$version"
git push origin "v$version"

