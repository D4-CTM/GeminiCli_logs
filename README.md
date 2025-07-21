# Gemini Logs

> [!NOTE]
> You can find the official releases in the [release page](https://github.com/D4-CTM/GeminiCli_logs/releases/tag/Releases)

This tool is used to extract and view Gemini CLI chat logs in a readable format. It's 
particularly useful for extracting checkpoints of your conversations with the Gemini 
CLI tool, allowing you to see the code that was written or modified.

## Building and Running

### Building the executable

To build the executable, run the following command:

```bash
go build .
```

This will create an executable named `GeminiLogs` in the current directory.

### Running the application

To run the application directly without building, you can use:

```bash
go run . [flags]
```

### Developing for this app

As any go aplication, it's as simple as using the:

```bash
go mod tidy
```

command to install the dependencies and then starting coding in your editor of choice.

## How to Use

This tool works by reading the Gemini CLI's chat logs. For the tool to find these logs, 
you need to have the chat saving feature enabled in the Gemini CLI. The logs are typically 
stored in `$HOME/.gemini/tmp`.

You can then run the tool and specify the directory where you want to save the output logs.

**Example:**

To generate markdown logs from the default location and save them to a `my-gemini-logs` directory, you would run:

```bash
go run . -o my-gemini-logs -md
```

To list all the available logs without generating the output files, use the `-l` flag:

```bash
go run . -l
```

## Flags:

| Flags                         |       | Effect            
|------------------------------ | ----- | ---
| --help                        | -h    | Show the list of possible arguments.
| --output-dir <directory path> | -o    | Creates a directory to place the logs.<br> By default it will create a directory named `logs`.       
| --target-dir <target dir>     | -t    | Especify in which directory are the logs.<br> by default it searches in `$HOME/.gemini/tmp`
| --output-markdown             | -md   | Output the log as a **markdown file**.
| --verbose                     | -vb   | Enabling this mode will show more information while creating the log file.
| --list                        | -l    | List the gemini logs specified in the `--target-dir`.

