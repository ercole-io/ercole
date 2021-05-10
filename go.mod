module github.com/ercole-io/ercole/v2

go 1.16

require (
	github.com/1set/gut v0.0.0-20201117175203-a82363231997
	github.com/360EntSecGroup-Skylar/excelize v1.4.1
	github.com/OpenPeeDeeP/xdg v1.0.0
	github.com/amreo/go-dyn-struct v1.2.0
	github.com/amreo/mu v0.0.0-20200710133257-efe27ae7258a
	github.com/aws/aws-sdk-go v1.38.36 // indirect
	github.com/bamzi/jobrunner v1.0.0
	github.com/goji/httpauth v0.0.0-20160601135302-2da839ab0f4d
	github.com/golang/gddo v0.0.0-20210115222349-20d68f94ee1f
	github.com/golang/mock v1.5.0
	github.com/google/go-github/v28 v28.1.1
	github.com/goraz/onion v0.1.3
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/go-version v1.2.1
	github.com/jinzhu/now v1.1.1
	github.com/jtblin/go-ldap-client v0.0.0-20170223121919-b73f66626b33
	github.com/klauspost/compress v1.12.2 // indirect
	github.com/leandro-lugaresi/hub v1.1.1
	github.com/lucasb-eyer/go-colorful v1.2.0
	github.com/pkg/sftp v1.13.0
	github.com/robertkrimen/otto v0.0.0-20200922221731-ef014fd054ac
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/rs/cors v1.7.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.1.3
	github.com/stretchr/testify v1.7.0
	github.com/xakep666/mongo-migrate v0.2.1 // indirect
	github.com/xdg/scram v1.0.3 // indirect
	github.com/xdg/stringprep v1.0.3 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	go.mongodb.org/mongo-driver v1.5.2
	golang.org/x/crypto v0.0.0-20210506145944-38f3c27a63bf
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/text v0.3.6 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/asn1-ber.v1 v1.0.0-20181015200546-f715ec2f112d // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gopkg.in/ldap.v2 v2.5.1 // indirect
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
	gopkg.in/square/go-jose.v2 v2.5.1
)

replace github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.5.0

replace github.com/pkg/sftp => github.com/amreo/sftp v1.10.2-0.20200107133605-5981645e4b3b

// replace github.com/pkg/sftp => ../sftp
// replace github.com/amreo/mu => ../mu
// replace github.com/amreo/go-dyn-struct => ../../amreo/go-dyn-struct
