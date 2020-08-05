#!/usr/bin/env bash

for f in `find . -iname "*.go"`; do # just for *.go files
  echo -e "/*\nCopyright 2020 Getup Cloud. All rights reserved.\n*/\n" > tmpfile
  cat $f >> tmpfile
  mv tmpfile $f
done