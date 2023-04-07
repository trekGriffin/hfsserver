set productname hfsserver$(git describe --tags --dirty)
rm bin/*
#go build -o bin/activation$(git describe --tags --dirty).exe -ldflags "-X main.appVersion=$(git describe --tags --dirty) -X 'main.appDate=$(date)' "  
go env -w GOOS=linux
go build -o bin/$productname -ldflags "-X main.appVersion=$(git describe --tags --dirty) -X 'main.appDate=$(date)' "  
go env -w GOOS=windows
go build -o bin/$productname.exe -ldflags "-X main.appVersion=$(git describe --tags --dirty) -X 'main.appDate=$(date)' "  

curl.exe -T ./bin/$productname http://10.10.10.3/upload/$productname
