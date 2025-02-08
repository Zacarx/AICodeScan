package AItools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"time"
)

// 配置结构体
type Config struct {
	API struct {
		URL string `yaml:"url"`
		Key string `yaml:"key"`
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

	// 构建请求的payload
	payloadMap := map[string]interface{}{
		"model": config.Model.Name,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": fmt.Sprintf(config.Prompt.Text, file, lineNumber, string(fileContent), lineContent),
			},
		},
		"stream":            false,
		"max_tokens":        2048,
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

	req, err := http.NewRequest("POST", config.API.URL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("Authorization", config.API.Key)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error performing request: %v", err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	//fmt.Println(string(body)) // 打印原始响应体以进行调试

	var responseMap map[string]interface{}
	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling response: %v", err)
	}

	if choices, ok := responseMap["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					time.Sleep(time.Duration(config.Settings.SleepSeconds) * time.Second)
					return string(content), nil
				}
			}
		}
	}

	return "", fmt.Errorf("unexpected response format, maybe api time out")
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

	return &config, nil
}
