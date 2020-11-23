# shaloc

**SHA**re files **LOC**ally !

(I didn't find anything better...)

## Getting started

There is multiple ways to get `shaloc` on your machine. You can either:

* Clone the repo and build it yourself

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
$ ./shaloc serve -f toto.txt 
Now serving on http://127.0.0.1:8080/myfile.txt
```

Note that you can choose the IP and the port (respectively -i and -p).

* To receive:

```
$ shaloc get -u http://127.0.0.1:8080/myfile.txt
```

Or whatever tool you want (`wget`, `curl`, your favorite browser...).

The content will be wrote in a file called `out`, but you can change the name with the flag -o.

### Share a folder

* To send:

```
$ shaloc serve -F /home/s3cr3t-f0ld3r
Now serving on http://127.0.0.1:8080/AHjdifpLMz.zip
```

In order to simplify the use, the zip file is renamed with a random string.
You can also specify the IP addresse to serve on, as well as the port with the same flags as before (-i and -p).

* To receive:

You can receive the zipped file using the same command as for a single file.

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

## License

[MIT](https://choosealicense.com/licenses/mit/)