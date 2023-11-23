CC=i686-w64-mingw32-gcc
CXX=i686-w64-mingw32-g++
AR=i686-w64-mingw32-ar

ifeq ($(ARCH),amd64)
	CC=x86_64-w64-mingw32-gcc
	CXX=x86_64-w64-mingw32-g++
	AR=x86_64-w64-mingw32-gcc-ar
endif

all:
	cd $(LOADER) && \
		rm -f *.o *.a *.dll && \
		$(CC) -Wall -c *.c $(asmFile) && \
		$(AR) rv lib$(LOADER).a *.o &&\
		go mod init &&\
		echo "replace github.com/b1gcat/anti-av/apps/av/utils => ../../../apps/av/utils" >> go.mod &&\
		go mod tidy

	go mod init
	echo "replace github.com/b1gcat/anti-av/dist/_tmp/$(LOADER)  => ./sc" >> go.mod
	echo "replace github.com/b1gcat/anti-av/apps/av/utils => ../../apps/av/utils" >> go.mod
	go mod tidy
	CGO_ENABLED=1 CC=$(CC) CXX=$(CXX) GOOS=$(OS) GOARCH=$(ARCH) \
		garble build  -trimpath -ldflags "-s -w" -o $(OUTFILE)