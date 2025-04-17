up:
	docker-compose up --build

down:
	docker-compose down -v

migrate-up:
	goose -dir ./migrations postgres "host=localhost port=54323 dbname=auth user=user password=1234" up

