package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/containers/image/v5/pkg/sysregistriesv2"
	"github.com/containers/storage/pkg/regexp"
	"github.com/docker/distribution/reference"
	apioperatorsv1alpha1 "github.com/openshift/api/operator/v1alpha1"
	"github.com/openshift/runtime-utils/pkg/registries"
)

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
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
	err = postProcessRegistries(config)
	check(err)
	data, err := json.MarshalIndent(config, "", "  ")
	check(err)
	log.Printf("Config: \n%s", data)

	name := "source-1.com/openshift-release-dev/ocp-v4.0-art-dev@sha256:7270ceb168750f0c4ae0afb0086b6dc111dd0da5a96ef32638e8c414b288d228"

	// get an appropriate registry
	reg, err := FindRegistryWithConfig(config, name)
	check(err)
	data, err = json.MarshalIndent(reg, "", "  ")
	check(err)
	log.Printf("Matched Registry: \n%s", data)

	// get the list of mirrors to pull from
	ref, err := reference.ParseNamed(name)
	check(err)
	srcs, err := reg.PullSourcesFromReference(ref)
	check(err)
	data, err = json.MarshalIndent(srcs, "", "  ")
	check(err)
	log.Printf("Pull Sources: \n%s", data)

	for i, src := range srcs {
		log.Printf("[%v] %v", i, src.Reference.String())
	}

	log.Print("end")
}

func FindRegistryWithConfig(config *sysregistriesv2.V2RegistriesConf, ref string) (*sysregistriesv2.Registry, error) {
	return findRegistryWithParsedConfig(config, ref)
}

// ========================================================================================================
// ========================================================================================================
// Methods below here were copied (with minor modifications) from:
// https://github.com/containers/image/blob/main/pkg/sysregistriesv2/system_registries_v2.go
// ========================================================================================================
// ========================================================================================================

func findRegistryWithParsedConfig(config *sysregistriesv2.V2RegistriesConf, ref string) (*sysregistriesv2.Registry, error) {
	reg := sysregistriesv2.Registry{}
	prefixLen := 0
	for _, r := range config.Registries {
		if refMatchingPrefix(ref, r.Prefix) != -1 {
			length := len(r.Prefix)
			if length > prefixLen {
				reg = r
				prefixLen = length
			}
		}
	}
	if prefixLen != 0 {
		return &reg, nil
	}
	return nil, nil
}

func refMatchingPrefix(ref, prefix string) int {
	switch {
	case strings.HasPrefix(prefix, "*."):
		return refMatchingSubdomainPrefix(ref, prefix)
	case len(ref) < len(prefix):
		return -1
	case len(ref) == len(prefix):
		if ref == prefix {
			return len(prefix)
		}
		return -1
	case len(ref) > len(prefix):
		if !strings.HasPrefix(ref, prefix) {
			return -1
		}
		c := ref[len(prefix)]
		// This allows "example.com:5000" to match "example.com",
		// which is unintended; that will get fixed eventually, DON'T RELY
		// ON THE CURRENT BEHAVIOR.
		if c == ':' || c == '/' || c == '@' {
			return len(prefix)
		}
		return -1
	default:
		panic("Internal error: impossible comparison outcome")
	}
}

func refMatchingSubdomainPrefix(ref, prefix string) int {
	index := strings.Index(ref, prefix[1:])
	if index == -1 {
		return -1
	}
	if strings.Contains(ref[:index], "/") {
		return -1
	}
	index += len(prefix[1:])
	if index == len(ref) {
		return index
	}
	switch ref[index] {
	case ':', '/', '@':
		return index
	default:
		return -1
	}
}

var anchoredDomainRegexp = regexp.Delayed("^" + reference.DomainRegexp.String() + "$")

func postProcessRegistries(config *sysregistriesv2.V2RegistriesConf) error {
	regMap := make(map[string][]*sysregistriesv2.Registry)

	for i := range config.Registries {
		reg := &config.Registries[i]
		// make sure Location and Prefix are valid
		var err error
		reg.Location, err = parseLocation(reg.Location)
		if err != nil {
			return err
		}

		if reg.Prefix == "" {
			if reg.Location == "" {
				return errors.New("invalid condition: both location and prefix are unset")
			}
			reg.Prefix = reg.Location
		} else {
			reg.Prefix, err = parseLocation(reg.Prefix)
			if err != nil {
				return err
			}
			// FIXME: allow config authors to always use Prefix.
			// https://github.com/containers/image/pull/1191#discussion_r610622495
			if !strings.HasPrefix(reg.Prefix, "*.") && reg.Location == "" {
				return errors.New("invalid condition: location is unset and prefix is not in the format: *.example.com")
			}
		}

		// validate the mirror usage settings does not apply to primary registry
		if reg.PullFromMirror != "" {
			return fmt.Errorf("pull-from-mirror must not be set for a non-mirror registry %q", reg.Prefix)
		}
		// make sure mirrors are valid
		for _, mir := range reg.Mirrors {
			mir.Location, err = parseLocation(mir.Location)
			if err != nil {
				return err
			}

			//FIXME: unqualifiedSearchRegistries now also accepts empty values
			//and shouldn't
			// https://github.com/containers/image/pull/1191#discussion_r610623216
			if mir.Location == "" {
				return errors.New("invalid condition: mirror location is unset")
			}

			if reg.MirrorByDigestOnly && mir.PullFromMirror != "" {
				return errors.New(fmt.Sprintf("cannot set mirror usage mirror-by-digest-only for the registry (%q) and pull-from-mirror for per-mirror (%q) at the same time", reg.Prefix, mir.Location))
			}
			if mir.PullFromMirror != "" && mir.PullFromMirror != sysregistriesv2.MirrorAll &&
				mir.PullFromMirror != sysregistriesv2.MirrorByDigestOnly && mir.PullFromMirror != sysregistriesv2.MirrorByTagOnly {
				return errors.New(fmt.Sprintf("unsupported pull-from-mirror value %q for mirror %q", mir.PullFromMirror, mir.Location))
			}
		}
		if reg.Location == "" {
			regMap[reg.Prefix] = append(regMap[reg.Prefix], reg)
		} else {
			regMap[reg.Location] = append(regMap[reg.Location], reg)
		}
	}

	// Given a registry can be mentioned multiple times (e.g., to have
	// multiple prefixes backed by different mirrors), we need to make sure
	// there are no conflicts among them.
	//
	// Note: we need to iterate over the registries array to ensure a
	// deterministic behavior which is not guaranteed by maps.
	for _, reg := range config.Registries {
		var others []*sysregistriesv2.Registry
		var ok bool
		if reg.Location == "" {
			others, ok = regMap[reg.Prefix]
		} else {
			others, ok = regMap[reg.Location]
		}
		if !ok {
			return fmt.Errorf("Internal error in V2RegistriesConf.PostProcess: entry in regMap is missing")
		}
		for _, other := range others {
			if reg.Insecure != other.Insecure {
				msg := fmt.Sprintf("registry '%s' is defined multiple times with conflicting 'insecure' setting", reg.Location)
				return errors.New(msg)
			}

			if reg.Blocked != other.Blocked {
				msg := fmt.Sprintf("registry '%s' is defined multiple times with conflicting 'blocked' setting", reg.Location)
				return errors.New(msg)
			}
		}
	}

	for i := range config.UnqualifiedSearchRegistries {
		registry, err := parseLocation(config.UnqualifiedSearchRegistries[i])
		if err != nil {
			return err
		}
		if !anchoredDomainRegexp.MatchString(registry) {
			return errors.New(fmt.Sprintf("Invalid unqualified-search-registries entry %#v", registry))
		}
		config.UnqualifiedSearchRegistries[i] = registry
	}

	// Registries are ordered and the first longest prefix always wins,
	// rendering later items with the same prefix non-existent. We cannot error
	// out anymore as this might break existing users, so let's just ignore them
	// to guarantee that the same prefix exists only once.
	//
	// As a side effect of parsedConfig.updateWithConfigurationFrom, the Registries slice
	// is always sorted. To be consistent in situations where it is not called (no drop-ins),
	// sort it here as well.
	prefixes := []string{}
	uniqueRegistries := make(map[string]sysregistriesv2.Registry)
	for i := range config.Registries {
		// TODO: should we warn if we see the same prefix being used multiple times?
		prefix := config.Registries[i].Prefix
		if _, exists := uniqueRegistries[prefix]; !exists {
			uniqueRegistries[prefix] = config.Registries[i]
			prefixes = append(prefixes, prefix)
		}
	}
	sort.Strings(prefixes)
	config.Registries = []sysregistriesv2.Registry{}
	for _, prefix := range prefixes {
		config.Registries = append(config.Registries, uniqueRegistries[prefix])
	}

	return nil
}

func parseLocation(input string) (string, error) {
	trimmed := strings.TrimRight(input, "/")

	// FIXME: This check needs to exist but fails for empty Location field with
	// wildcarded prefix. Removal of this check "only" allows invalid input in,
	// and does not prevent correct operation.
	// https://github.com/containers/image/pull/1191#discussion_r610122617
	//
	//	if trimmed == "" {
	//		return "", &InvalidRegistries{s: "invalid location: cannot be empty"}
	//	}
	//

	if strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://") {
		msg := fmt.Sprintf("invalid location '%s': URI schemes are not supported", input)
		return "", errors.New(msg)
	}

	return trimmed, nil
}
