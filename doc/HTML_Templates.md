#HTML Templates

Within the htmlTemplates field in webConfig.json there are 3 fields to consider.

	"htmlTemplates":{
		"enabled":false,
		"directory":"templates",
		"directoryLevels": 1
	}

####enabled  (bool)

This will enable the GIN router to route html requests to the templates directory.

####directory (string)

This is the directory located beneath the web directory of your application where the template files are to be stored.

####directoryLevels  (int)  range 0 to 2

The GIN HTML Templates can go multiple levels of directories deep.  GoCore supports 0 to 2.  

By default GoCore looks for index.tmpl in the templates directory with a 0 level for the default root file.

By default GoCore looks for root/index.tmpl in the templates directory with a 1 level for the default root file.

By default GoCore looks for root/root/index.tmpl in the templates directory with a 2 level for the default root file.

When you specify levels greater than zero you must template your files with the following begin and end tags for the router to properly process your template file.  At the top you must add a define the directory path tag and create an {{end}} tag, see below:

	{{ define "root/index.tmpl" }}
	<!DOCTYPE html>
	<html lang="en">
	  <head>
	    <!-- Required meta tags always come first -->
	    <meta charset="utf-8">
	    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
	    <meta http-equiv="x-ua-compatible" content="ie=edge">
	
	    <!-- Bootstrap CSS -->
	    <link rel="stylesheet" href="web/core/bootstrap/dist/css/bootstrap.min.css">
	    <link rel="stylesheet" href="web/core/tether-1.2.0/dist/css/tether.min.css">
	    <link rel="stylesheet" href="web/core/font-awesome-4.5.0/css/font-awesome.min.css">
	
	
	    <!-- jQuery first, then Bootstrap JS. -->
	    <script src="web/core/js/jquery-2.2.2.min.js"></script>
	    <script src="web/core/js/jquery-2.2.2.min.js"></script>
	    <script src="web/core/tether-1.2.0/dist/js/tether.min.js" ></script>
	
	    <!-- Favicon -->
	  <link rel="icon" href="favicon.ico"/>
	<script type="text/javascript">
	    $(document).ready(function (){
	      $.getJSON( "WebAPI", function( data ) {
	        console.log(data);
	      });
	    });
	</script>
	
	  <title>Hello World.com</title>
	  </head>
	  <body>
	  Hello World
	  </body>
	</html>
	{{ end }}