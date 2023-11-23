Version=1.0
ldFlag=-s -w -X main.Version=$(Version)
Topdir=$(CURDIR)

# must run in darwin
darwin:
	cd apps/ui && GOOS=darwin go build -ldflags "$(ldFlag)" -o $(Topdir)/dist/antiAv_darwin

# must run in windows
win:
	cd apps/ui && GOOS=windows go build -ldflags "$(ldFlag)" -o $(Topdir)/dist/antiAv.exe

# must run in linux
linux:
	cd apps/ui && GOOS=linux go build -ldflags "$(ldFlag)" -o $(Topdir)/dist/antiAv_linux