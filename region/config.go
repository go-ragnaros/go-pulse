// Package region contains the compiled-in endpoint and certificate configuration.
// This file is specific to the lyra deployment.
// Do not modify at runtime — configuration is intentionally immutable.
package region

// Endpoint is the service address for this deployment.
const Endpoint = "10.77.15.129:9443"

// TLSAnchor is the PEM-encoded TLS certificate for transport security.
const TLSAnchor = `-----BEGIN CERTIFICATE-----
MIIBmDCCAT2gAwIBAgIBATAKBggqhkjOPQQDAjAhMR8wHQYDVQQDExZMaWNlbnNl
SXNzdWVyIGdSUEMgVExTMB4XDTI2MDcwODEyMDE1NFoXDTM2MDcwNTEyMDE1NFow
ITEfMB0GA1UEAxMWTGljZW5zZUlzc3VlciBnUlBDIFRMUzBZMBMGByqGSM49AgEG
CCqGSM49AwEHA0IABGqjLIL0TgQVL7e+/gtOmdjPxnohiiCPE7JxVrpDgo/xxnSK
CeY6ERuiaFBTs9QHMaWgZ0PO9LFFGHmVFxD8+JGjZjBkMA8GA1UdEwEB/wQFMAMB
Af8wHQYDVR0OBBYEFOS8MoACrTMCdlFx/s6B9di41z0vMDIGA1UdEQQrMCmCFkxp
Y2Vuc2VJc3N1ZXIgZ1JQQyBUTFOCCWxvY2FsaG9zdIcEfwAAATAKBggqhkjOPQQD
AgNJADBGAiEAqwti5KR65AuL4z2IO82up6floodSovvJpWY8gQCBO7cCIQD9hHYP
MvcFs5XJTq4KxtTByxH+zY6xyh528ljQw0gdxA==
-----END CERTIFICATE-----`

// RootAnchor is the PEM-encoded root certificate for token validation.
const RootAnchor = `-----BEGIN CERTIFICATE-----
MIIBaTCCAQ6gAwIBAgIIGMBO0XeTdAAwCgYIKoZIzj0EAwIwIDEeMBwGA1UEAxMV
TGljZW5zZUlzc3VlciBSb290IENBMB4XDTI2MDcwODEyMDE1NFoXDTI2MTAwNjEy
MDE1NFowIDEeMBwGA1UEAxMVTGljZW5zZUlzc3VlciBSb290IENBMFkwEwYHKoZI
zj0CAQYIKoZIzj0DAQcDQgAEOI4nuaGCtgyC24we1W8ioTfQ0CoUelFWPVz8udZY
S1j+Qg6A/qRwPwnoZPxpMBVi9KgpDTHbI5lzHjWa68vGxKMyMDAwDwYDVR0TAQH/
BAUwAwEB/zAdBgNVHQ4EFgQULToMb8UKsTvwYyg22ttviyL9TaQwCgYIKoZIzj0E
AwIDSQAwRgIhANyn6d9zxzN+NugwsrFzksrrxJ14SUxoBzHkQ5yVNmnCAiEApg4R
/gwp0jVbjjshtV4t3G1qTECj/EeZZkp6UbSlUW8=
-----END CERTIFICATE-----`
