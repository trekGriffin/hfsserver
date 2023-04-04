rm bin/*
#go build -o bin/activation$(git describe --tags --dirty).exe -ldflags "-X main.appVersion=$(git describe --tags --dirty) -X 'main.appDate=$(date)' "  
go env -w GOOS=linux
go build -o bin/hfsserver$(git describe --tags --dirty) -ldflags "-X main.appVersion=$(git describe --tags --dirty) -X 'main.appDate=$(date)' "  
go env -w GOOS=windows
