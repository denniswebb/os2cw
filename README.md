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

## Configuration

Add configuration documentation

## Authors

* **Dennis Webb** - [dhwebb](https://github.com/dhwebb)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## Acknowledgments

* Thanks to the developers of [gopsutil](https://github.com/shirou/gopsutil) who did all the legwork for me.
