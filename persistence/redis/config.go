package redis

type Config struct {
	Host      string
	Port      uint16
	IsCluster bool
	Namespace string
}
