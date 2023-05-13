# Archivist

## Initial Source
Shamefully copied uretprob from:
https://github.com/cilium/ebpf/tree/master/examples/uretprobe

# Design
Attach an eBPF program to a uretprobe, getting the returned results of readline for the bash command.

# Issues / Limitations

* Command length limit to 400 chars
* BPF stack sized is limited to 512 bytes
