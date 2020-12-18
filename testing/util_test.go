package testing

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestWaitFor(t *testing.T) {
	err := gophercloud.WaitFor(2, func() (bool, error) {
		return true, nil
	})
	th.CheckNoErr(t, err)
}

func TestWaitForTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	err := gophercloud.WaitFor(1, func() (bool, error) {
		return false, nil
	})
	th.AssertEquals(t, "A timeout occurred", err.Error())
}

func TestWaitForError(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	err := gophercloud.WaitFor(2, func() (bool, error) {
		return false, errors.New("Error has occurred")
	})
	th.AssertEquals(t, "Error has occurred", err.Error())
}

func TestWaitForPredicateExceed(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	err := gophercloud.WaitFor(1, func() (bool, error) {
		time.Sleep(4 * time.Second)
		return false, errors.New("Just wasting time")
	})
	th.AssertEquals(t, "A timeout occurred", err.Error())
}

func TestNormalizeURL(t *testing.T) {
	urls := []string{
		"NoSlashAtEnd",
		"SlashAtEnd/",
	}
	expected := []string{
		"NoSlashAtEnd/",
		"SlashAtEnd/",
	}
	for i := 0; i < len(expected); i++ {
		th.CheckEquals(t, expected[i], gophercloud.NormalizeURL(urls[i]))
	}

}

func TestNormalizePathURL(t *testing.T) {
	baseDir, _ := os.Getwd()

	rawPath := "template.yaml"
	basePath, _ := filepath.Abs(".")
	result, _ := gophercloud.NormalizePathURL(basePath, rawPath)
	expected := strings.Join([]string{"file:/", filepath.ToSlash(baseDir), "template.yaml"}, "/")
	th.CheckEquals(t, expected, result)

	rawPath = "http://www.google.com"
	basePath, _ = filepath.Abs(".")
	result, _ = gophercloud.NormalizePathURL(basePath, rawPath)
	expected = "http://www.google.com"
	th.CheckEquals(t, expected, result)

	rawPath = "very/nested/file.yaml"
	basePath, _ = filepath.Abs(".")
	result, _ = gophercloud.NormalizePathURL(basePath, rawPath)
	expected = strings.Join([]string{"file:/", filepath.ToSlash(baseDir), "very/nested/file.yaml"}, "/")
	th.CheckEquals(t, expected, result)

	rawPath = "very/nested/file.yaml"
	basePath = "http://www.google.com"
	result, _ = gophercloud.NormalizePathURL(basePath, rawPath)
	expected = "http://www.google.com/very/nested/file.yaml"
	th.CheckEquals(t, expected, result)

	rawPath = "very/nested/file.yaml/"
	basePath = "http://www.google.com/"
	result, _ = gophercloud.NormalizePathURL(basePath, rawPath)
	expected = "http://www.google.com/very/nested/file.yaml"
	th.CheckEquals(t, expected, result)

	rawPath = "very/nested/file.yaml"
	basePath = "http://www.google.com/even/more"
	result, _ = gophercloud.NormalizePathURL(basePath, rawPath)
	expected = "http://www.google.com/even/more/very/nested/file.yaml"
	th.CheckEquals(t, expected, result)

	rawPath = "very/nested/file.yaml"
	basePath = strings.Join([]string{"file:/", filepath.ToSlash(baseDir), "only/file/even/more"}, "/")
	result, _ = gophercloud.NormalizePathURL(basePath, rawPath)
	expected = strings.Join([]string{"file:/", filepath.ToSlash(baseDir), "only/file/even/more/very/nested/file.yaml"}, "/")
	th.CheckEquals(t, expected, result)

	rawPath = "very/nested/file.yaml/"
	basePath = strings.Join([]string{"file:/", filepath.ToSlash(baseDir), "only/file/even/more"}, "/")
	result, _ = gophercloud.NormalizePathURL(basePath, rawPath)
	expected = strings.Join([]string{"file:/", filepath.ToSlash(baseDir), "only/file/even/more/very/nested/file.yaml"}, "/")
	th.CheckEquals(t, expected, result)

}

func TestRemainingKeys(t *testing.T) {
	type User struct {
		UserID    string `json:"user_id"`
		Username  string `json:"username"`
		Location  string `json:"-"`
		CreatedAt string `json:"-"`
		Status    string
		IsAdmin   bool
	}

	userResponse := map[string]interface{}{
		"user_id":      "abcd1234",
		"username":     "jdoe",
		"location":     "Hawaii",
		"created_at":   "2017-06-08T02:49:03.000000",
		"status":       "active",
		"is_admin":     "true",
		"custom_field": "foo",
	}

	expected := map[string]interface{}{
		"created_at":   "2017-06-08T02:49:03.000000",
		"is_admin":     "true",
		"custom_field": "foo",
	}

	actual := gophercloud.RemainingKeys(User{}, userResponse)

	isEqual := reflect.DeepEqual(expected, actual)
	if !isEqual {
		t.Fatalf("expected %s but got %s", expected, actual)
	}
}
