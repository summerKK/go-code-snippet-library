go build -ldflags \
 "-X main.buildTime=`date -u +"%Y-%m-%d,%H:%M:%S"` -X main.buildVersion=1.0.0 -X main.gitCommitId=`git rev-parse Head`"