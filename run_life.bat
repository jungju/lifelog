git pull origin master
docker stop run_life
docker rm  run_life
docker build -t life .
docker run -it -p 8373:8373 -e JAWBONE_KEY=XrfhRjqYY2w -e JAWBONE_SECRET=fec16e5370242d9b42ced89a87c81fb3f5ae42d5 -e HOST=life.jjgo.kr -v s:/apps/life/db:/go/src/bitbucket.org/jungju/life/db -d --name run_life life