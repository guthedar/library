This is based on youtube video playlist : https://www.youtube.com/watch?v=0qS7thrcFbc&list=PLh4KH3LtJvRRcSwTZgmW60lkhVrzKU1jA&index=6

Exposing  restapi to outside world using Port forward with below cmd:
kubectl port-forward -nrestapi svc/restapi 8080


Deploying using Helm:
1. Install helm chart: 
    helm install lib helm
2. Port forward on to svc restapi: 
    kubectl port-forward svc/restapi 8080
3. Exec into mysql po and run cmds from sql.sh file(password is in secret-mysql.yaml file): 
    kubectl exec -it pod-name -n namespace sh 
    mysql -u root -p 
    run cmds from sql.sh file 
4. Open separate terminal and curl get request: 
    curl localhost:8080/apis/v1/books



Publishing onto public repo:
1. Create a folder and package helm chart as tarbal into it: 
    helm package helm -d charts/
2. Create a index.yaml for the charts folder so as to serve this as a repo: 
    helm repo index charts
3. Push the changes to git repo and expose the repo to public so they can downlaod the helm chart and run it locally
4. Add the exposed URL as the repo: 
    helm repo add guthedar https://guthedar.github.io/library/charts/
5. Install the helm chart as below: library is the name of the chart from index.yaml 
    helm install lib guthedar/library
6. Check if 2 namespaces, pods, svcs, deployments are created or not. database and restapi and 2 new namespaces

