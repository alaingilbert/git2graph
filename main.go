package main

import (
	"git2graph/git2graph"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"os"
)

func bootstrap(c *cli.Context) error {
	var nodes []map[string]interface{}
	var err error
	fromFlag := c.Int("from")
	sizeFlag := c.Int("size")
	contextFlag := c.Bool("context")
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
		return err
	} else if jsonFlag != "" {
		nodes, err = git2graph.GetInputNodesFromJSON(jsonFlag)
	} else if fileFlag != "" {
		nodes, err = git2graph.GetInputNodesFromFile(fileFlag)
	} else {
		cli.ShowAppHelp(c)
		return err
	}
	if err != nil {
		log.Error(err)
		return err
	}

	myColors := git2graph.DefaultColors

	out, err := git2graph.BuildTree(nodes, myColors)
	if err != nil {
		log.Error(err)
		return err
	}
	for _, node := range out {
		delete(node, "parentsPaths")
	}

	var tmp []map[string]interface{}
	if fromFlag >= 0 && sizeFlag >= 1 {
		// TODO: include context (nodes before "from" that have parents inside or after the range)
		if contextFlag {
			for _, node := range out {
				hasParentsInContext := false
				for _, nodeParent := range node["parents_paths"].([]git2graph.Path) {
					if nodeParent.Path[len(nodeParent.Path)-1].Y >= fromFlag {
						hasParentsInContext = true
					}
				}
				if hasParentsInContext ||
					node["idx"].(int) >= fromFlag && node["idx"].(int) < fromFlag+sizeFlag {
					tmp = append(tmp, node)
				}
				if node["idx"].(int) >= fromFlag+sizeFlag-1 {
					break
				}
			}
		} else {
			tmp = out[fromFlag : fromFlag+sizeFlag]
		}
	} else {
		tmp = out
	}

	git2graph.SerializeOutput(tmp)

	return err
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
		// TODO: From should be a sha
		cli.IntFlag{
			Name:  "from",
			Usage: "From",
			Value: -1,
		},
		cli.IntFlag{
			Name:  "size",
			Usage: "Size",
			Value: -1,
		},
		cli.BoolFlag{
			Name:  "context",
			Usage: "Include context",
		},
	}
	app.Action = bootstrap
	app.Run(os.Args)
}
