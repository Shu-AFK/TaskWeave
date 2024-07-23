IF NOT EXIST db (
    mkdir db
)

set CGO_ENABLED=1
go run .\cmd\web\main.go