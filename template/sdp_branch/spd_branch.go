package sdpBranch

import (
	"log"

	"github.com/concourse-friends/concourse-builder/library"
	"github.com/concourse-friends/concourse-builder/project"
)

type Specification interface {
	BootstrapSpecification
	ModifyJobs(resourceRegistry *project.ResourceRegistry) (project.Jobs, error)
	VerifyJobs(resourceRegistry *project.ResourceRegistry) (project.Jobs, error)
}

func addPipelineResource(
	mainPipeline *project.Pipeline,
	selfUpdateJob *project.Job,
	pipelineJobResource *project.JobResource) error {

	jobs, err := mainPipeline.AllJobs()
	if err != nil {
		return err
	}

	excludeJobs, err := mainPipeline.JobsFor(project.JobsSet{selfUpdateJob: struct{}{}})
	if err != nil {
		return err
	}

	for _, job := range jobs {
		ok := true
		for _, exclude := range excludeJobs {
			if job.Name == exclude.Name {
				ok = false
				break
			}
		}
		if ok {
			log.Printf("Add pipeline resource to job %s", job.Name)
			job.ExtraResources = append(job.ExtraResources, pipelineJobResource)
		}
	}

	return nil
}

func GenerateProject(specification Specification) (*project.Project, error) {
	concourseBuilderGit, err := specification.ConcourseBuilderGit()
	if err != nil {
		return nil, err
	}

	concourseBuilderBranch := concourseBuilderGit.Source.(*library.GitSource).Branch

	imageRegistry, err := specification.DeployImageRegistry()
	if err != nil {
		return nil, err
	}

	concourse, err := specification.Concourse()
	if err != nil {
		return nil, err
	}

	environment, err := specification.Environment()
	if err != nil {
		return nil, err
	}

	concourseBuilderPipeline := project.NewPipeline()
	concourseBuilderPipeline.AllJobsGroup = project.AllJobsGroupFirst
	concourseBuilderPipeline.Name = project.ConvertToPipelineName(concourseBuilderBranch.FriendlyName() + "-cb")

	concourseBuilderLinuxImage, err := specification.LinuxImage(concourseBuilderPipeline.ResourceRegistry)
	if err != nil {
		return nil, err
	}

	concourseBuilderPipeline.Jobs = project.Jobs{
		library.AllImages(&library.AllImagesArgs{
			LinuxImageResource:  concourseBuilderLinuxImage,
			ConcourseBuilderGit: concourseBuilderGit,
			ImageRegistry:       imageRegistry,
			ResourceRegistry:    concourseBuilderPipeline.ResourceRegistry,
			Concourse:           concourse,
		}),
	}

	mainPipeline := project.NewPipeline()
	mainPipeline.AllJobsGroup = project.AllJobsGroupFirst
	mainPipeline.Name = project.ConvertToPipelineName(specification.Branch().FriendlyName() + "-sdpb")

	if !concourseBuilderBranch.IsImage() {
		mainPipeline.ReuseFrom = append(mainPipeline.ReuseFrom, concourseBuilderPipeline.ResourceRegistry)
	}

	linuxImage, err := specification.LinuxImage(mainPipeline.ResourceRegistry)
	if err != nil {
		return nil, err
	}

	goImage, err := specification.GoImage(mainPipeline.ResourceRegistry)
	if err != nil {
		return nil, err
	}

	generateProjectLocation, err := specification.GenerateProjectLocation(mainPipeline.ResourceRegistry)
	if err != nil {
		return nil, err
	}

	selfUpdateJob, pipelineResource := library.SelfUpdateJob(&library.SelfUpdateJobArgs{
		LinuxImageResource:      linuxImage,
		ConcourseBuilderGit:     concourseBuilderGit,
		ImageRegistry:           imageRegistry,
		GoImage:                 goImage,
		ResourceRegistry:        mainPipeline.ResourceRegistry,
		Concourse:               concourse,
		Environment:             environment,
		GenerateProjectLocation: generateProjectLocation,
	})

	mainPipeline.Jobs = project.Jobs{
		selfUpdateJob,
	}

	pipelineJobResource := mainPipeline.ResourceRegistry.JobResource(pipelineResource, true, nil)

	var modifyGroup = &project.JobGroup{
		Name: "modify",
	}

	var modifyJobs project.Jobs
	if specification.Branch().IsTask() {
		modifyJobs, err = specification.ModifyJobs(mainPipeline.ResourceRegistry)
		if err != nil {
			return nil, err
		}

		for _, job := range modifyJobs {
			job.AddToGroup(modifyGroup)
			job.AddJobToRunAfter(selfUpdateJob)
		}
		mainPipeline.Jobs = append(mainPipeline.Jobs, modifyJobs...)
	}

	var verifyGroup = &project.JobGroup{
		Name: "verify",
	}

	verifyJobs, err := specification.VerifyJobs(mainPipeline.ResourceRegistry)
	if err != nil {
		return nil, err
	}

	for _, job := range verifyJobs {
		job.AddToGroup(verifyGroup)
		job.AddJobToRunAfter(selfUpdateJob)
		job.AddJobToRunAfter(modifyJobs...)
	}

	mainPipeline.Jobs = append(mainPipeline.Jobs, verifyJobs...)

	err = addPipelineResource(mainPipeline, selfUpdateJob, pipelineJobResource)
	if err != nil {
		return nil, err
	}

	prj := &project.Project{
		Pipelines: project.Pipelines{
			concourseBuilderPipeline,
			mainPipeline,
		},
	}

	return prj, nil
}
