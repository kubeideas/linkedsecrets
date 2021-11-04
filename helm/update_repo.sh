helm package ./linkedsecrets -d ../docs/charts

helm repo index ../docs --url https://kubeideas.github.io/linkedsecrets --merge ../docs/index.yaml
