resource "helm_release" "kong" {
  name  = "kong"

  repository = "https://charts.konghq.com"
  chart = "ingress"
  namespace = "kong"
  create_namespace = true
  version = "0.17.0"

  values = [file("values/argocd.yaml")]
}
