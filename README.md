# shaloc

**SHA**re files **LOC**ally !

(I didn't find anything better...)

## Getting started

## Usage

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