package project

import (
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

var logger = log.New(os.Stdout, "", log.LstdFlags)

type Project struct {
	Pipelines Pipelines
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

// Removes every job from the project that does not match the pattern.
// This allows for the creating simpler pipelines that get to the jobs of interest faster.
// Purfect for debugging. Note that on save the piplines will self expand with all of the
// dependencies needed for the jobs to work, witch will make them valid and operational.
func (p *Project) Filter(jobRegex *regexp.Regexp) error {
	for _, pipeline := range p.Pipelines {
		jobs, err := pipeline.AllJobs()
		if err != nil {
			return err
		}

		for i := 0; i < len(jobs); {
			if jobRegex.Match([]byte(jobs[i].Name)) {
				i++
			} else {
				pipeline.Jobs = append(jobs[:i], jobs[i+1:]...)
			}
		}

		pipeline.Jobs = jobs
	}

	return nil
}
