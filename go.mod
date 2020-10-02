module github.com/drewbailey/nomad-deploy-notifier

go 1.15

replace github.com/hashicorp/nomad/api => /home/drew/work/go/nomad/api

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/hashicorp/go-hclog v0.14.1
	github.com/hashicorp/nomad/api v0.0.0-20201001180849-8238b9f86415
	github.com/mitchellh/mapstructure v1.3.3
	github.com/slack-go/slack v0.6.6
	github.com/stretchr/testify v1.5.1
)
