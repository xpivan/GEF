#!/bin/sh

java -jar /home/VOMS-proxy/target/voms-proxy-1.0-jar-with-dependencies.jar &&
java -jar /home/workflowEGI/target/jocci-create-resource-1.0-jar-with-dependencies.jar
