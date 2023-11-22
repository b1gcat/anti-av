module github.com/b1gcat/anti-av/apps/av

go 1.21.4

replace github.com/b1gcat/anti-av/apps/av/utils => ../av/utils

require (
	github.com/b1gcat/anti-av/apps/av/utils v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.9.3
)

require (
	github.com/akavel/rsrc v0.10.2 // indirect
	github.com/josephspurrier/goversioninfo v1.4.0 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
)
