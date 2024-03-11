.PHONY: deps build test

CARBON_VERSION := 1.33.0

deps:
	@echo "Installing dependencies..."
	@if [ "$(GOOS)" = "linux" ]; then \
		mkdir -p $$HOME/carbon && \
		cd $$HOME/carbon && \
        wget https://csspeechstorage.blob.core.windows.net/drop/$(CARBON_VERSION)/SpeechSDK-Linux-$(CARBON_VERSION).tar.gz && \
		tar xzf SpeechSDK-Linux-$(CARBON_VERSION).tar.gz && \
		rm SpeechSDK-Linux-$(CARBON_VERSION).tar.gz && \
		ln -s SpeechSDK-Linux-$(CARBON_VERSION) current; \
		sudo yum install -y alsa-lib && \
		sudo yum install -y gstreamer1 gstreamer1-plugins-good; \
    elif [ "$(GOOS)" = "darwin" ]; then \
		cd ~ && \
		wget https://csspeechstorage.blob.core.windows.net/drop/$(CARBON_VERSION)/MicrosoftCognitiveServicesSpeech-XCFramework-$(CARBON_VERSION).zip -O SpeechSDK.zip && \
		unzip SpeechSDK.zip -d speechsdk; \
   	fi

build:
	@echo "Building..."
	@if [ "$(GOOS)" = "linux" ]; then \
		export CGO_CFLAGS="-I$$HOME/carbon/current/include/c_api" && \
		export CGO_LDFLAGS="-L$$HOME/carbon/current/lib/x64 -lMicrosoft.CognitiveServices.Speech.core" && \
		export LD_LIBRARY_PATH="$$LD_LIBRARY_PATH:$$HOME/carbon/current/lib/x64" && \
		go build -v ./...; \
	elif [ "$(GOOS)" = "darwin" ]; then \
		export SDK_HOME="$$HOME/speechsdk/MicrosoftCognitiveServicesSpeech.xcframework/macos-arm64_x86_64" && \
		export CGO_CFLAGS="-I$$SDK_HOME/MicrosoftCognitiveServicesSpeech.framework/Headers" && \
		export CGO_LDFLAGS="-F$$SDK_HOME -framework MicrosoftCognitiveServicesSpeech" && \
		export DYLD_FRAMEWORK_PATH="$$DYLD_FRAMEWORK_PATH:$$SDK_HOME" && \
		go build -v ./...; \
	fi

test:
	@if [ "$(GOOS)" = "linux" ]; then \
		export CGO_CFLAGS="-I$$HOME/carbon/current/include/c_api" && \
		export CGO_LDFLAGS="-L$$HOME/carbon/current/lib/x64 -lMicrosoft.CognitiveServices.Speech.core" && \
		export LD_LIBRARY_PATH="$$LD_LIBRARY_PATH:$$HOME/carbon/current/lib/x64" && \
		go test -v ./...; \
	elif [ "$(GOOS)" = "darwin" ]; then \
		export SDK_HOME="$$HOME/speechsdk/MicrosoftCognitiveServicesSpeech.xcframework/macos-arm64_x86_64" && \
		export CGO_CFLAGS="-I$$SDK_HOME/MicrosoftCognitiveServicesSpeech.framework/Headers" && \
		export CGO_LDFLAGS="-F$$SDK_HOME -framework MicrosoftCognitiveServicesSpeech" && \
		export DYLD_FRAMEWORK_PATH="$$DYLD_FRAMEWORK_PATH:$$SDK_HOME" && \
		go test -v ./...; \
	fi