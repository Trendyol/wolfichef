package models

type Environment struct {
	Key   string `json:"key" form:"key" valid:"type(string)"`
	Value string `json:"value" form:"value" valid:"type(string)"`
}

type Builder struct {
	Packages        []string      `json:"packages" form:"packages" valid:"type([]string),required"`
	ImageRepository string        `json:"registry" form:"registry" valid:"type(string),required"`
	ImageName       string        `json:"image" form:"image" valid:"type(string),required"`
	ImageTag        string        `json:"tag" form:"tag" valid:"type(string),required"`
	Username        string        `json:"username" form:"username" valid:"type(string),required"`
	Token           string        `json:"token" form:"token" valid:"type(string),required"`
	Entrypoint      string        `json:"entrypoint" form:"entrypoint" valid:"type(string)"`
	WorkDir         string        `json:"cwd" form:"cwd" valid:"type(string)"`
	Environments    []Environment `json:"environments" form:"environments"`
}
