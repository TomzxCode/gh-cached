package version

import (
	"runtime/debug"
	"sync"
)

var (
	once    sync.Once
	version string
)

func Set(v string) {
	version = v
}

func Get() string {
	once.Do(func() {
		if version != "" {
			return
		}
		info, ok := debug.ReadBuildInfo()
		if !ok {
			version = "dev"
			return
		}
		for _, s := range info.Settings {
			if s.Key == "vcs.revision" {
				v := s.Value
				if len(v) > 7 {
					v = v[:7]
				}
				version = v
				return
			}
		}
		v := info.Main.Version
		if v != "" && v != "(devel)" {
			version = v
			return
		}
		version = "dev"
	})
	return version
}
