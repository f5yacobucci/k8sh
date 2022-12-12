package builtin

import (
	"context"
	"errors"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	k8spath "k8sh/internal/path"
	"k8sh/internal/symtab"
	"os"
	"path"
	"strings"
)

const pathMaxDepth = 3

const (
	listNamespaces = iota
	listAll        = iota
	listKind       = iota
	listResource   = iota
)

type BuiltInExec func(argc int, argv []string) error
type BuiltIn struct {
	name string
	exec BuiltInExec
}
type BuiltInMap map[string]BuiltIn

var BuiltIns = BuiltInMap{
	"echo": BuiltIn{
		"echo",
		echo,
	},
	"dump": BuiltIn{
		"dump",
		dump,
	},
	"ls": BuiltIn{
		"ls",
		ls,
	},
	"cd": BuiltIn{
		"cd",
		cd,
	},
	"pwd": BuiltIn{
		"pwd",
		pwd,
	},
	"exit": BuiltIn{
		"exit",
		exit,
	},
}

func (b BuiltIn) Exec(argc int, argv []string) error {
	return b.exec(argc, argv)
}

func echo(_ int, argv []string) error {
	fmt.Printf("%v\n", argv[len(argv)-1])
	return nil
}

func dump(_ int, _ []string) error {
	symtab.DumpLocalTable()
	return nil
}

func ls(argc int, argv []string) error {
	//var cmd string
	var p string

	entry := symtab.GetSymbolEntry("CWD")
	if entry == nil {
		return errors.New("ls: no current working directory")
	}

	//cmd = argv[0]
	if argc == 1 {
		p = entry.GetValue()
	} else {
		p = argv[len(argv)-1]

		if !strings.HasPrefix(p, "/") {
			entry := symtab.GetSymbolEntry("CWD")
			if entry == nil || entry.GetValue() == "" {
				return errors.New("ls: relpath error - no CWD")
			}
			p = entry.GetValue() + "/" + p
		}
	}

	p = path.Clean(p)
	elems := strings.Split(p[1:], "/")

	if len(elems) > pathMaxDepth {
		return errors.New("ls: path exceeded max depth")
	}

	/*
		var ns string
		var kind string
		var res string
		var query int = listNamespaces
		for i := range elems {
			switch i {
			case 0:
				ns = elems[i]
				if ns != "" {
					query = listAll
				}
			case 1:
				kind = elems[i]
				query = listKind
			case 2:
				res = elems[i]
				query = listResource
			}
		}
	*/

	/* XXX for now only list namespaces from '/' */
	kube := symtab.GetSymbolEntry("KUBECONFIG")
	if kube == nil {
		return errors.New("ls: no kubeconfig set")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kube.GetValue())
	if err != nil {
		return fmt.Errorf("ls: no client: %w", err)
	}

	/*
		client, err := discovery.NewDiscoveryClientForConfig(config)
		if err != nil {
			return fmt.Errorf("ls: no discovery client: %w", err)
		}

		groups, resources, err := client.ServerGroupsAndResources()
		fmt.Printf("%v\n", groups)
		fmt.Printf("%v\n", resources)
		fmt.Printf("%v\n", err)
	*/

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("ls: no dynamic client: %w", err)
	}
	resource := schema.GroupVersionResource{
		Version:  "v1",
		Resource: "namespaces",
	}
	obj, err := dynamicClient.Resource(resource).
		Namespace("").
		List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("ls: cannot query root: %w", err)
	}

	for _, o := range obj.Items {
		fmt.Printf("%s\n", o.Object["metadata"].(map[string]interface{})["name"])
	}
	return nil
}

func cd(argc int, argv []string) error {
	if argc > 2 {
		return errors.New("cd: too many arguments")
	}

	if argc == 1 {
		entry := symtab.GetSymbolEntry("CWD")
		entry.SetValue("/")
		return nil
	}

	if argv[1] == "-" {
		// cd to OLDCWD
	}

	abs := k8spath.MakeAbsolute(argv[1])
	fmt.Printf("ABSOLUTE PATH: %s\n", abs)
	return nil
}

func pwd(_ int, argv []string) error {
	entry := symtab.GetSymbolEntry("CWD")
	if entry == nil || entry.GetValue() == "" {
		return errors.New("no current working directory")
	}
	fmt.Fprintf(os.Stdout, "%s\n", entry.GetValue())
	return nil
}

func exit(_ int, _ []string) error {
	os.Exit(0)
	return nil
}
