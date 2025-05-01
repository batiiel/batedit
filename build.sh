export GOOS=linux
go build -o batedit main.go window.go editor.go
export GOOS=windows
go build -o batedit.exe main.go window.go editor.go
