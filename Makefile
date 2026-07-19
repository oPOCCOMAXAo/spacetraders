build:
	go build -o bin/stbot cmd/stbot/main.go

lint:
	golangci-lint-v2 run

format-templ:
	templ fmt .

generate-templ: format-templ
	templ generate

generate: generate-templ
	go generate ./...

# automatically regenerate templ files on changes
dev-watch-templ:
	templ generate --watch
