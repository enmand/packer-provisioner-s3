build:
	go build -v

vendor:
	GO111MODULE=on go mod vendor -v
	mkdir vendor/bin

gox: vendor
	GO111MODULE=on go build -o vendor/bin/gox github.com/mitchellh/gox/

install: build
	mkdir -p ~/.packer.d/plugins
	install ./packer-s3-provisioner ~/.packer.d/plugins/

release: gox
	mkdir -p releases
	vendor/bin/gox -osarch="darwin/amd64 darwin/386 linux/amd64 linux/386 windows/amd64 windows/386" --output 'dist/{{.OS}}_{{.Arch}}/{{.Dir}}'
	zip -j releases/packer-s3-provisioner_darwin_386.zip    dist/darwin_386/packer-s3-provisioner
	zip -j releases/packer-s3-provisioner_darwin_amd64.zip  dist/darwin_amd64/packer-s3-provisioner
	zip -j releases/packer-s3-provisioner_linux_386.zip     dist/linux_386/packer-s3-provisioner
	zip -j releases/packer-s3-provisioner_linux_amd64.zip   dist/linux_amd64/packer-s3-provisioner
	zip -j releases/packer-s3-provisioner_windows_386.zip   dist/windows_386/packer-s3-provisioner.exe
	zip -j releases/packer-s3-provisioner_windows_amd64.zip dist/windows_amd64/packer-s3-provisioner.exe

clean:
	rm -rf dist/
	rm -fR releases/

.PHONY: build vendor gox install release clean
