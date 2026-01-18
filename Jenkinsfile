pipeline {
	agent any

	stages {
		
		stage('Checkout') {
			steps{
				checkout scm
			}
		}

		stage('compose config check'){
			steps{
				sh 'docker compose config'
			}
		}

		stage('Docker build'){
			steps{
				sh 'docker build -t notes-api .'
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
				curl -f -X POST -d '{"email":"abc@gmail.com", "password":"supersecret"}' http://localhost:9000/register
				curl -f -X POST -d '{"email":"abc@gmail.com", "password":"supersecret"}' http://localhost:9000/login
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


