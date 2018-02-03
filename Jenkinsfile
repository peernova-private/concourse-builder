pipeline {
  agent {
    node {
      label 'jenkins-slave'
    }
    
  }
  stages {
    stage('Prepare') {
      steps {
        sh '${GOROOT}/bin/go version'
        script {
          gver = sh (returnStdout: true, script: "git --version").trim()
          echo "Git Version: ${gver}"
        }
        
        dir(path: env.GOBIN) {
          deleteDir()
        }
        
      }
    }
    stage('Checkout') {
      parallel {
        stage('Concourse-Builder') {
          steps {
            dir(path: 'src/github.com/peernova-private/concourse-builder') {
              checkout scm
            }
            
          }
        }
        stage('Integration') {
          steps {
            dir(path: "$REPODIR/integration") {
              checkout scm
            }
            
          }
        }
      }
    }
    stage('pre-check') {
      parallel {
        stage('Code-gen') {
          steps {
            dir(path: "$GOPATH/$REPODIR/canary/ci/scripts") {
              sh './check_generation.sh'
            }
            
            githubNotify(status: 'failed', description: 'fail', account: 'peernova-private', context: 'jenkins-ci@code_gen', credentialsId: 'peernova-go', repo: 'canary')
            catchError() {
              githubNotify(status: 'failed', description: 'failed', account: 'peernova-private', credentialsId: 'peernova-go', repo: 'canary')
            }
            
          }
        }
        stage('test') {
          steps {
            dir(path: "$GOPATH/$REPODIR/cuneiform/src/tools/ci/scripts") {
              echo 'Stage was successful'
            }
            
          }
        }
        stage('Lint-Test') {
          steps {
            dir(path: "$GOPATH/$REPODIR/canary") {
              sh 'go get -u github.com/golang/lint/golint'
              sh 'golint ./...'
            }
            
          }
        }
        stage('Copyright') {
          steps {
            dir(path: "$GOPATH/$REPODIR/canary/ci/scripts") {
              sh './check_copyright.sh'
            }
            
          }
        }
        stage('go_vendor') {
          steps {
            dir(path: "$GOPATH/$REPODIR/canary/ci/scripts") {
              sh './go_vendor.sh'
            }
            
          }
        }
        stage('Meta-Linter') {
          steps {
            dir(path: "$GOPATH/$REPODIR/canary/ci/scripts") {
              sh './meta_linter.sh'
            }
            
          }
        }
        stage('Vet') {
          steps {
            dir(path: "$GOPATH/$REPODIR/canary") {
              sh 'go vet  ./...'
            }
            
          }
        }
      }
    }
  }
  environment {
    GOROOT = '/var/opt/go'
    GRDLHOME = '/var/opt/gradle'
    GOPATH = "$WORKSPACE"
    GOBIN = "${env.GOPATH}/bin"
    REPODIR = 'src/github.com/peernova-private'
    PATH = "${env.GOROOT}/bin:${env.GRDLHOME}/bin:/usr/local/bin:/bin:/usr/bin:/usr/local/sbin:/usr/sbin:/sbin"
  }
}