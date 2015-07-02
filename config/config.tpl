{
  "cwd": "config",
  "standardEnv": true,
  "binary": "lanch_holed",
  "args": [
    "--addr", "{{.Scheme}}://{{.Host}}:{{.Port}}",
    "--ca", "certs/{{.Ca}}",
    "--key", "certs/{{.Cakey}}",
    "--use-tls"
  ]
}
