# easycla-api tests


This directory contains REST api tests built using [Travern](https://github.com/taverntesting/tavern). Its a RestAPI test framework based on [py.test](http://pytest.org/en/latest/). All the API tests are defined in file of format `test_*.tavern.yaml`.


## Install dependencies to run tests in local


- Install python dependencies using pip (in virtualenv).

```bash
pip install -r requirements.txt
```

## Run Tests

```bash

# EasyCLA API URL.
export API_URL="http://localhost:8080"

# Run tests
tavern-ci  test_easyclai_api.tavern.yaml -v 

```


## Run Tests in Docker

- Build docker image

```bash
docker build -t tavern-tests .
```
- Run tests in container

```bash
export API_URL="<URL>"

docker run --rm -it -e API_URL tavern-tests 

```


