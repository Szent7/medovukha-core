package types

type ContainerBaseInfo struct {
	Id          string   `json:"id"`
	Names       []string `json:"names"`
	ImageName   string   `json:"image"`
	Ports       []Port   `json:"ports"`
	Created     int64    `json:"created"`
	State       string   `json:"state"`
	IsMedovukha bool     `json:"isMedovukha"`
}

type Port struct {
	IP          string `json:"ip,omitempty"`
	PrivatePort uint16 `json:"privatePort"`
	PublicPort  uint16 `json:"publicPort,omitempty"`
	Type        string `json:"type"`
}

type ImageBaseInfo struct {
	Id      string   `json:"id"`
	Tags    []string `json:"tags"`
	Size    int64    `json:"size"`
	Created int64    `json:"created"`
}

type NetworkBaseInfo struct {
	Name          string   `json:"name"`
	Id            string   `json:"id"`
	Driver        string   `json:"driver"`
	EnableIPv6    bool     `json:"enableIPv6"`
	IPAMDriver    string   `json:"ipamDriver"`
	Subnet        []string `json:"subnet"`
	Gateway       []string `json:"gateway"`
	Attachable    bool     `json:"attachable"`
	DockerNetwork bool     `json:"dockerNetwork"`
}

type VolumeBaseInfo struct {
	Name       string `json:"name"`
	Driver     string `json:"driver"`
	Mountpoint string `json:"mountpoint"`
	Created    string `json:"created"`
}

type DeployFromGit struct {
	URL           string `json:"url"`
	Dockerfile    string `json:"dockerfile"`
	DockerCompose string `json:"docker_compose"`
	DockerRun     string `json:"docker_run"`
}
