### Notes before getting started

The full front-end back-end scripts only will work on a linux or mac machine. It provides a basic login and authentication schema and user management pages with a react front-end. As of right now, our webpack and react libraries are outdated in this template and are on react 15 and webpack 2.6.1. We plan on updating this soon to react 18 and webpack 5.

On a project that uses this we are fairly far behind node 18, so its recommended to have node 12 installed. We based this appGeneration off some of how we wired up GoCore. For this to work, you dont have to initially setup a go mod init as the main will do it for you. Note, that the main package name of your entire app will be github.com/yourUserName/appName as the module name.

You can get a feel for what your application will look like by previewing the source code [here](https://github.com/davidrenne/GoCoreAppTemplate). This is the template that would be copied.

If you would like to proceed and at least test it out, we would love some feedback on this so react out to [me](mailto:dnxglya4@duck.com)

### Install Homebrew

```
/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
```

### Install MongoDB

```
brew tap mongodb/brew
brew install mongodb-community
brew services start mongodb/brew/mongodb-community
```

### Install NPM and Node JS

https://nodejs.org/en/download/

### Install Webpack

    npm install --global webpack@v2.6.1

### Install Golang

[Setup](https://golang.org/doc/install/ "Setup")

### Add your gopath bin directory to your PATH and set NODE_ENV=development

```
vim ~/.bash_profile
```

Add lines to add NODE_ENV to your environment
```
export PATH="$PATH:$HOME/go/bin"
export NODE_ENV=development
```

Source it:

```
source ~/.bash_profile
```

### Build GoCore Front End/Backend App

A sample project generator is available to build a GoCore project.

```
go install github.com/DanielRenne/GoCore/goCoreCreateApp
```

And install the binary for fetching the project template files

```
go install github.com/DanielRenne/GoCore/getAppTemplate
```

Then run

```
goCoreCreateApp
```

Follow the prompts to generate your app. Ensure you install nodejs, npm before generating an app with this binary.

GoCore has built in functions to read json configuration files to generate SQL Tables, indexes, primary & foreign keys, NOSQL Collections & Buckets, auto-generated Golang model code files. See [BuildCore Readme](https://github.com/DanielRenne/GoCore/blob/master/buildCore/README.md) for details on how you will call your buildCore main to rebuild models typically when your json schema changes.

### Run Your GoCore App

1.  Run the following to start your app

```
bash bin/start_app
```

Open a web browser to: [http://127.0.0.1](http://127.0.0.1)

Login as admin/admin and setup application roles etc. More documentation to come later on application specific setup.
