package main

import "embed"

//go:embed assets/locales/*.toml
var LocaleFS embed.FS
