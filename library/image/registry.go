package image

type Registry struct {
	Domain             string
	AwsAccessKeyId     string
	AwsSecretAccessKey string
}

func (ir *Registry) Public() bool {
	return ir.AwsAccessKeyId == "" && ir.AwsSecretAccessKey == ""
}

var DockerHub = &Registry{
	Domain: "",
}
