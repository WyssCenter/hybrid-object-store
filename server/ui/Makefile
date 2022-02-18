BUILD_DATE=$(shell date -u "+%Y%m%d")
export BUILD_DATE


build-debug:
	GOOS=linux go build -gcflags "-N -l" -a --ldflags "-X main.version=%(BUILD_DATE) -extldflags  -static" -tags render -o ./render ./

build:
	GOOS=linux go build -a --ldflags "-X main.version=%(BUILD_DATE) -extldflags  -static" -tags render -o ./render ./
