#!/usr/bin/env groovy
registry = "eu.gcr.io/kyma-project/incubator"
isRelease = env.TAG_NAME
appVersion = isRelease?env.TAG_NAME:env.BRANCH_NAME

properties([
    buildDiscarder(logRotator(numToKeepStr: '30')),
])

pipeline {
    agent {
        node {
            label 'docker'
        }
    }

    stages{
        stage("setup")
        {
            steps {
                script {
                    checkout scm
                    withCredentials([usernamePassword(credentialsId: 'gcr-kyma-rw', passwordVariable: 'pwd', usernameVariable: 'uname')]) {
                        sh "docker login -u $uname -p '$pwd' $registry"
                    }
                }
            }
        }
        stage("build bundles")
        {
            steps{
                script {
                    parallel(
                        "mqtt-event-bridge": {
                            execute("mqtt-event-bridge")
                        },
                        "api-registration-job":{
                            execute("api-registration-job")
                        },
                        "oauth2server":{
                            execute("oauth2server")
                        }
                    )
                }
            }
        }
    }
}

def execute(application){
    stage ("build image $application"){
        sh "cd $application && docker build -t $application ."
    }

    stage ("push image $application"){
        sh "docker tag $application:latest $registry/$application:$appVersion"
        sh "docker push $registry/$application:$appVersion"
        if(isRelease){
            sh "docker tag $application:latest $registry/$application:latest"
            sh "docker push $registry/$application:latest"
        }
    }
}