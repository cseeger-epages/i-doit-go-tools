if [[ ! -d "bin" ]]; then
  mkdir bin
fi

for f in ./*/*.go;
  do
    echo "build $f"
    go build $f
    mv $(basename ${f::-3}) bin/
done
