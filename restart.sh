make delete-xapp
make docker
sudo docker save -o ~/Documents/rimedo_labs/ransim-api-demo/ransim-api.tar rimedo-labs/ransim-api-demo:v0.0.1
sudo ctr -n k8s.io i import ~/Documents/rimedo_labs/ransim-api-demo/ransim-api.tar
make install-xapp
