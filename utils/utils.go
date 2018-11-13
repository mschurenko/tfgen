package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func contains(s string, l []string) bool {

	for _, i := range l {

		if i == s {
			return true
		}
	}

	return false
}

// FileExists checks if file exits
func FileExists(file string) bool {
	_, err := os.Stat(file)

	if !os.IsNotExist(err) {
		return true
	}

	return false
}

func absPath() string {
	absPath, _ := filepath.Abs(".")
	return absPath
}

func splitPath(path string) []string {
	return strings.Split(path, string(filepath.Separator))
}

// GetStackPath returns path relative to last environment dir
func GetStackPath(pat string, environments []string) (string, error) {
	a := absPath()
	baseDir := filepath.Base(a)

	matched, err := regexp.MatchString(`^([a-z]|[A-Z]|[0-9]|-)+$`, baseDir)
	if err != nil {
		return "", err
	}

	if !matched {
		return "", fmt.Errorf("%v does not match regexp %v", baseDir, pat)
	}

	dirs := splitPath(a)
	foundPaths := []string{}

	for idx, dir := range dirs {
		if contains(dir, environments) {
			foundPaths = append(foundPaths, filepath.Join(dirs[idx:]...))
		}
	}

	if len(foundPaths) > 1 {
		return foundPaths[len(foundPaths)-1], nil
	} else if len(foundPaths) == 1 {
		return foundPaths[0], nil
	}

	return "", fmt.Errorf("no paths matched")
}

// ReplaceSlash replaces "/" with "_"
func ReplaceSlash(path string) string {
	return strings.Replace(path, string(filepath.Separator), "_", -1)
}

// KeyInS3 return true if key exits in s3, false if it doesn't
func KeyInS3(s3Config map[string]string, key string) (bool, error) {
	// https://docs.aws.amazon.com/sdk-for-go/api/service/s3/#S3.HeadObject
	sess := session.Must(
		session.NewSession(&aws.Config{
			Region: aws.String(s3Config["aws_region"]),
		}),
	)
	svc := s3.New(sess)

	input := &s3.HeadObjectInput{
		Bucket: aws.String(s3Config["bucket"]),
		Key:    aws.String(key),
	}

	_, err := svc.HeadObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return false, nil
			}
		} else {
			return false, err
		}
	}

	return true, nil
}

// FindTfGenPath finds the path to the last .tfgen.yml file in the cwd
func FindTfGenPath(tfgenConf string) (string, error) {
	a := absPath()
	path := a

	for range splitPath(a) {
		if FileExists(filepath.Join(path, tfgenConf)) {
			return path, nil
		}
		path = filepath.Dir(path)
	}

	return "", fmt.Errorf("no directories in %v have %v", a, tfgenConf)
}
