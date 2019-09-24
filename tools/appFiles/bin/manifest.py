import json
import sys
import os
git_commit = sys.argv[1]
package_name = sys.argv[2]
gopath = os.getenv("GOPATH")
manifest = json.load(open(gopath + '/src/github.com/DanielRenne/goCoreAppTemplate/vendorPackages/github-manifest.json', 'r'))
manifest[package_name] = git_commit
json.dump(manifest, open(gopath + '/src/github.com/DanielRenne/goCoreAppTemplate/vendorPackages/github-manifest.json', 'w'), indent=4)