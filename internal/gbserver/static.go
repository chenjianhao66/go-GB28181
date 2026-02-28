package gbserver

import "embed"

//go:embed ui/dist/*
var staticFiles embed.FS
