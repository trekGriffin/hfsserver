set productname hfsserver
rm bin/*
#go build -o bin/activation$(git describe --tags --dirty).exe -ldflags "-X main.appVersion=$(git describe --tags --dirty) -X 'main.appDate=$(date)' "  
go env -w GOOS=linux
go build -o bin/$productname$(git describe --tags --dirty) -ldflags "-X main.appVersion=$(git describe --tags --dirty) -X 'main.appDate=$(date)' "  
go env -w GOOS=windows
curl -T ./bin/$productname$(git describe --tags --dirty) http://10.10.10.3/upload/$productname
