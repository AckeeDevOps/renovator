module main

replace github.com/AckeeDevOps/renovator/config => ./config

replace github.com/AckeeDevOps/renovator/notifier => ./notifier

replace github.com/AckeeDevOps/renovator/client => ./client

go 1.12

require (
	github.com/AckeeDevOps/renovator/client v0.0.0-00010101000000-000000000000
	github.com/AckeeDevOps/renovator/config v0.0.0-00010101000000-000000000000
	github.com/AckeeDevOps/renovator/notifier v0.0.0-00010101000000-000000000000
	github.com/lusis/slack-test v0.0.0-20190426140909-c40012f20018 // indirect
	github.com/stretchr/testify v1.3.0 // indirect
)
