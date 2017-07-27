package library

type ImageRepository struct {
	Domain string
}

var DockerHub = &ImageRepository{
	Domain: "",
}
