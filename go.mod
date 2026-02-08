module github.com/MeowTux/drift-detector

go 1.21

require (
	github.com/aws/aws-sdk-go-v2 v1.24.1
	github.com/aws/aws-sdk-go-v2/config v1.26.6
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.141.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.48.1
	github.com/aws/aws-sdk-go-v2/service/rds v1.64.0
	github.com/hashicorp/terraform-exec v0.19.0
	github.com/spf13/cobra v1.8.0
	github.com/spf13/viper v1.18.2
	gopkg.in/yaml.v3 v3.0.1
	github.com/slack-go/slack v0.12.3
	github.com/go-mail/mail v2.3.1+incompatible
	github.com/fatih/color v1.16.0
	github.com/schollz/progressbar/v3 v3.14.1
)

require (
	cloud.google.com/go/compute v1.23.3
	cloud.google.com/go/storage v1.36.0
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.9.1
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.5.1
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5 v5.4.0
	github.com/sirupsen/logrus v1.9.3
	github.com/joho/godotenv v1.5.1
	golang.org/x/crypto v0.18.0
	github.com/google/uuid v1.5.0
)
