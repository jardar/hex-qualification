# usage: make ver=12345
# variable
server_std_exec = api-server.elf

# -ldflags
LDFLAGS=-ldflags "-X main.ver=${ver} -X main.minversion=`date +.%Y%m%d.%H%M%S`"

default: build

clean:
	rm -f $(server_std_exec)
	rm -f *.log

prebuild:
	# qtc -dir qt -ext qtpl
	# cd ../proto ;make;cd ../game-server;

build: clean
	go build ${LDFLAGS} -o $(server_std_exec) main.go

buildlinux: clean
	#gox -osarch="linux/amd64" -output="$(server_std_exec)" ${LDFLAGS}
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS_v2} -o $(server_std_exec) main.go

run: clean build
	./$(server_std_exec)



