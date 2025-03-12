//go:build !production

package bundle

import "embed"

//go:embed all:frontend/.nuxt
var Bundle embed.FS
