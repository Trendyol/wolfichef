package index

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/wolfi-dev/wolfictl/pkg/index"
	"golang.org/x/exp/maps"
	"golang.org/x/mod/semver"
	"trendyol.com/security/appsec/devsecops/wolfichef/internal/index/models"
)

type Packages []models.Package

func Prepare() Packages {
	packages := make(map[string]models.Package)
	idx, err := index.Index("x86_64", "https://packages.wolfi.dev/os")
	if err != nil {
		log.Fatal(err)
	}

	for _, pkg := range idx.Packages {
		p := models.Package{
			Name:          pkg.Name,
			Description:   pkg.Description,
			InstalledSize: pkg.InstalledSize,
			Version:       pkg.Version,
		}
		if _, ok := packages[p.Name]; ok {
			if semver.Compare(p.Version, packages[pkg.Name].Version) > 0 {
				packages[pkg.Name] = p
			}
		} else {
			packages[pkg.Name] = p
		}
	}

	return maps.Values(packages)
}

func (p *Packages) Save() {
	output, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}

	path, err := filepath.Abs("packages.json")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(path, output, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
