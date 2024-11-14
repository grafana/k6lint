Static analyzer for k6 extensions

k6lint analyzes the source of the k6 extension and try to build k6 with the extension.

By default, text output is generated. The `--json` flag can be used to generate the result in JSON format.

If the grade is `C` or higher, the command is successful, otherwise it returns an exit code larger than 0.
This passing grade can be modified using the `--passing` flag.
