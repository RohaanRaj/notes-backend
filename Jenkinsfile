pipeline {
	agent any

	stages {
		
		stage('Checkout') {
			steps{
				checkout scm
			}
		}

		stage('Docker build'){
			steps{
				sh 'docker build -t notes-api .'
			}
		}

		stage('compose config check'){
			steps{
				sh 'docker compose config'
			}
		}

		stage('Compose build'){
			steps{
				sh 'docker compose build'
			}
		}

		stage('Compose up'){
			steps{
				sh 'docker compose up -d'
				sh 'sleep 15'
			}
		}

		stage('Test'){
			steps{
				sh '''
				curl -X POST -d '{"email":"abc@gmail.com", "password":"supersecret"}' localhost:9000/register
				'''
			}
		}

		stage('Teardown'){
			steps{
				sh 'docker compose down -v'
			}
		}
}
}


