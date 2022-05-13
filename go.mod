module minecraft-control-panel

go 1.17


replace	github.com/gin-gonic/gin v1.6.3 => ./pkg/src/github.com/gin-gonic/gin@v1.6.3
replace	golang.org/x/oauth2 v0.0.0-20201109201403-9fd604954f58 => ./pkg/src/golang.org/x/oauth2@v0.0.0-20201109201403-9fd604954f58 

replace	cloud.google.com/go v0.65.0 => ./pkg/src/cloud.google.com/go@v0.65.0 
replace	github.com/gin-contrib/sse v0.1.0 => ./pkg/src/github.com/gin-contrib/sse@v0.1.0 
replace	github.com/go-playground/locales v0.13.0 => ./pkg/src/github.com/go-playground/locales@v0.13.0 
replace	github.com/go-playground/universal-translator v0.17.0 => ./pkg/src/github.com/go-playground/universal-translator@v0.17.0 
replace	github.com/go-playground/validator/v10 v10.2.0 => ./pkg/src/github.com/go-playground/validator/v10@v10.2.0
replace	github.com/golang/protobuf v1.4.2 => ./pkg/src/github.com/golang/protobuf@v1.4.2
replace	github.com/google/go-cmp v0.5.2 => ./pkg/src/github.com/google/go-cmp@v0.5.2
replace	github.com/json-iterator/go v1.1.9 => ./pkg/src/github.com/json-iterator/go@v1.1.9
replace	github.com/leodido/go-urn v1.2.0 => ./pkg/src/github.com/leodido/go-urn@v1.2.0 
replace	github.com/mattn/go-isatty v0.0.12 => ./pkg/src/github.com/mattn/go-isatty@v0.0.12 
replace	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 => ./pkg/src/ithub.com/modern-go/concurrent@v0.0.0-20180228061459-e0a39a4cb421
replace	github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742 => ./pkg/src/github.com/modern-go/reflect2@v0.0.0-20180701023420-4b7aa43c6742 
replace	github.com/ugorji/go/codec v1.1.7 => ./pkg/src/github.com/ugorji/go/codec@v1.1.7 
replace	golang.org/x/net v0.0.0-20200822124328-c89045814202 => ./pkg/src/golang.org/x/net@v0.0.0-20200822124328-c89045814202
replace	golang.org/x/sys v0.0.0-20200905004654-be1d3432aa8f => ./pkg/src/golang.org/x/sys@v0.0.0-20200905004654-be1d3432aa8f
replace	google.golang.org/appengine v1.6.6 => ./pkg/src/google.golang.org/appengine@v1.6.6
replace	google.golang.org/protobuf v1.25.0 => ./pkg/src/google.golang.org/protobuf@v1.25.0 
replace	gopkg.in/yaml.v2 v2.2.8 => ./pkg/src/gopkg.in/yaml.v2@v2.2.8

require (
	github.com/gin-gonic/gin v1.6.3
	golang.org/x/oauth2 v0.0.0-20201109201403-9fd604954f58
)

require (
	cloud.google.com/go v0.65.0 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator/v10 v10.2.0 // indirect
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/google/go-cmp v0.5.2 // indirect
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742 // indirect
	github.com/ugorji/go/codec v1.1.7 // indirect
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
	golang.org/x/sys v0.0.0-20200905004654-be1d3432aa8f // indirect
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)