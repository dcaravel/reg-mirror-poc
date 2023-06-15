A PoC using `containers/image` and `openshift/runtime-utils` to obtain a prioritized list of mirrors applicable to a given source image

```sh
$ go run .

2023/06/15 11:33:20 start
2023/06/15 11:33:20 Config: 
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
          "Location": "mirror-2.com",
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
2023/06/15 11:33:20 Matched Registry: 
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
      "Location": "mirror-2.com",
      "Insecure": false,
      "PullFromMirror": "digest-only"
    }
  ],
  "Blocked": false,
  "MirrorByDigestOnly": false
}
2023/06/15 11:33:20 Pull Sources: 
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
      "Location": "mirror-2.com",
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
2023/06/15 11:33:20 end
```