## updateip

A modern tool written in Go to update dynamic DNS records.

## Providers

* Namecheap Dynamic DNS

## Usage

Passwords are encoded (not encrypted), using base64. Save the output from the following into the `password` field in your `config.yaml`.

```bash
cat | base64
Password here
Control + D
```

See an [example config.yaml](/config.example.yaml)

## License

MIT