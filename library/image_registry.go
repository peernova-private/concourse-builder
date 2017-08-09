package library

type ImageRegistry struct {
	Domain string
}

var DockerHub = &ImageRegistry{
	Domain: "",
}
