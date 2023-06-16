package main

import (
	"encoding/json"
	"log"

	"github.com/containers/image/v5/docker/reference"
	"github.com/containers/image/v5/pkg/sysregistriesv2"
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

	// perform post processing / cleanup on the config
	err = config.PostProcessRegistries()
	check(err)
	log.Printf("Config: \n%s", asJson(config))

	name := "source-1.com/openshift-release-dev/ocp-v4.0-art-dev@sha256:7270ceb168750f0c4ae0afb0086b6dc111dd0da5a96ef32638e8c414b288d228"

	// get an appropriate registry
	reg, err := sysregistriesv2.FindRegistryWithConfig(config, name)
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
