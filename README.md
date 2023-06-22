A PoC using `containers/image` and `openshift/runtime-utils` to obtain a prioritized list of mirrors applicable to a given source image

The `withfile` version uses the `containers/image` library as is, with no API changes or code copied

The `withdupe` version required copying and modifying a few private functions from [here](https://github.com/containers/image/blob/04ba35d04e92eb4eba55563c8ccc404d3ffe7c19/pkg/sysregistriesv2/system_registries_v2.go) in order to operate purely in memory.

The `withoutdupe` version used [a fork]() of `containers/image` with the necessary API methods exposed


```sh
$ go run .

2023/06/16 15:44:01 start
2023/06/16 15:44:01 Config: 
{
  "Registries": [
    {
      "Prefix": "source-1.com",
      "Location": "source-1.com",
      "Insecure": false,
      "PullFromMirror": "",
      "Mirrors": [
        {
          "Location": "mirror-1.com",
          "Insecure": false,
          "PullFromMirror": "digest-only"
        },
        {
          "Location": "mirror-2.com/hello",
          "Insecure": false,
          "PullFromMirror": "digest-only"
        }
      ],
      "Blocked": false,
      "MirrorByDigestOnly": false
    },
    {
      "Prefix": "source-2.com",
      "Location": "source-2.com",
      "Insecure": false,
      "PullFromMirror": "",
      "Mirrors": [
        {
          "Location": "mirror-3.com",
          "Insecure": false,
          "PullFromMirror": "digest-only"
        },
        {
          "Location": "mirror-4.com",
          "Insecure": false,
          "PullFromMirror": "digest-only"
        }
      ],
      "Blocked": false,
      "MirrorByDigestOnly": false
    }
  ],
  "UnqualifiedSearchRegistries": null,
  "CredentialHelpers": null,
  "ShortNameMode": "",
  "Aliases": null
}
2023/06/16 15:44:01 Matched Registry: 
{
  "Prefix": "source-1.com",
  "Location": "source-1.com",
  "Insecure": false,
  "PullFromMirror": "",
  "Mirrors": [
    {
      "Location": "mirror-1.com",
      "Insecure": false,
      "PullFromMirror": "digest-only"
    },
    {
      "Location": "mirror-2.com/hello",
      "Insecure": false,
      "PullFromMirror": "digest-only"
    }
  ],
  "Blocked": false,
  "MirrorByDigestOnly": false
}
2023/06/16 15:44:01 Pull Sources: 
[
  {
    "Endpoint": {
      "Location": "mirror-1.com",
      "Insecure": false,
      "PullFromMirror": "digest-only"
    },
    "Reference": {}
  },
  {
    "Endpoint": {
      "Location": "mirror-2.com/hello",
      "Insecure": false,
      "PullFromMirror": "digest-only"
    },
    "Reference": {}
  },
  {
    "Endpoint": {
      "Location": "source-1.com",
      "Insecure": false,
      "PullFromMirror": ""
    },
    "Reference": {}
  }
]
2023/06/16 15:44:01 [0] mirror-1.com/openshift-release-dev/ocp-v4.0-art-dev@sha256:7270ceb168750f0c4ae0afb0086b6dc111dd0da5a96ef32638e8c414b288d228
2023/06/16 15:44:01 [1] mirror-2.com/hello/openshift-release-dev/ocp-v4.0-art-dev@sha256:7270ceb168750f0c4ae0afb0086b6dc111dd0da5a96ef32638e8c414b288d228
2023/06/16 15:44:01 [2] source-1.com/openshift-release-dev/ocp-v4.0-art-dev@sha256:7270ceb168750f0c4ae0afb0086b6dc111dd0da5a96ef32638e8c414b288d228
2023/06/16 15:44:01 end
```
