package utils

import (
  "fmt"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/awserr"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/s3"
  "os"
  "path/filepath"
  "strings"
)

func contains(s string, l []string) bool {
  for _, i := range l {
    if i == s {
      return true
    }
  }

  return false
}

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

func GetStackPath(environments []string) (string, error) {
  dirs := splitPath(absPath())
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

  return "", fmt.Errorf("Error: no paths matched")
}

func ReplaceSlash(path string) string {
  return strings.Replace(path, string(filepath.Separator), "_", -1)
}

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

func FindTfGenPath(tfgenConf string) (string, error) {
  a := absPath()
  path := a

  for _, _ = range splitPath(a) {
    if FileExists(filepath.Join(path, tfgenConf)) {
      return path, nil
    }
    path = filepath.Dir(path)
  }

  return "", fmt.Errorf("no directories in %v have %v", a, tfgenConf)
}
