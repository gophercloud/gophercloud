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

// Taken from: http://docs.openstack.org/api/openstack-identity-service/2.0/content/GET_listExtensions_v2.0_extensions_.html#GET_listExtensions_v2.0_extensions_-Request
const queryResults = `{
    "extensions":[{
            "name": "Reset Password Extension",
            "namespace": "http://docs.rackspacecloud.com/identity/api/ext/rpe/v2.0",
            "alias": "RS-RPE",
            "updated": "2011-01-22T13:25:27-06:00",
            "description": "Adds the capability to reset a user's password. The user is emailed when the password has been reset.",
            "links":[{
                    "rel": "describedby",
                    "type": "application/pdf",
                    "href": "http://docs.rackspacecloud.com/identity/api/ext/identity-rpe-20111111.pdf"
                },
                {
                    "rel": "describedby",
                    "type": "application/vnd.sun.wadl+xml",
                    "href": "http://docs.rackspacecloud.com/identity/api/ext/identity-rpe.wadl"
                }
            ]
        },
        {
            "name": "User Metadata Extension",
            "namespace": "http://docs.rackspacecloud.com/identity/api/ext/meta/v2.0",
            "alias": "RS-META",
            "updated": "2011-01-12T11:22:33-06:00",
            "description": "Allows associating arbritrary metadata with a user.",
            "links":[{
                    "rel": "describedby",
                    "type": "application/pdf",
                    "href": "http://docs.rackspacecloud.com/identity/api/ext/identity-meta-20111201.pdf"
                },
                {
                    "rel": "describedby",
                    "type": "application/vnd.sun.wadl+xml",
                    "href": "http://docs.rackspacecloud.com/identity/api/ext/identity-meta.wadl"
                }
            ]
        }
    ],
    "extensions_links":[]
}`

// Same as queryResults above, but with a bogus JSON envelop.
const bogusExtensionsResults = `{
    "explosions":[{
            "name": "Reset Password Extension",
            "namespace": "http://docs.rackspacecloud.com/identity/api/ext/rpe/v2.0",
            "alias": "RS-RPE",
            "updated": "2011-01-22T13:25:27-06:00",
            "description": "Adds the capability to reset a user's password. The user is emailed when the password has been reset.",
            "links":[{
                    "rel": "describedby",
                    "type": "application/pdf",
                    "href": "http://docs.rackspacecloud.com/identity/api/ext/identity-rpe-20111111.pdf"
                },
                {
                    "rel": "describedby",
                    "type": "application/vnd.sun.wadl+xml",
                    "href": "http://docs.rackspacecloud.com/identity/api/ext/identity-rpe.wadl"
                }
            ]
        },
        {
            "name": "User Metadata Extension",
            "namespace": "http://docs.rackspacecloud.com/identity/api/ext/meta/v2.0",
            "alias": "RS-META",
            "updated": "2011-01-12T11:22:33-06:00",
            "description": "Allows associating arbritrary metadata with a user.",
            "links":[{
                    "rel": "describedby",
                    "type": "application/pdf",
                    "href": "http://docs.rackspacecloud.com/identity/api/ext/identity-meta-20111201.pdf"
                },
                {
                    "rel": "describedby",
                    "type": "application/vnd.sun.wadl+xml",
                    "href": "http://docs.rackspacecloud.com/identity/api/ext/identity-meta.wadl"
                }
            ]
        }
    ],
    "extensions_links":[]
}`
