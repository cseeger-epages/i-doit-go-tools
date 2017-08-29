#!/bin/bash

if [[ ! -d "bin" ]]; then
  mkdir bin
fi

if [ "$(grep -R '<apikey>' ./* | wc -l)" -gt "1" ]; then
  echo "please set proper login information !"
  echo 'change all goidoit.NewLogin("<i-doit-url>/src/jsonrpc.php", "<apikey>", "<username>", "<password>") to use proper i-doit-url, apikey and username password)'
  read -p "proceed anyway ? [y/n] " -n 1 -r
  echo 
  if [[ $REPLY =~ ^[Yy]$ ]]; then
    for f in ./*/*.go;
      do
        echo "build $f"
        go build $f
        mv $(basename ${f::-3}) bin/
    done
  fi
  exit 1
fi

