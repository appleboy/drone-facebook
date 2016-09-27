package main

import (
	"github.com/stretchr/testify/assert"

	"os"
	"testing"
)

func TestMissingDefaultConfig(t *testing.T) {
	var plugin Plugin

	err := plugin.Exec()

	assert.NotNil(t, err)
}

func TestMissingUserConfig(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			PageToken:   "123456789",
			VerifyToken: "123456789",
		},
	}

	err := plugin.Exec()

	assert.NotNil(t, err)
}

func TestDefaultMessageFormat(t *testing.T) {
	plugin := Plugin{
		Repo: Repo{
			Name:  "go-hello",
			Owner: "appleboy",
		},
		Build: Build{
			Number: 101,
			Status: "success",
			Link:   "https://github.com/appleboy/go-hello",
			Author: "Bo-Yi Wu",
			Branch: "master",
			Commit: "e7c4f0a63ceeb42a39ac7806f7b51f3f0d204fd2",
		},
	}

	message := plugin.Message(plugin.Repo, plugin.Build)

	assert.Equal(t, []string{"[success] <https://github.com/appleboy/go-hello> (master) by Bo-Yi Wu"}, message)
}

func TestSendMessage(t *testing.T) {
	plugin := Plugin{
		Repo: Repo{
			Name:  "go-hello",
			Owner: "appleboy",
		},
		Build: Build{
			Number: 101,
			Status: "success",
			Link:   "https://github.com/appleboy/go-hello",
			Author: "Bo-Yi Wu",
			Branch: "master",
			Commit: "e7c4f0a63ceeb42a39ac7806f7b51f3f0d204fd2",
		},

		Config: Config{
			PageToken:   os.Getenv("FB_PAGE_TOKEN"),
			VerifyToken: os.Getenv("FB_VERIFY_TOKEN"),
			Verify:      false,
			To:          []string{os.Getenv("FB_TO"), "中文ID", "1234567890"},
			Message:     []string{"Test Facebook Bot From Travis or Local", " "},
		},
	}

	err := plugin.Exec()
	assert.Nil(t, err)

	// disable message
	plugin.Config.Message = []string{}
	err = plugin.Exec()
	assert.Nil(t, err)
}

func TestTrimElement(t *testing.T) {
	var input, result []string

	input = []string{"1", "     ", "3"}
	result = []string{"1", "3"}

	assert.Equal(t, result, trimElement(input))

	input = []string{"1", "2"}
	result = []string{"1", "2"}

	assert.Equal(t, result, trimElement(input))
}

func TestParseID(t *testing.T) {
	var input []string
	var result []int64

	input = []string{"1", "測試", "3"}
	result = []int64{int64(1), int64(3)}

	assert.Equal(t, result, parseID(input))

	input = []string{"1", "2"}
	result = []int64{int64(1), int64(2)}

	assert.Equal(t, result, parseID(input))
}
