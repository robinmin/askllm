## askllm

This is a tiny command line tool for you to execute LLM inquiry with prompt or prompt file. Be fore you can use it, please use the following command to copy the config file into your home directory:
```bash
cp config.example.yaml ~/.askllm/config.yaml
```

Please DO add you `api_key` to enable the inquiry. So far this tool only support the following LLM engines: chatgpt, gemini, ollama. The others will be added going forward.

Then you can use the following command to execute the inquiry:
```bash
askllm [-e chatgpt] [-m model] [-c config.yaml] [-p prompt_file.md] [-o output.md] [direct prompt instuctions]
```

For the details of command line options, please run `askllm --help` or `askllm [command] --help`.

Do not forget to install ollama, if you want to use it. Please refer to [here](https://github.com/ollama/ollama) for details.

#### Installation
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
