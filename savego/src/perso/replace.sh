#!/bin/bash
for file in *.go
do
  echo "Traitement de $file ..."
  sed -i -e "s/\$[0-9]/?/g" "$file"
done 