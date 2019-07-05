package templates

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/mschurenko/tfgen/utils"
)

const (
	terraformFile    = "terraform.tf.json"
	remoteStatesFile = "remote_states.tf.json"
	stateFileKey     = "stacks/%v/terraform.tfstate"
)

// top-level structs for templates
type tfDataSource struct {
	Data remoteState `json:"data"`
}

type tfBackend struct {
	Terraform terraform `json:"terraform"`
}

// sub structs used by top-level structs
type remoteState struct {
	TerraformRemoteState map[string]interface{} `json:"terraform_remote_state"`
}

type terraform struct {
	Backend         backend `json:"backend"`
	RequiredVersion string  `json:"required_version,omitempty"`
}

type backend struct {
	S3 s3 `json:"s3"`
}

type s3 struct {
	Bucket        string `json:"bucket"`
	DynamodbTable string `json:"dynamodb_table"`
	Encrypt       bool   `json:"encrypt"`
	Key           string `json:"key"`
	Region        string `json:"region"`
}

func writeFile(file string, d interface{}, force bool) error {
	fileExists := utils.FileExists(file)
	if (fileExists && force) || !fileExists {
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(d)
		if err != nil {
			return err
		}

		dst := &bytes.Buffer{}
		json.Indent(dst, []byte(buf.Bytes()), "", "  ")

		fmt.Println("Creating", file)
		if err := ioutil.WriteFile(file, dst.Bytes(), 0644); err != nil {
			return err
		}

	} else {
		fmt.Println(file, "alredy exists")
	}

	return nil
}

func createTerraform(s3Config map[string]string, path string, reqVer string, force bool) error {
	tb := tfBackend{
		Terraform: terraform{
			RequiredVersion: reqVer,
			Backend: backend{
				S3: s3{
					Bucket:        s3Config["bucket"],
					DynamodbTable: s3Config["dynamodb_table"],
					Encrypt:       true,
					Key:           fmt.Sprintf(stateFileKey, path),
					Region:        s3Config["aws_region"],
				},
			},
		},
	}

	if err := writeFile(terraformFile, tb, force); err != nil {
		return err
	}

	return nil
}

func createRemoteState(s3Config map[string]string, stack string, key string) error {
	var tfDs *tfDataSource
	if utils.FileExists(remoteStatesFile) {
		// first unmarshal existing json to data structure
		j, err := ioutil.ReadFile(remoteStatesFile)
		if err != nil {
			return err
		}

		tfDs = &tfDataSource{}

		if err := json.Unmarshal([]byte(j), tfDs); err != nil {
			return err
		}

		_, ok := tfDs.Data.TerraformRemoteState[stack]
		if ok {
			fmt.Println(stack, "already exits")
			return nil
		} else {
			tfDs.Data.TerraformRemoteState[stack] = map[string]interface{}{
				"backend": "s3",
				"config": s3{
					Bucket:        s3Config["bucket"],
					DynamodbTable: s3Config["dynamodb_table"],
					Encrypt:       true,
					Key:           key,
					Region:        s3Config["aws_region"],
				},
			}

		}
	} else {
		tfDs = &tfDataSource{
			Data: remoteState{
				TerraformRemoteState: map[string]interface{}{
					stack: map[string]interface{}{
						"backend": "s3",
						"config": s3{
							Bucket:        s3Config["bucket"],
							DynamodbTable: s3Config["dynamodb_table"],
							Encrypt:       true,
							Key:           key,
							Region:        s3Config["aws_region"],
						},
					},
				},
			},
		}

	}

	if err := writeFile(remoteStatesFile, tfDs, true); err != nil {
		return err
	}

	return nil
}

// InitStack sets up s3 backend
func InitStack(s3Config map[string]string, environments []string, stackRx string, reqVer string, force bool) error {
	path, err := utils.GetStackPath(stackRx, environments)
	if err != nil {
		return err
	}

	err = createTerraform(s3Config, path, reqVer, force)
	if err != nil {
		return err
	}

	return nil
}

// RemoteState sets up remote state data sources
func RemoteState(s3Config map[string]string, stackName string, noVerifyKey bool) error {
	stackNameSafe := utils.ReplaceSlash(stackName)

	key := fmt.Sprintf(stateFileKey, stackName)
	if !noVerifyKey {
		found, err := utils.KeyInS3(s3Config, key)
		if err != nil {
			return err
		}
		if !found {
			return fmt.Errorf("key: %v could not be found in %v", key, s3Config["bucket"])
		}
	}

	if err := createRemoteState(s3Config, stackNameSafe, key); err != nil {
		return err
	}

	return nil
}
