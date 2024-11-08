# server-shutdown
A simple Go app for Ubuntu server which runs as a service to check if there are any users logged in. If no users are logged in and a pre-determined period of time has elapsed, it will shut the server down.

It was created to save costs when using AWS EC2 as a development server and during your busy day you've forgot to stop the instance when not being used.

## Installation & Setup
### Go
Download the latest Go release from the [Go Download](https://go.dev/dl/) page and install using:
```
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.23.3.linux-amd64.tar.gz
```
Add /usr/local/go/bin to the PATH environment variable:
```
export PATH=$PATH:/usr/local/go/bin
```
Remember to add this to your .profile file if you want it to persist. 

More information can be found on the [Go Install](https://go.dev/doc/install) information page.

### Server Shutdown
1. Use the `make build` command to create the binary then `sudo make install` to copy binary and service files.

2. The server's idle time can be adjusted by setting the **SERVER_IDLE_TIME** environment variable. Enter the `sudo systemctl edit server-shutdown.service` command and add:
```
[Service]
Environment="SERVER_IDLE_TIME=xxx"
```
The default time is 3600 seconds (1 hour).

3. Enter `sudo systemctl start server-shutdown.service` to start the service, and `sudo systemctl enable server-shutdown.service` to automatically start the service after booting up.

4. `sudo systemctl status server-shutdown.service` can be used to check that everything is correct.

5. `journalctl -u server-shutdown.service` can be used to check logs.

### SSH
If desired, you can also modify the `/etc/ssh/sshd_config` file with the following to automatically end idle SSH sessions:
```
ClientAliveInterval 3600
ClientAliveCountMax 0
```
