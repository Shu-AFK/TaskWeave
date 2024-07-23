IF NOT EXIST out (
    mkdir out
)

IF NOT EXIST db (
    mkdir db
)

set CGO_ENABLED=1
go build -o .\out\TaskWeave .\cmd\web\main.go