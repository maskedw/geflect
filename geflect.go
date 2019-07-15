package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
)

type Options struct {
	Version           func()         `short:"v" long:"version" description:"print version and exit"`
	GitRepositoryPath flags.Filename `short:"g" long:"git-repo" description:"git repository path(default: automatically looks in the CWD)"`
	OutputFilePath    flags.Filename `short:"o" long:"out" description:"output file path(default: stdout)"`
	ForceOverwrite    bool           `short:"f" long:"force" description:"overwrite the destination with the same content(default: false)"`
	IgnoreGitErrors   bool           `long:"ignore-git-errors" description:"ignore git errors(default: false)"`
	Args              struct {
		TemplateFilename flags.Filename
	} `positional-args:"yes" required:"true"`
}

type GitInfo struct {
	Hash                   string
	ShortHash              string
	Branch                 string
	Tag                    string
	Describe               string
	IsClean                bool
	IsCleanNoUnTracedFiles bool
}

func trimNewLine(line string) string {
	return strings.TrimRight(line, "\r\n")
}

type GitParser struct {
	WorkingDirectory string
}

func (g *GitParser) HasTag() (bool, error) {
	output, err := g.call("tag")
	if err != nil {
		return false, err
	}
	outString := trimNewLine(string(output))

	return len(outString) != 0, nil
}

func (g *GitParser) Tag() (string, error) {
	hasTag, err := g.HasTag()
	if err != nil {
		return "", err
	}
	if !hasTag {
		return "", nil
	}

	output, err := g.call("describe", "--tags", "--abbrev=0")
	if err != nil {
		return "", err
	}

	return trimNewLine(output), nil
}

func (g *GitParser) Hash() (string, error) {
	output, err := g.call("rev-parse", "HEAD")
	if err != nil {
		return "", err
	}
	return trimNewLine(output), nil
}

func (g *GitParser) Branch() (string, error) {
	// https://stackoverflow.com/a/12142066/3270390
	output, err := g.call("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}
	return trimNewLine(output), nil
}

func (g *GitParser) Describe() (string, error) {
	hasTag, err := g.HasTag()
	if err != nil {
		return "", err
	}
	if !hasTag {
		return "", nil
	}

	// https://stackoverflow.com/a/12142066/3270390
	output, err := g.call("describe", "--tags")
	if err != nil {
		return "", err
	}
	return trimNewLine(output), nil
}

func (g *GitParser) IsClean() (bool, error) {
	output, err := g.call("status", "--porcelain")
	if err != nil {
		return false, err
	}

	return len(output) == 0, nil
}

func (g *GitParser) IsCleanNoUnTracedFiles() (bool, error) {
	output, err := g.call("status", "--porcelain", "--untracked-files=no")
	if err != nil {
		return false, err
	}

	return len(output) == 0, nil
}

func (g *GitParser) call(arg1 string, args ...string) (string, error) {
	cmdArgs := []string{"git", "-C", g.WorkingDirectory, arg1}
	cmdArgs = append(cmdArgs, args...)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmdLine := strings.Join(cmdArgs, " ")

	output, err := cmd.Output()
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("command-line: %s", cmdLine))
		return "", err
	}

	return string(output), err
}

func getGitInfo(workingDirectory string) (GitInfo, error) {
	var parser GitParser
	if workingDirectory == "" {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		parser.WorkingDirectory = cwd
	} else {
		parser.WorkingDirectory = workingDirectory
	}

	var gitInfo GitInfo
	var err error

	gitInfo.Hash, err = parser.Hash()
	if err != nil {
		return gitInfo, err
	}
	gitInfo.ShortHash = gitInfo.Hash[:7]

	gitInfo.Branch, err = parser.Branch()
	if err != nil {
		return gitInfo, err
	}

	gitInfo.Tag, err = parser.Tag()
	if err != nil {
		return gitInfo, err
	}

	gitInfo.Describe, err = parser.Describe()
	if err != nil {
		return gitInfo, err
	}

	gitInfo.IsClean, err = parser.IsClean()
	if err != nil {
		return gitInfo, err
	}

	gitInfo.IsCleanNoUnTracedFiles, err = parser.IsCleanNoUnTracedFiles()
	if err != nil {
		return gitInfo, err
	}

	return gitInfo, nil
}

type IsRepositoryCleanOptions struct {
	NoUntrackedFiles bool
}

func main() {
	log.SetFlags(log.Lshortfile)
	var opts Options
	opts.Version = func() {
		fmt.Fprintf(os.Stderr, "0.1.0\n")
		os.Exit(0)
	}

	_, err := flags.Parse(&opts)
	if err != nil {
		if flags.WroteHelp(err) {
			os.Exit(0)
		}
		log.Fatal(err)
		os.Exit(1)
	}

	gitInfo, err := getGitInfo(string(opts.GitRepositoryPath))
	if err != nil && !opts.IgnoreGitErrors {
		log.Fatal(fmt.Sprintf("%+v", err))
		os.Exit(1)
	}

	t := template.Must(template.ParseFiles(string(opts.Args.TemplateFilename)))
	var templateResult bytes.Buffer
	if err = t.Execute(&templateResult, gitInfo); err != nil {
		log.Fatal(fmt.Sprintf("%+v", err))
	}

	if opts.OutputFilePath == "" {
		fmt.Printf("%s", string(templateResult.String()))
	} else {
		overwrite := opts.ForceOverwrite
		if !overwrite {
			destinationContent, _ := ioutil.ReadFile(string(opts.OutputFilePath))
			overwrite = !bytes.Equal(templateResult.Bytes(), destinationContent)
		}
		if overwrite {
			file, err := os.Create(string(opts.OutputFilePath))
			if err != nil {
				log.Fatal(fmt.Sprintf("%+v", err))
			}
			defer file.Close()
			w := bufio.NewWriter(file)
			w.WriteString(templateResult.String())
			w.Flush()
		}
	}
}
