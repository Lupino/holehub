{
  "cwd": "config",
  "standardEnv": true,
  "binary": "lanch_hole_server",
  "args": [
    "--addr", "{{.Addr}}",
    "--ca", "certs/{{.Ca}}",
    "--key", "certs/{{.Cakey}}",
    "--use-tls"
  ]
}
