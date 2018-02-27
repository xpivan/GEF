#!/bin/bash

go run genEgiResource.go
go run scpTLSCert.go
go run dockerRemoteServer.go