package run

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/encoding"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	jwtv5 "github.com/golang-jwt/jwt/v5"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/aide-family/magicbox/load"
	"github.com/aide-family/magicbox/strutil/cnst"
	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/internal/server"
	"github.com/aide-family/rabbit/pkg/config"
	"github.com/aide-family/rabbit/pkg/middler"
)

const (
	defaultClientUsername = "root"
	defaultClientUserUID  = 1
)

// generateClientConfig 生成 ClientConfig 并写入 .rabbit/config.yaml
func generateClientConfig(
	bc *conf.Bootstrap,
	srvs server.Servers,
	helper *klog.Helper,
) error {
	if !bc.GetEnableClientConfig() {
		return nil
	}
	// 获取服务器 endpoint
	var httpEndpoint, grpcEndpoint string
	for _, srv := range srvs {
		switch s := srv.(type) {
		case *http.Server:
			endpoint, err := s.Endpoint()
			if err == nil {
				httpEndpoint = endpoint.String()
			}
		case *grpc.Server:
			endpoint, err := s.Endpoint()
			if err == nil {
				grpcEndpoint = endpoint.String()
			}
		}
	}

	// 构建 ClientConfig
	clientConfig := &config.ClientConfig{}

	// 1. 添加集群连接信息
	clusterConfig := bc.GetCluster()
	registryType := bc.GetRegistryType()
	serverName := bc.GetServer().GetName()

	clientConfig.RegistryType = registryType
	switch registryType {
	case config.RegistryType_ETCD:
		clientConfig.Etcd = bc.GetEtcd()
	case config.RegistryType_KUBERNETES:
		clientConfig.Kubernetes = bc.GetKubernetes()
	}

	if clusterConfig != nil {
		clientConfig.Cluster = clusterConfig
	} else {
		// 如果没有配置，使用默认值
		var endpoint string
		var protocol config.ClusterConfig_Protocol
		if registryType == config.RegistryType_ETCD || registryType == config.RegistryType_KUBERNETES {
			// 如果使用 etcd 或 k8s 注册，使用 discovery:///服务名称 格式
			endpoint = strings.Join([]string{"discovery://", serverName}, "/")
			protocol = config.ClusterConfig_GRPC
		} else {
			endpoint = normalizeAddress(grpcEndpoint)
			if endpoint == "" {
				endpoint = normalizeAddress(httpEndpoint)
				protocol = config.ClusterConfig_HTTP
			}
		}
		clientConfig.Cluster = &config.ClusterConfig{
			Name:      serverName,
			Endpoints: endpoint,
			Protocol:  protocol,
			Timeout:   durationpb.New(10 * time.Second),
		}
	}

	// 2. 生成 root jwtToken
	jwtConf := bc.GetJwt()
	if jwtConf != nil {
		rootClaims := middler.NewJwtClaims(jwtConf, middler.BaseInfo{
			UserID:   defaultClientUserUID,
			Username: defaultClientUsername,
		})
		// 设置较长的过期时间（1年）
		rootClaims.ExpiresAt = jwtv5.NewNumericDate(time.Now().Add(365 * 24 * time.Hour))
		if token, err := rootClaims.GenerateToken(); err == nil {
			clientConfig.JwtToken = strings.Join([]string{cnst.HTTPHeaderBearerPrefix, token}, " ")
		}
	}

	// 写入配置文件到当前目录的 .rabbit/config.yaml
	configDir := filepath.Dir(load.ExpandHomeDir(flags.clientConfigOutputPath))
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return fmt.Errorf("create config directory failed: %w", err)
	}

	configFile := filepath.Join(configDir, "config.yaml")

	// 转换为 YAML 格式
	yamlData, err := clientConfigToYAML(clientConfig)
	if err != nil {
		return fmt.Errorf("convert config to yaml failed: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(configFile, yamlData, 0o644); err != nil {
		return fmt.Errorf("write config file failed: %w", err)
	}

	helper.Infow("msg", "client config generated", "path", configFile)
	return nil
}

// normalizeAddress 将地址中的 0.0.0.0 转换为 localhost
func normalizeAddress(addr string) string {
	if strings.HasPrefix(addr, "0.0.0.0:") {
		return strings.Replace(addr, "0.0.0.0:", "localhost:", 1)
	}
	return addr
}

// clientConfigToYAML 将 ClientConfig 转换为 YAML 格式
func clientConfigToYAML(cfg *config.ClientConfig) ([]byte, error) {
	return encoding.GetCodec("yaml").Marshal(cfg)
}
