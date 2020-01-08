docker build -t tmp .
docker stop kexibq
docker rm kexibq
docker rmi kexibq
docker tag tmp kexibq
docker rmi tmp
docker run -d -p 5000:5000 --name kexibq -t kexibq
