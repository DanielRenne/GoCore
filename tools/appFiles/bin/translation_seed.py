import json
import os
import io
import glob
import sys
from language_contant import languages
pattern = sys.argv[1]
blastKey = sys.argv[2]

files = glob.glob(os.getenv("GOPATH") + '/src/github.com/DanielRenne/goCoreAppTemplate/web/app/globalization/translations/' + pattern + '/en/US.json')

for file in files:
    parts = file.split("/")
    for lang in languages:
        translations = json.load(io.open(file, 'r'))
        fname = file.replace("/en/US.json", "/" + lang + "/" + lang + ".json")
        iterate = True
        if os.path.isfile(fname):
            # print "Read" + fname
            i18n_translations = json.load(io.open(fname, 'r'))
        else:
            i18n_translations = translations
            for k, v in i18n_translations.iteritems():
                i18n_translations[k] = ""
            iterate = False

        if iterate:
            for k, v in translations.iteritems():

                ok = True
                try:
                    i18n_translations[k]
                    if blastKey == k:
                        i18n_translations[k] = ""
                except:
                    ok = False
        try:
            os.makedirs("/".join(parts[:-2]) + "/" + lang)
        except:
            pass
        json.dump(i18n_translations, open(file.replace("/en/US.json", "/" + lang + "/" + lang + ".json"), 'w'), indent=4, sort_keys=True)