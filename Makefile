.PHONY: build test vet lint fmt check clean serve migrate

build:
	go build -o bin/supost .

test:
	go test ./... -race -coverprofile=coverage.out

vet:
	go vet ./...

lint:
	golangci-lint run

fmt:
	gofmt -w .

check: fmt vet build test
	@echo "All checks passed."

serve:
	go run . serve

migrate:
	@echo "Apply migrations to your database with Supabase CLI:"
	@echo "  supabase db push --db-url \"$$DATABASE_URL\""

clean:
	rm -rf bin/ coverage.out
