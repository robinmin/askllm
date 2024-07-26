# ASKLLM

This is a tiny command line tool for you to execute LLM inquiry with prompt or prompt file. The goal is to provide a convinent way to share and reuse prompts, and provide a way to observe the different results from different LLM engines or models with the same prompt.

## Features

- [x] Command AI tool.
- [x] LangChain Support, so far suppor chatgpt, gemini, ollama, groq and claude.
- [x] Prompt Template Support with Golang text/template syntax(yaml file only), or plaint text as old version.

## Installation

On macOS/Linux you can install it via [HomeBrew](https://brew.sh/) as shown below:

```bash
# add tap, only for the first time
brew tap robinmin/tap

# install
brew install askllm

# upgrade new version
brew upgrade askllm
```

You also can install it via Scoop on Windows as shown below:

```bash
scoop bucket add robinmin https://github.com/robinmin/scoop-bucket.git
scoop install robinmin/askllm
```

## Useage

Be fore you can use it, please use the following command to copy the config file into your home directory:

```bash
mkdir ~/.askllm/ && cp config.example.yaml ~/.askllm/config.yaml
```

Please DO add you `api_key` to enable the inquiries. If you want to use ollama (by default), please do not forget to install it. Please refer to [here](https://github.com/ollama/ollama) for details. So far askllm only supports the following LLM engines:

- [chatgpt](https://chatgpt.com/)
- [gemini](https://gemini.google.com/)
- [ollama](https://ollama.com/)
- [claude](https://claude.ai/)
- [groq](https://groq.com/)

The others will be added soon.

Once everything is ready, then you can use the following command to ask whatever you want to know:

```bash
askllm [-a action] [-e chatgpt] [-m model] [-c config.yaml] [-p prompt_file.md] [-o output.md] [direct prompt instuctions]

# use the default model (gemma2) to ask local ollama
askllm "hello, llm"

# use model gpt-3.5-turbo to ask openai chatgpt
askllm -e chatgpt -m gpt-3.5-turbo "hello, llm"

# use model gemini-1.5-flash to ask google gemini
askllm -e gemini -m gemini-1.5-flash "hello, llm"

# use model claude-3-sonnet-20240229 to ask anthropic claude
askllm -e claude -m claude-3-sonnet-20240229 "hello, llm"

```

For the details of command line options, please run `askllm --help`.

## Prompt template file

Askllm defined a file layout for the relevant prompt information in YAML format. It composed with three parts: metadata section, variable section and prompt template section. Once you defined variables in the variable section, then you can use them in the template section in golang text template syntax. It will give you the capability to design the reuseable prompt. Here comes a sample.

```yaml
id: prompt_yaml_golang_struct
name: "Prompt To Generate Golang Struct Definition for YAML"
description: "A tiny tool to generate golang struct definition based on YAML file content."
author: "Robin Min"
default_engine: "chatgpt"
default_model: "gpt-4o"
variables:
  - name: "yaml_file"
    vtype: "file"
    otype: "text"
    default: ""
    validation: ""
template: |
  I have the following utility generic function to load information from YAML file and dump the content into YAML file. 

  // LoadConfig: load information from YAML file
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
  
  // SaveConfig: dump the content into YAML file
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
  
  Now, I have a YAML file with following content. You need to generate a single and compound golang struct definition. So that I can call LoadConfig[T]() and SaveConfig[T]() with this new type for further action.
  DO NOT forget to add inline comment for each field. Here's the yaml content:

  {{ .yaml_file }}
```

## Reference

- [5 simple tips and tricks for writing unit tests in #golang](https://medium.com/@matryer/5-simple-tips-and-tricks-for-writing-unit-tests-in-golang-619653f90742)
- [Meet Moq: Easily mock interfaces in Go](https://medium.com/@matryer/meet-moq-easily-mock-interfaces-in-go-476444187d10)
