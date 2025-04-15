package genstruct

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/dave/jennifer/jen"
)

// generateInitFunction creates an init() function to set up cross-references
// This is used to avoid initialization cycles in export mode
func (g *Generator) generateInitFunction() {
	// Collect all the references that need to be initialized
	// Map from struct type name to a slice of struct field references
	type fieldRef struct {
		fieldName   string // Field to be populated (e.g., "Tags")
		srcField    string // Source field containing identifiers (e.g., "TagSlugs")
		targetType  string // Type of the references (e.g., "Tag")
		targetSlice string // All-slice name for the type (e.g., "AllTags")
		isPointer   bool   // Whether the reference is a pointer
	}

	refMap := make(map[string][]fieldRef)

	// Analyze all types for structgen tags
	for typeName, refDataObj := range g.Refs {
		refDataValue := reflect.ValueOf(refDataObj)
		if refDataValue.Kind() != reflect.Slice && refDataValue.Kind() != reflect.Array || refDataValue.Len() == 0 {
			continue
		}

		// Get the first element to analyze its type
		refElem := refDataValue.Index(0)
		var structType reflect.Type

		// Handle both direct struct and pointer-to-struct cases
		if refElem.Kind() == reflect.Struct {
			structType = refElem.Type()
		} else if refElem.Kind() == reflect.Pointer && refElem.Elem().Kind() == reflect.Struct {
			structType = refElem.Elem().Type()
		} else {
			continue
		}

		// Look for structgen tags in fields
		for i := 0; i < structType.NumField(); i++ {
			field := structType.Field(i)

			// Skip unexported fields
			if !field.IsExported() {
				continue
			}

			// Check for structgen tag
			structgenVal, hasStructgenTag := field.Tag.Lookup("structgen")
			if !hasStructgenTag || structgenVal == "" {
				continue
			}

			// Get source field name from tag
			srcField := structgenVal

			// Only process fields that are slices of structs or pointers to structs
			if field.Type.Kind() != reflect.Slice {
				continue
			}

			elemType := field.Type.Elem()
			isPointer := elemType.Kind() == reflect.Pointer

			// Get the actual struct type
			var targetType reflect.Type
			if isPointer {
				targetType = elemType.Elem()
			} else {
				targetType = elemType
			}

			// Skip if not a struct
			if targetType.Kind() != reflect.Struct {
				continue
			}

			// Add to reference map
			targetTypeName := targetType.Name()
			targetSlice := "All" + targetTypeName + "s" // Default to plural form

			// Some common irregular plurals
			switch targetTypeName {
			case "Entry":
				targetSlice = "AllEntries"
			case "Category":
				targetSlice = "AllCategories"
			}

			refMap[typeName] = append(refMap[typeName], fieldRef{
				fieldName:   field.Name,
				srcField:    srcField,
				targetType:  targetTypeName,
				targetSlice: targetSlice,
				isPointer:   isPointer,
			})
		}
	}

	// If no cross-references found, don't generate init function
	if len(refMap) == 0 {
		return
	}

	// Start the init function
	init := jen.Func().Id("init").Params()

	// Create the function body
	body := jen.Block().Add(
		jen.Comment("Set up cross-references in structs to avoid initialization cycles"),
		jen.Line(),
	)

	// For each struct type that has references
	for typeName, fieldRefs := range refMap {
		// Generate code to update "All[TypeName]s" for each field reference
		allSliceName := "All" + typeName + "s" // Default to plural form

		// Some common irregular plurals
		switch typeName {
		case "Entry":
			allSliceName = "AllEntries"
		case "Category":
			allSliceName = "AllCategories"
		}

		// Add a comment for this type
		body.Add(jen.Comment(fmt.Sprintf("Initialize references for %s", typeName)))

		// Skip if we don't have the slice
		ifCheck := jen.If(jen.Id(allSliceName).Op("!=").Nil())
		ifBody := jen.Block()

		// For loop over the slice
		forLoop := jen.For(jen.Id("_").Op(",").Id(strings.ToLower(typeName)).Op(":=").Range().Id(allSliceName))
		forBody := jen.Block()

		// For each reference field in this type
		for _, ref := range fieldRefs {
			// Add a comment for this field
			forBody.Add(jen.Comment(fmt.Sprintf("Initialize %s field with %s references", ref.fieldName, ref.targetType)))

			// Skip if source slice is empty
			ifSrcCheck := jen.If(jen.Len(jen.Id(strings.ToLower(typeName)).Dot(ref.srcField)).Op(">").Lit(0))
			ifSrcBody := jen.Block()

			// Reset the target slice to empty but with same capacity
			ifSrcBody.Add(jen.Id(strings.ToLower(typeName)).Dot(ref.fieldName).Op("=").Id(strings.ToLower(typeName)).Dot(ref.fieldName).Op("[:0]"))

			// Loop through each ID in source field
			idLoop := jen.For(jen.Id("_").Op(",").Id("id").Op(":=").Range().Id(strings.ToLower(typeName)).Dot(ref.srcField))
			idLoopBody := jen.Block()

			// Loop through targets to find match
			targetLoop := jen.For(jen.Id("_").Op(",").Id("target").Op(":=").Range().Id(ref.targetSlice))
			targetLoopBody := jen.Block()

			// Check if target.Slug matches the ID
			ifMatch := jen.If(jen.Id("target").Dot("Slug").Op("==").Id("id"))
			ifMatchBody := jen.Block()

			// Add the matching reference to the slice
			if ref.isPointer {
				// For pointer references ([]*Type)
				ifMatchBody.Add(jen.Id(strings.ToLower(typeName)).Dot(ref.fieldName).Op("=").Append(
					jen.Id(strings.ToLower(typeName)).Dot(ref.fieldName),
					jen.Id("target"),
				))
			} else {
				// For value references ([]Type)
				ifMatchBody.Add(jen.Id(strings.ToLower(typeName)).Dot(ref.fieldName).Op("=").Append(
					jen.Id(strings.ToLower(typeName)).Dot(ref.fieldName),
					jen.Op("*").Id("target"),
				))
			}

			// Break out of inner loop when matched
			ifMatchBody.Add(jen.Break())

			// Build up the nested loops
			targetLoopBody.Add(ifMatch.Block(ifMatchBody))
			idLoopBody.Add(targetLoop.Block(targetLoopBody))
			ifSrcBody.Add(idLoop.Block(idLoopBody))
			forBody.Add(ifSrcCheck.Block(ifSrcBody))
		}

		// Add the loops to the block
		ifBody.Add(forLoop.Block(forBody))
		body.Add(ifCheck.Block(ifBody))
		body.Add(jen.Line())
	}

	// Add the init function to the file
	g.File.Add(init.Block(body))
}
