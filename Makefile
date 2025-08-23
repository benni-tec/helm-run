.PHONY: all clean

files := bin/helm-run_windows_amd64.exe bin/helm-run_windows_arm64.exe bin/helm-run_linux_amd64 bin/helm-run_linux_arm64 bin/helm-run_darwin_amd64 bin/helm-run_darwin_arm64

all: $(files)

bin/helm-run_windows_amd64.exe:
	GOOS=windows GOARCH=amd64 go build -o bin/helm-run_windows_amd64.exe

bin/helm-run_windows_arm64.exe:
	GOOS=windows GOARCH=arm64 go build -o bin/helm-run_windows_arm64.exe

bin/helm-run_linux_amd64:
	GOOS=linux GOARCH=amd64 go build -o bin/helm-run_linux_amd64

bin/helm-run_linux_arm64:
	GOOS=linux GOARCH=arm64 go build -o bin/helm-run_linux_arm64

bin/helm-run_darwin_amd64:
	GOOS=darwin GOARCH=amd64 go build -o bin/helm-run_darwin_amd64

bin/helm-run_darwin_arm64:
	GOOS=darwin GOARCH=arm64 go build -o bin/helm-run_darwin_arm64

clean:
	go clean || true
	rm -r bin

