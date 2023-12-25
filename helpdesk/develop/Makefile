build:
	go build -o bin/helpdesk ./cmd/helpdesk/helpdesk.go

run: build
	./bin/helpdesk

dbinit:
	go run ./cmd/database/database.go init

dbdrop:
	go run ./cmd/database/database.go drop

dbdelete:
	rm app.db
