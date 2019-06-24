package testing

import (
	"fmt"
)

const FipDNS = `{
        "floating_network_id": "6d67c30a-ddb4-49a1-bec3-a65b286b4170",
        "router_id": null,
        "fixed_ip_address": null,
        "floating_ip_address": "192.0.0.4",
        "tenant_id": "017d8de156df4177889f31a9bd6edc00",
        "created_at": "2019-06-30T04:15:37",
        "updated_at": "2019-06-30T05:18:49",
        "status": "DOWN",
        "port_id": null,
        "id": "2f95fd2b-9f6a-4e8e-9e9a-2cbe286cbf9e",
        "router_id": "1117c30a-ddb4-49a1-bec3-a65b286b4170",
        "dns_domain": "local.",
        "dns_name": "test-fip"
    }`

const FipNoDNS = `{
        "floating_network_id": "90f742b1-6d17-487b-ba95-71881dbc0b64",
        "router_id": "0a24cb83-faf5-4d7f-b723-3144ed8a2167",
        "fixed_ip_address": "192.0.0.2",
        "floating_ip_address": "10.0.0.3",
        "tenant_id": "017d8de156df4177889f31a9bd6edc00",
        "created_at": "2019-06-30T04:15:37Z",
        "updated_at": "2019-06-30T05:18:49Z",
        "status": "DOWN",
        "port_id": "74a342ce-8e07-4e91-880c-9f834b68fa25",
        "id": "ada25a95-f321-4f59-b0e0-f3a970dd3d63",
        "router_id": "2227c30a-ddb4-49a1-bec3-a65b286b4170",
        "dns_domain": "",
        "dns_name": ""
    }`

const PoFw = `{
      "protocol": "tcp",
      "internal_ip_address": "10.0.0.24",
      "internal_port": 25,
      "internal_port_id": "070ef0b2-0175-4299-be5c-01fea8cca522",
      "external_port": 2229,
      "id": "1798dc82-c0ed-4b79-b12d-4c3c18f90eb2"
    }`

const PoFw_second = `{
	  "protocol": "tcp",
      "internal_ip_address": "10.0.0.11",
      "internal_port": 25,
      "internal_port_id": "1238be08-a2a8-4b8d-addf-fb5e2250e480",
      "external_port": 2230,
      "id": "e0a0274e-4d19-4eab-9e12-9e77a8caf3ea"
	}`

var ListResponse = fmt.Sprintf(`
{
    "floatingips": [
%s,
%s
    ]
}
`, FipDNS, FipNoDNS)

var ListResponseDNS = fmt.Sprintf(`
{
    "floatingips": [
%s
    ]
}
`, FipDNS)

var ListPortForwardingsResponse = fmt.Sprintf(`
{
    "port_forwardings": [
%s,
%s
    ]
}
`, PoFw, PoFw_second)
