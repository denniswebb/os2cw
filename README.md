[![CircleCI](https://circleci.com/gh/dhwebb/os2cw/tree/master.svg?style=svg&circle-token=4927e733b73bdffb9a70375ff5b54db416e44b60)](https://circleci.com/gh/dhwebb/os2cw/tree/master)
# os2cw
Push disk and memory metrics to cloudwatch from Linux, Windows, and Mac.

## Installing/Building

The easiest way to build the binaries is to run `make on-docker`.
This uses a docker build container to build the application and only requires that you
have docker running on your local system.  No dev environment needed.

If you have Go installed on your system, simply run `make`.

Everything is outputted into the `build/` dirctory locally.

## Deployment

Simply copy the executable to the system and run.

## Running

`os2cw send mem-free mem-used vol-free vol-used -v all`

## Configuration

Configuration is handled either via command-line arguments, the configuration file
[`os2cw.yaml`](os2cw.yaml.sample), or a combination of both with command line taking precedent.

Application will search for `os2cw.yaml` in `/etc` and the current working directory.

AWS Region by default will try to be inferred by instance metadata.
It can be overriden, in order of precedence, by:

* command-line flag `--region`
* enviroment variable `AWS_REGION`
* configuration file `region` value

Authentication to AWS is handled by the following methods, in order or precedence:

* command-line flags `--access-key` and `--secret0key`
* environment variables `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`
* configuration file `accessKey` and `secretKey` values
* IAM Role assigned to EC2 instance

The program requires write permission to CloudWatch.  Use the following policy as a guide.
```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "cloudwatch:DescribeAlarmHistory",
                "cloudwatch:DescribeAlarms",
                "cloudwatch:DescribeAlarmsForMetric",
                "cloudwatch:GetMetricStatistics",
                "cloudwatch:ListMetrics",
                "cloudwatch:PutMetricAlarm",
                "cloudwatch:PutMetricData"
            ],
            "Resource": [
                "*"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "logs:CreateLogGroup",
                "logs:CreateLogStream",
                "logs:DescribeLogGroups",
                "logs:DescribeLogStreams",
                "logs:PutLogEvents"
            ],
            "Resource": [
                "*"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "ec2:DescribeInstanceAttribute",
                "ec2:DescribeInstanceStatus",
                "ec2:DescribeInstances",
                "ec2:DescribeTags"
            ],
            "Resource": [
                "*"
            ]
        }
    ]
}
```

## Help

```
Usage:
     os2cw send [flags]

   Examples:
     os2cw send -u gb -m mb -v / -v /home mem-avail mem-used vol-free uptime

   Flags:
         --dryrun                output metrics without sending to CloudWatch
     -i, --id string             system id to store metrics
     -m, --mem-unit string       memory size unit (b, kb, mb, gb)
     -n, --namespace string      CloudWatch namespace
     -u, --vol-unit string       volume size unit (b, kb, mb, gb, tb)
     -v, --volumes stringSlice   volumes to report (examples: /,/home,C:)



   Global Flags:
         --access-key string   AWS access key id
         --config string       config file (default is os2cw.yaml)
         --region string       AWS region
         --secret-key string   AWS secret access key

   Available Metrics:
         mem-avail  MemoryFreePercentage
         mem-free   MemoryFree
         mem-used   MemoryUsed
         mem-util   MemoryUsedPercentage
         uptime     Uptime
         vol-avail  VolumeFreePercentage
         vol-free   VolumeFree
         vol-used   VolumeUsed
         vol-util   VolumeUsedPercentage
```
## Authors

* **Dennis Webb** - [dhwebb](https://github.com/dhwebb)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## Acknowledgments

* Thanks to the developers of [gopsutil](https://github.com/shirou/gopsutil) who did all the legwork for me.
