import json
import os
import io
import glob
import sys

languages = ["es", "fr", "ru", "de", "it", "sv", "ro", "pt", "hu", "nl", "ar", "ko", "ja", "zh"]

bullshitTestFlag = "!"
#
pattern = sys.argv[1]

files = glob.glob(os.getenv("GOPATH") + '/src/github.com/DanielRenne/goCoreAppTemplate/web/app/globalization/translations/' + pattern + '/en/US.json')

for file in files:
    if "appError" in file or "transaction" in file:
        continue
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
            if lang in ["es", "ru", "fr", "de", "it", "sv", "pt", "hu", "nl"]:
                for k, v in i18n_translations.iteritems():
                    i18n_translations[k] = ""
                iterate = False

        if iterate:
            for k, v in translations.iteritems():
                ok = True
                try:
                    i18n_translations[k]
                except:
                    ok = False
                if ok and i18n_translations[k].startswith(bullshitTestFlag) and lang in ["es", "ru", "fr", "de", "it", "sv", "pt", "hu", "nl"]:
                    i18n_translations[k] = ""
                else:
                    if not ok and lang in ["es", "ru", "fr", "de", "it", "sv", "pt", "hu", "nl"] or ok and i18n_translations[k].startswith(bullshitTestFlag) and lang in ["es", "ru", "fr", "de", "it", "sv", "pt", "hu", "nl"]:
                        # Leave blank so translate_all can run translations on new items
                        i18n_translations[k] = ""
                    else:
                        if ok and not i18n_translations[k].startswith(bullshitTestFlag) and lang in ["es", "ru", "fr", "de", "it", "sv", "pt", "hu", "nl"]:
                            pass
                        else:
                            i18n_translations[k] = bullshitTestFlag + lang + bullshitTestFlag + v
        try:
            os.makedirs("/".join(parts[:-2]) + "/" + lang)
        except:
            pass
        json.dump(i18n_translations, open(file.replace("/en/US.json", "/" + lang + "/" + lang + ".json"), 'w'), indent=4, sort_keys=True)