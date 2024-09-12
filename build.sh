set GOOS=darwin
set GOARCH=amd64
go build -o myprogram_darwin_amd64
set GOOS=darwin
set GOARCH=arm64
go build -o myprogram_darwin_arm64
set GOOS=linux
set GOARCH=amd64
go build -o myprogram_linux_amd64
set GOOS=linux
set GOARCH=arm
go build -o myprogram_linux_arm
set GOOS=linux
set GOARCH=arm64
go build -o myprogram_linux_arm64
#