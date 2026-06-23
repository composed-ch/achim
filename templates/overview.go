package templates

type OverviewData struct {
	Selection string
	Entries   []OverviewEntry
}

type OverviewEntry struct {
	Owner      string
	SSHCommand string
	HostName   string
	IPAddress  string
}

const Overview = `
<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <style type="text/css">
      body {
        margin: 0 auto;
        font-family: sans-serif;
        font-size: 24px;
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
      Instance Overview{{ if (gt (len .Selection) 0) }} for <tt>{{ .Selection }}</tt>{{ end }}
    </h1>
    <table>
      <tr>
        <th>Owner</th>
        <th>SSH Command</th>
        <th>Host Name</th>
        <th>IP</th>
      </tr>
      {{ range .Entries }}
      <tr>
        <td><tt>{{ .Owner }}</tt></td>
        <td>
          <tt><strong>{{ .SSHCommand }}</strong></tt>
        </td>
        <td><tt>{{ .HostName }}</tt></td>
        <td><tt>{{ .IPAddress }}</tt></td>
      </tr>
      {{ end }}
    </table>
  </body>
</html>
`
