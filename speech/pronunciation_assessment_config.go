package speech

import (
	"fmt"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"unsafe"
)

// #include <stdlib.h>
// #include <speechapi_c_common.h>
// #include <speechapi_c_property_bag.h>
// #include <speechapi_c_pronunciation_assessment_config.h>
import "C"

type PronunciationAssessmentGradingSystemType int

const (
	PronunciationAssessmentGradingSystemFivePoint   PronunciationAssessmentGradingSystemType = 1
	PronunciationAssessmentGradingSystemHundredMark PronunciationAssessmentGradingSystemType = 2
)

func (p PronunciationAssessmentGradingSystemType) ToCType() C.Pronunciation_Assessment_Grading_System {
	switch p {
	case PronunciationAssessmentGradingSystemFivePoint:
		return C.PronunciationAssessmentGradingSystem_FivePoint
	case PronunciationAssessmentGradingSystemHundredMark:
		return C.PronunciationAssessmentGradingSystem_HundredMark
	default:
		panic(fmt.Sprintf("Unknown PronunciationAssessmentGradingSystemType: %d", p))
	}
}

type PronunciationAssessmentGranularity int

const (
	PronunciationAssessmentGranularityPhoneme  PronunciationAssessmentGranularity = 1
	PronunciationAssessmentGranularityWord     PronunciationAssessmentGranularity = 2
	PronunciationAssessmentGranularityFullText PronunciationAssessmentGranularity = 3
)

func (p PronunciationAssessmentGranularity) ToCType() C.Pronunciation_Assessment_Granularity {
	switch p {
	case PronunciationAssessmentGranularityPhoneme:
		return C.PronunciationAssessmentGranularity_Phoneme
	case PronunciationAssessmentGranularityWord:
		return C.PronunciationAssessmentGranularity_Word
	case PronunciationAssessmentGranularityFullText:
		return C.PronunciationAssessmentGranularity_FullText
	default:
		panic(fmt.Sprintf("Unknown PronunciationAssessmentGranularity: %d", p))
	}
}

// PronunciationAssessmentConfig represents pronunciation assessment configuration.
type PronunciationAssessmentConfig struct {
	handle     C.SPXPRONUNCIATIONASSESSMENTCONFIGHANDLE
	properties *common.PropertyCollection
}

// GetHandle gets the handle to the resource (for internal use)
func (c *PronunciationAssessmentConfig) GetHandle() common.SPXHandle {
	return handle2uintptr(c.handle)
}

// Close releases the underlying resources
func (c *PronunciationAssessmentConfig) Close() {
	c.properties.Close()
	if C.pronunciation_assessment_config_is_handle_valid(c.handle) {
		C.pronunciation_assessment_config_release(c.handle)
	}
}

func newPronunciationAssessmentConfigFromHandle(handle C.SPXPRONUNCIATIONASSESSMENTCONFIGHANDLE) (*PronunciationAssessmentConfig, error) {
	var propBagHandle C.SPXPROPERTYBAGHANDLE
	ret := uintptr(C.pronunciation_assessment_config_get_property_bag(handle, &propBagHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	config := new(PronunciationAssessmentConfig)
	config.handle = handle
	config.properties = common.NewPropertyCollectionFromHandle(handle2uintptr(propBagHandle))
	return config, nil
}

// NewPronunciationAssessmentConfig creates a new pronunciation assessment configuration.
func NewPronunciationAssessmentConfig(
	referenceText string,
	gradingSystem PronunciationAssessmentGradingSystemType,
	granularity PronunciationAssessmentGranularity,
	enableMiscue bool,
) (*PronunciationAssessmentConfig, error) {
	var handle C.SPXPRONUNCIATIONASSESSMENTCONFIGHANDLE
	cRefText := C.CString(referenceText)
	defer C.free(unsafe.Pointer(cRefText))
	ret := uintptr(C.create_pronunciation_assessment_config(&handle, cRefText, gradingSystem.ToCType(), granularity.ToCType(), (C.bool)(enableMiscue)))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	return newPronunciationAssessmentConfigFromHandle(handle)
}

// ApplyTo applies the pronunciation assessment configuration to a speech recognizer.
func (c *PronunciationAssessmentConfig) ApplyTo(recognizer *SpeechRecognizer) error {
	ret := uintptr(C.pronunciation_assessment_config_apply_to_recognizer(c.handle, recognizer.handle))
	if ret != C.SPX_NOERROR {
		return common.NewCarbonError(ret)
	}
	return nil
}

func (c *PronunciationAssessmentConfig) ReferenceText() string {
	return c.GetProperty(common.PronunciationAssessmentReferenceText)
}

func (c *PronunciationAssessmentConfig) SetReferenceText(referenceText string) error {
	return c.SetProperty(common.PronunciationAssessmentReferenceText, referenceText)
}

func (c *PronunciationAssessmentConfig) SetPhonemeAlphabet(alphabet string) error {
	return c.SetProperty(common.PronunciationAssessmentPhonemeAlphabet, alphabet)
}

func (c *PronunciationAssessmentConfig) SetNBestPhonemeCount(count int) error {
	return c.SetProperty(common.PronunciationAssessmentNBestPhonemeCount, fmt.Sprintf("%d", count))
}

func (c *PronunciationAssessmentConfig) EnableProsodyAssessment() error {
	return c.SetProperty(common.PronunciationAssessmentEnableProsodyAssessment, "true")
}

func (c *PronunciationAssessmentConfig) EnableContentAssessmentWithTopic(topic string) error {
	return c.SetProperty(common.PronunciationAssessmentContentTopic, topic)
}

// SetProperty sets a property value by ID.
func (c *PronunciationAssessmentConfig) SetProperty(id common.PropertyID, value string) error {
	return c.properties.SetProperty(id, value)
}

// GetProperty gets a property value by ID.
func (c *PronunciationAssessmentConfig) GetProperty(id common.PropertyID) string {
	return c.properties.GetProperty(id, "")
}

// SetPropertyByString sets a property value by name.
func (c *PronunciationAssessmentConfig) SetPropertyByString(name string, value string) error {
	return c.properties.SetPropertyByString(name, value)
}

// GetPropertyByString gets a property value by name.
func (c *PronunciationAssessmentConfig) GetPropertyByString(name string) string {
	return c.properties.GetPropertyByString(name, "")
}
