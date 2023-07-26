package models

type Package struct {
	Description   string `json:"description"`
	InstalledSize uint64 `json:"installed_size"`
	Name          string `json:"name"`
	Version       string `json:"version"`
}
