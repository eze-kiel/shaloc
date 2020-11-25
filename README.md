<h1 align="center">shaloc</h1>

<p align="center"><b>SHA</b>re files <b>LOC</b>ally !</p>

`shaloc` is a LAN-scoped sharing tool. With only 2 main commands, it is designed to be intuitive, easy and fast to use.

It has some cool features: URI randomization, AES-256 encryption/decryption, archive creation...

<p align="center">
  <a href="https://github.com/eze-kiel/shaloc/releases">
    <img src="https://img.shields.io/github/v/release/eze-kiel/shaloc" alt="Releases">
  </a>
  <a href="https://github.com/eze-kiel/shaloc/actions">
    <img src="https://img.shields.io/github/workflow/status/eze-kiel/shaloc/Release%20Go%20project" alt="Build">
  </a>
</p>

- [Getting started](#getting-started)
- [Usage](#usage)
  - [Share a single file](#share-a-single-file)
  - [Share a folder](#share-a-folder)
  - [Share something a limited number of times](#share-something-a-limited-number-of-times)
  - [Share an encrypted file/folder](#share-an-encrypted-filefolder)
  - [Clean /tmp](#clean-tmp)
  - [Update shaloc](#update-shaloc)
- [Completion](#completion)
  - [Bash](#bash)
  - [Zsh](#zsh)
  - [Fish](#fish)
- [Security note](#security-note)
- [License](#license)

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

This section will cover typical use cases. If you need a more in-depth explaination about commands and flags, check [the wiki](https://github.com/eze-kiel/shaloc/wiki).

### Share a single file

The minimal command to share a single file is the following:

```
$ shaloc share -f myfile.txt
Sharing myfile.txt on http://0.0.0.0:8080/myfile.txt
```

Note that you can choose the IP and the port (respectively `-i` and `-p`). With the flag `-r`, you can randomize the URI with a given length. For example :

```
$ shaloc share -f picture.png -i 192.168.25.33 -p 1337 -r 15
Sharing picture.png http://192.168.25.33:1337/sbChTqWQqPOiFqz
```

To receive a file, you can enter:

```
$ shaloc get -u http://127.0.0.1:8080/myfile.txt
Downloaded: myfile.txt from http://127.0.0.1:8080/myfile.txt
```

Or use whatever tool you want (`wget`, `curl`, your favorite browser...).

The content will be wrote in a file called as the file name in the url, but you can change the name with the flag `-o`:

```
$ shaloc get -u http://127.0.0.1:8080/myfile.txt -o better-name.txt
Downloaded: better-name.txt from http://127.0.0.1:8080/myfile.txt
```

### Share a folder

This command is the minimal command to share a folder:

```
$ shaloc share -F /home/user/sup3r-f0ld3r
INFO[0000] Zipping /home/user/sup3r-f0ld3r into /tmp/sup3r-f0ld3r.zip... 
Sharing /tmp/sup3r-f0ld3r.zip on http://0.0.0.0:8080/sup3r-f0ld3r.zip
```

You can also specify the IP addresse to share on, as well as the port with the same flags as before (`-i` and `-p`), and randomize the URI as well with `-r`.

You can receive the zip file using the same command as for a single file.

### Share something a limited number of times

By default, the file can be downloaded an unlimited amout of times. If you want your file to be downloaded only a certain number of times, you can specify it thanks to the `-m` flag. If it is a negative value (which is the default case), your file will be available until server shutdown. Elsewhere, the value of the flag defines the number of times it can be downloaded. Here is an example:

```
$ ./shaloc share -f foobar.txt -m 2
Sharing foobar.txt on http://0.0.0.0:8080/foobar.txt
INFO[0003] Downloads remaining: 1                       
INFO[0006] Downloads remaining: 0                       
INFO[0006] Max number of downloads reached, shutting down the server.
```

It works for both `-f` and `-F` flags.

### Share an encrypted file/folder

You can easily share an encrypted file/folder :

```
$ shaloc share -F /home/user/folder --aes
Type encryption key:
INFO[0001] Zipping /home/user/folder into /tmp/folder.zip... 
Sharing /tmp/folder.zip on http://0.0.0.0:8080/folder.zip
```

To receive it, just launch:

```
$ shaloc get -u http://127.0.0.1:8080/folder.zip --aes
Downloaded: out from http://127.0.0.1:8080/folder.zip
Type decryption key:
Decrypted out in out.dec
```

`shaloc` uses AES-256 encryption. To generate the 32 bytes key, it hashes the provided password with SHA256.

If you forgot to use `--aes` to download the file, don't worry ! You can still decrypt your file using this command:

```
$ shaloc decrypt file.txt
```

### Clean /tmp

If you do not shutdown your computer often like me, the .zip created by `shaloc` while compressing folders will stay for a long time in /tmp. So there is the `clean` command that will wipe everything that ends by ".zip" in /tmp. It is super easy to use:

```
$ shaloc clean
WARN[0000] Wiped /tmp/FgdYhsOI.zip
```

### Update shaloc

The command `update` allow you to easily keep `shaloc` up to date.

* Update to the latest version:

```
$ shaloc update latest
```

* Update to a specified version (for example v1.4.1):

```
$ shaloc update v1.4.1
```

* List all the available versions:

```
$ shaloc update list
```

## Completion

Completion is supported on multiple shells.

### Bash

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

### Zsh

If shell completion is not already enabled in your environment you will need to enable it.  You can execute the following once:

```
$ echo "autoload -U compinit; compinit" >> ~/.zshrc
```

To load completions for each session, execute once:

```
$ shaloc completion zsh > "${fpath[1]}/_shaloc"
```

You will need to start a new shell for this setup to take effect.

### Fish

```
$ shaloc completion fish | source
```

To load completions for each session, execute once:

```
$ shaloc completion fish > ~/.config/fish/completions/shaloc.fish
```

## Security note

By default, nothing is encrypted in `shaloc` which make it vulnerable to eavesdropping attacks such as [MITM](https://en.wikipedia.org/wiki/Man-in-the-middle_attack). Also, anyone with the link to your file can download it. If you want to send encrypted files, **please use the flag `--aes`**. It will ask you for a passphrase that will be needed by the receiver to decrypt the file.

If you plan to share something that should not be guessed, use the `-r` flag to randomize the URI with a random string of the length you want.

## License

[MIT](https://choosealicense.com/licenses/mit/)
