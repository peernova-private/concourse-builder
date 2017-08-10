package library

type ImageRegistry struct {
	Domain             string
	AwsAccessKeyId     string
	AwsSecretAccessKey string
}

var DockerHub = &ImageRegistry{
	Domain: "",
}
