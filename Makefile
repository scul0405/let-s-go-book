server:
	go run ./cmd/web
db:
	psql -h localhost -d snippetbox -p 5432 -U web 

.PHONY: server db