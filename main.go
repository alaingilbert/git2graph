package main

import (
	"git2graph/git2graph"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func startAction(c *cli.Context) error {
	var nodes []git2graph.Node
	var err error
	fromFlag := c.Int("from")
	sizeFlag := c.Int("size")
	contextFlag := c.Bool("context")
	jsonFlag := c.String("json")
	fileFlag := c.String("file")
	repoFlag := c.Bool("repo")
	git2graph.NoOutput = c.Bool("no-output")
	repoLinearFlag := c.Bool("repo-linear")
	seqIds := c.Bool("seq-ids")
	logLevel := c.String("log")
	setLogLevel(logLevel)

	if repoFlag {
		nodes, err = git2graph.GetInputNodesFromRepo(seqIds)
	} else if repoLinearFlag {
		nodes, err = git2graph.GetInputNodesFromRepo(seqIds)
		git2graph.SerializeOutput(nodes)
		return err
	} else if jsonFlag != "" {
		nodes, err = git2graph.GetInputNodesFromJSON([]byte(jsonFlag))
	} else if fileFlag != "" {
		nodes, err = git2graph.GetInputNodesFromFile(fileFlag)
	} else {
		return cli.ShowAppHelp(c)
	}
	if err != nil {
		log.Error(err)
		return err
	}

	out, err := git2graph.Get(nodes)
	if err != nil {
		log.Error(err)
		return err
	}
	for _, node := range out {
		delete(node, "parentsPaths")
	}

	var tmp []git2graph.Node
	if fromFlag >= 0 && sizeFlag >= 1 {
		// TODO: include context (nodes before "from" that have parents inside or after the range)
		if contextFlag {
			for _, node := range out {
				nodeIdx := node["idx"].(int)
				parentsPaths := node["parents_paths"].([]git2graph.Path)
				hasParentsInContext := false
				for _, nodeParent := range parentsPaths {
					if nodeParent.Path[len(nodeParent.Path)-1].Y >= fromFlag {
						hasParentsInContext = true
					}
				}
				if hasParentsInContext ||
					nodeIdx >= fromFlag && nodeIdx < fromFlag+sizeFlag {
					tmp = append(tmp, node)
				}
				if nodeIdx >= fromFlag+sizeFlag-1 {
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
			Name:  "r, repo",
			Usage: "Repository",
		},
		cli.BoolFlag{
			Name:  "l, repo-linear",
			Usage: "Repository linear history",
		},
		cli.BoolFlag{
			Name:  "s, seq-ids",
			Usage: "Use sequential ids instead of sha for linear history",
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
	app.Action = startAction
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
