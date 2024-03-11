.PHONY: setup-go get-deps build-linux build-macos test-linux test-macos

CARBON_VERSION := 1.33.0

setup-go:
	go version

get-deps:
	go get -v -t -d ./...
	if [ -f Gopkg.toml ]; then \
		curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh; \
		dep ensure; \
	fi

build-linux:
	mkdir -p $$HOME/carbon && \
	cd $$HOME/carbon && \
	wget https://csspeechstorage.blob.core.windows.net/drop/$(CARBON_VERSION)/SpeechSDK-Linux-$(CARBON_VERSION).tar.gz && \
	tar xzf SpeechSDK-Linux-$(CARBON_VERSION).tar.gz && \
	rm SpeechSDK-Linux-$(CARBON_VERSION).tar.gz && \
	ln -s SpeechSDK-Linux-$(CARBON_VERSION) current
	sudo apt-get update
	sudo apt-get install libasound2
	sudo apt-get install -y libgstreamer1.0-0 gstreamer1.0-plugins-good
	export CGO_CFLAGS="-I$$HOME/carbon/current/include/c_api" && \
	export CGO_LDFLAGS="-L$$HOME/carbon/current/lib/x64 -lMicrosoft.CognitiveServices.Speech.core" && \
	export LD_LIBRARY_PATH="$$LD_LIBRARY_PATH:$$HOME/carbon/current/lib/x64" && \
	go build -v ./...

build-macos:
	cd ~ && \
	wget https://csspeechstorage.blob.core.windows.net/drop/$(CARBON_VERSION)/MicrosoftCognitiveServicesSpeech-XCFramework-$(CARBON_VERSION).zip -O SpeechSDK.zip && \
	unzip SpeechSDK.zip -d speechsdk
	export SDK_HOME="$$HOME/speechsdk/MicrosoftCognitiveServicesSpeech.xcframework/macos-arm64_x86_64" && \
	export CGO_CFLAGS="-I$$SDK_HOME/MicrosoftCognitiveServicesSpeech.framework/Headers" && \
	export CGO_LDFLAGS="-F$$SDK_HOME -framework MicrosoftCognitiveServicesSpeech" && \
	export DYLD_FRAMEWORK_PATH="$$DYLD_FRAMEWORK_PATH:$$SDK_HOME" && \
	go build -v ./...

test-linux:
	export CGO_CFLAGS="-I$$HOME/carbon/current/include/c_api" && \
	export CGO_LDFLAGS="-L$$HOME/carbon/current/lib/x64 -lMicrosoft.CognitiveServices.Speech.core" && \
	export LD_LIBRARY_PATH="$$LD_LIBRARY_PATH:$$HOME/carbon/current/lib/x64" && \
	go test -v ./...

test-macos:
	export SDK_HOME="$$HOME/speechsdk/MicrosoftCognitiveServicesSpeech.xcframework/macos-arm64_x86_64" && \
	export CGO_CFLAGS="-I$$SDK_HOME/MicrosoftCognitiveServicesSpeech.framework/Headers" && \
	export CGO_LDFLAGS="-F$$SDK_HOME -framework MicrosoftCognitiveServicesSpeech" && \
	export DYLD_FRAMEWORK_PATH="$$DYLD_FRAMEWORK_PATH:$$SDK_HOME" && \
	go test -v ./...
