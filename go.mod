module go-docker

go 1.14

require (
	github.com/gorilla/mux v1.7.4
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	k8s.io/apimachinery v0.18.3
	k8s.io/client-go v11.0.0+incompatible
)

replace (
	k8s.io/apimachinery v0.18.3 => k8s.io/apimachinery v0.17.4
	k8s.io/client-go v11.0.0+incompatible => k8s.io/client-go v0.17.4
)
