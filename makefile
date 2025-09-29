up:
	docker compose up -d

up-build:
	docker compose up -d --build app

down:
	docker compose down

logs:
	docker compose logs -f app
