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
source ./onprem-demo/scripts/variables.sh
mkdir -p $LOGFOLDER
touch $LOGFILE
bash ./onprem-demo/scripts/generate-logs.sh $LOGFILE
```

#### webcam

In another terminal - we need to setup the webcam so it puts images into a specific folder.

```bash
source ./onprem-demo/scripts/variables.sh
mkdir -p $IMAGEFOLDER
while true; do
  streamer -s 1920x1080 -c /dev/video0 -b 16 -o $IMAGEFOLDER/cam01-$(date +%s).jpeg
  sleep 1
done
```

#### devstack

In another terminal - we allow list our log folder and start devstack:

```bash
source ./onprem-demo/scripts/variables.sh
make devstack
```

#### log parser job

This job will run for a long time - it will mount and tail the log file and trigger the bacalhau streaming results http endpoint for each line it finds.

First we build the image for the job:

```bash
source ./onprem-demo/scripts/variables.sh
docker build -t $SOURCE_LOGS_DOCKER_IMAGE -f Dockerfile.onprem-source-logs .
cat ./onprem-demo/onprem-demo-job-logs.yaml | envsubst > /tmp/onprem-demo-job-logs.yaml
bacalhau create /tmp/onprem-demo-job-logs.yaml
```

#### webcam job

```bash
source ./onprem-demo/scripts/variables.sh
docker build -t $SOURCE_IMAGES_DOCKER_IMAGE -f Dockerfile.onprem-source-images .
cat ./onprem-demo/onprem-demo-job-images.yaml | envsubst > /tmp/onprem-demo-job-images.yaml
bacalhau create /tmp/onprem-demo-job-images.yaml
```

Next: have workload call docker executor's streaming http server