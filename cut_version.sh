set -e
currentDir=$(pwd)
echo "New version number (with out the v)"
read version
echo "Whats your title/description for the release?"
read title
gh release create "v$version" --title "$title" --generate-notes
set +e
say "Release created, rebuilding documentation via go get on new tmp package.  The docs should be available in one minute"
cd /tmp
rm -rf goCore
mkdir goCore
cd goCore
go mod init rebuildDocumentation
go get github.com/DanielRenne/GoCore@v$version
cd ..
rm -rf goCore
cd $currentDir
