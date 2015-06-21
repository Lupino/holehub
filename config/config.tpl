{
  "cwd": "config",
  "standardEnv": true,
  "binary": "lauch_hole_server",
  "args": [
    "--addr", "{{.Addr}}",
    "--ca", "certs/{{.Ca}}",
    "--key", "certs/{{.Cakey}}",
    "--use-tls"
  ]
}
