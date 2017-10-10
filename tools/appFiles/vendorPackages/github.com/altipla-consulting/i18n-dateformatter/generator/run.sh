#!/bin/bash

set -eu

echo " [*] Install generator..."
go install github.com/altipla-consulting/i18n-dateformatter/generator

if [ ! -e /tmp/core.zip ]; then
  echo " [*] Download CLDR data..."
  wget http://www.unicode.org/Public/cldr/27.0.1/core.zip -O /tmp/core.zip
fi

echo " [*] Generate "
#english, spanish, french, russian, german, italian, swedish, romanian, portuguese, hungarian, netherlands (Dutch), Arabic, korean, japanese, chinese
generator -locales en,es,fr,ru,de,it,sv,ro,pt,hu,nl,ar,ko,ja,zh
gofmt -w symbols
