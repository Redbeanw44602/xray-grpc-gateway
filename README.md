## Build
```
git clone https://github.com/XTLS/Xray-core.git
cd Xray-core
git checkout v26.3.27  # Don't forget to change this.
cd ..
./protoc-gen.sh
go build
```

## Usage
```
./xray-grpc-gateway --help
```
