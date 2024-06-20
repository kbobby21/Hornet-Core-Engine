# UI-Backend

## Dependencies
Install all these dependencies
1. [Golang](https://go.dev/doc/install).  
2. [Postgres](https://www.postgresql.org/download/).  
3. [Migrate](https://github.com/golang-migrate/migrate).   
4. [Redis](https://redis.io/docs/getting-started/installation/install-redis-on-linux/)

## Run the server
1. Create a new user `hornet` with password `H0rnSt@r` from terminal and add it to sudoers list. [Learn more.](https://www.digitalocean.com/community/tutorials/how-to-create-a-new-sudo-enabled-user-on-ubuntu-20-04-quickstart)
2. Create a new role named `hornet`. [Learn more](https://www.digitalocean.com/community/tutorials/how-to-install-postgresql-on-ubuntu-20-04-quickstart)
3. Create the db if not already created, using `psql -h localhost -U postgres -W -c "create database hornet;"` or run `createdb hornet` from terminal when logged in a user `hornet`.
4. Provide configuration in the file `~/.hornet/config-{dev or prod}.json`in the root of the new user. 
5. Clone this repository, then `cd` to the root of this repository, then run  
	1. Export the postgres connection URL: `export POSTGRESQL_URL='postgres://hornet:H0rnSt@r@localhost:5432/hornet?sslmode=disable'`  
	2. Finally apply the migrations by: `migrate -database ${POSTGRESQL_URL} -path db/migrations up`  

6. `go run main.go env={dev or prod}`.  



