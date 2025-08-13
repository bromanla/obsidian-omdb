---
poster: {{.Poster}}
imdb_score: {{.ImdbRating}}
imdb_url: https://www.imdb.com/title/{{.ImdbID}}
length: {{.Runtime}}
genre:
{{- range $g := split .Genre "," }}
  - "[[{{ trim $g }}]]"
{{- end }}
year: {{.Year}}
cast:
{{- range $c := split .Actors "," }}
  - "[[{{ trim $c }}]]"
{{- end }}
director:
{{- range $d := split .Director "," }}
  - "[[{{ trim $d }}]]"
{{- end }}
created: {{ now }}
watched:
score:
---

![]({{.Poster}})

## Description

{{.Plot}}

## Comment

...

