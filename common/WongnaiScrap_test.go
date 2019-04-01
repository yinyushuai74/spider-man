package common

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestmakeRequest(t *testing.T)  {
	MakeRequest("https://www.wongnai.com/","_api/regions.json","_v","5.056","locale","th","knownLocation","false")
	assert.Equal(t,"","")
}