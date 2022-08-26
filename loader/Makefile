CC=i686-w64-mingw32-gcc
CXX=i686-w64-mingw32-g++
AR=i686-w64-mingw32-ar

ifeq ($(ARCH),amd64)
	CC=x86_64-w64-mingw32-gcc
	CXX=x86_64-w64-mingw32-g++
	AR=x86_64-w64-mingw32-ar
endif

all:
	@if [ $(ASM) = true -a $(ARCH) = amd64 ]; then \
		echo "Not Support Now amd64 asm!"; \
		exit 1;\
	fi;

	@if [ $(ASM) = true -a $(ARCH) = 386 ]; then \
		asmFile=""; \
		cp asm/386/*  $(LOADER); \
		asmFile="*.s"; \
	fi;

	cd $(LOADER) && \
		rm -f *.o *.a *.dll && \
		$(CC) -Wall -c *.c $(asmFile) && \
		$(AR) rv lib$(LOADER).a *.o

	CGO_ENABLED=1 CC=$(CC) CXX=$(CXX) GOOS=$(OS) GOARCH=$(ARCH) \
		garble build  -trimpath -ldflags "-s -w" -o $(OUTFILE)