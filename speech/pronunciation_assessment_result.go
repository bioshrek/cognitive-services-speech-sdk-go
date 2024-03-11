package speech

import (
	"encoding/json"
	"fmt"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

type PronunciationAssessmentResult struct {
	RecognitionResult *SpeechRecognitionJsonResult
}

func NewPronunciationAssessmentResultFromRecognitionResult(recognitionResult SpeechRecognitionResult) (*PronunciationAssessmentResult, error) {
	jsonResultStr := recognitionResult.Properties.GetProperty(common.SpeechServiceResponseJSONResult, "{}")
	fmt.Printf("jsonResultStr: %s\n", jsonResultStr)
	jsonResult := &SpeechRecognitionJsonResult{}
	err := json.Unmarshal([]byte(jsonResultStr), jsonResult)
	if err != nil {
		return nil, err
	}

	return &PronunciationAssessmentResult{
		RecognitionResult: jsonResult,
	}, nil
}

type SpeechRecognitionJsonResult struct {
	Id                string                          `json:"Id"`
	RecognitionStatus string                          `json:"RecognitionStatus"`
	Offset            int64                           `json:"Offset"`
	Duration          int64                           `json:"Duration"`
	Channel           int                             `json:"Channel"`
	DisplayText       string                          `json:"DisplayText"`
	SNR               float64                         `json:"SNR"`
	NBest             []*PronunciationAssessmentNBest `json:"NBest"`
}

type PronunciationAssessmentNBest struct {
	Confidence              float64                             `json:"Confidence"`
	Lexical                 string                              `json:"Lexical"`
	ITN                     string                              `json:"ITN"`
	MaskedITN               string                              `json:"MaskedITN"`
	Display                 string                              `json:"Display"`
	PronunciationAssessment *PronunciationAssessment            `json:"PronunciationAssessment"`
	Words                   []PronunciationAssessmentWordResult `json:"Words"`
	ContentAssessment       *ContentAssessmentResult            `json:"ContentAssessment"`
}

type PronunciationAssessment struct {
	AccuracyScore     float64 `json:"AccuracyScore"`
	FluencyScore      float64 `json:"FluencyScore"`      // only available for full text granularity
	ProsodyScore      float64 `json:"ProsodyScore"`      // only available for full text granularity
	CompletenessScore float64 `json:"CompletenessScore"` // only available for full text granularity
	PronScore         float64 `json:"PronScore"`         // only available for full text granularity

	ErrorType string                           `json:"ErrorType"`
	Feedback  *PronunciationAssessmentFeedback `json:"Feedback"`

	NBestPhonemes []*PronunciationAssessmentNBestPhoneme `json:"NBestPhonemes"`
}

type PronunciationAssessmentWordResult struct {
	Word                    string                                              `json:"Word"`
	Offset                  int64                                               `json:"Offset"`
	Duration                int64                                               `json:"Duration"`
	PronunciationAssessment *PronunciationAssessment                            `json:"PronunciationAssessment"`
	Syllables               []*PronunciationAssessmentSyllableLevelTimingResult `json:"Syllables"`
	Phonemes                []*PronunciationAssessmentPhonemeResult             `json:"Phonemes"`
}

type ContentAssessmentResult struct {
	GrammarScore    float64 `json:"GrammarScore"`
	VocabularyScore float64 `json:"VocabularyScore"`
	TopicScore      float64 `json:"TopicScore"`
}

type PronunciationAssessmentFeedback struct {
	Prosody *PronunciationAssessmentProsodyFeedback `json:"Prosody"`
}

type PronunciationAssessmentProsodyFeedback struct {
	Break      *PronunciationAssessmentProsodyFeedbackBreak      `json:"Break"`
	Intonation *PronunciationAssessmentProsodyFeedbackIntonation `json:"Intonation"`
}

type PronunciationAssessmentProsodyFeedbackBreak struct {
	ErrorTypes      []string `json:"ErrorTypes"`
	UnexpectedBreak *struct {
		Confidence float64 `json:"Confidence"`
	} `json:"UnexpectedBreak"`
	MissingBreak *struct {
		Confidence float64 `json:"Confidence"`
	} `json:"MissingBreak"`
	BreakLength int64 `json:"BreakLength"`
}

type PronunciationAssessmentProsodyFeedbackIntonation struct {
	ErrorTypes []string `json:"ErrorTypes"`
	Monotone   *struct {
		SyllablePitchDeltaConfidence float64 `json:"SyllablePitchDeltaConfidence"`
	} `json:"Monotone"`
}

type PronunciationAssessmentSyllableLevelTimingResult struct {
	Syllable                string                   `json:"Syllable"`
	Grapheme                string                   `json:"Grapheme"`
	PronunciationAssessment *PronunciationAssessment `json:"PronunciationAssessment"`
	Offset                  int64                    `json:"Offset"`
	Duration                int64                    `json:"Duration"`
}

type PronunciationAssessmentPhonemeResult struct {
	Phoneme                 string                   `json:"Phoneme"`
	PronunciationAssessment *PronunciationAssessment `json:"PronunciationAssessment"`
	Offset                  int64                    `json:"Offset"`
	Duration                int64                    `json:"Duration"`
}

type PronunciationAssessmentNBestPhoneme struct {
	Phoneme string  `json:"Phoneme"`
	Score   float64 `json:"Score"`
}
