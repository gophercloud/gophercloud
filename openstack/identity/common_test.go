package identity

// Taken from: http://docs.openstack.org/api/openstack-identity-service/2.0/content/POST_authenticate_v2.0_tokens_.html
const authResultsOK = `{
    "access":{
        "token":{
            "id": "ab48a9efdfedb23ty3494",
            "expires": "2010-11-01T03:32:15-05:00",
            "tenant":{
                "id": "t1000",
                "name": "My Project"
            }
        },
        "user":{
            "id": "u123",
            "name": "jqsmith",
            "roles":[{
                    "id": "100",
                    "name": "compute:admin"
                },
                {
                    "id": "101",
                    "name": "object-store:admin",
                    "tenantId": "t1000"
                }
            ],
            "roles_links":[]
        },
        "serviceCatalog":[{
                "name": "Cloud Servers",
                "type": "compute",
                "endpoints":[{
                        "tenantId": "t1000",
                        "publicURL": "https://compute.north.host.com/v1/t1000",
                        "internalURL": "https://compute.north.internal/v1/t1000",
                        "region": "North",
                        "versionId": "1",
                        "versionInfo": "https://compute.north.host.com/v1/",
                        "versionList": "https://compute.north.host.com/"
                    },
                    {
                        "tenantId": "t1000",
                        "publicURL": "https://compute.north.host.com/v1.1/t1000",
                        "internalURL": "https://compute.north.internal/v1.1/t1000",
                        "region": "North",
                        "versionId": "1.1",
                        "versionInfo": "https://compute.north.host.com/v1.1/",
                        "versionList": "https://compute.north.host.com/"
                    }
                ],
                "endpoints_links":[]
            },
            {
                "name": "Cloud Files",
                "type": "object-store",
                "endpoints":[{
                        "tenantId": "t1000",
                        "publicURL": "https://storage.north.host.com/v1/t1000",
                        "internalURL": "https://storage.north.internal/v1/t1000",
                        "region": "North",
                        "versionId": "1",
                        "versionInfo": "https://storage.north.host.com/v1/",
                        "versionList": "https://storage.north.host.com/"
                    },
                    {
                        "tenantId": "t1000",
                        "publicURL": "https://storage.south.host.com/v1/t1000",
                        "internalURL": "https://storage.south.internal/v1/t1000",
                        "region": "South",
                        "versionId": "1",
                        "versionInfo": "https://storage.south.host.com/v1/",
                        "versionList": "https://storage.south.host.com/"
                    }
                ]
            },
            {
                "name": "DNS-as-a-Service",
                "type": "dnsextension:dns",
                "endpoints":[{
                        "tenantId": "t1000",
                        "publicURL": "https://dns.host.com/v2.0/t1000",
                        "versionId": "2.0",
                        "versionInfo": "https://dns.host.com/v2.0/",
                        "versionList": "https://dns.host.com/"
                    }
                ]
            }
        ]
    }
}`
