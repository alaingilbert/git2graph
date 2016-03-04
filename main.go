package main

import (
	"git2graph/git2graph"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"os"
)

func bootstrap(c *cli.Context) {
	var nodes []map[string]interface{}
	var err error
	jsonFlag := c.String("json")
	fileFlag := c.String("file")
	git2graph.DebugMode = c.Bool("debug")
	repoFlag := c.Bool("repo")
	git2graph.NoOutput = c.Bool("no-output")
	repoLinearFlag := c.Bool("repo-linear")
	logLevel := c.String("log")
	setLogLevel(logLevel)

	if repoFlag {
		nodes, err = git2graph.GetInputNodesFromRepo()
	} else if repoLinearFlag {
		nodes, err = git2graph.GetInputNodesFromRepo()
		git2graph.SerializeOutput(nodes)
		return
	} else if jsonFlag != "" {
		nodes, err = git2graph.GetInputNodesFromJson(jsonFlag)
	} else if fileFlag != "" {
		nodes, err = git2graph.GetInputNodesFromFile(fileFlag)
	} else {
		cli.ShowAppHelp(c)
		return
	}
	if err != nil {
		log.Error(err)
		return
	}

	myColors := git2graph.DefaultColors

	out, err := git2graph.BuildTree(nodes, myColors)
	if err != nil {
		log.Error(err)
		return
	}
	for _, node := range out {
		delete(node, "parentsPaths")
	}

	git2graph.SerializeOutput(out)
}

func setLogLevel(logLevel string) {
	switch logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.WarnLevel)
	}
}

func init() {
	log.SetLevel(log.WarnLevel)
}

func main() {
	var authors []cli.Author
	// Collaborators, add your name here :)
	authors = append(authors, cli.Author{"Alain Gilbert", "alain.gilbert.15@gmail.com"})

	app := cli.NewApp()
	app.Authors = authors
	app.Version = "0.0.0"
	app.Name = "git2graph"
	app.Usage = "Take a git tree, make a graph structure"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "f, file",
			Usage: "File",
		},
		cli.StringFlag{
			Name:  "j, json",
			Usage: "Json input",
		},
		cli.StringFlag{
			Name:  "L, log",
			Usage: "Log level",
		},
		cli.BoolFlag{
			Name:  "d, debug",
			Usage: "Debug mode",
		},
		cli.BoolFlag{
			Name:  "r, repo",
			Usage: "Repository",
		},
		cli.BoolFlag{
			Name:  "l, repo-linear",
			Usage: "Repository linear history",
		},
		cli.BoolFlag{
			Name:  "n, no-output",
			Usage: "No output",
		},
	}
	app.Action = bootstrap
	app.Run(os.Args)
}
