module github.com/ruelala/arconn

go 1.22.4

require (
	github.com/aws/aws-sdk-go-v2 v1.30.1
	github.com/aws/aws-sdk-go-v2/config v1.27.24
	github.com/aws/aws-sdk-go-v2/credentials v1.17.24
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.168.0
	github.com/aws/aws-sdk-go-v2/service/ecs v1.44.1
	github.com/aws/aws-sdk-go-v2/service/ssm v1.52.1
	github.com/aws/session-manager-plugin v1.2.633
	github.com/aws/smithy-go v1.20.3
	github.com/buger/jsonparser v1.1.1
	github.com/integrii/flaggy v1.5.2
	github.com/manifoldco/promptui v0.9.0
	golang.org/x/mod v0.19.0
)

require (
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.9 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.13 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.13 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/kms v1.35.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.22.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.26.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.30.1 // indirect
	github.com/chzyer/readline v1.5.1 // indirect
	github.com/eiannone/keyboard v0.0.0-20220611211555-0d226195f203 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/twinj/uuid v1.0.0 // indirect
	github.com/xtaci/smux v1.5.24 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/term v0.22.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/aws/session-manager-plugin => github.com/ruelala/session-manager-plugin v1.6.1
