// Package region contains the compiled-in endpoint and certificate configuration.
// This file is specific to the test deployment.
// Do not modify at runtime — configuration is intentionally immutable.
package region

// Endpoint is the service address for this deployment.
const Endpoint = "10.116.105.118:9443"

// TLSAnchor is the PEM-encoded TLS certificate for transport security.
const TLSAnchor = `-----BEGIN CERTIFICATE-----
MIIBmDCCAT2gAwIBAgIBATAKBggqhkjOPQQDAjAhMR8wHQYDVQQDExZMaWNlbnNl
SXNzdWVyIGdSUEMgVExTMB4XDTI2MDcwNjExMzc1OVoXDTM2MDcwMzExMzc1OVow
ITEfMB0GA1UEAxMWTGljZW5zZUlzc3VlciBnUlBDIFRMUzBZMBMGByqGSM49AgEG
CCqGSM49AwEHA0IABDTcm+5ke/GV4Of10WOOTvg3uJnOpnEpbhDmMGuOZOI/tCFy
mvT7c+HQA8cIsEFn7E/0eKpr3KSBdVUfpJEpphqjZjBkMA8GA1UdEwEB/wQFMAMB
Af8wHQYDVR0OBBYEFM/zU0mVxnd+zp3/1UDrLIk22YbCMDIGA1UdEQQrMCmCFkxp
Y2Vuc2VJc3N1ZXIgZ1JQQyBUTFOCCWxvY2FsaG9zdIcEfwAAATAKBggqhkjOPQQD
AgNJADBGAiEA3uLZrJCGR2NN3uLgdP8gISPq4ciqF68dfVVHNzZO9SYCIQDJS/Pk
zaPeB6WVzTanGdBPuwmaJy+ohLUFA7M66DDXFw==
-----END CERTIFICATE-----`

// RootAnchor is the PEM-encoded root certificate for token validation.
const RootAnchor = `-----BEGIN CERTIFICATE-----
MIIBaTCCAQ6gAwIBAgIIGL+wWjhLJgAwCgYIKoZIzj0EAwIwIDEeMBwGA1UEAxMV
TGljZW5zZUlzc3VlciBSb290IENBMB4XDTI2MDcwNjExMzc1OVoXDTI2MTAwNDEx
Mzc1OVowIDEeMBwGA1UEAxMVTGljZW5zZUlzc3VlciBSb290IENBMFkwEwYHKoZI
zj0CAQYIKoZIzj0DAQcDQgAE37PZS1B88N2d/K7Y+7F3P7C/1vrvJsISthdVqZJV
ZURoTElWI2EvPmHWrWFkxDGUD0ya26vMceAtKosf7Y4EM6MyMDAwDwYDVR0TAQH/
BAUwAwEB/zAdBgNVHQ4EFgQU+O2bO6OmyQhL2y5/xWRhGM1nBggwCgYIKoZIzj0E
AwIDSQAwRgIhAMCoYV0skGgN+G8hG8StmvYpFbr5EUdn9cgWyi+BRGHSAiEA5RX+
bOIa1gf/iZvaLsJ6Cvv+20VHde/VYnezsIG/9qs=
-----END CERTIFICATE-----`
