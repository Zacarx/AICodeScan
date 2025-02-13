package AItools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// 配置结构体
type Config struct {
	API struct {
		URL  string   `yaml:"url"`
		Keys []string `yaml:"keys"`
	} `yaml:"api"`
	Settings struct {
		SleepSeconds int `yaml:"sleep_seconds"`
	} `yaml:"settings"`
	Model struct {
		Name string `yaml:"name"`
	} `yaml:"model"`
	Prompt struct {
		Text string `yaml:"text"`
	} `yaml:"prompt"`
}

// AWA 分析整个文件，而不仅仅是一行代码
func AWA(file string, lineNumber int, lineContent string) (string, error) {
	// 读取配置文件
	config, err := loadConfig("config.yaml")
	if err != nil {
		return "", fmt.Errorf("error loading config: %v", err)
	}

	// 读取整个文件内容
	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}

	// 构建请求的 payload
	payloadMap := map[string]interface{}{
		"model": config.Model.Name,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": fmt.Sprintf(config.Prompt.Text, file, lineNumber, string(fileContent), lineContent),
			},
		},
		"stream":            false,
		"max_tokens":        4096,
		"stop":              []string{"null"},
		"temperature":       0.7,
		"top_p":             0.7,
		"top_k":             50,
		"frequency_penalty": 0.5,
		"n":                 1,
		"response_format":   map[string]string{"type": "text"},
	}

	payloadBytes, err := json.Marshal(payloadMap)
	if err != nil {
		return "", fmt.Errorf("error marshaling payload: %v", err)
	}

	// 打印详细的请求 payload
	// fmt.Printf("Request Payload for %s: %s\n", config.API.URL, string(payloadBytes))

	// 用于跟踪当前使用的 API 密钥索引
	apiIndex := 0
	tries := len(config.API.Keys)

	for i := 0; i < tries; i++ {
		apiKey := config.API.Keys[apiIndex]

		req, err := http.NewRequest("POST", config.API.URL, bytes.NewBuffer(payloadBytes))
		if err != nil {
			apiIndex = (apiIndex + 1) % len(config.API.Keys)
			fmt.Printf("Error creating request for API %s with key %s: %v, trying next key\n", config.API.URL, apiKey, err)
			continue
		}

		req.Header.Add("Authorization", "Bearer "+apiKey)
		req.Header.Add("Content-Type", "application/json")

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			apiIndex = (apiIndex + 1) % len(config.API.Keys)
			fmt.Printf("Error performing request for API %s with key %s: %v, trying next key\n", config.API.URL, apiKey, err)
			continue
		}

		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			apiIndex = (apiIndex + 1) % len(config.API.Keys)
			fmt.Printf("Error reading response body for API %s with key %s: %v, trying next key\n", config.API.URL, apiKey, err)
			continue
		}

		// 打印详细的响应 body
		// fmt.Printf("Response Body for API %s with key %s: %s\n", config.API.URL, apiKey, string(body))

		if res.StatusCode != http.StatusOK {
			apiIndex = (apiIndex + 1) % len(config.API.Keys)
			fmt.Printf("Received non-OK status code %d from API %s with key %s, trying next key\n", res.StatusCode, config.API.URL, apiKey)
			if res.StatusCode == 429 {
				fmt.Println("冷却5s继续")
				time.Sleep(5 * time.Second)
				continue
			}
			if res.StatusCode == http.StatusBadRequest {
				var responseMap map[string]interface{}
				err := json.Unmarshal(body, &responseMap)
				if err == nil {
					if message, ok := responseMap["message"].(string); ok {
						if strings.Contains(message, "less than max_seq_len") {
							fmt.Printf("Skipping request due to prompt length limit: %s\n", message)
							return "", nil // 直接返回 nil，跳过审计
						}
					}
				}
				fmt.Printf("Response Details for 400 from API %s with key %s: %s\n", config.API.URL, apiKey, string(body))
			}
			continue
		}

		var responseMap map[string]interface{}
		err = json.Unmarshal(body, &responseMap)
		if err != nil {
			apiIndex = (apiIndex + 1) % len(config.API.Keys)
			fmt.Printf("Error unmarshaling response for API %s with key %s: %v, trying next key\n", config.API.URL, apiKey, err)
			continue
		}

		if choices, ok := responseMap["choices"].([]interface{}); ok && len(choices) > 0 {
			if choice, ok := choices[0].(map[string]interface{}); ok {
				if message, ok := choice["message"].(map[string]interface{}); ok {
					if content, ok := message["content"].(string); ok {
						time.Sleep(time.Duration(config.Settings.SleepSeconds) * time.Second)
						return content, nil
					}
				}
			}
		}

		apiIndex = (apiIndex + 1) % len(config.API.Keys)
		fmt.Printf("Unexpected response format from API %s with key %s, trying next key\n", config.API.URL, apiKey)
	}

	return "", fmt.Errorf("all API keys failed")
}

// 加载配置文件
func loadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	// 检查密钥列表是否为空
	if len(config.API.Keys) == 0 {
		return nil, fmt.Errorf("no API keys provided")
	}

	return &config, nil
}
