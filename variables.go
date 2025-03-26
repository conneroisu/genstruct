package genstruct

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/dave/jennifer/jen"
)

// generateVariables creates variables for each struct
func (g *Generator) generateVariables(dataValue reflect.Value) {
	// Generate a variable for each struct
	for i := range dataValue.Len() {
		elem := dataValue.Index(i)

		// Determine the variable name using the identifier function
		identValue := g.getStructIdentifier(elem)
		varName := g.Config.VarPrefix + slugToIdentifier(identValue)

		// Create the variable with its value
		g.File.Var().Id(varName).Op("=").Id(g.Config.TypeName).ValuesFunc(func(group *jen.Group) {
			g.generateStructValues(group, elem)
		})
	}
}

// generateSlice creates a slice containing all struct instances
func (g *Generator) generateSlice(dataValue reflect.Value) {
	// Determine the slice name - handle both regular and irregular plurals
	var sliceName string
	if g.Config.TypeName[len(g.Config.TypeName)-1] == 's' ||
		g.Config.TypeName[len(g.Config.TypeName)-1] == 'x' ||
		g.Config.TypeName[len(g.Config.TypeName)-1] == 'z' ||
		strings.HasSuffix(g.Config.TypeName, "sh") ||
		strings.HasSuffix(g.Config.TypeName, "ch") {
		sliceName = fmt.Sprintf(
			"All%ses",
			g.Config.TypeName,
		)
	} else if g.Config.TypeName[len(g.Config.TypeName)-1] == 'y' {
		sliceName = fmt.Sprintf(
			"All%sies",
			g.Config.TypeName[:len(g.Config.TypeName)-1],
		)
	} else {
		sliceName = fmt.Sprintf("All%ss", g.Config.TypeName)
	}

	g.File.Var().Id(
		sliceName,
	).Op(
		"=",
	).Index().Id(
		g.Config.TypeName,
	).ValuesFunc(func(group *jen.Group) {
		for i := range dataValue.Len() {
			elem := dataValue.Index(i)

			// Get the variable name using the same method as in generateVariables
			identValue := g.getStructIdentifier(elem)
			varName := g.Config.VarPrefix + slugToIdentifier(identValue)

			group.Id(varName)
		}
	})
}
