# Circonus Logwatch service configurations

## Systemd

### Basic installation

If circonus-logwatch is installed in default location `/opt/circonus/logwatch`. Otherwise, edit the `circonus-logwatch.service` file and change the path in the `ExecStart` command.

```
# cp circonus-logwatch.service /usr/lib/systemd/system/circonus-logwatch.service
# systemctl enable circonus-logwatch
# systemctl start circonus-logwatch
# systemctl status circonus-logwatch
```
