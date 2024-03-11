package speech

import (
	"encoding/json"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"testing"
	"time"
)

func TestPronunciationAssessment(t *testing.T) {
	format, err := audio.GetDefaultInputFormat()
	if err != nil {
		t.Error("Got an error ", err.Error())
	}
	defer format.Close()
	stream, err := audio.CreatePushAudioInputStreamFromFormat(format)
	if err != nil {
		t.Error("Got an error ", err.Error())
	}
	defer stream.Close()
	audioConfig, err := audio.NewAudioConfigFromStreamInput(stream)
	if err != nil {
		t.Error("Got an error ", err.Error())
	}
	defer audioConfig.Close()
	recognizer := createSpeechRecognizerFromAudioConfig(t, audioConfig)
	if recognizer == nil {
		return
	}
	defer recognizer.Close()

	pronunciationAssessmentConfig, err := NewPronunciationAssessmentConfig(
		"turn on the lamp",
		PronunciationAssessmentGradingSystemHundredMark,
		PronunciationAssessmentGranularityPhoneme,
		true,
	)
	pronunciationAssessmentConfig.EnableProsodyAssessment()
	pronunciationAssessmentConfig.EnableContentAssessmentWithTopic("light")
	if err != nil {
		t.Error("Got an error ", err.Error())
	}
	defer pronunciationAssessmentConfig.Close()
	if err := pronunciationAssessmentConfig.ApplyTo(recognizer); err != nil {
		t.Error("Got an error ", err.Error())
	}

	firstResult := true
	recognizedFuture := make(chan string)
	recognizingFuture := make(chan string)
	sessionStopedChan := make(chan string)
	recognizedHandler := func(event SpeechRecognitionEventArgs) {
		defer event.Close()
		firstResult = true
		t.Log("Recognized: ", event.Result.Text)
		pronunciationAssessmentResult, err := NewPronunciationAssessmentResultFromRecognitionResult(event.Result)
		if err != nil {
			t.Error("Got an error ", err.Error())
		}
		pronunciationAssessmentJson, _ := json.Marshal(pronunciationAssessmentResult.RecognitionResult)
		t.Log("PronunciationAssessmentResult: ", string(pronunciationAssessmentJson))

		recognizedFuture <- "Recognized"
	}
	recognizingHandle := func(event SpeechRecognitionEventArgs) {
		defer event.Close()
		t.Log("Recognizing: ", event.Result.Text)
		if firstResult {
			firstResult = false
			recognizingFuture <- "Recognizing"
		}
	}
	recognizer.SessionStopped(func(event SessionEventArgs) {
		sessionStopedChan <- event.SessionID
	})
	recognizer.Recognized(recognizedHandler)
	recognizer.Recognizing(recognizingHandle)
	err = <-recognizer.StartContinuousRecognitionAsync()
	if err != nil {
		t.Error("Got error: ", err)
	}
	pumpFileIntoStream(t, "../test_files/turn_on_the_lamp.wav", stream)
	pumpFileIntoStream(t, "../test_files/turn_on_the_lamp.wav", stream)
	pumpSilenceIntoStream(t, stream)
	stream.CloseStream()

	err = <-recognizer.StopContinuousRecognitionAsync()
	if err != nil {
		t.Error("Got error: ", err)
	}

	var sessionStopped bool
	for {
		select {
		case <-recognizingFuture:
			t.Log("Received Recognizing event.")
		case <-recognizedFuture:
			t.Log("Received Recognized event.")
		case <-sessionStopedChan:
			t.Log("Received SessionStopped event.")
			sessionStopped = true
		case <-time.After(5 * time.Second):
			t.Error("Didn't receive Recognizing or Recognized event.")
			sessionStopped = true
		}
		if sessionStopped {
			break
		}
	}

}
