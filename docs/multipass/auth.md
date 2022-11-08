Hey All,

We recently updated Multipass to version 1.9 and with it, there is [now a way](https://discourse.ubuntu.com/t/authenticating-clients-with-the-multipass-service/26183)  to [authorize](https://discourse.ubuntu.com/t/multipass-authenticate-command/26500) clients to connect to the Multipass service for better security and allow users not part of an admin group to use the `multipass` client.

We tried to put in logic to make it seamless for previous installs that upgrade to 1.9 without needing any intervention, but this seems to not be the case as evidence by [this](https://github.com/canonical/multipass/issues/2552), [this](https://github.com/canonical/multipass/issues/2554), and [this](https://github.com/canonical/multipass/issues/2549).

If you stumble here trying to find a solution, first, apologies for your troubles, and second, this post should hopefully help you fix the issue and be able to continue to use Multipass.  I've been looking at the code over and over and have yet to figure out what is triggering this issue...

### How to recover on Linux
```plain
$ sudo snap stop multipass
$ sudo killall multipass.gui
$ sudo rm /var/snap/multipass/common/data/multipassd/authenticated-certs/multipass_client_certs.pem
$ sudo cp ~/snap/multipass/current/data/multipass-client-certificate/multipass_cert.pem /var/snap/multipass/common/data/multipassd/authenticated-certs/multipass_client_certs.pem
$ sudo snap start multipass
```
### How to recover on macOS
```plain
$ sudo launchctl unload /Library/LaunchDaemons/com.canonical.multipassd.plist
$ sudo killall multipass.gui
$ sudo killall Multipass
$ sudo rm /var/root/Library/Application\ Support/multipassd/authenticated-certs/multipass_client_certs.pem
$ sudo cp ~/Library/Application\ Support/multipass-client-certificate/multipass_cert.pem /var/root/Library/Application\ Support/multipassd/authenticated-certs/multipass_client_certs.pem
$ sudo launchctl load /Library/LaunchDaemons/com.canonical.multipassd.plist
```
After the above steps for either platform, the `multipass` client running under your current user *should* connect without and authorization errors.

Again, apologies for the troubles and hopefully Multipass is working for you again!  Thanks for using Multipass!$ sudo rm /var/root/Library/Application\ Support/multipassd/authenticated-certs/multipass_client_certs.pem
$ sudo cp ~/Library/Application\ Support/multipass-client-certificate/multipass_cert.pem /var/root/Library/Application\ Support/multipassd/authenticated-certs/multipass_client_certs.pem
$ sudo launchctl load /Library/LaunchDaemons/com.canonical.multipassd.plist
```
After the above steps for either platform, the `multipass` client running under your current user *should* connect without and authorization errors.

Again, apologies for the troubles and hopefully Multipass is working for you again!  Thanks for using Multipass!
