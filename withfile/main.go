package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/containers/image/v5/docker/reference"
	"github.com/containers/image/v5/pkg/sysregistriesv2"
	"github.com/containers/image/v5/types"
	apioperatorsv1alpha1 "github.com/openshift/api/operator/v1alpha1"
	"github.com/openshift/runtime-utils/pkg/registries"
)

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func asJson(obj interface{}) []byte {
	data, err := json.MarshalIndent(obj, "", "  ")
	check(err)
	return data
}

func main() {
	log.Print("start")

	// create ICSP rules
	icspRules := []*apioperatorsv1alpha1.ImageContentSourcePolicy{
		{
			Spec: apioperatorsv1alpha1.ImageContentSourcePolicySpec{
				RepositoryDigestMirrors: []apioperatorsv1alpha1.RepositoryDigestMirrors{
					{Source: "source-1.com", Mirrors: []string{"mirror-1.com", "mirror-2.com/hello"}},
					{Source: "source-2.com", Mirrors: []string{"mirror-3.com", "mirror-4.com"}},
				},
			},
		},
	}

	// create empty config
	config := new(sysregistriesv2.V2RegistriesConf)

	// populate the config
	err := registries.EditRegistriesConfig(config, nil, nil, icspRules, nil, nil)
	check(err)
	log.Printf("Config: \n%s", asJson(config))

	// write the config to disk
	filePath := "registries.conf"
	fullPath, err := filepath.Abs(filePath)
	check(err)
	err = createRegistriesFile(config, fullPath)
	check(err)
	defer func() { check(os.Remove(filePath)) }()

	name := "source-1.com/openshift-release-dev/ocp-v4.0-art-dev@sha256:7270ceb168750f0c4ae0afb0086b6dc111dd0da5a96ef32638e8c414b288d228"

	// get an appropriate registry
	ctx := &types.SystemContext{SystemRegistriesConfPath: fullPath}
	reg, err := sysregistriesv2.FindRegistry(ctx, name)
	check(err)
	log.Printf("Matched Registry: \n%s", asJson(reg))

	// get the list of mirrors to pull from
	ref, err := reference.ParseNamed(name)
	check(err)
	srcs, err := reg.PullSourcesFromReference(ref)
	check(err)
	log.Printf("Pull Sources: \n%s", asJson(srcs))

	for i, src := range srcs {
		log.Printf("[%v] %v", i, src.Reference.String())
	}

	log.Print("end")
}

func createRegistriesFile(config *sysregistriesv2.V2RegistriesConf, path string) error {
	var newData bytes.Buffer
	encoder := toml.NewEncoder(&newData)
	if err := encoder.Encode(config); err != nil {
		return err
	}

	err := os.WriteFile(path, newData.Bytes(), 0744)
	if err != nil {
		return err
	}

	return nil
}
