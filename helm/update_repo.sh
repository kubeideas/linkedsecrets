helm package ./linkedsecrets -d ../docs/charts

helm repo index ../docs --url https://kubeideas.github.io/linkedsecrets/charts --merge ../docs/index.yaml
