// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package manifest

import (
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

// Node represents a generic mnifest node
type Node struct {
	ManifType `yaml:",inline"`

	FileType `yaml:",inline"`

	DirType `yaml:",inline"`

	FilesTreeType `yaml:",inline"`

	// Properties of the node
	Properties map[string]interface{} `yaml:"properties,omitempty"`
	// Frontmatter of the node
	Frontmatter map[string]interface{} `yaml:"frontmatter,omitempty"`
	// Type of node
	Type string `yaml:"type,omitempty"`
	// Path of node
	Path string `yaml:"path,omitempty"`
	// Parent of node
	parent *Node
}

// Name is the name of the node
func (n *Node) Name() string {
	switch n.Type {
	case "file":
		return n.File
	case "dir":
		return n.Dir
	default:
		return ""
	}
}

// NodePath returns fully qualified name of this node
// i.e. Node.Path + Node.Name
func (n *Node) NodePath() string {
	return path.Join(n.Path, n.Name())
}

// HugoPrettyPath returns hugo pretty path
func (n *Node) HugoPrettyPath() string {
	name := n.Name()
	name = strings.TrimSuffix(name, ".md")
	name = strings.TrimSuffix(name, "_index")
	return path.Join(n.Path, name) + "/"
}

// HasContent returns true if the node is a document node
func (n *Node) HasContent() bool {
	return len(n.MultiSource) > 0 || len(n.Source) > 0
}

// Parent is the node parent
func (n *Node) Parent() *Node {
	return n.parent
}

func (n *Node) String() string {
	node, err := yaml.Marshal(n)
	if err != nil {
		return ""
	}
	return string(node)
}
