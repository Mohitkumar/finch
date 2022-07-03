package config

type StorageImplementation string

const STORAGE_IMPL_REDIS StorageImplementation = "redis"
const STORAGE_IMPL_INMEM StorageImplementation = "memory"

type Config struct {
	RedisConfig    RedisStorageConfig
	InMemoryConfig InmemoryStorageConfig
	HttpPort       int
	GrpcPort       int
	StorageImpl    StorageImplementation
}

type RedisStorageConfig struct {
	Host      string
	Port      int
	Namespace string
}

type InmemoryStorageConfig struct {
}
