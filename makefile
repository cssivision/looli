test:
	go test ./... -v -race

cover:
	rm -f *.coverprofile
	go test -coverprofile=looli.coverprofile
	go test -coverprofile=cors.coverprofile ./cors
	go test -coverprofile=csrf.coverprofile ./csrf
	go test -coverprofile=session.coverprofile ./session
	gover
	go tool cover -html=gover.coverprofile
	rm -f *.coverprofile

.PHONY: test cover