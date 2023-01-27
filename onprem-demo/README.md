## on prem demo

### local dev

This is how to simulate the on-prem demo locally.

First - let's make sure our bacalhau client CLI is built from the local code:

```bash
make build-dev
```

This will place the `bacalhau` binary in the `/usr/local/bin` folder.

#### log file

First we make a folder & log file on the host that is where the logs will stream to and then run the script that is generating logs - this will be replaced by a journatlctl command running on the WIFI access point node:

```bash
export FOLDER=/tmp/bacalhau-onprem-demo/data
export LOGFILE="$FOLDER/accesspoint.log"
mkdir -p $FOLDER
touch $LOGFILE
bash ./onprem-demo/scripts/generate-logs.sh $LOGFILE
```

#### devstack

In another terminal - we allow list our log folder and start devstack:

```bash
export FOLDER=/tmp/bacalhau-onprem-demo/data
export PREDICTABLE_API_PORT=1
export BACALHAU_LOCAL_DIRECTORY_ALLOW_LIST=$FOLDER
make devstack
```

#### log parser job

This job will run for a long time - it will mount and tail the log file and trigger the bacalhau streaming results http endpoint for each line it finds.

```bash
export BACALHAU_API_HOST=127.0.0.1
export BACALHAU_API_PORT=20000
export FOLDER=/tmp/bacalhau-onprem-demo/data
export LOGFILE="$FOLDER/accesspoint.log"
cat ./onprem-demo/job.yaml | envsubst > /tmp/onprem-job.yaml
bacalhau create /tmp/onprem-job.yaml
```