build tests:
from root: 
```
docker build -t ubuntu . -f ./tests/ubuntu/Dockerfile
docker run -it --rm ubuntu
```