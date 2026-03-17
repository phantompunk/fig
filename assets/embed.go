package assets

import "embed"

//go:embed *.flf
var FontFS embed.FS

//go:embed fonts.yaml
var FontsYAML []byte
