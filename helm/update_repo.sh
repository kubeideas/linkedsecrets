helm package ${MY_PROJECT_DIR}/helm/linkedsecrets -d ${MY_PROJECT_DIR}/docs/charts

helm repo index ${MY_PROJECT_DIR}/docs --url https://kubeideas.github.io/linkedsecrets --merge ${MY_PROJECT_DIR}/docs/index.yaml
