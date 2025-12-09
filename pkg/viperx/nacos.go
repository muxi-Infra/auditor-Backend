package viperx

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
)

func NewVipperSettingFromNacos(path string) *VipperSetting {
	content, err := getConfigFromNacos()

	if err != nil {
		log.Printf("从 Nacos 拉取配置失败: %v, 尝试读取本地文件", err)
		fileContent, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("无法读取本地配置文件 %s: %v", path, err)
		}
		content = string(fileContent)
	}

	vp := viper.New()
	vp.SetConfigType("yaml")
	if err := vp.ReadConfig(bytes.NewBufferString(content)); err != nil {
		log.Fatalf("解析配置失败: %v", err)
	}
	return &VipperSetting{Viper: vp}
}

func getConfigFromNacos() (string, error) {
	server, port, namespace, user, pass, group, dataId := parseNacosDSN()
	fmt.Println(server, port, namespace, user, pass, group, dataId)

	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: server,
			Port:   port,
			Scheme: "http",
		},
	}

	clientConfig := constant.ClientConfig{
		NamespaceId:         namespace,
		Username:            user,
		Password:            pass,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		CacheDir:            "./data/configCache",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		return "", fmt.Errorf("初始化 Nacos 客户端失败: %w", err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	if err != nil {
		return "", fmt.Errorf("拉取 Nacos 配置失败: %w", err)
	}

	return content, nil
}

func parseNacosDSN() (server string, port uint64, ns, user, pass, group, dataId string) {
	dsn := os.Getenv("NACOS_DSN")
	if dsn == "" {
		log.Fatal("环境变量 NACOS_DSN 未设置")
	}

	parts := strings.SplitN(dsn, "?", 2)
	host := parts[0]
	params := url.Values{}

	if len(parts) == 2 {
		params, _ = url.ParseQuery(parts[1])
	}

	hostParts := strings.Split(host, ":")
	server = hostParts[0]
	if len(hostParts) > 1 {
		p, _ := strconv.Atoi(hostParts[1])
		port = uint64(p)
	} else {
		port = 8848
	}

	ns = params.Get("namespace")
	if ns == "" {
		ns = "public"
	}

	user = params.Get("username")
	pass = params.Get("password")
	group = params.Get("group")
	dataId = params.Get("dataId")
	return
}
