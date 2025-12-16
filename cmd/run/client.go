package run

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aide-family/magicbox/load"
	"github.com/aide-family/magicbox/strutil"
	"github.com/aide-family/magicbox/strutil/cnst"
	"github.com/go-kratos/kratos/v2/encoding"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/internal/server"
	"github.com/aide-family/rabbit/pkg/config"
)

// GenerateClientConfig 生成 ClientConfig 并写入 指定路径 默认.rabbit/config.yaml
func GenerateClientConfig(
	bc *conf.Bootstrap,
	srvs server.Servers,
	helper *klog.Helper,
) error {
	if bc.GetEnableClientConfig() != "true" {
		helper.Debugw("msg", "client config is not enabled")
		return nil
	}
	var clusterEndpoint string
	for _, srv := range srvs {
		if grpcSrv, ok := srv.(*grpc.Server); ok {
			endpoint, err := grpcSrv.Endpoint()
			if err == nil {
				clusterEndpoint = endpoint.String()
				break
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

	if clusterConfig != nil && strutil.IsNotEmpty(clusterConfig.GetEndpoints()) {
		clientConfig.Cluster = clusterConfig
	} else {
		// 如果没有配置，使用默认值
		var endpoint string
		if registryType == config.RegistryType_ETCD || registryType == config.RegistryType_KUBERNETES {
			// 如果使用 etcd 或 k8s 注册，使用 discovery:///服务名称 格式
			endpoint = strings.Join([]string{"discovery://", serverName}, "/")
		} else {
			endpoint = normalizeAddress(clusterEndpoint)
		}
		clientConfig.Cluster = &config.ClusterConfig{
			Name:      serverName,
			Endpoints: endpoint,
			Timeout:   durationpb.New(10 * time.Second),
		}
	}

	clientConfig.JwtToken = strings.Join([]string{cnst.HTTPHeaderBearerPrefix, "<jwt-token>"}, " ")

	// 写入配置文件到当前目录的 .rabbit/client_config.yaml
	configDir := filepath.Dir(load.ExpandHomeDir(runFlags.RabbitConfigPath))
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return fmt.Errorf("create config directory failed: %w", err)
	}

	configFile := filepath.Join(configDir, "client_config.yaml")

	// 转换为 YAML 格式
	yamlData, err := clientConfigToYAML(clientConfig)
	if err != nil {
		return fmt.Errorf("convert config to yaml failed: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(configFile, yamlData, 0o644); err != nil {
		return fmt.Errorf("write config file failed: %w", err)
	}

	helper.Debugw("msg", "config generated successfully", "path", configFile)
	return nil
}

// normalizeAddress 将地址中的 0.0.0.0 转换为 localhost，并移除协议前缀
func normalizeAddress(addr string) string {
	if addr == "" {
		return addr
	}
	// 移除协议前缀：去除 // 前面的所有字符
	if idx := strings.Index(addr, "//"); idx != -1 {
		addr = addr[idx+2:]
	}
	// 移除末尾的 / 字符
	addr = strings.TrimSpace(strings.TrimSuffix(addr, "/"))
	// 将 0.0.0.0 转换为 localhost
	if strings.HasPrefix(addr, "0.0.0.0") {
		addr = strings.Replace(addr, "0.0.0.0", "localhost", 1)
	}
	return addr
}

// clientConfigToYAML 将 ClientConfig 转换为 YAML 格式
func clientConfigToYAML(cfg *config.ClientConfig) ([]byte, error) {
	return encoding.GetCodec("yaml").Marshal(cfg)
}
