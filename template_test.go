package main

import (
	"github.com/stretchr/testify/assert"

	"fmt"
	"os"
	"testing"
	"time"
)

var plugin = Plugin{
	Repo: Repo{
		Name:  "go-hello",
		Owner: "appleboy",
	},
	Build: Build{
		Number:  101,
		Status:  "success",
		Link:    "https://github.com/appleboy/go-hello",
		Author:  "Bo-Yi Wu",
		Branch:  "master",
		Message: "update travis",
		Commit:  "e7c4f0a63ceeb42a39ac7806f7b51f3f0d204fd2",
	},
}

func TestTruncate(t *testing.T) {
	assert.Equal(t, "e7c4f0a63c", truncate(plugin.Build.Commit, 10))
	assert.Equal(t, plugin.Build.Commit, truncate(plugin.Build.Commit, 200))
}

func TestUppercaseFirst(t *testing.T) {
	assert.Equal(t, "Success", uppercaseFirst(plugin.Build.Status))
}

func TestToDuration(t *testing.T) {
	assert.Equal(t, "3m20s\n", toDuration(float64(1477550550), float64(1477550750)))
}

func TestToDatetime(t *testing.T) {
	localTime := time.Unix(int64(1477550550), 0).Local().Format("3:04PM")
	assert.Equal(t, "6:42AM", toDatetime(float64(1477550550), "3:04PM", "UTC"))

	// missing zone
	assert.Equal(t, localTime, toDatetime(float64(1477550550), "3:04PM", ""))
	// wrong zone
	assert.Equal(t, localTime, toDatetime(float64(1477550550), "3:04PM", "ABCDEFG"))
}

func TestUrlEncode(t *testing.T) {
	res, err := RenderTrim("{{#urlencode}}build successfully{{/urlencode}}", plugin)

	assert.Nil(t, err)
	assert.Equal(t, "build+successfully", res)
}

func TestRender(t *testing.T) {
	// test parse from string
	res, err := RenderTrim("Trigger from {{ build.author }}", plugin)

	assert.Nil(t, err)
	assert.Equal(t, "Trigger from Bo-Yi Wu", res)

	// test parse from url
	res, err = RenderTrim("https://goo.gl/EAivJP", plugin)

	assert.Nil(t, err)
	assert.Equal(t, "Trigger from Bo-Yi Wu", res)

	// test parse from file
	res, err = RenderTrim(fmt.Sprintf("file://%s/handlebar/template.handlebar", os.Getenv("PWD")), plugin)

	assert.Nil(t, err)
	assert.Equal(t, "Trigger from Bo-Yi Wu", res)

	// success build
	res, err = RenderTrim("{{#success build.status}}{{ build.author }} successfully pushed to {{ build.branch}}{{/success}}", plugin)

	assert.Nil(t, err)
	assert.Equal(t, "Bo-Yi Wu successfully pushed to master", res)

	// Inverse success build
	res, err = RenderTrim("{{#failure build.status}}{{ build.author }} successfully pushed to {{ build.branch}}{{/failure}}", plugin)

	assert.Nil(t, err)
	assert.Equal(t, "", res)

	plugin.Build.Status = "failure"
	// failure build
	res, err = RenderTrim("{{#failure build.status}}Something is busted{{/failure}}", plugin)

	assert.Nil(t, err)
	assert.Equal(t, "Something is busted", res)

	// Inverse failure build
	res, err = RenderTrim("{{#success build.status}}Something is busted{{/success}}", plugin)

	assert.Nil(t, err)
	assert.Equal(t, "", res)
}
