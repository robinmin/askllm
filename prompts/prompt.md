You are a senior system architect especially good at designing system in golang and other relevent technology skills.

I need your help to design a command line tool in golang with LangChain-go to help me inquiry with LLM engins. The full command line should like the following:

```bash
askllm [-e chatgpt] [-m model] [-c config.yaml] [-p prompt_file.md] [-o output.md] [direct prompt instuctions]
```

In above command line, askllm is the name of the tool. There are 5 options for this tool:

| option | default value | optional | remark                                                       |
|:------:|:------ |:------:| ------ |
| -e     | ollama        | Yes | specify the LLM engine, so far we can use chatgpt, gemini and ollama which defined in langchain-go and config file |
| -m     | gemma2        | Yes |  specify the model for current LLM engine |
| -c     | config.yaml | Yes |  specify config file for this tool, all configuration items will be load from this file |
| -p     |  | Yes |  file name of the LLM prompt |
| -o     |  | Yes |  file name to dump the LLM inquiry result, in case of blank, will dump the result into stdout |

If you have some simple prompt instruction you do not want to create a new prompt file, you can direct add them as the extra argument of the command line(a.k.a direct prompt instuctions). In case you already provide the -p option, all direct prompt instuctions will be ignored.

Additional but very important things you need to take care:

- use Makefile as build tool, use goreleaser, errcheck and golangci-lint support in Makefile for a better quality;
- add complete comments and log for better tracking on system behaviors;
- add enough error handling and exception handling logic for a better quality;
- add unit test for each .go file so that it will be maintainable going forward.
