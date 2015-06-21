{
  "cwd": "{{.Cwd}}",
  "standardEnv": true,
  "env": {
    "WANT_USER": ["_env", "want-${USER}"]
  },
  "binary": "{{.Command}}",
  "args": [
    "--addr", "{{.Addr}}",
    "--ca", "{{.Ca}}",
    "--key", "{{.Cakey}}",
    "--use-tls"
  ]
}
