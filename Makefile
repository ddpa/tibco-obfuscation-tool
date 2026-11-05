APP_NAME = TiBCO Obfuscation Tool
APP_ID   = com.diogo.tibcoobfuscationtool

.PHONY: run app universal clean
run:                      ## run the GUI without packaging
	go run .

app:                      ## build a native .app
	./build-macos.sh

universal:                ## build an Intel+Apple-Silicon universal .app
	CGO_ENABLED=1 GOARCH=arm64 go build -o /tmp/tibco-arm64 .
	CGO_ENABLED=1 GOARCH=amd64 go build -o /tmp/tibco-amd64 .
	lipo -create -output /tmp/tibco-universal /tmp/tibco-arm64 /tmp/tibco-amd64
	fyne package -os darwin -name "$(APP_NAME)" -appID "$(APP_ID)" -icon Icon.png -executable /tmp/tibco-universal

clean:
	rm -rf "$(APP_NAME).app" /tmp/tibco-arm64 /tmp/tibco-amd64 /tmp/tibco-universal
