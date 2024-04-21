Postgres Container:

```
docker run \
  --name coerhat-db \
  --network coerhat \
  --volume coerhat-db-data:/var/lib/postgresql/data \
  -p 5432:5432 \
  --rm \
  --detach \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=coerhat \
  postgres
```

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
