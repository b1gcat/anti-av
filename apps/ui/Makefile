Output=dist

all:

darwin:
	mkdir -p $(Output)/darwin
	fyne package -os darwin -icon resources/logo.png -name "免杀生成器" -appBuild 1
	cp -rf resources 免杀生成器.app/Contents/MacOS/
	rm -rf $(Output)/darwin/免杀生成器.app
	mv 免杀生成器.app $(Output)/darwin/免杀生成器.app
	rm -rf helper

win:
	fyne package -os windows -icon resources/logo.png -name "免杀生成器" -appBuild 1