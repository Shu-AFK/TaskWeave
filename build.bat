IF NOT EXIST
    out mkdir out
set CGO_ENABLED=1
go build -o .\out\TaskWeave .\cmd\web\main.go