# note: Makefile uses tabs as syntax, GoLand replaces tabs with spaces, so it cannot be edited with it

# local
AgentService:
	go build -o ./build/AgentService .

# cross-platform compile Linux
AgentService-amd64:
	GOARCH=amd64 GOOS=linux go build -o ./build/AgentService-linux .

# local all
all-proc: AgentService

# Linux all
all-amd64: AgentService-amd64

# all
all: all-proc all-amd64

# test
run:
	./build/AgentService --config ./build/conf.ini start

version:
	./build/AgentService --config ./build/conf.ini version

tool:
	./build/AgentService --config ./build/conf.ini tool $(var)

# clean all
clean:
	rm -rf ./build/AgentService
	rm -rf ./build/AgentService-linux
