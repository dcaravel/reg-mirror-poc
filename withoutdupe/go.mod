module github.com/dcaravel/reg-mirror-poc/withoutdupe

go 1.20

require (
	github.com/containers/image/v5 v5.22.0
	github.com/openshift/api v0.0.0-20220901185337-0b39f81154fa
	github.com/openshift/runtime-utils v0.0.0-20220926190846-5c488b20a19f
)

require (
	github.com/BurntSushi/toml v1.3.2 // indirect
	github.com/containers/storage v1.46.1 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/moby/sys/mountinfo v0.6.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0-rc3 // indirect
	github.com/opencontainers/runc v1.1.7 // indirect
	github.com/opencontainers/runtime-spec v1.1.0-rc.2 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/syndtr/gocapability v0.0.0-20200815063812-42c35b437635 // indirect
	golang.org/x/exp v0.0.0-20230522175609-2e198f4a06a1 // indirect
	golang.org/x/net v0.11.0 // indirect
	golang.org/x/sys v0.9.0 // indirect
	golang.org/x/text v0.10.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/api v0.25.0 // indirect
	k8s.io/apimachinery v0.25.0 // indirect
	k8s.io/klog/v2 v2.70.1 // indirect
	k8s.io/utils v0.0.0-20220728103510-ee6ede2d64ed // indirect
	sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
)

// github.com/containers/image/v5 => github.com/dcaravel/containers-image v0.0.0-20230615224121-814720beff09
// github.com/containers/image/v5 => github.com/dcaravel/containers-image/v5 v5.0.0-20230615224121-814720beff09
replace github.com/containers/image/v5 => ../../containers-image
