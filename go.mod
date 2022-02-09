module github.com/ercole-io/ercole/v2

go 1.17

require (
	github.com/1set/gut v0.0.0-20201117175203-a82363231997
	github.com/360EntSecGroup-Skylar/excelize v1.4.1
	github.com/OpenPeeDeeP/xdg v1.0.0
	github.com/amreo/go-dyn-struct v1.2.0
	github.com/amreo/mu v0.0.0-20200710133257-efe27ae7258a
	github.com/bamzi/jobrunner v1.0.0
	github.com/fatih/color v1.13.0
	github.com/goji/httpauth v0.0.0-20160601135302-2da839ab0f4d
	github.com/golang-jwt/jwt/v4 v4.2.0
	github.com/golang/gddo v0.0.0-20210115222349-20d68f94ee1f
	github.com/golang/mock v1.6.0
	github.com/google/go-github/v28 v28.1.1
	github.com/goraz/onion v0.1.3
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/go-version v1.4.0
	github.com/jinzhu/now v1.1.4
	github.com/jtblin/go-ldap-client v0.0.0-20170223121919-b73f66626b33
	github.com/leandro-lugaresi/hub v1.1.1
	github.com/lucasb-eyer/go-colorful v1.2.0
	github.com/microcosm-cc/bluemonday v1.0.17
	github.com/oracle/oci-go-sdk v24.3.0+incompatible
	github.com/oracle/oci-go-sdk/v45 v45.2.0
	github.com/pmezard/go-difflib v1.0.0
	github.com/robertkrimen/otto v0.0.0-20200922221731-ef014fd054ac
	github.com/rs/cors v1.8.2
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.3.0
	github.com/stretchr/testify v1.7.0
	github.com/xakep666/mongo-migrate v0.2.1
	github.com/xeipuuv/gojsonschema v1.2.0
	go.mongodb.org/mongo-driver v1.8.3
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
)

require (
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/felixge/httpsnoop v1.0.1 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/golang/snappy v0.0.3 // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/gorilla/css v1.0.0 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mitchellh/mapstructure v1.4.3 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/robfig/cron/v3 v3.0.0 // indirect
	github.com/skarademir/naturalsort v0.0.0-20150715044055-69a5d87bef62 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.0.2 // indirect
	github.com/xdg-go/stringprep v1.0.2 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20180127040702-4e3ac2762d5f // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	golang.org/x/net v0.0.0-20210813160813-60bc85c4be6d // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20211205182925-97ca703d548d // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/asn1-ber.v1 v1.0.0-20181015200546-f715ec2f112d // indirect
	gopkg.in/ldap.v2 v2.5.1 // indirect
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.5.0

// replace github.com/amreo/mu => ../mu
// replace github.com/amreo/go-dyn-struct => ../../amreo/go-dyn-struct
