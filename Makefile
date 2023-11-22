Version=1.0
ldFlag=-s -w -X main.Version=$(Version)
Topdir=$(CURDIR)

all: darwin linux win

darwin:res
	cd apps/ui && GOOS=darwin go build -ldflags "$(ldFlag)" -o $(Topdir)/dist/antiAv_darwin

win:res
	cd apps/ui && GOOS=windows go build -ldflags "$(ldFlag)" -o $(Topdir)/dist/antiAv.exe

linux:res
	cd apps/ui && GOOS=linux go build -ldflags "$(ldFlag)" -o $(Topdir)/dist/antiAv_linux

res:
	cp -rf apps/ui/resources $(Topdir)/dist/