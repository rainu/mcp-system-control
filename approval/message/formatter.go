package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/mark3labs/mcp-go/mcp"
)

// Formatter handles message formatting in different languages
type Formatter struct {
	language  Language
	templates map[string]*template.Template
}

// NewFormatter creates a new message formatter
func NewFormatter(lang Language) (*Formatter, error) {
	f := &Formatter{
		language:  lang,
		templates: make(map[string]*template.Template),
	}

	// Compile all templates
	for toolName, langTemplates := range templates {
		templateStr, ok := langTemplates[lang]
		if !ok {
			// Fallback to English if language not found
			templateStr = langTemplates[LanguageEnglish]
		}

		tmpl, err := template.New(toolName).Parse(templateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template for %s: %w", toolName, err)
		}
		f.templates[toolName] = tmpl
	}

	return f, nil
}

// NewFormatterAuto creates a new formatter with auto-detected language
func NewFormatterAuto() (*Formatter, error) {
	return NewFormatter(DetectLanguage())
}

// Format formats a tool request into a human-readable message
func (f *Formatter) Format(request *mcp.CallToolRequest) string {
	toolName := request.Params.Name
	tmpl, ok := f.templates[toolName]
	if !ok {
		// Use generic template
		tmpl = f.templates["generic"]
	}

	data := f.prepareData(request)

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		// Fallback to generic message on error
		return f.formatGeneric(request)
	}

	return buf.String()
}

// prepareData prepares the data for template execution
func (f *Formatter) prepareData(request *mcp.CallToolRequest) map[string]interface{} {
	data := make(map[string]interface{})

	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return data
	}

	// Copy all arguments to data
	for k, v := range args {
		data[k] = v
	}

	// Add special handling for content preview
	if content, ok := args["content"].(string); ok {
		if len(content) > 100 {
			data["content_preview"] = content[:100]
			data["content_truncated"] = true
		} else {
			data["content_preview"] = content
			data["content_truncated"] = false
		}
		data["content_size"] = len(content)
	}

	// Add special handling for body preview (HTTP calls)
	if body, ok := args["body"].(string); ok {
		if len(body) > 100 {
			data["body_preview"] = body[:100]
			data["body_truncated"] = true
		} else {
			data["body_preview"] = body
			data["body_truncated"] = false
		}
	}

	// Set default method for HTTP calls
	if request.Params.Name == "callHttp" {
		if _, ok := data["method"]; !ok {
			data["method"] = "GET"
		}
	}

	// Format numeric values properly
	if uid, ok := args["user_id"].(float64); ok {
		data["user_id"] = fmt.Sprintf("%.0f", uid)
	}
	if gid, ok := args["group_id"].(float64); ok {
		data["group_id"] = fmt.Sprintf("%.0f", gid)
	}
	if lo, ok := args["lo"].(float64); ok {
		data["lo"] = fmt.Sprintf("%.0f", lo)
	}
	if ll, ok := args["ll"].(float64); ok {
		data["ll"] = fmt.Sprintf("%.0f", ll)
	}

	return data
}

// formatGeneric formats a generic message when template is not found or fails
func (f *Formatter) formatGeneric(request *mcp.CallToolRequest) string {
	data := map[string]interface{}{
		"tool_name": request.Params.Name,
		"arguments": "",
	}

	if request.Params.Arguments != nil {
		if argsBytes, err := json.MarshalIndent(request.Params.Arguments, "", "  "); err == nil {
			data["arguments"] = string(argsBytes)
		} else {
			data["arguments"] = fmt.Sprintf("%v", request.Params.Arguments)
		}
	}

	tmpl := f.templates["generic"]
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		// Last resort fallback
		return fmt.Sprintf("Tool: %s\nArguments: %v", request.Params.Name, request.Params.Arguments)
	}

	return buf.String()
}
