oss=(darwin linux windows)
archs=(386 amd64)
for i in ${oss[@]}
do
  for j in ${archs[@]}
  do
    echo "CGO_ENABLED=0 GOOS=${i} GOARCH=${j} go install github.com/lastsweetop/gitdown"
    CGO_ENABLED=0 GOOS=${i} GOARCH=${j} go install github.com/lastsweetop/gitdown
  done
done

cp -R ../../../../bin/* bin/
