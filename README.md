```
export DATABASE_URL="postgres://admin:secret@localhost:5433/marketplace"

docker run --name marketplace-db -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=marketplace -p 5433:5432 -d postgres

docker run --name marketplace-redis -p 6370:6379 -d redis

go run main.go
```
or
```
export DATABASE_URL="postgres://admin:secret@localhost:5433/marketplace"

make run
```