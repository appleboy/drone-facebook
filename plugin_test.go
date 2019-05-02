package main

import (
	"github.com/stretchr/testify/assert"

	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMissingDefaultConfig(t *testing.T) {
	var plugin Plugin

	err := plugin.Exec()

	assert.NotNil(t, err)
}

func TestDefaultMessageFormat(t *testing.T) {
	plugin := Plugin{
		Repo: Repo{
			FullName:  "appleboy/go-hello",
			Name:      "go-hello",
			Namespace: "appleboy",
		},
		Commit: Commit{
			Sha:     "e7c4f0a63ceeb42a39ac7806f7b51f3f0d204fd2",
			Author:  "Bo-Yi Wu",
			Branch:  "master",
			Message: "update travis by drone plugin",
			Email:   "test@gmail.com",
		},
		Build: Build{
			Tag:    "1.0.0",
			Number: 101,
			Status: "success",
			Link:   "https://github.com/appleboy/go-hello",
		},
	}

	message := plugin.Message()

	assert.Equal(t, []string{"[success] <https://github.com/appleboy/go-hello> (master)『update travis by drone plugin』by Bo-Yi Wu"}, message)
}

func TestErrorTemplate(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			PageToken:   os.Getenv("FB_PAGE_TOKEN"),
			VerifyToken: os.Getenv("FB_VERIFY_TOKEN"),
			Verify:      false,
			To:          []string{"1234567890"},
			Message:     []string{"file://xxxxx/xxxxx"},
		},
	}

	err := plugin.Exec()
	assert.NoError(t, err)
}

func TestSendMessage(t *testing.T) {
	plugin := Plugin{
		Repo: Repo{
			FullName:  "appleboy/go-hello",
			Name:      "go-hello",
			Namespace: "appleboy",
		},
		Commit: Commit{
			Sha:     "e7c4f0a63ceeb42a39ac7806f7b51f3f0d204fd2",
			Author:  "Bo-Yi Wu",
			Branch:  "master",
			Message: "update travis by drone plugin",
			Email:   "test@gmail.com",
		},
		Build: Build{
			Tag:    "1.0.0",
			Number: 101,
			Status: "success",
			Link:   "https://github.com/appleboy/go-hello",
		},

		Config: Config{
			PageToken:   os.Getenv("FB_PAGE_TOKEN"),
			VerifyToken: os.Getenv("FB_VERIFY_TOKEN"),
			Verify:      false,
			To:          []string{os.Getenv("FB_TO")},
			Message:     []string{"Test Facebook Bot From Travis or Local from {{ build.author }}", " "},
			Image:       []string{"https://cdn3.iconfinder.com/data/icons/picons-social/57/16-apple-256.png"},
			// Audio:       []string{"https://ia802508.us.archive.org/5/items/testmp3testfile/mpthreetest.mp3"},
			// Video:       []string{"https://www.sample-videos.com/video123/mp4/720/big_buck_bunny_720p_1mb.mp4"},
			// File: []string{"http://open.qiniudn.com/where-can-you-use-golang.pdf"},
		},
	}

	err := plugin.Exec()
	assert.Nil(t, err)

	// disable message
	plugin.Config.Message = []string{}
	err = plugin.Exec()
	assert.Nil(t, err)
}

func TestDefaultMessageFormatFromGitHub(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			GitHub: true,
		},
		Repo: Repo{
			FullName:  "appleboy/go-hello",
			Name:      "go-hello",
			Namespace: "appleboy",
		},
		GitHub: GitHub{
			Workflow:  "test-workflow",
			Action:    "send notification",
			EventName: "push",
		},
	}

	message := plugin.Message()

	assert.Equal(t, []string{"appleboy/go-hello/test-workflow triggered by appleboy (push)"}, message)
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

func TestParseTo(t *testing.T) {
	input := []string{"0", "1:1@gmail.com", "2:2@gmail.com", "3:3@gmail.com", "4", "5"}

	ids := parseTo(input, "1@gmail.com", false)
	assert.Equal(t, []int64{0, 4, 5, 1}, ids)

	ids = parseTo(input, "1@gmail.com", true)
	assert.Equal(t, []int64{1}, ids)

	ids = parseTo(input, "a@gmail.com", false)
	assert.Equal(t, []int64{0, 4, 5}, ids)

	ids = parseTo(input, "a@gmail.com", true)
	assert.Equal(t, []int64{0, 4, 5}, ids)

	// test empty ids
	ids = parseTo([]string{"", " ", "   "}, "a@gmail.com", true)
	assert.Equal(t, 0, len(ids))
}

func performRequest(r http.Handler, method, url string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, url, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestDefaultRouter(t *testing.T) {
	p := Plugin{
		Config: Config{
			PageToken:   os.Getenv("FB_PAGE_TOKEN"),
			VerifyToken: os.Getenv("FB_VERIFY_TOKEN"),
			Verify:      false,
		},
	}

	router := p.serveMux()
	w := performRequest(router, "GET", "/")
	assert.Equal(t, "Welcome to facebook webhook page.\n", w.Body.String())

	w = performRequest(router, "GET", "/metrics")
	assert.Equal(t, 200, w.Code)
}
