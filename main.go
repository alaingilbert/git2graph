package main

import (
	"github.com/alaingilbert/git2graph/git2graph"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func startAction(c *cli.Context) error {
	var nodes []*git2graph.Node
	var err error
	fromFlag := c.String("from")
	limitFlag := c.Int("limit")
	//contextFlag := c.Bool("context")
	jsonFlag := c.String("json")
	fileFlag := c.String("file")
	repoFlag := c.Bool("repo")
	topoOrderFlag := c.Bool("topo-order")
	dateOrderFlag := c.Bool("date-order")
	git2graph.NoOutput = c.Bool("no-output")
	repoLinearFlag := c.Bool("repo-linear")
	seqIds := c.Bool("seq-ids")
	rowsFlag := c.Bool("rows")
	logLevel := c.String("log")
	setLogLevel(logLevel)

	if repoFlag || repoLinearFlag {
		order := git2graph.DefaultOrder
		if topoOrderFlag {
			order = git2graph.TopoOrder
		} else if dateOrderFlag {
			order = git2graph.DateOrder
		}
		if seqIds {
			nodes, err = git2graph.GetInputNodesFromRepoSeq("", order, limitFlag)
		} else {
			nodes, err = git2graph.GetInputNodesFromRepo("", order, limitFlag)
		}
		if repoLinearFlag {
			git2graph.SerializeOutput(&git2graph.Out{Nodes: nodes})
			return err
		}
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

	var out *git2graph.Out
	if rowsFlag {
		out, err = git2graph.GetPaginatedRows(nodes, fromFlag, limitFlag)
	} else {
		out, err = git2graph.GetPaginated(nodes, fromFlag, limitFlag)
	}
	if err != nil {
		log.Error(err)
		return err
	}

	git2graph.SerializeOutput(out)

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
	authors = append(authors, cli.Author{Name: "Alain Gilbert", Email: "alain.gilbert.15@gmail.com"})

	app := cli.NewApp()
	app.Authors = authors
	app.Version = "0.0.0"
	app.Name = "git2graph"
	app.Usage = "Take a git tree, make a graph structure"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "f, file", Usage: "File"},
		cli.StringFlag{Name: "j, json", Usage: "Json input"},
		cli.StringFlag{Name: "L, log", Usage: "Log level"},
		cli.BoolFlag{Name: "r, repo", Usage: "Repository"},
		cli.BoolFlag{Name: "topo-order", Usage: "Topological order"},
		cli.BoolFlag{Name: "l, repo-linear", Usage: "Repository linear history"},
		cli.BoolFlag{Name: "s, seq-ids", Usage: "Use sequential ids instead of sha for linear history"},
		cli.BoolFlag{Name: "n, no-output", Usage: "No output"},
		cli.BoolFlag{Name: "rows", Usage: "Rows graph"},
		cli.StringFlag{Name: "from", Usage: "From"},
		cli.IntFlag{Name: "limit", Usage: "Limit", Value: -1},
		cli.BoolFlag{Name: "context", Usage: "Include context"},
	}
	app.Action = startAction
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
