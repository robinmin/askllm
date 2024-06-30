I have the following generic function to load information from YAML file or dump into YAML file. 
```golang
// LoadConfig 从指定的YAML文件中加载配置信息
func LoadConfig[T any](yamlFile string) (*T, error) {
	data, err := os.ReadFile(yamlFile)
	if err != nil {
		return nil, err
	}

	var config T
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig 将配置信息保存到指定的YAML文件中
func SaveConfig[T any](cfg *T, yamlFile string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(yamlFile, data, 0o644)
	if err != nil {
		return err
	}

	return nil
}
```

Meanwhile, I have a YAML file with content as shown bellow. Please help to generic a relevant struct definition. So that I can call LoadConfig and SaveConfig with this new type for further action. The type of struct must be a single and compound struct. Meanwhile, DO NOT forget to add inline comment for each field. Here's the yaml content:
```YAML
id: firt_ai_template
name: "Personality Analyzer"
description: "Analyzes a person's personality based on their name and mood"
author: "John Doe"
variables:
  - name: "user_name"
    vtype: "string"
    otype: "text"
    default: "User"
    validation: "^[a-zA-Z ]{2,30}$"
  - name: "mood"
    vtype: "string"
    otype: "select"
    options: ["happy", "sad", "excited"]
    default: "happy"
    validation: "^(happy|sad|excited)$"
template: |
  Analyze the personality of a person with the following characteristics:
  Name: {{ user_name }}
  Current Mood: {{ mood }}

  Please provide a brief personality analysis based on these factors.
```
