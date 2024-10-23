# Load environment variables from the .env file
ifneq (,$(wildcard ./.env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

# Database connection variables
DB_URL=mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST))/$(DB_NAME)


# Migration commands
migration_create:
	@test -n "$(name)" || (echo "Error: migration name is not set. Use 'make migration_create name=<migration_name>'"; exit 1)
	migrate create -ext sql -dir database/migration/ -seq $(name)

migration_up:
	migrate -path database/migration/ -database "$(DB_URL)" -verbose up

migration_down:
	migrate -path database/migration/ -database "$(DB_URL)" -verbose down

migration_fix:
	migrate -path database/migration/ -database "$(DB_URL)" force VERSION

migration_force:
	migrate -path database/migration/ -database "$(DB_URL)" force 1 

# Run the server
run_server:
	air