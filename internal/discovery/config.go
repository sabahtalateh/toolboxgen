package discovery

import (
	"fmt"

	"github.com/sabahtalateh/toolboxgen/internal/discovery/mod"
)

type Conf struct {
	RootDir      string
	Mod          *mod.Module
	ProviderDirs []string
	ExcludeDirs  []string
}

type ConfBuilder struct {
	withRoot        string
	withProviders   []string
	withExcludeDirs []string
}

func ConfigBuilder() *ConfBuilder {
	return &ConfBuilder{}
}

func (b *ConfBuilder) Root(x string) *ConfBuilder {
	b.withRoot = x
	return b
}

func (b *ConfBuilder) Provider(x string) *ConfBuilder {
	b.withProviders = append(b.withProviders, x)
	return b
}

func (b *ConfBuilder) ExcludeDirs(x []string) *ConfBuilder {
	b.withExcludeDirs = x
	return b
}

func (b *ConfBuilder) Build() (*Conf, error) {
	module, err := mod.LookupAtDir(b.withRoot, true)
	if err != nil {
		return nil, err
	}

	for _, provider := range b.withProviders {
		providerMod, err := mod.LookupAtDir(provider, false)
		if err != nil {
			return nil, err
		}
		if providerMod.Path != module.Path || providerMod.Dir != module.Dir {
			return nil, fmt.Errorf(
				"validateProvider module: `%s` at `%s`, different from root module: `%s` at `%s. all directories should be within same module",
				providerMod.Path, providerMod.Dir, module.Path, module.Dir,
			)
		}
	}

	return &Conf{
		Mod:          module,
		RootDir:      b.withRoot,
		ProviderDirs: b.withProviders,
		ExcludeDirs:  b.withExcludeDirs,
	}, nil
}
