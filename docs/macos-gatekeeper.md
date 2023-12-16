# macOS Gatekeeper

macOS may prevent you from running the pre-compiled binary due to the built-in security feature called Gatekeeper.
This is because the binary isn't signed with an Apple Developer ID certificate.

**If you get an error message saying that the binary is from an unidentified developer or something similar, you can
allow it to run by doing one of the following:**

1. :material-apple-finder: **Finder:** right-click the binary and select "Open" from the context menu and confirm that
   you want to run the binary. Gatekeeper remembers your choice and allows you to run the binary in the future.
2. :material-console: **Terminal:** add the binary to the list of allowed applications by running the following command:

```console
user@host:~$ spctl --add /path/to/tmpl
```
