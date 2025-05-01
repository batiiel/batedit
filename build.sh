export GOOS=linux
go build -o ego main.go window.go editor.go
export GOOS=windows
go build -o ego.exe main.go window.go editor.go
