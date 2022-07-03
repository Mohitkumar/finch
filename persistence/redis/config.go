package redis

type Config struct {
	Host      string
	Port      int
	IsCluster bool
	Namespace string
}
