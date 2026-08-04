package templates

type ScenarioData struct {
	Selection string
	Entries   []ScenarioEntry
}

type ScenarioEntry struct {
	Owner      string
	Host       string
	Image      string
	IP         string
	Username   string
	Password   string
	Connection string
}

const Scenario = `
<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <style type="text/css">
      body {
        margin: 0 auto;
        font-family: sans-serif;
        font-size: 12pt;
        max-width: 1800px;
      }
      table {
        border: none;
      }
      td,
      th {
        border: none;
        padding: 0.5em 1em;
        text-align: left;
      }
      table tr:nth-child(even) td {
        background: black;
        color: white;
      }
      table tr:nth-child(odd) td {
        background: white;
        color: black;
      }
    </style>
    <title>Instance Overview</title>
  </head>
  <body>
    <h1>
      Scenario Overview{{ if (gt (len .Selection) 0) }} for <tt>{{ .Selection }}</tt>{{ end }}
    </h1>
    <table>
      <tr>
        <th>Owner</th>
        <th>Host</th>
        <th>Image</th>
        <th>IP</th>
        <th>Username</th>
        <th>Password</th>
        <th>Connection</th>
      </tr>
      {{ range .Entries }}
      <tr>
        <td><tt>{{ .Owner }}</tt></td>
        <td><tt>{{ .Host }}</tt></td>
        <td><tt>{{ .Image }}</tt></td>
        <td><tt>{{ .IP }}</tt></td>
        <td><tt>{{ .Username }}</tt></td>
        <td><tt>{{ .Password }}</tt></td>
        <td><tt>{{ .Connection }}</tt></td>
      </tr>
      {{ end }}
    </table>
  </body>
</html>
`
