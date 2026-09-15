// Package region contains the compiled-in endpoint and certificate configuration.
// This file is specific to the test deployment.
// Do not modify at runtime — configuration is intentionally immutable.
package region

// Endpoint is the service address for this deployment.
const Endpoint = "10.68.36.108:9443"

// TLSAnchor is the PEM-encoded TLS certificate for transport security.
const TLSAnchor = `-----BEGIN CERTIFICATE-----
MIIBlzCCAT2gAwIBAgIBATAKBggqhkjOPQQDAjAhMR8wHQYDVQQDExZMaWNlbnNl
SXNzdWVyIGdSUEMgVExTMB4XDTI2MDcxNjExNDczNloXDTM2MDcxMzExNDczNlow
ITEfMB0GA1UEAxMWTGljZW5zZUlzc3VlciBnUlBDIFRMUzBZMBMGByqGSM49AgEG
CCqGSM49AwEHA0IABBbGCxLoBYezd0Fza10SHhGjVjTZ+Y8LnpUD1QzTxCi0EUeC
Z/b7kJooteaLIw7CnrroP7LdXXzBVM4hyw+JXtKjZjBkMA8GA1UdEwEB/wQFMAMB
Af8wHQYDVR0OBBYEFD00oSFpBggAEMy8Z2gAqBUNHaXKMDIGA1UdEQQrMCmCFkxp
Y2Vuc2VJc3N1ZXIgZ1JQQyBUTFOCCWxvY2FsaG9zdIcEfwAAATAKBggqhkjOPQQD
AgNIADBFAiB+dPQwu6r2YEwYDyPgyDH1cjm56OQVryRuSnEJS/1mhQIhAOGOmUzM
ReRsf8ZR3cHvRZSXEeL7WC4uIAWWKCuMrmxa
-----END CERTIFICATE-----`

// RootAnchor is the PEM-encoded root certificate for token validation.
const RootAnchor = `-----BEGIN CERTIFICATE-----
MIIBaTCCAQ6gAwIBAgIIGMLCrj1CcAAwCgYIKoZIzj0EAwIwIDEeMBwGA1UEAxMV
TGljZW5zZUlzc3VlciBSb290IENBMB4XDTI2MDcxNjExNDczNloXDTI2MDcxNjEz
NDczNlowIDEeMBwGA1UEAxMVTGljZW5zZUlzc3VlciBSb290IENBMFkwEwYHKoZI
zj0CAQYIKoZIzj0DAQcDQgAE2dc5hzdQZJesIfZt7ki3ocf6DRXRBBb6SB38i7xk
KHD6lP1IRVR9J5AjDGVUi+K3T7n2WMR8X1ll3XTLRVvxS6MyMDAwDwYDVR0TAQH/
BAUwAwEB/zAdBgNVHQ4EFgQUYIzJcVCZ2PhIxUqs8+R3vHGTBaIwCgYIKoZIzj0E
AwIDSQAwRgIhAImZx29aSh71mjj1DrRYjok1OR5sAaGC3+q2ZZNLB/snAiEAmxOV
YtXX62cHN3jTGg0QBM6HJtXAb9OQJUDPUOnu0ew=
-----END CERTIFICATE-----`
