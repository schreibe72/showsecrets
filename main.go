package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"

	"gopkg.in/yaml.v3"
)

type secrets struct {
	Data map[string]string
}

type outputKV map[string]string

func getDecryptedSecretYaml(filename string) []byte {
	cmd := exec.Command("sops", "-d", filename)
	stdout, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	return stdout
}

func getSecretsKV(secretYaml []byte) outputKV {
	s := secrets{}
	err := yaml.Unmarshal(secretYaml, &s)
	if err != nil {
		panic(err)
	}
	output := outputKV{}
	for k, v := range s.Data {
		decodedvalue, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			panic(err)
		}
		output[k] = string(decodedvalue)
	}
	return output
}

func (kv *outputKV) toYaml() string {
	o, err := yaml.Marshal(kv)
	if err != nil {
		panic(err)
	}
	return string(o)
}

func (kv *outputKV) toJson() string {
	o, err := json.MarshalIndent(kv, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(o)
}

func main() {
	j := flag.Bool("j", false, "format output as json")
	flag.Parse()
	filename := flag.Args()[0]

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		panic(fmt.Errorf("%s - file does not exists", filename))
	}

	s := getDecryptedSecretYaml(filename)

	kv := getSecretsKV(s)
	if *j {
		fmt.Println(kv.toJson())
	} else {
		fmt.Println(kv.toYaml())
	}

}
