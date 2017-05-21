pipeline {
    agent { docker 'golang:1.8' }
    stages {
        stage('build') {
            steps {
                sh build.sh
            }
        }
    }
}