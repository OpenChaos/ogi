
### Dev Env Help


* git clone `kafka-ogi` repo

```
git clone https://github.com/OpenChaos/kafka-ogi
```


* fetch dependencies: uses go modules `GO111MODULE=on`


* prepare environment config file

```
cp env.sample env
## now replace values for all keys there with required one
```


* run tests

```
source tests/tests-env.cfg ; make test
```

---
