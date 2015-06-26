{
  "cwd": "config",
  "standardEnv": true,
  "binary": "lanch_holed",
  "args": [
    "--addr", "{{.Addr}}",
    "--ca", "certs/{{.Ca}}",
    "--key", "certs/{{.Cakey}}",
    "--use-tls"
  ]
}
