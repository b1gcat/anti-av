Version=1.0
ldFlag=-s -w -X main.Version=$(Version)
Topdir=$(CURDIR)

all: darwin linux win

darwin:res
	cd apps/ui && GOOS=darwin go build -ldflags "$(ldFlag)" -o $(Topdir)/dist/antiAv

win:res
	cd apps/ui && GOOS=windows go build -ldflags "$(ldFlag)" -o $(Topdir)/dist/antiAv.exe

linux:res
	cd apps/ui && GOOS=linux go build -ldflags "$(ldFlag)" -o $(Topdir)/dist/antiAv

res:
	cp -rf apps/ui/resources -o $(Topdir)/dist/