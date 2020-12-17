module github.com/ercole-io/ercole/v2

go 1.15

require (
	github.com/1set/gut v0.0.0-20200225162230-3995492b8589
	github.com/360EntSecGroup-Skylar/excelize v1.4.1
	github.com/OpenPeeDeeP/xdg v0.2.0
	github.com/amreo/go-dyn-struct v1.2.0
	github.com/amreo/mu v0.0.0-20200710133257-efe27ae7258a
	github.com/bamzi/jobrunner v0.0.0-20190810144113-852b7ca4d475
	github.com/goji/httpauth v0.0.0-20160601135302-2da839ab0f4d
	github.com/golang/gddo v0.0.0-20200219175727-df439dd5819e
	github.com/golang/mock v1.4.4
	github.com/google/go-cmp v0.4.0 // indirect
	github.com/google/go-github/v28 v28.1.1
	github.com/goraz/onion v0.1.3
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.3
	github.com/hashicorp/go-version v1.2.0
	github.com/jinzhu/now v1.1.1
	github.com/jtblin/go-ldap-client v0.0.0-20170223121919-b73f66626b33
	github.com/leandro-lugaresi/hub v1.1.0
	github.com/lucasb-eyer/go-colorful v1.0.3
	github.com/pkg/sftp v1.10.1
	github.com/robertkrimen/otto v0.0.0-20191219234010-c382bd3c16ff
	github.com/rs/cors v1.7.0
	github.com/sirupsen/logrus v1.5.0
	github.com/spf13/cobra v0.0.6
	github.com/stretchr/testify v1.5.1
	github.com/xdg/stringprep v1.0.0 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0
	go.mongodb.org/mongo-driver v1.3.3
	golang.org/x/crypto v0.0.0-20191112222119-e1110fd1c708
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/asn1-ber.v1 v1.0.0-20181015200546-f715ec2f112d // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gopkg.in/ldap.v2 v2.5.1 // indirect
	gopkg.in/robfig/cron.v3 v3.0.0-00010101000000-000000000000 // indirect
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
	gopkg.in/square/go-jose.v2 v2.4.1
)

replace github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.5.0

replace gopkg.in/robfig/cron.v3 => github.com/robfig/cron/v3 v3.0.0

replace github.com/pkg/sftp => github.com/amreo/sftp v1.10.2-0.20200107133605-5981645e4b3b

// replace github.com/pkg/sftp => ../sftp
// replace github.com/amreo/mu => ../mu
// replace github.com/amreo/go-dyn-struct => ../../amreo/go-dyn-struct
