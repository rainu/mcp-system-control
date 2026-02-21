package message

// Template definitions for each tool in different languages
var templates = map[string]map[Language]string{
	"executeCommand": {
		LanguageEnglish: `🖥️  Execute Command: {{.command}}{{if .working_directory}}
Working Directory: {{.working_directory}}{{end}}{{if .environment}}
Environment Variables:{{range $key, $value := .environment}}
  {{$key}}={{$value}}{{end}}{{end}}`,
		LanguageGerman: `🖥️  Befehl ausführen: {{.command}}{{if .working_directory}}
Arbeitsverzeichnis: {{.working_directory}}{{end}}{{if .environment}}
Umgebungsvariablen:{{range $key, $value := .environment}}
  {{$key}}={{$value}}{{end}}{{end}}`,
	},
	"createFile": {
		LanguageEnglish: `📝 Create File: {{.path}}{{if .permission}}
Permission: {{.permission}}{{end}}{{if .content}}
Content: {{.content_preview}}{{if .content_truncated}}...{{end}}
Size: {{.content_size}} Bytes{{end}}`,
		LanguageGerman: `📝 Datei erstellen: {{.path}}{{if .permission}}
Berechtigung: {{.permission}}{{end}}{{if .content}}
Inhalt: {{.content_preview}}{{if .content_truncated}}...{{end}}
Größe: {{.content_size}} Bytes{{end}}`,
	},
	"deleteFile": {
		LanguageEnglish: `🗑️  Delete File: {{.path}}`,
		LanguageGerman:  `🗑️  Datei löschen: {{.path}}`,
	},
	"appendFile": {
		LanguageEnglish: `➕ Append to File: {{.path}}{{if .content}}
Content: {{.content_preview}}{{if .content_truncated}}...{{end}}
Size: {{.content_size}} Bytes{{end}}`,
		LanguageGerman: `➕ An Datei anhängen: {{.path}}{{if .content}}
Inhalt: {{.content_preview}}{{if .content_truncated}}...{{end}}
Größe: {{.content_size}} Bytes{{end}}`,
	},
	"readTextFile": {
		LanguageEnglish: `📖 Read File: {{.path}}{{if .lm}}
Limit Mode: {{.lm}}{{if .lo}}
Offset: {{.lo}}{{end}}{{if .ll}}
Limit: {{.ll}}{{end}}{{end}}`,
		LanguageGerman: `📖 Datei lesen: {{.path}}{{if .lm}}
Limit-Modus: {{.lm}}{{if .lo}}
Offset: {{.lo}}{{end}}{{if .ll}}
Limit: {{.ll}}{{end}}{{end}}`,
	},
	"createTempFile": {
		LanguageEnglish: `📄 Create Temporary File{{if .suffix}} (Suffix: {{.suffix}}){{end}}{{if .permission}}
Permission: {{.permission}}{{end}}{{if .content}}
Content: {{.content_preview}}{{if .content_truncated}}...{{end}}
Size: {{.content_size}} Bytes{{end}}`,
		LanguageGerman: `📄 Temporäre Datei erstellen{{if .suffix}} (Suffix: {{.suffix}}){{end}}{{if .permission}}
Berechtigung: {{.permission}}{{end}}{{if .content}}
Inhalt: {{.content_preview}}{{if .content_truncated}}...{{end}}
Größe: {{.content_size}} Bytes{{end}}`,
	},
	"createDirectory": {
		LanguageEnglish: `📁 Create Directory: {{.path}}{{if .permission}}
Permission: {{.permission}}{{end}}`,
		LanguageGerman: `📁 Verzeichnis erstellen: {{.path}}{{if .permission}}
Berechtigung: {{.permission}}{{end}}`,
	},
	"deleteDirectory": {
		LanguageEnglish: `🗑️  Delete Directory: {{.path}}
⚠️  All files and subdirectories will be deleted!`,
		LanguageGerman: `🗑️  Verzeichnis löschen: {{.path}}
⚠️  Alle Dateien und Unterverzeichnisse werden gelöscht!`,
	},
	"createTempDirectory": {
		LanguageEnglish: `📂 Create Temporary Directory`,
		LanguageGerman:  `📂 Temporäres Verzeichnis erstellen`,
	},
	"changeMode": {
		LanguageEnglish: `🔐 Change Permission: {{.path}}{{if .permission}} → {{.permission}}{{end}}`,
		LanguageGerman:  `🔐 Berechtigung ändern: {{.path}}{{if .permission}} → {{.permission}}{{end}}`,
	},
	"changeOwner": {
		LanguageEnglish: `👤 Change Owner: {{.path}}{{if .user_id}}
User ID: {{.user_id}}{{end}}{{if .group_id}}
Group ID: {{.group_id}}{{end}}`,
		LanguageGerman: `👤 Eigentümer ändern: {{.path}}{{if .user_id}}
Benutzer-ID: {{.user_id}}{{end}}{{if .group_id}}
Gruppen-ID: {{.group_id}}{{end}}`,
	},
	"changeTimes": {
		LanguageEnglish: `🕐 Change Timestamps: {{.path}}{{if .access_time}}
Access Time: {{.access_time}}{{end}}{{if .modification_time}}
Modification Time: {{.modification_time}}{{end}}`,
		LanguageGerman: `🕐 Zeitstempel ändern: {{.path}}{{if .access_time}}
Zugriffszeit: {{.access_time}}{{end}}{{if .modification_time}}
Änderungszeit: {{.modification_time}}{{end}}`,
	},
	"getStats": {
		LanguageEnglish: `ℹ️  Get File Information: {{.path}}`,
		LanguageGerman:  `ℹ️  Dateiinformationen abrufen: {{.path}}`,
	},
	"getSystemTime": {
		LanguageEnglish: `🕐 Get System Time`,
		LanguageGerman:  `🕐 Systemzeit abrufen`,
	},
	"getEnvironment": {
		LanguageEnglish: `🌍 Get Environment Variables`,
		LanguageGerman:  `🌍 Umgebungsvariablen abrufen`,
	},
	"getSystemInformation": {
		LanguageEnglish: `💻 Get System Information`,
		LanguageGerman:  `💻 Systeminformationen abrufen`,
	},
	"callHttp": {
		LanguageEnglish: `🌐 HTTP Call: {{.method}} {{.url}}{{if .header}}
Headers:{{range $key, $value := .header}}
  {{$key}}: {{$value}}{{end}}{{end}}{{if .body}}
Body: {{.body_preview}}{{if .body_truncated}}...{{end}}{{end}}`,
		LanguageGerman: `🌐 HTTP-Aufruf: {{.method}} {{.url}}{{if .header}}
Header:{{range $key, $value := .header}}
  {{$key}}: {{$value}}{{end}}{{end}}{{if .body}}
Body: {{.body_preview}}{{if .body_truncated}}...{{end}}{{end}}`,
	},
	"generic": {
		LanguageEnglish: `Tool: {{.tool_name}}
Arguments:
{{.arguments}}`,
		LanguageGerman: `Tool: {{.tool_name}}

Argumente:
{{.arguments}}`,
	},
}
