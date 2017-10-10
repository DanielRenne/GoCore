from yandex_translate import YandexTranslate
import json
import os
import re
import sys
import glob
translate = YandexTranslate('trnsl.1.1.20170204T061220Z.e60d4406971eb449.b68a2aea336927188afd82eecc8cbe54455a3695')

languages = ["es", "ru", "fr", "de", "it", "sv", "pt", "hu", "nl"]
# languages = [ "ru" ]
if len(sys.argv) == 1:
    pattern = "*"
else:
    pattern = sys.argv[1]

# delete ro "ar", "ko", "ja", "zh"

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
                            if len(matchesEn) > 0:
                                matchesLang = re.findall(r"\{([^}]+)\}", translations[k])
                                for kk, vv in enumerate(matchesEn):
                                    i18n_translations[k] = translations[k].replace("{" + matchesLang[kk] + "}", "{" + vv + "}")
                            print lang + ":" + response["text"][0]
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
