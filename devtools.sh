go install github.com/spf13/cobra-cli@latest
go install github.com/goreleaser/goreleaser@latest
go install github.com/caarlos0/svu@latest

if [ -f setuplinks.sh ]; then
    . ./setuplinks.sh
fi
