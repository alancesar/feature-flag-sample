dep:
	go mod download
build: dep
	go build -o service
grep:
	git grep ".*if.*Eval("