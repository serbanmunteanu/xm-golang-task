package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type WebServerConfig struct {
	JwtConfig          JwtConfig      `yaml:"jwtConfig"`
	MongoConfig        MongoConfig    `yaml:"mongoConfig"`
	PostgresConfig     PostgresConfig `yaml:"postgresConfig"`
	KafkaConfig        KafkaConfig    `yaml:"kafkaConfig"`
	ServerPort         int            `yaml:"serverPort"`
	ConcurrentRequests int            `yaml:"concurrentRequests"`
}

type JwtConfig struct {
	TTL        int    `yaml:"ttl"`
	PrivateKey string `yaml:"privateKey"`
	PublicKey  string `yaml:"publicKey"`
}

type MongoConfig struct {
	URL              string      `yaml:"url"`
	TimeOutInSeconds int         `yaml:"timeOutInSeconds"`
	Database         string      `yaml:"database"`
	Collections      Collections `yaml:"collections"`
}

type PostgresConfig struct {
	URL              string `yaml:"url"`
	TimeOutInSeconds int    `yaml:"timeOutInSeconds"`
	Database         string `yaml:"database"`
}

type KafkaConfig struct {
	Addr          string `yaml:"addr"`
	Topic         string `yaml:"topic"`
	ReaderGroupID string `yaml:"readerGroupId"`
}

type Collections struct {
	UserCollection    string `yaml:"userCollection"`
	CompanyCollection string `yaml:"companyCollection"`
}

func LoadConfig(filePath string) (*WebServerConfig, error) {
	cfg := &WebServerConfig{}
	err := LoadYAML(cfg, filePath)
	if err != nil {
		return nil, err
	}

	return cfg, err
}

func OpenFile(relativePath string) (*os.File, error) {
	path, err := filepath.Abs(relativePath)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	return f, nil
}

func LoadYAML(target interface{}, path string) error {
	f, err := OpenFile(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return yaml.NewDecoder(f).Decode(target)
}
