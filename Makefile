LD_FLAGS ="-s -w"
BUILD_DIR =./build

clean: 
		rm -rf $(BUILD_DIR)
		rm -f raidstat.tar.gz
build: 
		mkdir $(BUILD_DIR)
		go build -ldflags=$(LD_FLAGS) -buildmode=plugin -o $(BUILD_DIR)/adaptec.so plugins/adaptec/main.go
		go build -ldflags=$(LD_FLAGS) -buildmode=plugin -o $(BUILD_DIR)/hp.so plugins/hp/main.go
		go build -ldflags=$(LD_FLAGS) -o $(BUILD_DIR)/raidstat main.go
install: build
	install -d /opt/raidstat
	install -m 644 $(BUILD_DIR)/adaptec.so /opt/raidstat/adaptec.so
	install -m 644 $(BUILD_DIR)/hp.so /opt/raidstat/hp.so
	install -m 755 $(BUILD_DIR)/raidstat /opt/raidstat/raidstat
tar: build
	tar cfz raidstat.tar.gz build --transform 's/build/raidstat/'

.DEFAULT_GOAL = build
