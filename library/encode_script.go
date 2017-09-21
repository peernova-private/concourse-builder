package library

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"

	"time"

	"github.com/concourse-friends/concourse-builder/library/primitive"
	"github.com/concourse-friends/concourse-builder/project"
)

func GZipBase64Lines(script string, delimiter string) string {
	var buff bytes.Buffer
	gz, err := gzip.NewWriterLevel(&buff, gzip.BestCompression)
	if err != nil {
		panic(err)
	}
	gz.Header.ModTime, _ = time.Parse(time.RFC3339, time.RFC3339)

	if _, err := gz.Write([]byte(script)); err != nil {
		panic(err)
	}
	if err := gz.Close(); err != nil {
		panic(err)
	}

	str := base64.StdEncoding.EncodeToString(buff.Bytes())

	result := ""
	i := 0
	for ; i+80 < len(str); i += 80 {
		result += str[i:i+80] + delimiter
	}

	result += str[i:]
	return result
}

func EncodeScript(script string) (project.IRun, []interface{}) {
	bash := &primitive.Location{
		Volume: &primitive.Directory{
			Root: "/bin",
		},
		RelativePath: "bash",
	}

	arguments := []interface{}{
		"-c",
		fmt.Sprintf(`mkdir -p /tmp \
&& echo \
%s |\
    base64 --decode |\
    gzip -cfd > /tmp/script.sh \
&& cat /tmp/script.sh \
&& echo \
&& chmod 777 /tmp/script.sh \
&& /tmp/script.sh`,
			GZipBase64Lines(script, "\\\n")),
	}

	return bash, arguments
}
