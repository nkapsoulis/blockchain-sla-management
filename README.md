
## Development environment

Build and run the development environment.

```bash
make -C $FPC_PATH/utils/docker build-dev
make -C $FPC_PATH/utils/docker run-dev
```

Inside the dev environment:

```bash
make -C $FPC_PATH/protos
make -C $FPC_PATH/common
make -C $FPC_PATH/cc
make -C $FPC_PATH/application
```
