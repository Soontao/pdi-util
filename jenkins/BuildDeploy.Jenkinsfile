// PDI Solution Continuous Deployment Jenkinsfile
// Jenkins must install the 'fileOperations' and 'withCredentials' plugins
// For the 'Deploy' stage, the source & target tenant must have the same 'Release' version

// Recommendation, please use this CD pipeline after first-time deployment
// first-time deploy is not testes by 'pdi-util project' by now

// Recommendation, please import this jenkinsfile as normal 'Pipeline' (instead of Multi Branch Pipeline), 
// so that you could control build/deployment flexible

// This pipeline will assemble a new package in the source tenant and popup a new version (create patch)
// This pipeline will deploy the assembled package to target tenant and activate it after uploading

def cliVersion = "v1.9.10"

def utilDownloadURI = "https://oss-theo.oss-cn-shenzhen.aliyuncs.com/download/pdiutil-${cliVersion}-darwin-amd64.zip"

def cliName = "pdiutil-${cliVersion}-darwin-amd64"

def utilZipFileName = "pdiutil.zip"

// please use the readable 'Solution Description' instead of the technical 'Solution ID'
// so that pdi-util tool can auto detect the correct solution id after patch solution created
def solution = ""

// source tenant
def devTenant = ""

// target tenant
def targetTenant = ""

// Jenkins user/password credential ID
def devUserCredentialId = ''

// Jenkins user/password credential ID
def targetUserCredentialId = ''

// specific the current ByD/C4C release version
def releaseVersion = ''

pipeline {

	agent any

	environment { 
		// tenant release version
		TENANT_RELEASE = "${releaseVersion}"
		// for assemble 
		PDI_TENANT_HOST = "${devTenant}"
		// for assemble
		SOLUTION_NAME = "${solution}"
		// for deploy
		TARGET_TENANT = "${targetTenant}"
		// for deploy
		SOURCE_SOLUTION_NAME = "${solution}"

	}

	stages {

		// download pdi tool from CDN
		stage('Setup') {
			steps {
				script {
					fileOperations([
						fileDownloadOperation(targetLocation: '', userName: '', password: '', targetFileName: utilZipFileName, url: utilDownloadURI),
						fileUnZipOperation(filePath: utilZipFileName, targetLocation: '')
					])
				}
			}
		}

		stage('Build') {

			steps {
				lock(devTenant) {
					withCredentials([usernamePassword(credentialsId: devUserCredentialId, passwordVariable: 'PDI_PASSWORD', usernameVariable: 'PDI_USER')]) {
						sh script: "./${cliName} package assemble", label: "Assemble & Download"
						archiveArtifacts artifacts: '*.zip', excludes: utilZipFileName
					}
				}
			}

		}

		stage('Deploy') {

			steps {
				lock(targetTenant) {
					withCredentials([usernamePassword(credentialsId: devUserCredentialId, passwordVariable: 'PDI_PASSWORD', usernameVariable: 'PDI_USER'), usernamePassword(credentialsId: targetUserCredentialId, passwordVariable: 'TARGET_TENANT_PASSWORD', usernameVariable: 'TARGET_TENANT_USER')]) {
						sh script: "./${cliName} solution deploy", label: "Deploy to target tenant"
					}
				}
			}

		}

		stage('Clean') {

			// clean workspace 
			steps{
				deleteDir()
			}

		}

	}

}