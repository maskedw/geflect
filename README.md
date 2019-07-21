# geflect

Have you ever wanted to embed Git's HASH value into a C/C + + application?  
This tool supports it.

# Example
Prepare the template source.

```
#include <stdbool.h>

const char GIT_HASH[] = "{{if.Hash}} {{- .Hash -}} {{- else -}} unknown {{- end}}";
const char GIT_BRANCH[] = "{{.Branch}}";
const char GIT_TAG[] = "{{.Tag}}";
const char GIT_DESCRIBE[] = "{{.Describe}}";
const char GIT_SHOART_HASH[] = "{{.ShortHash}}";
const bool GIT_IS_CLEAN = {{if.IsClean}} true {{- else -}} false {{- end}};
const bool GIT_IS_CLEAN_NO_UNTRACED_FILES = {{if.IsCleanNoUnTracedFiles}} true {{- else -}} false {{- end}};
```

Run geflect.


```sh
geflect -o gitmeta.c gitmeta.template
```

Results.

```c
#include <stdbool.h>

const char GIT_HASH[] = "2ab01c92157560bb0a18623e3cf10399a80699d0";
const char GIT_BRANCH[] = "master";
const char GIT_TAG[] = "v1.0.0";
const char GIT_DESCRIBE[] = "v1.0.0-1-g2ab01c9";
const char GIT_SHOART_HASH[] = "2ab01c9";
const bool GIT_IS_CLEAN = false;
const bool GIT_IS_CLEAN_NO_UNTRACED_FILES =  true;
```

List of available template arguments.

|  |  |
| --- | --- |
| Hash | `git rev-parse HEAD` |
| ShortHash | First seven characters of Hash |
| Branch | `git rev-parse --abbrev-ref HEAD` |
| Tag | `git describe --tags --abbrev=0` |
| Describe | `git describe --tags` |
| IsClean | `git status --porcelain` |
| IsCleanNoUnTracedFiles | `git status --porcelain --untraced-files=no` |


Please refer to [/example](https://github.com/maskedw/geflect/tree/master/example) for more details.
