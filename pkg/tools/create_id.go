package tools

import "github.com/bwmarrin/snowflake"

func CreateID() string {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}

	id := node.Generate()

	return id.String()
}
