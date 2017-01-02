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

## Help

Please refer to runtime help using `os2cw help` and `os2cw help send`

## Configuration

Configuration is handled either via command-line arguments, the configuration file
[`os2cw.yaml`](os2cw.yaml.sample), or a combination of both with command line taking precedent.

Authentication to AWS is handled using the standard methods of leveraging the IAM role or
`AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables.

## Authors

* **Dennis Webb** - [dhwebb](https://github.com/dhwebb)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## Acknowledgments

* Thanks to the developers of [gopsutil](https://github.com/shirou/gopsutil) who did all the legwork for me.
