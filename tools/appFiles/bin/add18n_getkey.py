import sys

controller = sys.argv[1]
key = sys.argv[2].replace(" ","").replace("/","").replace("\\","").replace("`","").replace("~","").replace("@","").replace("#","").replace("$","").replace("%","").replace("^","").replace("&","").replace("*","").replace("(","").replace(")","").replace("-","").replace("+","").replace("=","").replace("[","").replace("]","").replace("{","").replace("}","").replace("|","").replace(":","").replace(";","").replace("'","").replace("\"","").replace("<","").replace(">","").replace(",","").replace(".","").replace("?","").replace("/","").replace("!","")
key = key[:35]

translation = sys.argv[3]
print key