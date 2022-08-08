set -e
echo "New version number (with out the v)"
read version
# git tag "v$version"
# git push origin "v$version"
echo "Tag pushed, whats your title/description for the release?"
read title
gh release create "$version" --title "$title" --generate-notes
