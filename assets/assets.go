package assets

import "embed"

//go:embed css/* css/** fonts/* fonts/**
var Assets embed.FS
