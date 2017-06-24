package project

import (
	"io/ioutil"
	"log"
	"os"
)

var logger = log.New(os.Stdout, "", log.LstdFlags)

type Project struct {
	Pipelines []*Pipeline
}

func (p *Project) Deploy() error {
	for _, pipeline := range p.Pipelines {
		tmpfile, err := ioutil.TempFile("", "pipeline")
		if err != nil {
			log.Panic(err)
		}

		logger.Printf("Saving pipeline %s in %s", pipeline.Name, tmpfile.Name())

		defer func() {
			if tmpfile != nil {
				err = tmpfile.Close()
				if err != nil {
					logger.Printf("Warining: %s", err.Error())
				}
			}
		}()
		defer os.Remove(tmpfile.Name())

		err = pipeline.Save(tmpfile)
		if err != nil {
			return err
		}
	}

	return nil
}
