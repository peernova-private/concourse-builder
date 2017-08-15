package library

type ImageRegistry struct {
	Domain             string
	AwsAccessKeyId     string
	AwsSecretAccessKey string
}

func (ir *ImageRegistry) Public() bool {
	return ir.AwsAccessKeyId == "" && ir.AwsSecretAccessKey == ""
}

var DockerHub = &ImageRegistry{
	Domain: "",
}
