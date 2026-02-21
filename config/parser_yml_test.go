package config

import (
	"mcp-system-control/config/model"
	"strings"
	"testing"

	"github.com/rainu/go-yacl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_processYaml(t *testing.T) {
	c := &model.Config{}

	yamlContent := `
ui:
  window:
    title: Test Window
    init-width:
      expression: "800"
      value: 0
    max-height:
      expression: "600"
      value: 0
    init-pos-x:
      expression: "100"
      value: 0
    init-pos-y:
      expression: "100"
      value: 0
    init-zoom:
      expression: "1.0"
      value: 0
    bg-color:
      r: 255
      g: 255
      b: 255
      a: 255
    start-state: 1
    always-on-top: true
    frameless: true
    resizeable: true
    translucent: never
  prompt:
    min-rows: 1
    max-rows: 10
    submit:
      binding:
        - "alt+ctrl+meta+shift+enter"
  file-dialog:
    default-dir: /root
    show-hidden: true
    can-create-dirs: true
    resolve-aliases: true
    treat-packages-as-dirs: true
    filter-display:
      - Image
    filter-pattern:
      - '*.png'
  stream: true
  quit:
    binding:
      - "alt+ctrl+meta+shift+escape"
  theme: dark
  code-style: default
  lang: en
llm:
  backend: anthropic
  localai:
    api-key: 
      plain: APIKey
    model: model
    base-url: baseurl
  openai:
    api-key:
      command: 
        name: echo
        args:
          - APIKey
    api-type: APIType
    api-version: APIVersion
    model: Model
    base-url: BaseUrl
    organization: Organization
  anythingllm:
    base-url: BaseURL
    token:
      command: 
        name: echo
        args:
          - '-n'
          - Token
        no-trim: true
    workspace: Workspace
  ollama:
    server-url: ServerURL
    model: Model
  mistral:
    api-key:
      command: 
        name: echo
        args:
          - '-n'
          - ApiKey
    endpoint: Endpoint
    model: Model
  anthropic:
    api-key: 
      plain: Token
    base-url: BaseUrl
    model: Model
  call:
    prompt:
      system: Your system prompt
      init-message:
        - role: human
          content: This is an initial message.
      init-tool-call:
        - server: _builtin
          name: getSystemInformation
        - server: _custom
          name: test
          args: 
            arg1: value1
            arg2: 42
      init-value: Initial Prompt
      init-attachment:
        - attachment1
        - attachment2
    max-token: 1000
    temperature: 0.7
    top-k: 50
    top-p: 0.9
    min-length: 10
    max-length: 200
  tool:
    builtin:
      command-execution:
        disable: true
    custom:
      test:
        description: This is a test function.
        parameters:
          type: object
          properties:
            arg1:
              type: string
              description: The first argument.
            arg2:
              type: number
              description: The second argument.
          required:
            - arg1
        command: doTest.sh
        "approval": true
    mcpServers:
      command1:
        command: docker
        args:
          - run
          - --rm
          - -i
          - -e
          - GITHUB_PERSONAL_ACCESS_TOKEN=github_
          - ghcr.io/github/github-mcp-server
        env:
          TEST: test
        timeout:
          init: 500ms
          list: 10s
          execution: 1m30s
      command2:
        command: echo
      http:
        baseUrl: http://localhost:8080/api/v1
        headers:
          Authorization: Bearer TOKEN
print:
  format: json
  targets:
    - stdout
log-level: debug
pprof-address: ":1312"
vue-dev-tools:
  host: "localhost"
  port: 1312
webkit:
  open-inspector: true
  http-server: "127.0.0.1:5000"
themes:
  dark:
    colors:
      background: "#FFFFFF"
      surface: "#FFFFFF"
  light:
    colors:
      chat-tool-call: "#FF0000"
  custom:
    test:
      colors:
         background: "#00FF00"
         chat-tool-call: "#0000FF"
`
	// add profile "test" with the same values as default
	yamlContent += "\nprofiles:\n  test:\n" + strings.ReplaceAll(yamlContent, "\n", "\n    ")

	sr := strings.NewReader(yamlContent)
	config := yacl.NewConfig(c, yacl.WithAutoApplyDefaults(false))

	require.NoError(t, processYaml(config, sr))

	expDefConf := model.Config{}

	expConfig := expDefConf

	assert.Equal(t, &expConfig, c)
}
