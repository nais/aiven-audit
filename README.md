# Aiven Audit (Go) 📝🕵️
Fetches project event logs from Aiven API, and logs them to logs.adeo.no.

## Configuration

| environment variable  | description |
| ------------- | ------------- |
| AIVEN_AUDIT_PAT  | Access token for Aiven API  |

## Verifying the aiven-audit images and their contents

The images are signed "keylessly" using [Sigstore cosign](https://github.com/sigstore/cosign).
To verify their authenticity run
```
cosign verify \
--certificate-identity "https://github.com/nais/aiven-audit/.github/workflows/main.yml@refs/heads/main" \
--certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
europe-north1-docker.pkg.dev/nais-io/nais/images/aiven-audit@sha256:<shasum>
```

The images are also attested with SBOMs in the [CycloneDX](https://cyclonedx.org/) format.
You can verify these by running
```
cosign verify-attestation --type cyclonedx  \
--certificate-identity "https://github.com/nais/aiven-audit/.github/workflows/main.yml@refs/heads/main" \
--certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
europe-north1-docker.pkg.dev/nais-io/nais/images/aiven-audit@sha256:<shasum>
```