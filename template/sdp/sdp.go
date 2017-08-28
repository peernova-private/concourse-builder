package sdp

import (
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/concourse-friends/concourse-builder/library"
	"github.com/concourse-friends/concourse-builder/library/primitive"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/template/sdp_branch"
)

type Specification interface {
	Concourse() (*primitive.Concourse, error)
	DeployImageRegistry() (*library.ImageRegistry, error)
	ConcourseBuilderGitSource() (*library.GitSource, error)
	GenerateProjectLocation(resourceRegistry *project.ResourceRegistry, branch *primitive.GitBranch) (project.IRun, error)
	TargetGitRepo() (*library.GitRepo, error)
	Environment() (map[string]interface{}, error)
}

const BranchesFileEnvVar = "BRANCHES_FILE"

func BootstrapBranches() ([]string, error) {
	branchesFile, exist := os.LookupEnv(BranchesFileEnvVar)
	if !exist {
		return nil, nil
	}
	branches, err := ioutil.ReadFile(branchesFile)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(branches), "\n"), err
}

func GenerateProject(specification Specification) (*project.Project, error) {
	prj := &project.Project{}

	branches, err := BootstrapBranches()
	if err != nil {
		return nil, err
	}

	reInclude := regexp.MustCompile(`master|release/.*|feature/.*|task/.*`)
	reExclude := regexp.MustCompile(`task/.*?-pr$`)

	for _, branch := range branches {
		if !reInclude.MatchString(branch) || reExclude.MatchString(branch) {
			log.Printf("Branch %s does not fit the expected name protocol", branch)
			continue
		}
		log.Printf("Preparing pipeline for branch %s", branch)

		branchSpecification := &BranchBootstrapSpecification{
			Specification: specification,
			TargetBranch:  &primitive.GitBranch{Branch: branch},
		}

		project, err := sdpBranch.GenerateBootstrapProject(branchSpecification)
		if err != nil {
			return nil, err
		}

		prj.Pipelines = append(prj.Pipelines, project.Pipelines...)
	}

	mainPipeline := project.NewPipeline()
	mainPipeline.AllJobsGroup = project.AllJobsGroupFirst

	concourseBuilderGitSource, err := specification.ConcourseBuilderGitSource()
	if err != nil {
		return nil, err
	}

	mainPipeline.Name = project.ConvertToPipelineName(concourseBuilderGitSource.Branch + "-sdp")

	imageRegistry, err := specification.DeployImageRegistry()
	if err != nil {
		return nil, err
	}

	concourse, err := specification.Concourse()
	if err != nil {
		return nil, err
	}

	generateProjectLocation, err := specification.GenerateProjectLocation(mainPipeline.ResourceRegistry, nil)
	if err != nil {
		return nil, err
	}

	environment, err := specification.Environment()
	if err != nil {
		return nil, err
	}

	selfUpdateJob := library.SelfUpdateJob(&library.SelfUpdateJobArgs{
		ConcourseBuilderGitSource: concourseBuilderGitSource,
		ImageRegistry:             imageRegistry,
		ResourceRegistry:          mainPipeline.ResourceRegistry,
		Tag:                       library.ConvertToImageTag(concourseBuilderGitSource.Branch),
		Concourse:                 concourse,
		Environment:               environment,
		GenerateProjectLocation:   generateProjectLocation,
	})

	targetGit, err := specification.TargetGitRepo()
	if err != nil {
		return nil, err
	}

	branchesJob := BranchesJob(&BranchesJobArgs{
		ConcourseBuilderGitSource: concourseBuilderGitSource,
		ImageRegistry:             imageRegistry,
		ResourceRegistry:          mainPipeline.ResourceRegistry,
		Tag:                       library.ConvertToImageTag(concourseBuilderGitSource.Branch),
		Concourse:                 concourse,
		TargetGitRepo:             targetGit,
		Environment:               environment,
		GenerateProjectLocation:   generateProjectLocation,
	})

	mainPipeline.Jobs = project.Jobs{
		selfUpdateJob,
		branchesJob,
	}

	prj.Pipelines = append(prj.Pipelines, mainPipeline)
	return prj, nil
}
