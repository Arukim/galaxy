#!/bin/sh
go build .
rsync galaxy ubuntu@ec2-52-57-3-148.eu-central-1.compute.amazonaws.com:/var/www/galaxy/
rsync -r templates/ ubuntu@ec2-52-57-3-148.eu-central-1.compute.amazonaws.com:/var/www/galaxy/templates/
rsync -r static/ ubuntu@ec2-52-57-3-148.eu-central-1.compute.amazonaws.com:/var/www/galaxy/static/
ssh ubuntu@ec2-52-57-3-148.eu-central-1.compute.amazonaws.com supervisorctl restart galaxy
