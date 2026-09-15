// Package region contains the compiled-in endpoint and certificate configuration.
// This file is specific to the test deployment.
// Do not modify at runtime — configuration is intentionally immutable.
package region

// Endpoint is the service address for this deployment.
const Endpoint = "10.68.36.108:9443"

// TLSAnchor is the PEM-encoded TLS certificate for transport security.
const TLSAnchor = `-----BEGIN CERTIFICATE-----
MIIBlzCCAT2gAwIBAgIBATAKBggqhkjOPQQDAjAhMR8wHQYDVQQDExZMaWNlbnNl
SXNzdWVyIGdSUEMgVExTMB4XDTI2MDkxNTAzMTUzMVoXDTM2MDkxMjAzMTUzMVow
ITEfMB0GA1UEAxMWTGljZW5zZUlzc3VlciBnUlBDIFRMUzBZMBMGByqGSM49AgEG
CCqGSM49AwEHA0IABG6AWThwms4rL/u/g1pKEof5l+C1tFcm1KIAixAZJ41KmBr6
8sbzAqxfVLFjakT1wksP1uHtSDh40rT/m/Iv1hSjZjBkMA8GA1UdEwEB/wQFMAMB
Af8wHQYDVR0OBBYEFCSQK1T2KarQ4TepH/dx4LOUEckxMDIGA1UdEQQrMCmCFkxp
Y2Vuc2VJc3N1ZXIgZ1JQQyBUTFOCCWxvY2FsaG9zdIcEfwAAATAKBggqhkjOPQQD
AgNIADBFAiAPd/nQb+dMDJf8tXR4uuyf2OYVw3k4PtT1yvAGW86OWwIhAOQNXJtt
k9N0tKjSWqA7bss4fTAwqbcHpYR4hux4PWj+
-----END CERTIFICATE-----`

// RootAnchor is the PEM-encoded root certificate for token validation.
const RootAnchor = `-----BEGIN CERTIFICATE-----
MIIBaTCCAQ6gAwIBAgIIGNVgIyRffgAwCgYIKoZIzj0EAwIwIDEeMBwGA1UEAxMV
TGljZW5zZUlzc3VlciBSb290IENBMB4XDTI2MDkxNTAzMTUzMVoXDTI2MTIxNDAz
MTUzMVowIDEeMBwGA1UEAxMVTGljZW5zZUlzc3VlciBSb290IENBMFkwEwYHKoZI
zj0CAQYIKoZIzj0DAQcDQgAEJs/55MGRaNkumaPzdsbABFLirCx4hInsIN2BSgnp
krB6wzmIvnkoxUPEVbf4mk1Gd31jGiSnu7Knr8orZ8c+uqMyMDAwDwYDVR0TAQH/
BAUwAwEB/zAdBgNVHQ4EFgQUsUrHbjQSn1j2Ysh39NpIAbtBNMQwCgYIKoZIzj0E
AwIDSQAwRgIhAKKHcC1L+Erx9O668wGwLmjE8EKd+qH4qx7P1HlTcsqhAiEAl1gh
soFGh5TRUG022WyLjO0vFrc9P19PpJh3axq5Q4g=
-----END CERTIFICATE-----`
