# BFBaaS  
With the concept of Block-chain s a Service, BFBaas Platform provides block chain serives to clients, assits industries to set business on block chain, and excavates commercial values of block chain.  
For detailed information: https://www.bfbaas.com/home  

# Feature
BFBaaS provides one-stop block chain development service with low costs and fast deployment.  

## Fast deployment  
As a light-weighting framework, BFBaaS Platform can build block chain at high speed, and runs client services at much lower costs.  

## Safety  
Industries own main block chain, with particular organization data share, multiple node endorses and channel options, keeps data realiable and security.  

## Visible Management  
Provide visible background mangement, easier usage of block chain, help industries control main block chain.  

# Basement  
BFBaaS is based on Kubernetes (k8s) to manage containerized workloads and services.  
More information: https://kubernetes.io/  

# Main Function  
## Dynamic fabric creation with different consensus  
- solo  
- etcdraft  
## Block chain control  
- Block chain dashboard statiscs & analysis  
- Block chain browser  
## Block chain resources  
- Dynamic block chain expansion  
- Optimization  

# File directory structure  
## baas-fabricengine  
- execute fabric operations  
## baas-kubecluster  
- k8s cluster, based on flannel network, dashboard and some other plugin installed to build a simple k8s cluster  
## baas-kubeengine  
- kubeconfig/config to replace $HOME/.kube/config of k8s master, use for k8s client linking to k8s cluster  
## baas-gateway  
- unified api gateway managment, call gateway  
