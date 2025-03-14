# Default url: http://localhost:7331

live/server:
	go run github.com/air-verse/air@latest

live/tailwind:
	npx @tailwindcss/cli -i ./static/css/input.css -o ./static/css/output.css -w -m >/dev/null 2>&1

live: 
	make -j2 live/tailwind live/server



build/templ:
	go run github.com/a-h/templ/cmd/templ@latest generate

build/sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate

build/tailwindcss:
	npm install
	npx @tailwindcss/cli -i ./static/css/input.css -o ./build/static/css/output.css -m

build/go:
	go build -ldflags "-s -w" -o ./build/carshift main.go

build/mini:
	upx ./build/carshift

build:
	mkdir build
	cp static build -r
	make build/templ build/sqlc build/tailwindcss build/go build/mini
