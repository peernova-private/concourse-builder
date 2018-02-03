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
				        
                            dir("$GOPATH/$REPODIR/canary/ci/scripts") {
                                    sh './check_generation.sh'
                            }
            
                    },
                   "test": {
                         dir("$GOPATH/$REPODIR/cuneiform/src/tools/ci/scripts") {
                                    echo "Stage was successful"
                                
                        }
                        
                    },    
                          
                    
                    "Lint-Test": {
			            
                                dir("$GOPATH/$REPODIR/canary") {
                                    sh 'go get -u github.com/golang/lint/golint'
                                    sh 'golint ./...'
                                }
            
                        
			        },
                    "Copyright": {
				        
                            dir("$GOPATH/$REPODIR/canary/ci/scripts") {
                                    sh './check_copyright.sh'
                            }
            
                    },
                    "go_vendor": {
				        
                            dir("$GOPATH/$REPODIR/canary/ci/scripts") {
                                    sh './go_vendor.sh'
                            }
            
                    },
                    "Meta-Linter": {
				        
                            dir("$GOPATH/$REPODIR/canary/ci/scripts") {
                                    sh './meta_linter.sh'
                            }
            
                    },
                    "Vet" : {
                        dir("$GOPATH/$REPODIR/canary") {
		                sh 'go vet  ./...'
                        }
                    }
                )
	    }	
        }    
    } 
}   
