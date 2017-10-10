#!/usr/bin/env node
var program = require('commander');
var inquirer = require("inquirer");
var Promise = require("bluebird");
var fs = require('fs');
var fsp = require('fs-extra-promise');
var path = require('path');
var glob = require("multi-glob").glob;
var extend = require('extend');
var Applause = require('applause');
var pp = require('preprocess');
var treeify = require('treeify');
var loadJsonFile = require('load-json-file');
var writeJsonFile = require('write-json-file');
var dive = require('dive');

program
 .version('1.0.0')
 .parse(process.argv);

var defaults = {
  directory: 'generated',
  includeSkins: true,
  includeSource: true,
  includeExamples: true,
  includeBrowserSync: true
};

var questions = [
  {
    type: "input",
    name: "directory",
    message: "Which directory do you want generated to",
    default: defaults.directory,
    validate: function( value ) {
      if(fsp.existsSync(value)){
        return 'The directory ' + value + ' exists. Please remove the directory. Or typping another directory.';
      }
      return true;
    },
  },
  {
    type: "list",
    name: "style",
    message: "Which style?",
    choices: [ {
      name: "Classic",
      value: 'classic'
    },{
      name: "Material",
      value: 'material',
      checked: true
    }],
    validate: function( answer ) {
      if ( answer.length < 1 ) {
        return "You must choose at least one.";
      }
      return true;
    }
  },
  {
    type: "list",
    name: "layout",
    message: "Which layout?",
    choices: [
      {
        name: "Base",
        value: 'base',
        checked: true
      },
      {
        name: "Center",
        value: 'center'
      },
      {
        name: "Iconbar",
        value: 'iconbar'
      },
      {
        name: "Mmenu",
        value: 'mmenu'
      },
      {
        name: "Topbar",
        value: 'topbar'
      },
      {
        name: "Topicon",
        value: 'topicon'
      }
    ]
  },
  {
    type: "list",
    name: "buildSystem",
    message: "Which build system",
    choices: [
      {
        name: "Grunt",
        value: 'grunt',
        checked: true
      },
      {
        name: "Gulp",
        value: 'gulp'
      },
      {
        name: "None",
        value: 'none'
      }
    ]
  },
  {
    type: "confirm",
    name: "includeSkins",
    message: "Use skins?",
    default: defaults.includeSkins
  },
  {
    type: "confirm",
    name: "includeExamples",
    message: "Include examples files?",
    default: defaults.includeExamples
  },
  {
    type: "confirm",
    name: "includeSource",
    message: "Include Source files?",
    when: function( answers ) {
      return answers.buildSystem === 'none';
    },
    default: defaults.includeSource
  },
  {
    type: "checkbox",
    name: "includeSources",
    message: "What sources to include?",
    choices: function(answers){
      var choices = [
        {
          name: "Javascript",
          value: 'js',
          checked: true
        },
        {
          name: "Less",
          value: 'css',
          checked: true
        },
        {
          name: "Vendor's less",
          value: 'vendor',
          checked: false
        },
        {
          name: "Fonts's less",
          value: 'fonts',
          checked: false
        },
        {
          name: "Html layout system",
          value: 'html',
          checked: true
        }
      ];

      if(answers.includeSkins) {
        choices.push({
          name: "Skins's less",
          value: 'skins',
          checked: true
        });
      }

      if(answers.includeExamples) {
        choices.push({
          name: "Examples's javascript & less",
          value: 'examples',
          checked: true
        });
      }

      return choices;
    },
    when: function( answers ) {
      return answers.includeSource !== false;
    },
    validate: function( answer ) {
      if ( answer.indexOf('css') == -1 && answer.indexOf('skins') != -1 ) {
        return "You must choose 'Less' as Skins's less dependency.";
      }

      if ( answer.indexOf('css') == -1 && answer.indexOf('examples') != -1 ) {
        return "You must choose 'Less' as Examples's dependency.";
      }

      if ( answer.indexOf('js') == -1 && answer.indexOf('examples') != -1 ) {
        return "You must choose 'Javascript' as Examples's dependency.";
      }

      if ( answer.indexOf('css') == -1 && answer.indexOf('vendor') != -1 ) {
        return "You must choose 'Less' as Vendor's dependency.";
      }

      if ( answer.indexOf('css') == -1 && answer.indexOf('fonts') != -1 ) {
        return "You must choose 'Less' as Fonts's dependency.";
      }

      if ( answer.indexOf('html') == -1 && answer.indexOf('examples') != -1 ) {
        return "You must choose 'Html layout system' as Examples's dependency.";
      }

      return true;
    }
  },
  {
    type: "confirm",
    name: "includeBrowserSync",
    message: "Use browsersync?",
    default: defaults.includeBrowserSync
  }
];

inquirer.prompt( questions).then(function( answers ) {

  answers = extend(defaults, answers);

  var sourcePath = path.join(answers.style, answers.layout);
  var destPath = answers.directory;
  var globalPath = path.join(answers.style, 'global');
  var buildSourcePath = '.src';

  var methods = {
    replaceCopy: function(file, applause) {
      var sourceFile = path.join(sourcePath, file);
      var destFile = path.join(destPath, file);

      if(!fsp.lstatSync(sourceFile).isDirectory()){
        var content = fsp.readFileSync(sourceFile, 'utf8');
        var result = applause.replace(content).content;

        if (result === false ){
          result = content;
        }

        return fsp.outputFileAsync(destFile, result);
      }
    },
    preprocessCopy: function(context, srcFile, destFile, options){
      if(typeof destFile === 'object') {
        options = destFile;
        destFile = srcFile;
      }

      if(typeof destFile === 'undefined') {
        destFile = srcFile;
      }

      return fsp.copyAsync(path.join(buildSourcePath, srcFile), path.join(destPath, destFile)).then(function(){
        pp.preprocessFileSync(path.join(destPath, destFile), path.join(destPath, destFile), context, options);
      });
    },

    justCopy: function(srcFile, destFile) {
      if(typeof destFile === 'undefined') {
        destFile = srcFile;
      }
      return fsp.copyAsync(path.join(buildSourcePath, srcFile), path.join(destPath, destFile));
    }
  };


  var tasks = {
    copyGlobalDist: function(){
      return new Promise(function(resolve, reject){
        var wait = [];

        var patterns = [
          'css',
          'fonts',
          'js',
          'vendor',
        ];

        if(answers.includeExamples) {
          patterns = patterns.concat([
            'photos',
            'portraits'
          ]);
        }

        glob(patterns, {
          cwd: globalPath
        }, function(err, files){
          if(err) {
            console.log(err);
          } else {
            files.forEach(function(file){
              wait.push(fsp.copyAsync(path.join(globalPath, file), path.join(destPath, 'assets', file)));
            });
          }
        });

        Promise.all(wait).then(function(){
          console.log('Copy global assets files successful.');

          resolve();
        });
      });
    },
    copyGlobalSource: function(){
      return new Promise(function(resolve, reject){
        var patterns = [];
        var wait = [];

        if(answers.includeSource) {
          if(answers.includeSources.indexOf('js') == -1) {
            patterns = patterns.concat([
              '!src/js/**/*',
            ]);
          } else {
            patterns = patterns.concat([
              'src/js/**/*',
              'components.json'
            ]);
          }

          if(answers.includeSources.indexOf('css') == -1) {
            patterns = patterns.concat([
              '!src/less/**/*',
            ]);
          } else {
            patterns = patterns.concat([
              'src/less/**/*',
            ]);
          }

          if(answers.includeSources.indexOf('skins') == -1) {
            patterns = patterns.concat([
              '!src/skins/**/*',
            ]);
          } else {
            patterns = patterns.concat([
              'src/skins/**/*',
            ]);
          }

          if(answers.includeSources.indexOf('vendor') == -1) {
            patterns = patterns.concat([
              '!src/vendor/**/*',
            ]);
          } else {
            patterns = patterns.concat([
              'src/vendor/**/*',
            ]);
          }

          if(answers.includeSources.indexOf('fonts') == -1) {
            patterns = patterns.concat([
              '!src/fonts/**/*',
            ]);
          } else {
            patterns = patterns.concat([
              'src/fonts/**/*',
            ]);
          }

          glob(patterns, {
            cwd: globalPath
          }, function(err, files){
            if(err) {
              console.log(err);
            } else {
              files.forEach(function(file){
                //if(fsp.existsSync(path.join(globalPath, file))){
                var filePath = path.join(globalPath, file);
                if(!fsp.lstatSync(filePath).isDirectory()){
                  wait.push(fsp.copyAsync(filePath, path.join(destPath, file), {
                    clobber: true
                  }));
                }
              });
            }
          });
        }

        Promise.all(wait).then(function(){
          setTimeout(function() {
            console.log('Copy global source files successful.');
            resolve();
          }, 1000);

        }, function(error){
          console.info(error);
        });
      });
    },
    copyLayout: function(){
      return new Promise(function(resolve, reject){
        var wait = [];

        var patterns = [
          '.editorconfig',
          '.gitattributes',
          '.gitignore',
          'assets/css/**/*',
          'assets/js/**/*',
          'assets/images/**/*'
        ];

        if(answers.includeSource) {
          patterns = patterns.concat([
            '.csscomb.json',
            '.csslintrc',
            '.jshintrc',
            'bower.json',
            'color.yml',
          ]);
        }

        if(answers.includeExamples) {
          patterns = patterns.concat([
            'assets/examples/css/**/*',
            'assets/examples/images/**/*',
            //'assets/data/**/*',
          ]);
        } else {
          patterns = patterns.concat([
            '!html/**/*',
            '!assets/examples/css/**/*',
            '!assets/examples/images/**/*',
            '!assets/data/**/*',
          ]);
        }

        if(answers.includeSkins) {
          patterns = patterns.concat([
            'assets/skins/**/*',
          ]);
        } else {
          patterns = patterns.concat([
            '!assets/skins/**/*',
          ]);
        }

        if(answers.includeSource) {
          if(answers.includeSources.indexOf('js') == -1) {
            patterns = patterns.concat([
              '!src/js/**/*',
            ]);
          } else {
            patterns = patterns.concat([
              'src/js/**/*',
            ]);
          }

          if(answers.includeSources.indexOf('css') == -1) {
            patterns = patterns.concat([
              '!src/less/**/*',
            ]);
          } else {
            patterns = patterns.concat([
              'src/less/**/*',
            ]);
          }

          if(answers.includeSources.indexOf('skins') == -1) {
            patterns = patterns.concat([
              '!src/skins/**/*',
            ]);
          } else {
            patterns = patterns.concat([
              'src/skins/**/*',
            ]);
          }

          if(answers.includeSources.indexOf('html') == -1) {
            patterns = patterns.concat([
              '!src/templates/**/*',
            ]);
          } else {
            // patterns = patterns.concat([
            //   'src/templates/**/*',
            // ]);
          }

          if(answers.includeSources.indexOf('examples') == -1) {
            patterns = patterns.concat([
              '!src/examples/less/**/*',
              '!src/templates/pages/*/*',
            ]);
          } else {
            patterns = patterns.concat([
              'src/examples/less/**/*',
              'src/examples/js/**/*'
              //'src/templates/pages/*/*',
            ]);
          }
        }

        glob(patterns, {
          cwd: sourcePath
        }, function(err, files){
          if(err) {
            console.log(err);
          } else {
            files.forEach(function(file){
              var filePath = path.join(sourcePath, file);

              if(!fsp.lstatSync(filePath).isDirectory()){
                wait.push(fsp.copyAsync(filePath, path.join(destPath, file), {
                  clobber: true
                }));
              }
            });
          }
        });

        Promise.all(wait).then(function(){
          setTimeout(function() {
            console.log('Copy layout files successful.');
            resolve();
          }, 1000);
        }, function(error){
          console.info(error);
        });
      });
    },
    updateDist: function(){
      return new Promise(function(resolve, reject){
        var patterns = ['html/**/*'];
        var wait = [];

        if(answers.includeExamples) {
          patterns = patterns.concat([
            'assets/data/**/*',
            'assets/examples/js/**/*',
          ]);
        }

        if(answers.includeSource) {
          if(answers.includeSources.indexOf('examples') == -1) {
            patterns = patterns.concat([
              'src/examples/js/**/*',
            ]);
          }
        }

        var applause = Applause.create({
          patterns: [
            {
              match: '../global/css',
              replacement: 'assets/css'
            },
            {
              match: '../global/fonts',
              replacement: 'assets/fonts'
            },
            {
              match: '../global/photos',
              replacement: 'assets/photos'
            },
            {
              match: '../global/js',
              replacement: 'assets/js'
            },
            {
              match: '../global/portraits',
              replacement: 'assets/portraits'
            },
            {
              match: '../global/vendor',
              replacement: 'assets/vendor'
            }
          ],
          usePrefix: false
        });

        glob(patterns, {
          cwd: sourcePath
        }, function(err, files){
          if(err) {
            console.log(err);
          } else {
            files.forEach(function(file){
              wait.push(methods.replaceCopy(file, applause));
            });
          }
        });

        Promise.all(wait).then(function(){
          setTimeout(function() {
            console.log('Update dist files successful.');
            resolve();
          }, 1000);
        }, function(error){
          console.info(error);
        });
      });
    },
    updateSource: function(){
      return new Promise(function(resolve, reject){
        var wait = [];

        if(answers.includeSource) {
          var applause = Applause.create({
            patterns: [
              {
                match: '{{global}}',
                replacement: '{{assets}}'
              }
            ],
            usePrefix: false
          });

          var patterns = [];

          if(answers.includeSources.indexOf('html') == -1) {
            patterns = patterns.concat([
              '!src/templates/**/*',
            ]);
          } else {
            patterns = patterns.concat([
              'src/templates/**/*',
            ]);
          }

          if(answers.includeSources.indexOf('examples') == -1) {
            patterns = patterns.concat([
              '!src/templates/pages/*/*',
            ]);
          } else {
            patterns = patterns.concat([
              'src/templates/pages/*/*',
            ]);
          }

          glob(patterns, {
            cwd: sourcePath
          }, function(err, files){
            if(err) {
              console.log(err);
            } else {
              files.forEach(function(file){
                wait.push(methods.replaceCopy(file, applause));
              });
            }
          });
        }

        Promise.all(wait).then(function(){
          setTimeout(function() {
            console.log('Update source files successful.');
            resolve();
          }, 1000);
        }, function(error){
          console.info(error);
        });
      });
    },
    copyBuildSystem: function(){
      return new Promise(function(resolve, reject){
        var wait = [];

        if(answers.buildSystem !== 'none') {
          var context = {};
          var wait = [];

          if(answers.includeSources.indexOf('js') !== -1) {
            context.processJs = true;
            context.processLint = true;
          }
          if(answers.includeSources.indexOf('css') !== -1) {
            context.processCss = true;
            context.processLint = true;
          }

          if(answers.includeSources.indexOf('skins') !== -1) {
            context.processSkins = true;
          }
          if(answers.includeSources.indexOf('html') !== -1) {
            context.processHtml = true;
          }
          if(answers.includeSources.indexOf('examples') !== -1) {
            context.processExamples = true;
          }
          if(answers.includeSources.indexOf('vendor') !== -1) {
            context.processVendor = true;
          }
          if(answers.includeSources.indexOf('fonts') !== -1) {
            context.processFonts = true;
          }

          if(answers.includeBrowserSync) {
            context.processBrowserSync = true;
          }

          switch(answers.buildSystem){
            case 'grunt':
              wait.push(methods.preprocessCopy(context, 'package.json.grunt', 'package.json', {type: 'js'}));
              wait.push(methods.preprocessCopy(context, 'Gruntfile.js'));

              wait.push(methods.preprocessCopy(context, path.join('grunt', 'clean.js')));

              if(context.processBrowserSync) {
                wait.push(methods.preprocessCopy(context, path.join('grunt', 'browserSync.js')));
              }

              if(context.processHtml) {
                wait.push(methods.preprocessCopy(context, path.join('grunt', 'bootlint.js')));
                wait.push(methods.preprocessCopy(context, path.join('grunt', 'hb.js')));
                wait.push(methods.preprocessCopy(context, path.join('grunt', 'htmllint.js')));
                wait.push(methods.preprocessCopy(context, path.join('grunt', 'prettify.js')));
              }

              if(context.processJs || context.processExamples){
                wait.push(methods.preprocessCopy(context, path.join('grunt', 'concat.js')));
                wait.push(methods.preprocessCopy(context, path.join('grunt', 'uglify.js')));
                wait.push(methods.preprocessCopy(context, path.join('grunt', 'jshint.js')));
              }

              if(context.processCss || context.processExamples || context.processSkins || context.processVendor || context.processFonts) {
                wait.push(methods.preprocessCopy(context, path.join('grunt', 'autoprefixer.js')));
                wait.push(methods.preprocessCopy(context, path.join('grunt', 'csscomb.js')));
                wait.push(methods.preprocessCopy(context, path.join('grunt', 'csslint.js')));
                wait.push(methods.preprocessCopy(context, path.join('grunt', 'cssmin.js')));
                wait.push(methods.preprocessCopy(context, path.join('grunt', 'less.js')));
              }

              wait.push(methods.preprocessCopy(context, path.join('grunt', 'notify.js')));
              break;
            case 'gulp':
              wait.push(methods.preprocessCopy(context, 'package.json.gulp', 'package.json', {type: 'js'}));
              wait.push(methods.preprocessCopy(context, 'gulpfile.js'));
              wait.push(methods.justCopy(path.join('gulp', 'utils')));

              if(context.processBrowserSync) {
                wait.push(methods.justCopy(path.join('gulp', 'options', 'serve.js')));
                wait.push(methods.justCopy(path.join('gulp', 'recipes', 'serve.js')));
              }

              if(context.processHtml) {
                wait.push(methods.justCopy(path.join('gulp', 'options', 'html.js')));
                wait.push(methods.justCopy(path.join('gulp', 'recipes', 'html')));
              }

              if(context.processExamples) {
                wait.push(methods.justCopy(path.join('gulp', 'options', 'examples.js')));
                wait.push(methods.justCopy(path.join('gulp', 'recipes', 'examples')));
              }

              if(context.processSkins) {
                wait.push(methods.justCopy(path.join('gulp', 'options', 'skins.js')));
                wait.push(methods.justCopy(path.join('gulp', 'recipes', 'skins')));
              }

              if(context.processFonts) {
                wait.push(methods.justCopy(path.join('gulp', 'options', 'fonts.js')));
                wait.push(methods.justCopy(path.join('gulp', 'recipes', 'fonts')));
              }

              if(context.processVendor) {
                wait.push(methods.justCopy(path.join('gulp', 'options', 'vendor.js')));
                wait.push(methods.justCopy(path.join('gulp', 'recipes', 'vendor')));
              }

              if(context.processJs) {
                wait.push(methods.justCopy(path.join('gulp', 'options', 'scripts.js')));
                wait.push(methods.justCopy(path.join('gulp', 'recipes', 'scripts')));

                wait.push(methods.justCopy(path.join('gulp', 'options', 'jshint.js')));
                wait.push(methods.justCopy(path.join('gulp', 'recipes', 'jshint.js')));
              }

              if(context.processCss) {
                wait.push(methods.justCopy(path.join('gulp', 'options', 'styles.js')));
                wait.push(methods.justCopy(path.join('gulp', 'recipes', 'styles')));

                wait.push(methods.justCopy(path.join('gulp', 'options', 'csslint.js')));
                wait.push(methods.justCopy(path.join('gulp', 'recipes', 'csslint.js')));
              }

              break;
          }

          var configObj = {
            "assets": "assets",
            "destination": {},
            "source": {}
          };
          if(context.processHtml) {
            configObj.html = "html";
            configObj.templates = {
              "data": "src/templates/data",
              "helpers": "src/templates/helpers",
              "partials": "src/templates/partials",
              "pages": "src/templates/pages"
            }

          }

          if(context.processCss) {
            configObj.destination.css = "assets/css";
            configObj.source.less = "src/less";
            configObj.bootstrap = {
              "less": "src/less/bootstrap",
              "mixins": "src/less/mixins"
            };
            configObj.autoprefixerBrowsers = [
              "Android 2.3",
              "Android >= 4",
              "Chrome >= 20",
              "Firefox >= 24",
              "Explorer >= 8",
              "iOS >= 6",
              "Opera >= 12",
              "Safari >= 6"
            ];
          }

          if(context.processJs) {
            configObj.destination.js = "assets/js";
            configObj.source.js = "src/js";
          }

          if(context.processSkins) {
            configObj.destination.skins = "assets/skins";
            configObj.source.skins = "src/skins";
          }

          if(context.processExamples) {
            configObj.destination.examples = "assets/examples";
            configObj.source.examples = "src/examples";
          }

          if(context.processFonts) {
            configObj.destination.fonts = "assets/fonts";
            configObj.source.fonts = "src/fonts";
          }

          if(context.processVendor) {
            configObj.destination.vendor = "assets/vendor";
            configObj.source.vendor = "src/vendor";
          }

          wait.push(writeJsonFile(path.join(destPath, 'config.json'), configObj));
        }

        Promise.all(wait).then(function(){
          console.log('Copy build system files successful.');
          resolve();
        }, function(error){
          console.info(error);
        });
      });
    }
  };

  Promise.all([tasks.copyGlobalDist()]).then(function(){
    return tasks.copyGlobalSource();
  }).then(function(){
    return tasks.copyLayout();
  }).then(function(){
    return tasks.updateDist();
  }).then(function(){
    return tasks.updateSource();
  }).then(function(){
    return tasks.copyBuildSystem();
  }).then(function(){
    console.info('All successful.');

    var tree = {};

    dive(destPath, { all: true, directories: true }, function(err, thisPath) {
      var relativePath = path.relative(destPath, thisPath),
          node = tree;

      if (relativePath.indexOf('..') !== 0) {
         relativePath.split(path.sep).forEach(function(part) {
            typeof node[part] !== 'object' && (node[part] = {});
            node = node[part];
         });
      }
   }, function(){
      console.log(treeify.asTree(tree, true));

      console.log('All files are saved to "'+destPath+'" folder.');
   });

  }).catch(function(error){
    console.log(error);
  });
});
