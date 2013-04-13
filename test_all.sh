go test -i ./... &&
go test -v &&
cd cmd && go test -v && cd .. &&
cd example && go test -v && cd .. &&
cd interpreter && go test -v && cd ..
