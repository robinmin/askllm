# ASKLLM

This is a tiny command line tool for you to execute LLM inquiry with prompt or prompt file. The goal is to provide a convinent way to share and reuse prompts, and provide a way to observe the different results from different LLM engines and models with the same prompt.

## Features

- [x] Command AI tool.
- [x] LangChain Support, so far suppor chatgpt, gemini, ollama and claude.
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

- chatgpt
- gemini
- ollama
- claude.

The others will be added soon.

Once everything is ready, then you can use the following command to ask whatever you want to know:

```bash
askllm [-e chatgpt] [-m model] [-c config.yaml] [-p prompt_file.md] [-o output.md] [direct prompt instuctions]

# use the default model (gemma2) to ask local ollama
askllm "hello, llm"

# use model gpt-3.5-turbo to ask openai chatgpt
askllm -e chatgpt -m gpt-3.5-turbo "hello, llm"

# use model gemini-1.5-flash to ask google gemini
askllm -e gemini -m gemini-1.5-flash "hello, llm"

# use model claude-3-sonnet-20240229 to ask anthropic claude
askllm -e claude -m claude-3-sonnet-20240229 "hello, llm"

```

For the details of command line options, please run `askllm --help` or `askllm [command] --help`.
