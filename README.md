# Kong with metrics

This is example directory for setting up ArgoCd with Kong and api.
The argoCd deployments for Kong and apis in different repo: [kong-with-metrics-deployment](https://github.com/KestutisKazlauskas/kong-with-metrics-deployment)

To the code:

- First you need to have local kubernetes cluster
- Second runt the terraform to create the kong and api deployments
```bash
cd terraform
terraform init
terraform apply
```
```bash
kubectl apply -f argocd-intial-app/application.yaml
```

Now you could change the deployment in the https://github.com/KestutisKazlauskas/kong-with-metrics-deployment and argo will redeploy tha application.

## TODO
- [ ] Deploy go api (add container to image repo, add deployment to argo as api deployment for kong)
- [ ] Remove Kong other services leave only evetns api
- [ ] Add api key authentication to kong
- [ ] Separet CI/CD process - use argocd image puller - to pull the newest images of the go app - when commit the newest version git repo for kong-with-metrics deployments
- [ ] Add metrics for kong api
- [ ] Add database for kong deployment
- [ ] Test everything on EKS.
- [ ] Think about using aws prometheus and grafana services?
- [ ] Add logs for the kubernetes cluster
- [ ] Add application architecture diagram
