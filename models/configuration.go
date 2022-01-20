package models

// Configuration cache server configuration
type Configuration struct {
	RedisServer DataStorageInfo `yaml:"redisServer" json:"redisServer"`
}

// DataStorageInfo data storage connection info
type DataStorageInfo struct {
	IP         string `yaml:"ip" json:"ip"`
	Port       int    `yaml:"port" json:"port"`
	Database   string `yaml:"database" json:"database"`
	ID         string `yaml:"id" json:"id"`
	PW         string `yaml:"pw" json:"pw"`
	MaxIdleCon int    `yaml:"maxIdleCon" json:"maxIdleCon"`
	MaxOpenCon int    `yaml:"MaxOpenCon" json:"MaxOpenCon"`
}
