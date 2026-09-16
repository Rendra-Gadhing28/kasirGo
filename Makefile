.PHONY: all up down restart build logs seed test clean

all: up

up:
	docker compose up -d

down:
	docker compose down

restart:
	docker compose restart

build:
	docker compose build

logs:
	docker compose logs -f

clean:
	docker compose down -v

seed:
	docker compose exec backend ./server --seed
