package ai

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	err := godotenv.Overload("../../.env")
	if err != nil {
		log.Fatal("ai test: error loading .env file:", err)
	}

	Init(os.Getenv("OPENAI_KEY"))

	m.Run()
}

func TestImage(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	image, err := Image("Old man")
	require.Nil(err, "ai test: error generating image:", err)
	assert.NotNil(image, "ai test: image is nil")
}
