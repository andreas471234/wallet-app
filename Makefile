# Version - this is optionally used on goto command
V?=

# Number of migrations - this is optionally used on up and down commands
N?=

# In a real world scenario, these environment variables
# would be injected by your build tool, like Drone for example (https://drone.io/)
MYSQL_USER ?= root
MYSQL_PASSWORD ?= 
MYSQL_HOST ?= localhost
MYSQL_DATABASE ?= walletapp
MYSQL_PORT ?= 3306

MYSQL_DSN ?= $(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DATABASE)

local-db:
	@ until mysql --host=$(MYSQL_HOST) --port=$(MYSQL_PORT) --user=$(MYSQL_USER) -p$(MYSQL_PASSWORD) --protocol=tcp -e 'SELECT 1' >/dev/null 2>&1 && exit 0; do \
	  >&2 echo "MySQL is unavailable - sleeping"; \
	  sleep 5 ; \
	done

	@ echo "MySQL is up and running!"

migrate-setup:
	@if [ -z "$$(which migrate)" ]; then echo "Installing migrate command..."; go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate; fi

migrate-up: migrate-setup
	@ migrate -database 'mysql://$(MYSQL_DSN)?multiStatements=true' -path migrations up $(N)

migrate-down: migrate-setup
	@ migrate -database 'mysql://$(MYSQL_DSN)?multiStatements=true' -path migrations down $(N)

migrate-to-version: migrate-setup
	@ migrate -database 'mysql://$(MYSQL_DSN)?multiStatements=true' -path migrations goto $(V)

drop-db: migrate-setup
	@ migrate -database 'mysql://$(MYSQL_DSN)?multiStatements=true' -path migrations drop

force-version: migrate-setup
	@ migrate -database 'mysql://$(MYSQL_DSN)?multiStatements=true' -path migrations force $(V)

migration-version: migrate-setup
	@ migrate -database 'mysql://$(MYSQL_DSN)?multiStatements=true' -path migrations version

build:
	@ go build tests/inspect_database.go

run: build
	@ ./inspect_database