Create Migration File

Install `go-migrate`

https://github.com/golang-migrate/migrate

```
migrate create -seq -ext=.sql -dir=./database/migrations <migration_name>
```

Create Migration for Postgresql database:

```
migrate -source <migration_path> -database "postgres://postgres:postgres@localhost/postgres?sslmode=disable" up

```

for starting postgresql database:
```
sudo docker compose up -d
```
