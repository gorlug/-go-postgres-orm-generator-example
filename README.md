# Go Postgres ORM generator example

## Instructions

This example uses a Postgres as the database backend.

To get up running quickly use Docker to start a Postgres instance:

```
cd docker/go-postgres-orm-generator
docker compose up -d
```

Add this .env File:

```
DATABASE_URL="postgres://postgres:local@127.0.0.1:5432/example"
```

Run

```
go generate
```

To generate the Prisma schema and the ORM code.

The ORM generation code can be modified in the generator package.

To set up the database, run this script:

```
./db_push.sh
```

Running

```
go run main.go
```

Runs the simple CRUD operations from the generated code.

## Resources

* https://goprisma.org/
