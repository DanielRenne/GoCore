from yandex_translate import YandexTranslate
import json
import os
import re
import sys
import glob
from language_contant import languages
translate = YandexTranslate('trnsl.1.1.20171006T195350Z.3f62cfa260827868.0ebb36cd5081b66bd7f63cb172d785a5bce5bbee')

if len(sys.argv) == 1:
    pattern = "*"
else:
    pattern = sys.argv[1]

initialDump = False
log = False

bullshitTestFlag = "!"
#
files = glob.glob(os.getenv("GOPATH") + '/src/github.com/DanielRenne/goCoreAppTemplate/web/app/globalization/translations/' + pattern + '/en/US.json')

for file in files:
    if log:
        print "Working on " + file
    parts = file.split("/")
    for lang in languages:
        translations = json.load(open(file, 'r'))
        i18n_translations = json.load(open(file.replace("/en/US.json", "/" + lang + "/" + lang + ".json"), 'r'))
        for k, v in translations.iteritems():
            try:
                i18n_translations[k]
            except:
                i18n_translations[k] = ""
            if initialDump or (not initialDump and i18n_translations[k] == ""):

                if log:
                    print "Make request for (" + lang + "): " + v
                try:
                    response = translate.translate(v, 'en-' + lang)
                    if response["code"] == 200:
                        try:
                            i18n_translations[k] = response["text"][0]
                            matchesEn = re.findall(r"\{([^}]+)\}", v)
                            matchesI18n = re.findall(r"\{([^}]+)\}", i18n_translations[k])
                            if len(matchesEn) > 0:
                                matchesLang = re.findall(r"\{([^}]+)\}", i18n_translations[k])
                                for kk, vv in enumerate(matchesEn):
                                    i18n_translations[k] = i18n_translations[k].replace("{" + matchesLang[kk] + "}", "{" + vv + "}")
                            #print lang + ":" + i18n_translations[k]
                        except:
                            i18n_translations[k] = ""

                            if log:
                                print "Index error"
                                print response
                    else:
                        i18n_translations[k] = ""

                        if log:
                            print "err!!!!!!!!!!"
                            print response
                except Exception:
                    i18n_translations[k] = ""
                    if log:
                        print "exception!!!"
        try:
            os.makedirs("/".join(parts[:-2]) + "/" + lang)
        except:
            pass
        ff = file.replace("/en/US.json", "/" + lang + "/" + lang + ".json")
        if log:
            print "writing to" + ff
        json.dump(i18n_translations, open(ff, 'w'), indent=4, sort_keys=True)
