Create Migration File

Install `go-migrate`

https://github.com/golang-migrate/migrate

```
migrate create -seq -ext=.sql -dir=./database/migrations <migration_name>
```

or

```
make migration/create name=alter_user_table_password_type
make migration/create name=<migration_name>
```

Create Migration for Postgresql database:

```
migrate -source <migration_path> -database "postgres://postgres:postgres@localhost/postgres?sslmode=disable" up

```

for starting postgresql database:
```
sudo docker compose up -d
```
