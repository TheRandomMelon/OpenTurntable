//go:build production

package bundle

import "embed"

//go:embed all:frontend/.output/public
var Bundle embed.FS
