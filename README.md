# bash-history
Tools designed to share and improve bash history across sessions and nodes

## Current Issues

- history-exporter: When run in a container it does not see commands
  executed on the host OS. Only the commands executed inside the
container are picked up. This might be related to a file system and
docker issue. Current the binary needs to be executed on the host OS.
https://github.com/iovisor/bcc/issues/2363?notification_referrer_id=MDE4Ok5vdGlmaWNhdGlvblRocmVhZDQ5ODg1MTgwNzo3ODQxOQ%3D%3D#issuecomment-1503050884
- history-exporter: The generated ebpf .o files are still being copied over from
  cilium/ebpf repo for now.
- history-search: Can run as a container but it might add too much
  complexity because of how the commands need to be executed within the
current shell and not a subprocess.

# monorepo structure
This repo was structure was based off:

https://eli.thegreenplace.net/2019/simple-go-project-layout-with-modules/


# eBPF setup
The following describes the eBPF setup used:

- libbpf is used (https://github.com/libbpf?type=source)
- cilium/ebpf (https://github.com/cilium/ebpf/tree/master/examples)

## Suggested Reading

- https://nakryiko.com/posts/bpf-portability-and-co-re/
- https://facebookmicrosites.github.io/bpf/blog/2020/02/19/bpf-portability-and-co-re.html
