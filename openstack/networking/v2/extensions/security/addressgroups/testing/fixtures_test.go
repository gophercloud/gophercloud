package testing

const AddressGroupListResponse = `
{
    "address_groups": [
        {
            "id": "8722e0e0-9cc9-4490-9660-8c9a5732fbb0",
            "project_id": "45977fa2dbd7482098dd68d0d8970117",
            "name": "ADDR_GP_1",
            "addresses": [
            "132.168.4.12/24"
            ]
        }
    ]
}
`

const AddressGroupGetResponse = `
{
    "address_group": {
        "description": "",
        "id": "8722e0e0-9cc9-4490-9660-8c9a5732fbb0",
        "name": "ADDR_GP_1",
        "project_id": "45977fa2dbd7482098dd68d0d8970117",
        "addresses": [
           "132.168.4.12/24"
        ]
    }
}
`

const AddressGroupCreateRequest = `
{
    "address_group": {
        "name": "ADDR_GP_1",
        "addresses": [
           "132.168.4.12/24"
        ]
    }
}
`

const AddressGroupCreateResponse = `
{
   "address_group": {
        "description": "",
        "id": "8722e0e0-9cc9-4490-9660-8c9a5732fbb0",
        "name": "ADDR_GP_1",
        "project_id": "45977fa2dbd7482098dd68d0d8970117",
        "addresses": [
           "132.168.4.12/24"
        ]
    }
}
`

const AddressGroupUpdateRequest = `
{
   "address_group": {
        "description": "new description",
        "name": "ADDR_GP_2"
    }
}
`

const AddressGroupUpdateResponse = `
{
   "address_group": {
        "description": "new description",
        "id": "8722e0e0-9cc9-4490-9660-8c9a5732fbb0",
        "name": "ADDR_GP_2",
        "project_id": "45977fa2dbd7482098dd68d0d8970117",
        "addresses": [
           "192.168.4.1/32"
        ]
    }
}
`

const AddressGroupAddAddressesRequest = `
{
    "addresses": ["192.168.4.1/32"]
}
`

const AddressGroupAddAddressesResponse = `
{
   "address_group": {
        "description": "original description",
        "id": "8722e0e0-9cc9-4490-9660-8c9a5732fbb0",
        "name": "ADDR_GP_1",
        "project_id": "45977fa2dbd7482098dd68d0d8970117",
        "addresses": [
           "132.168.4.12/24",
           "192.168.4.1/32"
        ]
    }
}
`

const AddressGroupRemoveAddressesRequest = `
{
    "addresses": ["192.168.4.1/32"]
}
`

const AddressGroupRemoveAddressesResponse = `
{
   "address_group": {
        "description": "original description",
        "id": "8722e0e0-9cc9-4490-9660-8c9a5732fbb0",
        "name": "ADDR_GP_1",
        "project_id": "45977fa2dbd7482098dd68d0d8970117",
        "addresses": [
           "132.168.4.12/24"
        ]
    }
}
`
