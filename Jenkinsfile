pipeline {
   
    agent { 
        node { 
           label 'jenkins-slave' 
        } 
    } 
    
    environment {
        GOROOT = "/var/opt/go"
        GRDLHOME = "/var/opt/gradle"
        GOPATH = "$WORKSPACE"
        GOBIN = "${env.GOPATH}/bin"
        REPODIR = "src/github.com/peernova-private"
        PATH = "${env.GOROOT}/bin:${env.GRDLHOME}/bin:/usr/local/bin:/bin:/usr/bin:/usr/local/sbin:/usr/sbin:/sbin"
    }
    
    stages {
        
        stage('Prepare') { // Setup Check
            steps {
                sh '${GOROOT}/bin/go version'
                script {    
                    gver = sh (returnStdout: true, script: "git --version").trim()
                    echo "Git Version: ${gver}"
                }
                dir (env.GOBIN) {
                    deleteDir()
                }
            }
        }
        stage('Checkout') {
            steps {
                parallel (
                    "Concourse-Builder" : {
                        dir('src/github.com/peernova-private/concourse-builder') {
                            checkout scm
                        } 
                    },
		    "Integration" : {
                        dir("$REPODIR/integration") {
 			    checkout scm
                        }        
                    } 
                )
            }
    	}
	stage('pre-check') {
	   steps {
                parallel (
		     "Code-gen": {
				        
                            env.context="jenkins-ci/@code-gen"
                            setBuildStatus ("${context}", 'Jenkins CI build pending', 'PENDING')
                            dir("$GOPATH/$REPODIR/canary/ci/scripts") {
                                    sh './check_generation.sh'
                            }
            
                    },
                   "test": {
                        env.context="jenkins-ci/@test"
                        setBuildStatus ("${context}", 'Jenkins CI build pending', 'PENDING')
                         dir("$GOPATH/$REPODIR/cuneiform/src/tools/ci/scripts") {
                                    echo "Stage was successful"
                                
                        }
                        
                    },    
                          
                    
                    "Lint-Test": {
			            
                                env.context="jenkins-ci/@lint"
                                setBuildStatus ("${context}", 'Jenkins CI build pending', 'PENDING')
                                dir("$GOPATH/$REPODIR/canary") {
                                    sh 'go get -u github.com/golang/lint/golint'
                                    sh 'golint ./...'
                                }
            
                        
			        },
                    "Copyright": {
				        
                            env.context="jenkins-ci/@copyright"
                            setBuildStatus ("${context}", 'Jenkins CI build pending', 'PENDING')
                            dir("$GOPATH/$REPODIR/canary/ci/scripts") {
                                    sh './check_copyright.sh'
                            }
            
                    },
                    "go_vendor": {
				        
                            env.context="jenkins-ci/@go-vendor"
                            setBuildStatus ("${context}", 'Jenkins CI build pending', 'PENDING')
                            dir("$GOPATH/$REPODIR/canary/ci/scripts") {
                                    sh './go_vendor.sh'
                            }
            
                    },
                    "Meta-Linter": {
				        
                            env.context="jenkins-ci/@meta-linter"
                            setBuildStatus ("${context}", 'Jenkins CI build pending', 'PENDING')
                            dir("$GOPATH/$REPODIR/canary/ci/scripts") {
                                    sh './meta_linter.sh'
                            }
            
                    },
                    "Vet" : {
                        env.context="jenkins-ci/@vet"
                        setBuildStatus ("${context}", 'Jenkins CI build pending', 'PENDING')
                        dir("$GOPATH/$REPODIR/canary") {
		                sh 'go vet  ./...'
                        }
                    }
                )
	    }	
        }    
    } 
}   
