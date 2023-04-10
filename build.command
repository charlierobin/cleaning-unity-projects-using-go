env GOOS=darwin GOARCH=arm64 go build -o ./build/arm64/clean-unity-project ./clean-unity-project/main.go
env GOOS=darwin GOARCH=amd64 go build -o ./build/x86_64/clean-unity-project ./clean-unity-project/main.go

env GOOS=darwin GOARCH=arm64 go build -o ./build/arm64/find-unity-project-folders ./find-unity-project-folders/main.go
env GOOS=darwin GOARCH=amd64 go build -o ./build/x86_64/find-unity-project-folders ./find-unity-project-folders/main.go

env GOOS=darwin GOARCH=arm64 go build -o ./build/arm64/get-unity-project ./get-unity-project/main.go
env GOOS=darwin GOARCH=amd64 go build -o ./build/x86_64/get-unity-project ./get-unity-project/main.go

env GOOS=darwin GOARCH=arm64 go build -o ./build/arm64/list-mounted-volumes ./list-mounted-volumes/main.go
env GOOS=darwin GOARCH=amd64 go build -o ./build/x86_64/list-mounted-volumes ./list-mounted-volumes/main.go

env GOOS=darwin GOARCH=arm64 go build -o ./build/arm64/trash ./trash/main.go
env GOOS=darwin GOARCH=amd64 go build -o ./build/x86_64/trash ./trash/main.go

env GOOS=darwin GOARCH=arm64 go build -o ./build/arm64/unity-projects-cleaner ./unity-projects-cleaner/main.go
env GOOS=darwin GOARCH=amd64 go build -o ./build/x86_64/unity-projects-cleaner ./unity-projects-cleaner/main.go

mv -v ./build/x86_64/* ./build

