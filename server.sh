env GOOS=linux GOARCH=amd64 go build
mv identity-network /Users/connor/Downloads
cd /Users/connor/Downloads
ssh -i "gateway-node-1.pem" ubuntu@ec2-54-197-20-70.compute-1.amazonaws.com "rm identity-network"
scp -i ./gateway-node-1.pem identity-network ubuntu@ec2-54-197-20-70.compute-1.amazonaws.com:/home/ubuntu
cd /Users/connor/Locke/identity-network