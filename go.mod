module github.com/amreo/ercole-services

go 1.13

require (
	github.com/bamzi/jobrunner v0.0.0-20190810144113-852b7ca4d475
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/goji/httpauth v0.0.0-20160601135302-2da839ab0f4d
	github.com/golang/mock v1.3.1
	github.com/golang/snappy v0.0.1 // indirect
	github.com/google/go-cmp v0.3.0 // indirect
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.3
	github.com/leandro-lugaresi/hub v1.1.0
	github.com/robfig/cron v1.2.0 // indirect
	github.com/rs/cors v1.7.0
	github.com/stretchr/testify v1.4.0
	github.com/tidwall/pretty v1.0.0 // indirect
	github.com/xdg/scram v0.0.0-20180814205039-7eeb5667e42c // indirect
	github.com/xdg/stringprep v1.0.0 // indirect
	go.mongodb.org/mongo-driver v1.1.1
	golang.org/x/crypto v0.0.0-20190911031432-227b76d455e7 // indirect
	golang.org/x/text v0.3.2 // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22 // indirect
	gopkg.in/robfig/cron.v3 v3.0.0-00010101000000-000000000000 // indirect
)

replace github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.4.2

replace gopkg.in/robfig/cron.v3 => github.com/robfig/cron/v3 v3.0.0
