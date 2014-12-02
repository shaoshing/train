go test -i ./... &&
go test -v &&
pushd cmd/train && go test -v && popd &&
pushd example && go test -v && popd &&
pushd interpreter && go test -v && popd
