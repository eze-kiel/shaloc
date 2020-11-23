# shaloc

**SHA**re files **LOC**ally !

(I didn't find anything better...)

`shaloc` is a LAN-scoped sharing tool. With only 2 main commands, it is designed to be intuitive, easy and fast to use.

## Getting started

There is multiple ways to get `shaloc` on your machine. You can either:

* Clone the repo and build it yourself:

```
$ git clone https://github.com/eze-kiel/shaloc.git
$ cd shaloc
$ go build .
```

And then move `shaloc` somewhere in the range of your PATH.

* Download the latest release depending on your architecture.

The releases are [here](https://github.com/eze-kiel/shaloc/releases).

## Usage

This section will cover typical use cases.

### Share a single file

* To send:

```
$ shaloc serve -f myfile.txt
Serving myfile.txt on http://127.0.0.1:8080/myfile.txt
```

Note that you can choose the IP and the port (respectively `-i` and `-p`). With the flag `-r`, you can randomize the URI with a given length. For example :

```
$ shaloc serve -f picture.png -i 192.168.25.33 -p 1337 -r 15
Serving picture.png http://192.168.25.33:1337/sbChTqWQqPOiFqz
```

* To receive:

```
$ shaloc get -u http://127.0.0.1:8080/myfile.txt
Downloaded: out from http://127.0.0.1:8080/toto.txt
```

Or whatever tool you want (`wget`, `curl`, your favorite browser...).

The content will be wrote in a file called `out`, but you can change the name with the flag `-o`:

```
$ shaloc get -u http://127.0.0.1:8080/myfile.txt -o better-name.txt
Downloaded: better-name.txt from http://127.0.0.1:8080/toto.txt
```

### Share a folder

* To send:

```
$ shaloc serve -F /home/sup3r-f0ld3r
Serving /home/sup3r-f0ld3r on http://127.0.0.1:8080/AHjdifpLMz.zip
```

You can also specify the IP addresse to serve on, as well as the port with the same flags as before (`-i` and `-p`), and randomize the URI as well with `-r`.

* To receive:

You can receive the zip file using the same command as for a single file.

### Share a limited number of files

By default, the file can be downloaded an unlimited amout of times. If you want you file to be downloaded only a certain number of times, you can specify it thanks to the `-m` flag. If it is a negative value (which is the default case), your file will be available until server shutdown. Elsewhere, the value of the flag defines the number of times it can be downloaded. Here is an example:

```
$ ./shaloc serve -f foobar.txt -m 2
Serving foobar.txt on http://127.0.0.1:8080/foobar.txt
INFO[0003] Downloads remaining: 1                       
INFO[0006] Downloads remaining: 0                       
INFO[0006] Max number of downloads reached, shutting down the server.
```

It works for both `-f` and `-F` flags.

### Clean /tmp

If you do not shutdown your computer often like me, the .zip created by `shaloc` while compressing folders will stay for a long time in /tmp. So there is the `clean` command that will wipe everything that ends by ".zip" in /tmp. It is super easy to use:

```
$ shaloc clean
WARN[0000] Wiped /tmp/FgdYhsOI.zip
```

## Completion

Completion is supported on multiple shells.

### Bash:

```
$ source <(shaloc completion bash)
```

To load completions for each session, execute once:

Linux:

```
$ shaloc completion bash > /etc/bash_completion.d/shaloc
```

MacOS:

```
$ shaloc completion bash > /usr/local/etc/bash_completion.d/shaloc
```

### Zsh:

If shell completion is not already enabled in your environment you will need to enable it.  You can execute the following once:

```
$ echo "autoload -U compinit; compinit" >> ~/.zshrc
```

To load completions for each session, execute once:

```
$ shaloc completion zsh > "${fpath[1]}/_shaloc"
```

You will need to start a new shell for this setup to take effect.

### Fish:

```
$ shaloc completion fish | source
```

To load completions for each session, execute once:

```
$ shaloc completion fish > ~/.config/fish/completions/shaloc.fish
```

## Security note

By design, nothing is encrypted in `shaloc` which make it vulnerable to eavesdropping attacks such as [MITM](https://en.wikipedia.org/wiki/Man-in-the-middle_attack). Also, anyone with the link to your file can download it. You should not use `shaloc` outside your private network, or with sensitive files/folders. If you plan to share something that should not be guessed, use the `-r` flag to randomize the URI with a random string of the length you want.

## License

[MIT](https://choosealicense.com/licenses/mit/)