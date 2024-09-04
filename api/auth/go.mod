module auth

go 1.22.0

replace utils => ../../utils

replace services => ../../services

require (
	github.com/aws/aws-lambda-go v1.47.0
	utils v0.0.0-00010101000000-000000000000
)
