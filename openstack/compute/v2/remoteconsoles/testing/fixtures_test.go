package testing

// RemoteConsoleCreateRequest represents a request to create a remote console.
const RemoteConsoleCreateRequest = `
{
    "remote_console": {
        "protocol": "vnc",
        "type": "novnc"
    }
}
`

// RemoteConsoleCreateResult represents a raw server response to the RemoteConsoleCreateRequest.
const RemoteConsoleCreateResult = `
{
    "remote_console": {
        "protocol": "vnc",
        "type": "novnc",
        "url": "http://192.168.0.4:6080/vnc_auto.html?token=9a2372b9-6a0e-4f71-aca1-56020e6bb677"
    }
}
`

// ConsoleGetResult represents a raw server response to using the console auth token
const ConsoleGetResult = `
{
    "console": {
        "instance_uuid": "933c8963-8a83-43bb-9618-98ee8025044d",
        "host": "10.0.2.224",
        "port": 5900,
        "tls_port": 5901,
        "internal_access_path": "null"}
}`
